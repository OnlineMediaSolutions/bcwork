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
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/storage/db"
	s3storage "github.com/m6yf/bcwork/storage/s3_storage"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
		databaseEnvDefault    = "local_prod"
		cronDefault           = "0 0 * * *"
		bucketDefault         = "new-platform-data-migration"
		prefixDefault         = "publishers/"
		managerMapPathDefault = "/etc/oms/managers_map.json"
		logSeverityDefault    = 2
		daysBeforeDefault     = -2
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

	w.DB = db.New(bcdb.DB())

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
		modPub, modDomains, blacklist := loadedPub.ToModel(w.ManagersMap)
		log.Debug().Interface("pub", modPub).Interface("domain", modDomains).Msg("Updating pub and domains")
		err = w.DB.UpsertPublisher(ctx, modPub, blacklist)
		if err != nil {
			return eris.Wrapf(err, "Failed to upsert row [%v] in publisher table", modPub.PublisherID)
		}
		for _, modDomain := range modDomains {
			err = w.DB.InsertPublisherDomain(ctx, modDomain)
			if err != nil {
				return eris.Wrapf(
					err,
					"Failed to insert row [%v] to publisher domain table for publisherId [%v]",
					modDomain.Domain, modDomain.PublisherID,
				)
			}
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

type LoadedPublisherSlice []*LoadedPublisher

type LoadedPublisher struct {
	Id                 string        `json:"_id"`
	Name               string        `json:"name"`
	AccountManager     *field        `json:"accountManager"`
	MediaBuyer         *field        `json:"mediaBuyer"`
	CampaignManager    *field        `json:"campaignManager"`
	OfficeLocation     string        `json:"officeLocation"`
	PausedDate         int64         `json:"pausedDate"`
	StartDate          int64         `json:"startDate"`
	ReactivatedDate    int64         `json:"reactivatedDate"`
	BlacklistedDomains []interface{} `json:"blacklistedDomains"`
	ProtectorDomains   []interface{} `json:"protectorDomains"`
	Site               []string      `json:"site"`
}

type field struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (loaded *LoadedPublisher) ToModel(managersMap map[string]string) (*models.Publisher, models.PublisherDomainSlice, boil.Columns) {
	columnBlacklist := []string{
		models.PublisherColumns.CreatedAt,
		models.PublisherColumns.Status,
		models.PublisherColumns.IntegrationType,
	}

	mod := models.Publisher{
		PublisherID: loaded.Id,
		Name:        loaded.Name,
	}

	if loaded.AccountManager != nil && loaded.AccountManager.Id != "" {
		mod.AccountManagerID = null.StringFrom(getManagerID(loaded.AccountManager.Id, managersMap))
	} else {
		columnBlacklist = append(columnBlacklist, models.PublisherColumns.AccountManagerID)
	}

	if loaded.CampaignManager != nil && loaded.CampaignManager.Id != "" {
		mod.CampaignManagerID = null.StringFrom(getManagerID(loaded.CampaignManager.Id, managersMap))
	} else {
		columnBlacklist = append(columnBlacklist, models.PublisherColumns.CampaignManagerID)
	}

	if loaded.OfficeLocation != "" {
		mod.OfficeLocation = null.StringFrom(loaded.OfficeLocation)
	} else {
		columnBlacklist = append(columnBlacklist, models.PublisherColumns.OfficeLocation)
	}

	if loaded.ReactivatedDate > 0 {
		mod.ReactivateTimestamp = null.Int64From(loaded.ReactivatedDate)
	} else {
		columnBlacklist = append(columnBlacklist, models.PublisherColumns.ReactivateTimestamp)
	}

	if loaded.StartDate > 0 {
		mod.StartTimestamp = null.Int64From(loaded.StartDate)
	} else {
		columnBlacklist = append(columnBlacklist, models.PublisherColumns.StartTimestamp)
	}

	if loaded.MediaBuyer != nil && loaded.MediaBuyer.Id != "" {
		mod.MediaBuyerID = null.StringFrom(getManagerID(loaded.MediaBuyer.Id, managersMap))
	}
	if loaded.PausedDate > 0 {
		mod.PauseTimestamp = null.Int64From(loaded.PausedDate)
	}

	var modDomains models.PublisherDomainSlice
	for _, site := range loaded.Site {
		modDomains = append(modDomains, &models.PublisherDomain{
			Domain:      site,
			PublisherID: mod.PublisherID,
		})
	}

	return &mod, modDomains, boil.Blacklist(columnBlacklist...)
}

func getManagerID(id string, managersMap map[string]string) string {
	npID, ok := managersMap[id]
	if ok {
		return npID
	}
	return id
}
