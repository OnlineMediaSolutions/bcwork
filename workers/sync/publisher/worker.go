package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	adstxt "github.com/m6yf/bcwork/modules/ads_txt"
	"github.com/m6yf/bcwork/storage/db"
	s3storage "github.com/m6yf/bcwork/storage/s3_storage"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "github.com/sirupsen/logrus"
)

type Worker struct {
	File        string            `json:"file"`
	DatabaseEnv string            `json:"dbenv"`
	LogSeverity int               `json:"logsev"`
	Cron        string            `json:"cron"`
	Bucket      string            `json:"bucket"`
	Prefix      string            `json:"prefix"`
	DaysBefore  int               `json:"days_before"` // From how many days before now objects will be processed
	ManagersMap map[string]string // ids mapping from compass to NP

	S3 s3storage.S3
	DB db.PublisherSyncStorage
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	const (
		databaseEnvDefault                 = "local_prod"
		cronDefault                        = "0 0 * * *"
		bucketDefault                      = "new-platform-data-migration"
		prefixDefault                      = "publishers/"
		managerMapPathDefault              = "/etc/oms/managers_map.json"
		logSeverityDefault                 = 2
		daysBeforeDefault                  = -2
		isNeededToCreateAdsTxtLinesDefault = false
	)

	w.DatabaseEnv = conf.GetStringValueWithDefault(config.DBEnvKey, databaseEnvDefault)
	w.Cron = conf.GetStringValueWithDefault(config.CronExpressionKey, cronDefault)
	w.Bucket = conf.GetStringValueWithDefault(config.BucketKey, bucketDefault)
	w.Prefix = conf.GetStringValueWithDefault(config.PrefixKey, prefixDefault)

	logSeverity, err := conf.GetIntValueWithDefault(config.LogSeverityKey, logSeverityDefault)
	if err != nil {
		return eris.Wrapf(err, "failed to parse log severity")
	}

	w.LogSeverity = logSeverity
	zerolog.SetGlobalLevel(zerolog.Level(w.LogSeverity))

	daysBefore, err := conf.GetIntValueWithDefault(config.DaysBeforeKey, daysBeforeDefault)
	if err != nil {
		return eris.Wrapf(err, "failed to initalize DB")
	}

	if daysBefore > 0 {
		return errors.New("variable 'days_before' must be negative")
	}

	w.DaysBefore = daysBefore

	managersMapPath := conf.GetStringValueWithDefault(config.ManagersMapPathKey, managerMapPathDefault)
	file, err := os.Open(managersMapPath)
	if err != nil {
		return eris.Wrapf(err, "failed to open file with managers map")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return eris.Wrapf(err, "failed to read file with managers map")
	}

	var managersMap map[string]string
	err = json.Unmarshal(data, &managersMap)
	if err != nil {
		return eris.Wrapf(err, "failed to unmarshal managers map")
	}

	w.ManagersMap = managersMap

	err = bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "failed to initalize DB")
	}

	w.DB = db.New(
		bcdb.DB(),
		adstxt.NewAdsTxtModule(),
		conf.GetBoolValueWithDefault(config.CreateAdsTxtLineKey, isNeededToCreateAdsTxtLinesDefault),
	)

	s3Client, err := s3storage.New()
	if err != nil {
		return eris.Wrapf(err, "failed to initalize s3 client")
	}

	w.S3 = s3Client

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	log.Info().Msg("Starting publisher automation process")

	list, err := w.S3.ListS3Objects(w.Bucket, w.Prefix)
	if err != nil {
		return eris.Wrapf(err, "failed to list objects")
	}

	for _, obj := range list.Contents {
		var (
			key       string
			hasErrors bool
		)
		if obj.Key != nil {
			key = *obj.Key
		}

		if !w.isNeededToUpdate(ctx, key, obj.LastModified) {
			log.Debug().Msgf("skipping object [%v]", key)
			continue
		}

		err = w.processObject(ctx, key)
		if err != nil {
			hasErrors = true
			log.Debug().Msgf("failed to process object [%v]: %v", key, err.Error())
		}

		err = w.DB.SaveResultOfLastSync(ctx, key, hasErrors)
		if err != nil {
			log.Debug().Msgf(
				"Failed to save result of syncing [%v:%v]: %v",
				key, hasErrors, err.Error(),
			)
		}
	}

	log.Info().Msg("Finished publisher automation process")

	return nil
}

func (w *Worker) GetSleep() int {
	next := bccron.Next(w.Cron)
	log.Info().Msg(fmt.Sprintf("next run in: %v", time.Duration(next)*time.Second))
	if w.Cron != "" {
		return next
	}

	return 0
}

func (w *Worker) processObject(ctx context.Context, key string) error {
	pubJson, err := w.S3.GetObjectInput(w.Bucket, key)
	if err != nil {
		return eris.Wrapf(err, "failed to read publisher")
	}

	var loadedPubs LoadedPublisherSlice
	err = json.Unmarshal(pubJson, &loadedPubs)
	if err != nil {
		return eris.Wrapf(err, "failed to unmarshal publisher list(file=%s)", key)
	}

	for _, loadedPub := range loadedPubs {
		err := w.processPublisher(ctx, loadedPub)
		if err != nil {
			return eris.Wrapf(err, "failed to update publisher [%v]", loadedPub.Id)
		}
	}

	return nil
}

func (w *Worker) isNeededToUpdate(ctx context.Context, key string, lastModified *time.Time) bool {
	// If there were errors last time we try to sync object from bucket,
	// we need to try it again in order not to miss updates
	hadLoadingErrorLastTime := w.DB.HadLoadingErrorLastTime(ctx, key)

	// We are updating all the publishers everyday - it should be done only
	// for the ones that were updated in the last 2 days (were updated on Compass)
	var period = time.Now().AddDate(0, 0, w.DaysBefore)
	wasUpdatedInLastNDays := lastModified.After(period)

	return wasUpdatedInLastNDays || hadLoadingErrorLastTime
}

func (w *Worker) processPublisher(ctx context.Context, loadedPub *LoadedPublisher) error {
	modPub, modDomains, blacklist := loadedPub.ToModel(w.ManagersMap)
	log.Debug().Interface("pub", modPub).Interface("domain", modDomains).Msg("Updating pub and domains")

	err := w.DB.UpsertPublisherAndDomains(ctx, modPub, modDomains, blacklist)
	if err != nil {
		return err
	}

	return nil
}
