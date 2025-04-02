package dto

import (
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
)

type PublisherDomainUpdateRequest struct {
	PublisherID       string   `json:"publisher_id" validate:"required"`
	Domain            string   `json:"domain" validate:"required"`
	GppTarget         *float64 `json:"gpp_target"`
	IntegrationType   []string `json:"integration_type"` // validate:"integrationType"
	Automation        bool     `json:"automation"`
	MirrorPublisherID *string  `json:"mirror_publisher_id"`
	IsDirect          bool     `json:"is_direct"`
}

type PublisherDomainRequest struct {
	DemandParnerId string                `json:"demand_partner_id"`
	Data           []PublisherDomainData `json:"data"`
}

type PublisherDomainData struct {
	PubId        string `json:"pubId"`
	Domain       string `json:"domain"`
	AdsTxtStatus bool   `json:"ads_txt_status"`
}

type PublisherDomain struct {
	PublisherID       string         `json:"publisher_id"`
	PublisherName     string         `json:"publisher_name"`
	Domain            string         `json:"domain,omitempty"`
	Automation        bool           `json:"automation"`
	GppTarget         float64        `json:"gpp_target"`
	IntegrationType   []string       `json:"integration_type"`
	CreatedAt         time.Time      `json:"created_at"`
	Confiant          Confiant       `json:"confiant,omitempty"`
	Pixalate          Pixalate       `json:"pixalate,omitempty"`
	BidCaching        []BidCaching   `json:"bid_caching"`
	RefreshCache      []RefreshCache `json:"refresh_cache"`
	UpdatedAt         *time.Time     `json:"updated_at,omitempty"`
	MirrorPublisherID *string        `json:"mirror_publisher_id,omitempty"`
	IsDirect          bool           `json:"is_direct"`
	IsDirectPublisher bool           `json:"is_direct_publisher"`
}

func (pubDom *PublisherDomain) FromModel(
	mod *models.PublisherDomain,
	confiant models.Confiant,
	pixalate models.Pixalate,
	bidCache []models.BidCaching,
	refreshCache []models.RefreshCache,
) error {
	pubDom.PublisherID = mod.PublisherID
	pubDom.CreatedAt = mod.CreatedAt
	pubDom.UpdatedAt = mod.UpdatedAt.Ptr()
	pubDom.Domain = mod.Domain
	pubDom.GppTarget = mod.GPPTarget.Float64
	pubDom.Automation = mod.Automation
	if mod.R != nil && mod.R.Publisher != nil {
		pubDom.PublisherName = mod.R.Publisher.Name
		pubDom.IsDirectPublisher = mod.R.Publisher.IsDirect
	}
	if len(mod.IntegrationType) == 0 {
		pubDom.IntegrationType = []string{}
	} else {
		pubDom.IntegrationType = mod.IntegrationType
	}
	pubDom.Confiant = Confiant{}
	pubDom.Pixalate = Pixalate{}
	pubDom.RefreshCache = make([]RefreshCache, 0)
	pubDom.BidCaching = make([]BidCaching, 0)

	if len(confiant.ConfiantKey) > 0 {
		pubDom.Confiant.createConfiant(confiant)
	}
	if len(pixalate.ID) > 0 {
		pubDom.Pixalate.createPixalate(pixalate)
	}

	pubDom.addBidCaching(bidCache)
	pubDom.addRefreshCaching(refreshCache)

	pubDom.MirrorPublisherID = mod.MirrorPublisherID.Ptr()
	pubDom.IsDirect = mod.IsDirect

	return nil
}

func (pubDom *PublisherDomain) addBidCaching(cache []models.BidCaching) {
	pubDom.BidCaching = []BidCaching{}
	for _, bidCaching := range cache {
		if bidCaching.Active {
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
			pubDom.BidCaching = append(pubDom.BidCaching, newBidCache)
		}
	}
}

func (pubDom *PublisherDomain) addRefreshCaching(cache []models.RefreshCache) {
	pubDom.RefreshCache = []RefreshCache{}

	for _, refresh := range cache {
		if refresh.Active {
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
			pubDom.RefreshCache = append(pubDom.RefreshCache, newRefresh)
		}
	}
}

func (newConfiant *Confiant) createConfiant(confiant models.Confiant) {
	newConfiant.PublisherID = confiant.PublisherID
	newConfiant.CreatedAt = &confiant.CreatedAt
	newConfiant.UpdatedAt = confiant.UpdatedAt.Ptr()
	newConfiant.Domain = &confiant.Domain
	newConfiant.Rate = &confiant.Rate
	newConfiant.ConfiantKey = &confiant.ConfiantKey
}

type PublisherDomainSlice []*PublisherDomain

func (cs *PublisherDomainSlice) FromModel(
	slice models.PublisherDomainSlice,
	confiantMap map[string]models.Confiant,
	pixalateMap map[string]models.Pixalate,
	bidCacheMap map[string][]models.BidCaching,
	refreshCacheMap map[string][]models.RefreshCache,
) error {
	for _, mod := range slice {
		c := PublisherDomain{}
		key := mod.PublisherID + ":" + mod.Domain
		confiant := confiantMap[key]
		pixalate := pixalateMap[key]
		bidCache := bidCacheMap[key]
		refreshCache := refreshCacheMap[key]
		err := c.FromModel(mod, confiant, pixalate, bidCache, refreshCache)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}
