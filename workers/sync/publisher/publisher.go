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

const insert_dpo_rule_query = `INSERT INTO dpo_rule (rule_id, demand_partner_id, publisher, domain, country, browser, os, device_type, placement_type, factor,created_at, updated_at) VALUES `

const on_conflict_query = `ON CONFLICT (rule_id) DO UPDATE SET country = EXCLUDED.country,
	factor = EXCLUDED.factor, device_type = EXCLUDED.device_type, domain = EXCLUDED.domain, placement_type = EXCLUDED.placement_type, updated_at = EXCLUDED.updated_at`

var query = `INSERT INTO publisher_domain (domain, publisher_id, automation, gpp_target, created_at, updated_at)
VALUES ('example.com', '1231234', CURRENT_TIMESTAMP)
ON CONFLICT (domain, publisher_id)
DO UPDATE SET
automation = EXCLUDED.automation,
gpp_target = EXCLUDED.gpp_target,
updated_at = EXCLUDED.updated_at,
created_at = publisher_domain.created_at`

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
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
	list, err := utils.ListS3Objects(w.Bucket, "publishers/")
	if err != nil {
		return eris.Wrapf(err, "failed to list objects")
	}

	for _, obj := range list.Contents {
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
					return eris.Wrapf(err, "failed to update domain (file=%s)", *obj.Key)
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

	var modDomains models.PublisherDomainSlice
	for _, s := range loaded.Site {
		modDomains = append(modDomains, &models.PublisherDomain{
			Domain:      s,
			PublisherID: mod.PublisherID,
		})
	}

	return &mod, modDomains
}
