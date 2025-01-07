package core

import (
	"context"
	"database/sql"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DomainService struct {
	historyModule history.HistoryModule
}

func NewDomainService(historyModule history.HistoryModule) *DomainService {
	return &DomainService{
		historyModule: historyModule,
	}
}

type PublisherDomainUpdateRequest struct {
	PublisherID     string   `json:"publisher_id" validate:"required"`
	Domain          string   `json:"domain"`
	GppTarget       *float64 `json:"gpp_target,omitempty"`
	IntegrationType []string `json:"integration_type"`
	Automation      bool     `json:"automation"`
}

type GetPublisherDomainOptions struct {
	Filter     PublisherDomainFilter  `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type PublisherDomainFilter struct {
	Domain      filter.StringArrayFilter `json:"domain,omitempty"`
	PublisherID filter.StringArrayFilter `json:"publisher_id,omitempty"`
	Automation  *filter.BoolFilter       `json:"automation,omitempty"`
	GppTarget   filter.StringArrayFilter `json:"gpp_target,omitempty"`
}

func (d *DomainService) GetPublisherDomain(ctx context.Context, ops *GetPublisherDomainOptions) (PublisherDomainSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.PublisherDomainColumns.PublisherID).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.PublisherDomains(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve publisher domains values")
	}

	confiantMap, err := LoadConfiantByPublisherAndDomain(ctx, mods)
	pixalateMap, err := LoadPixalateByPublisherAndDomain(ctx, mods)
	bidCachingMap, err := LoadBidCacheByPublisherAndDomain(ctx, mods)
	refreshCacheMap, err := LoadRefreshCacheByPublisherAndDomain(ctx, mods)

	if err != nil {
		return nil, eris.Wrap(err, "Error while retreving additional Data for publisher domains values")
	}
	res := make(PublisherDomainSlice, 0)
	res.FromModel(mods, confiantMap, pixalateMap, bidCachingMap, refreshCacheMap)

	return res, nil
}

type PublisherDomainSlice []*PublisherDomain

func (cs *PublisherDomainSlice) FromModel(slice models.PublisherDomainSlice, confiantMap map[string]models.Confiant, pixalateMap map[string]models.Pixalate, bidCacheMap map[string][]models.BidCaching, refreshCacheMap map[string][]models.RefreshCache) error {

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

func (filter *PublisherDomainFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.PublisherDomainColumns.PublisherID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.PublisherDomainColumns.Domain))
	}

	if len(filter.GppTarget) > 0 {
		mods = append(mods, filter.GppTarget.AndIn(models.PublisherDomainColumns.GPPTarget))
	}

	if filter.Automation != nil {
		mods = append(mods, filter.Automation.Where(models.PublisherDomainColumns.Automation))
	}

	return mods
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
		if bidCaching.Active == true {
			var newBidCache = BidCaching{}
			newBidCache.Publisher = bidCaching.Publisher
			newBidCache.CreatedAt = bidCaching.CreatedAt
			newBidCache.UpdatedAt = bidCaching.UpdatedAt
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
		if refresh.Active == true {
			var newRefresh = RefreshCache{}
			newRefresh.Publisher = refresh.Publisher
			newRefresh.CreatedAt = refresh.CreatedAt
			newRefresh.UpdatedAt = refresh.UpdatedAt
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

func (d *DomainService) UpdatePublisherDomain(ctx context.Context, data *PublisherDomainUpdateRequest) error {
	var oldModPointer any
	mod, err := models.PublisherDomains(
		models.PublisherDomainWhere.PublisherID.EQ(data.PublisherID),
		models.PublisherDomainWhere.Domain.EQ(data.Domain),
	).One(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if mod == nil {
		mod = &models.PublisherDomain{
			Domain:          data.Domain,
			PublisherID:     data.PublisherID,
			Automation:      data.Automation,
			GPPTarget:       null.Float64FromPtr(data.GppTarget),
			IntegrationType: data.IntegrationType,
		}

		err := mod.Insert(ctx, bcdb.DB(), boil.Infer())
		if err != nil {
			return err
		}
	} else {
		oldMod := *mod
		oldModPointer = &oldMod

		mod.Automation = data.Automation
		mod.GPPTarget = null.Float64FromPtr(data.GppTarget)
		mod.IntegrationType = data.IntegrationType
		mod.UpdatedAt = null.TimeFrom(time.Now().UTC())

		_, err := mod.Update(ctx, bcdb.DB(), boil.Infer())
		if err != nil {
			return err
		}
	}

	d.historyModule.SaveAction(ctx, oldModPointer, mod, nil)

	return nil
}
