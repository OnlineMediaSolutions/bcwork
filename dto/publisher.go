package dto

import (
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
)

type Publisher struct {
	PublisherID             string         `json:"publisher_id"`
	CreatedAt               time.Time      `json:"created_at"`
	Name                    string         `json:"name"`
	AccountManagerID        string         `json:"account_manager_id"`
	AccountManagerFullName  string         `json:"account_manager_full_name"`
	MediaBuyerID            string         `json:"media_buyer_id"`
	MediaBuyerFullName      string         `json:"media_buyer_full_name"`
	CampaignManagerID       string         `json:"campaign_manager_id"`
	CampaignManagerFullName string         `json:"campaign_manager_full_name"`
	OfficeLocation          string         `json:"office_location,omitempty"`
	PauseTimestamp          int64          `json:"pause_timestamp,omitempty"`
	StartTimestamp          int64          `json:"start_timestamp,omitempty"`
	ReactivateTimestamp     int64          `json:"reactivate_timestamp,omitempty"`
	Domains                 []string       `json:"domains,omitempty"`
	IntegrationType         []string       `json:"integration_type"` // validate:"integrationType"
	MediaType               []string       `json:"media_type"`       // validate:"mediaType"
	Status                  string         `json:"status"`
	Confiant                Confiant       `json:"confiant,omitempty"`
	Pixalate                Pixalate       `json:"pixalate,omitempty"`
	BidCaching              []BidCaching   `json:"bid_caching"`
	RefreshCache            []RefreshCache `json:"refresh_cache"`
	LatestTimestamp         int64          `json:"latest_timestamp,omitempty"`
	IsDirect                bool           `json:"is_direct"`
}

func (pub *Publisher) FromModel(mod *models.Publisher, usersMap map[string]string) error {
	integrationType := []string{}
	if len(mod.IntegrationType) > 0 {
		integrationType = mod.IntegrationType
	}

	mediaType := []string{}
	if len(mod.MediaType) > 0 {
		mediaType = mod.MediaType
	}

	pub.PublisherID = mod.PublisherID
	pub.CreatedAt = mod.CreatedAt
	pub.Name = mod.Name
	pub.Status = mod.Status.String
	pub.AccountManagerID = mod.AccountManagerID.String
	pub.AccountManagerFullName = usersMap[mod.AccountManagerID.String]
	pub.MediaBuyerID = mod.MediaBuyerID.String
	pub.MediaBuyerFullName = usersMap[mod.MediaBuyerID.String]
	pub.CampaignManagerID = mod.CampaignManagerID.String
	pub.CampaignManagerFullName = usersMap[mod.CampaignManagerID.String]
	pub.OfficeLocation = mod.OfficeLocation.String
	pub.PauseTimestamp = mod.PauseTimestamp.Int64
	pub.StartTimestamp = mod.StartTimestamp.Int64
	pub.ReactivateTimestamp = mod.ReactivateTimestamp.Int64
	pub.LatestTimestamp = max(pub.StartTimestamp, pub.ReactivateTimestamp)
	pub.IntegrationType = integrationType
	pub.MediaType = mediaType
	pub.IsDirect = mod.IsDirect

	if mod.R != nil {
		if len(mod.R.PublisherDomains) > 0 {
			for _, dom := range mod.R.PublisherDomains {
				pub.Domains = append(pub.Domains, dom.Domain)
			}
		}
		if len(mod.R.Confiants) > 0 {
			pub.Confiant = Confiant{}
			err := pub.Confiant.FromModelToCOnfiantWIthoutDomains(mod.R.Confiants)
			if err != nil {
				return eris.Wrap(err, "failed to add Confiant data for publisher")
			}
		}
		if len(mod.R.Pixalates) > 0 {
			pub.Pixalate = Pixalate{}
			err := pub.Pixalate.FromModelToPixalateWIthoutDomains(mod.R.Pixalates)
			if err != nil {
				return eris.Wrap(err, "failed to add Pixalate data for publisher")
			}
		}
		pub.BidCaching = make([]BidCaching, 0)
		if len(mod.R.BidCachings) > 0 {
			pub.addBidCachingData(mod)
		}

		pub.RefreshCache = make([]RefreshCache, 0)
		if len(mod.R.RefreshCaches) > 0 {
			pub.addRefreshCacheData(mod)
		}
	}

	return nil
}

