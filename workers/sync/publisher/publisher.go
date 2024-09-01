package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Worker struct {
	Bucket      string
	File        string `json:"file"`
	DatabaseEnv string `json:"dbenv"`
	LogSeverity int    `json:"logsev"`
	Cron        string `json:"cron"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local_prod")
	w.LogSeverity, _ = conf.GetIntValueWithDefault("logsev", int(2))
	w.Cron = conf.GetStringValueWithDefault("cron", "0 0 * * *")
	zerolog.SetGlobalLevel(zerolog.Level(w.LogSeverity))

	err := bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "failed to initalize DB")
	}

	w.Bucket = "new-platform-data-migration"
	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	log.Info().Msg("Starting publisher automation process")
	//var twoDaysAgo = time.Now().AddDate(0, 0, -2).UnixNano()
	list, err := utils.ListS3Objects(w.Bucket, "publishers/")
	if err != nil {
		return eris.Wrapf(err, "failed to list objects")
	}

	for _, obj := range list.Contents {
		//if obj.LastModified.UnixNano() < twoDaysAgo {
		//	continue
		//}
		pubJson, err := utils.GetObjectInput(w.Bucket, *obj.Key)
		if err != nil {
			return eris.Wrapf(err, "failed to read publisher")
		}

		var loadedPubs LoadedPublisherSlice
		err = json.Unmarshal([]byte(pubJson), &loadedPubs)
		if err != nil {
			return eris.Wrapf(err, "failed to unmarshal publisher list(file=%s)", *obj.Key)
		}

		for _, loadedPub := range loadedPubs {
			modPub, modDomains := loadedPub.ToModel()
			log.Debug().Interface("pub", modPub).Interface("domain", modDomains).Msg("Updating pub and domains")
			err = modPub.Upsert(ctx, bcdb.DB(), true, []string{models.PublisherColumns.PublisherID}, boil.Infer(), boil.Infer())
			if err != nil {
				return eris.Wrapf(err, "failed to update publisher (file=%s)", *obj.Key)
			}
			for _, modDomain := range modDomains {
				err = modDomain.Insert(ctx, bcdb.DB(), boil.Infer())
				if err != nil {
					log.Debug().Msgf("Failed to insert row to publisher domain table for publisherId: '%s'", modDomain.PublisherID)
				}
			}
		}
	}
	log.Info().Msg("Finished publisher automation process")
	return nil

}

func (w *Worker) GetSleep() int {
	log.Info().Msg(fmt.Sprintf("next run in: %d seconds", bccron.Next(w.Cron)))
	if w.Cron != "" {
		return bccron.Next(w.Cron)
	}
	return 0
}

type LoadedPublisherSlice []*LoadedPublisher

type LoadedPublisher struct {
	Id             string `json:"_id"`
	Name           string `json:"name"`
	AccountManager *struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"accountManager"`
	MediaBuyer struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"mediaBuyer"`
	CampaignManager struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"campaignManager"`
	OfficeLocation     string        `json:"officeLocation"`
	PausedDate         int64         `json:"pausedDate"`
	StartDate          int64         `json:"startDate"`
	ReactivatedDate    int64         `json:"reactivatedDate"`
	BlacklistedDomains []interface{} `json:"blacklistedDomains"`
	ProtectorDomains   []interface{} `json:"protectorDomains"`
	Site               []string      `json:"site"`
}

func (loaded *LoadedPublisher) ToModel() (*models.Publisher, models.PublisherDomainSlice) {
	mod := models.Publisher{}
	mod.PublisherID = loaded.Id
	mod.Name = loaded.Name
	if loaded.AccountManager.Id != "" {
		mod.AccountManagerID = null.StringFrom(loaded.AccountManager.Id)
	}
	if loaded.MediaBuyer.Id != "" {
		mod.MediaBuyerID = null.StringFrom(loaded.MediaBuyer.Id)
	}
	if loaded.CampaignManager.Id != "" {
		mod.CampaignManagerID = null.StringFrom(loaded.CampaignManager.Id)
	}
	if loaded.OfficeLocation != "" {
		mod.OfficeLocation = null.StringFrom(loaded.OfficeLocation)
	}
	if loaded.ReactivatedDate > 0 {
		mod.ReactivateTimestamp = null.Int64From(loaded.ReactivatedDate)
	}
	if loaded.PausedDate > 0 {
		mod.PauseTimestamp = null.Int64From(loaded.PausedDate)
	}
	if loaded.StartDate > 0 {
		mod.StartTimestamp = null.Int64From(loaded.StartDate)
	}
	var modDomains models.PublisherDomainSlice
	for _, s := range loaded.Site {
		modDomains = append(modDomains, &models.PublisherDomain{
			Domain:      s,
			PublisherID: mod.PublisherID,
		})
	}

	return &mod, modDomains
}
