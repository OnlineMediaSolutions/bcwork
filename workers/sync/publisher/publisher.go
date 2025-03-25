package publisher

import (
	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type LoadedPublisher struct {
	Id                 string            `json:"_id"`
	Name               string            `json:"name"`
	AccountManager     *field            `json:"accountManager"`
	MediaBuyer         *field            `json:"mediaBuyer"`
	CampaignManager    *field            `json:"campaignManager"`
	OfficeLocation     string            `json:"officeLocation"`
	PausedDate         int64             `json:"pausedDate"`
	StartDate          int64             `json:"startDate"`
	ReactivatedDate    int64             `json:"reactivatedDate"`
	BlacklistedDomains []interface{}     `json:"blacklistedDomains"`
	ProtectorDomains   []interface{}     `json:"protectorDomains"`
	Site               []string          `json:"site"`
	DomainOptions      []*domainsOptions `json:"domainsOptions"`
}

type LoadedPublisherSlice []*LoadedPublisher

type field struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type domainsOptions struct {
	Domain          string `json:"domain"`
	IntegrationType string `json:"integrationType"`
	MirrorPublisher string `json:"mirrorPublisher,omitempty"`
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

	mirroredDomainsMap := make(map[string]*string)
	for _, domainOption := range loaded.DomainOptions {
		if domainOption.MirrorPublisher != "" {
			mirroredDomainsMap[domainOption.Domain] = &domainOption.MirrorPublisher
		}
	}

	var modDomains models.PublisherDomainSlice
	for _, site := range loaded.Site {
		modDomains = append(modDomains, &models.PublisherDomain{
			Domain:            site,
			PublisherID:       mod.PublisherID,
			MirrorPublisherID: null.StringFromPtr(mirroredDomainsMap[site]),
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