func (pub *Publisher) addRefreshCacheData(mod *models.Publisher) {
	pub.RefreshCache = []RefreshCache{}

	for _, refresh := range mod.R.RefreshCaches {
		if len(refresh.Domain.String) == 0 && refresh.Active {
			var newRefresh = RefreshCache{}
			newRefresh.Publisher = refresh.Publisher
			newRefresh.CreatedAt = refresh.CreatedAt
			newRefresh.UpdatedAt = refresh.UpdatedAt.Ptr()
			newRefresh.Domain = refresh.Domain.String
			newRefresh.Device = refresh.Device.String
			newRefresh.Country = refresh.Country.String
			newRefresh.RefreshCache = refresh.RefreshCache
			newRefresh.RuleID = refresh.RuleID
			newRefresh.Active = true
			pub.RefreshCache = append(pub.RefreshCache, newRefresh)
		}
	}
}

func (pub *Publisher) addBidCachingData(mod *models.Publisher) {
	pub.BidCaching = []BidCaching{}

	for _, bidCaching := range mod.R.BidCachings {
		if len(bidCaching.Domain.String) == 0 && bidCaching.Active {
			var newBidCache = BidCaching{}
			newBidCache.Publisher = bidCaching.Publisher
			newBidCache.CreatedAt = bidCaching.CreatedAt
			newBidCache.UpdatedAt = bidCaching.UpdatedAt.Ptr()
			newBidCache.Domain = bidCaching.Domain.String
			newBidCache.Device = bidCaching.Device.String
			newBidCache.Country = bidCaching.Country.String
			newBidCache.BidCaching = bidCaching.BidCaching
			newBidCache.RuleID = bidCaching.RuleID
			newBidCache.Active = true
			pub.BidCaching = append(pub.BidCaching, newBidCache)
		}
	}
}

type PublisherSlice []*Publisher

func (cs *PublisherSlice) FromModel(mods models.PublisherSlice, usersMap map[string]string) error {
	for _, mod := range mods {
		c := Publisher{}
		err := c.FromModel(mod, usersMap)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

const UpdatePublisherValuesStructName = "UpdatePublisherValues"

type UpdatePublisherValues struct {
	Name                *string  `json:"name"`
	AccountManagerID    *string  `json:"account_manager_id,omitempty"`
	MediaBuyerID        *string  `json:"media_buyer_id,omitempty"`
	CampaignManagerID   *string  `json:"campaign_manager_id,omitempty"`
	OfficeLocation      *string  `json:"office_location,omitempty"`
	PauseTimestamp      *int64   `json:"pause_timestamp,omitempty"`
	StartTimestamp      *int64   `json:"start_timestamp,omitempty"`
	ReactivateTimestamp *int64   `json:"reactivate_timestamp,omitempty"`
	Status              *string  `json:"status,omitempty"`
	IntegrationType     []string `json:"integration_type,omitempty"` // validate:"integrationType"
	MediaType           []string `json:"media_type,omitempty"`       // validate:"mediaType"
	IsDirect            *bool    `json:"is_direct,omitempty"`
}

type PublisherCreateValues struct {
	Name              string   `json:"name" validate:"required"`
	AccountManagerID  string   `json:"account_manager_id"`
	MediaBuyerID      string   `json:"media_buyer_id"`
	CampaignManagerID string   `json:"campaign_manager_id"`
	OfficeLocation    string   `json:"office_location"`
	Status            string   `json:"status"`
	IntegrationType   []string `json:"integration_type"` // validate:"integrationType"
	MediaType         []string `json:"media_type"`       // validate:"mediaType"
	IsDirect          bool     `json:"is_direct"`
}
