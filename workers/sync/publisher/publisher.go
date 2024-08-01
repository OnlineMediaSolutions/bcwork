package publisher

import (
	"context"
	"encoding/json"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Worker struct {
	Bucket      string
	File        string `json:"file"`
	DatabaseEnv string `json:"dbenv"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local_prod")
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
			log.Info().Interface("pub", modPub).Interface("domain", modDomains).Msg("Updating pub and domains")

			err = modPub.Upsert(ctx, bcdb.DB(), true, []string{models.PublisherColumns.PublisherID}, boil.Infer(), boil.Infer())
			if err != nil {
				return eris.Wrapf(err, "failed to update publisher (file=%s)", *obj.Key)
			}
			for _, modDomain := range modDomains {
				err = modDomain.Upsert(ctx, bcdb.DB(), true, []string{models.PublisherDomainColumns.Name, models.PublisherDomainColumns.PublisherID}, boil.Infer(), boil.Infer())
				if err != nil {
					return eris.Wrapf(err, "failed to update domain (file=%s)", *obj.Key)
				}
			}
		}
	}
	log.Info().Msg("Finished publisher automation")
	return nil

}

func (w *Worker) GetSleep() int {
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
			Name:        s,
			PublisherID: mod.PublisherID,
		})
	}

	return &mod, modDomains
}
