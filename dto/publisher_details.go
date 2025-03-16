package dto

import (
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
)

type PublisherDetailModel struct {
	Publisher       models.Publisher       `boil:"publisher,bind"`
	PublisherDomain models.PublisherDomain `boil:"publisher_domain,bind"`
	User            UserModelCompact       `boil:"user,bind"`
}

// UserModelCompact to support possible null string in first and last names
type UserModelCompact struct {
	ID        int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	FirstName null.String `boil:"first_name" json:"first_name" toml:"first_name" yaml:"first_name"`
	LastName  null.String `boil:"last_name" json:"last_name" toml:"last_name" yaml:"last_name"`
}

type PublisherDetail struct {
	Name                   string         `json:"name"`
	PublisherID            string         `json:"publisher_id"`
	Domain                 string         `json:"domain"`
	AccountManagerID       string         `json:"account_manager_id"`
	AccountManagerFullName string         `json:"account_manager_full_name"`
	Automation             bool           `json:"automation"`
	GPPTarget              float64        `json:"gpp_target"`
	ActivityStatus         string         `json:"activity_status"`
	Confiant               Confiant       `json:"confiant,omitempty"`
	Pixalate               Pixalate       `json:"pixalate,omitempty"`
	BidCaching             []BidCaching   `json:"bid_caching"`
	RefreshCache           []RefreshCache `json:"refresh_cache" `
}

func (pd *PublisherDetail) FromModel(mod *PublisherDetailModel, activityStatus map[string]map[string]ActivityStatus,
	confiant models.Confiant,
	pixalate models.Pixalate,
	bidCache []models.BidCaching,
	refreshCache []models.RefreshCache) error {
	pd.Name = mod.Publisher.Name
	pd.PublisherID = mod.Publisher.PublisherID
	pd.Domain = mod.PublisherDomain.Domain
	pd.AccountManagerID = mod.Publisher.AccountManagerID.String
	pd.AccountManagerFullName = buildFullName(mod.User)
	pd.Automation = mod.PublisherDomain.Automation
	pd.GPPTarget = mod.PublisherDomain.GPPTarget.Float64
	pd.ActivityStatus = activityStatus[pd.Domain][pd.PublisherID].String()
	pd.Confiant = Confiant{}
	pd.Pixalate = Pixalate{}
	pd.RefreshCache = make([]RefreshCache, 0)
	pd.BidCaching = make([]BidCaching, 0)

	if len(confiant.ConfiantKey) > 0 {
		pd.Confiant.createConfiant(confiant)
	}
	if len(pixalate.ID) > 0 {
		pd.Pixalate.createPixalate(pixalate)
	}

	pd.addBidCaching(bidCache)
	pd.addRefreshCaching(refreshCache)

	return nil
}

type PublisherDetailsSlice []*PublisherDetail

func (pds *PublisherDetailsSlice) FromModel(mods []*PublisherDetailModel, activityStatus map[string]map[string]ActivityStatus,
	confiantMap map[string]models.Confiant,
	pixalateMap map[string]models.Pixalate,
	cachingMap map[string][]models.BidCaching,
	refreshCacheMap map[string][]models.RefreshCache) error {
	for _, mod := range mods {
		key := mod.Publisher.PublisherID + ":" + mod.PublisherDomain.Domain
		confiant := confiantMap[key]
		pixalate := pixalateMap[key]
		bidCache := cachingMap[key]
		refreshCache := refreshCacheMap[key]

		pd := PublisherDetail{}
		err := pd.FromModel(mod, activityStatus, confiant, pixalate, bidCache, refreshCache)
		if err != nil {
			return eris.Cause(err)
		}
		*pds = append(*pds, &pd)
	}

	return nil
}

func buildFullName(user UserModelCompact) string {
	if user.FirstName.Valid && user.LastName.Valid {
		return user.FirstName.String + " " + user.LastName.String
	}

	return ""
}

func (pubDetail *PublisherDetail) addBidCaching(cache []models.BidCaching) {
	pubDetail.BidCaching = []BidCaching{}
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
			pubDetail.BidCaching = append(pubDetail.BidCaching, newBidCache)
		}
	}
}

func (pubDetail *PublisherDetail) addRefreshCaching(cache []models.RefreshCache) {
	pubDetail.RefreshCache = []RefreshCache{}

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
			pubDetail.RefreshCache = append(pubDetail.RefreshCache, newRefresh)
		}
	}
}
