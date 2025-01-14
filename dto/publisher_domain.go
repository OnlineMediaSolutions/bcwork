package dto

import (
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
)

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
	PublisherID     string         `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain          string         `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Automation      bool           `boil:"automation" json:"automation" toml:"automation" yaml:"automation"`
	GppTarget       float64        `boil:"gpp_target" json:"gpp_target" toml:"gpp_target" yaml:"gpp_target"`
	IntegrationType []string       `boil:"integration_type" json:"integration_type" toml:"integration_type" yaml:"integration_type"`
	CreatedAt       time.Time      `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	Confiant        Confiant       `boil:"confiant" json:"confiant,omitempty" toml:"confiant" yaml:"confiant"`
	Pixalate        Pixalate       `boil:"pixalate" json:"pixalate,omitempty" toml:"pixalate" yaml:"pixalate"`
	BidCaching      []BidCaching   `boil:"bid_caching" json:"bid_caching" toml:"bid_caching" yaml:"bid_caching"`
	RefreshCache    []RefreshCache `boil:"refresh_cache" json:"refresh_cache" toml:"refresh_cache" yaml:"refresh_cache"`
	UpdatedAt       *time.Time     `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

func (pubDom *PublisherDomain) FromModel(mod *models.PublisherDomain, confiant models.Confiant, pixalate models.Pixalate, bidCache []models.BidCaching, refreshCache []models.RefreshCache) error {
	pubDom.PublisherID = mod.PublisherID
	pubDom.CreatedAt = mod.CreatedAt
	pubDom.UpdatedAt = mod.UpdatedAt.Ptr()
	pubDom.Domain = mod.Domain
	pubDom.GppTarget = mod.GPPTarget.Float64
	pubDom.Automation = mod.Automation

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
