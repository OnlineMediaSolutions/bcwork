package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/net/context"
	"strings"
)

type RefreshCacheService struct {
	historyModule history.HistoryModule
}

func NewRefreshCacheService(historyModule history.HistoryModule) *RefreshCacheService {
	return &RefreshCacheService{
		historyModule: historyModule,
	}
}

type RefreshCache struct {
	RuleId        string `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string `boil:"device" json:"device" toml:"device" yaml:"device"`
	RefreshCache  int16  `boil:"refresh_cache" json:"refresh_cache,omitempty" toml:"refresh_cache" yaml:"refresh_cache,omitempty"`
	Browser       string `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
}

type RefreshCacheSlice []*RefreshCache

type RefreshCacheRealtimeRecord struct {
	Rule         string `json:"rule"`
	RefreshCache int16  `json:"refresh_cache"`
	RuleID       string `json:"rule_id"`
}

type GetRefreshCacheOptions struct {
	Filter     RefreshCacheFilter     `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type RefreshCacheFilter struct {
	Publisher filter.StringArrayFilter `json:"publisher,omitempty"`
	Domain    filter.StringArrayFilter `json:"domain,omitempty"`
	Country   filter.StringArrayFilter `json:"country,omitempty"`
	Device    filter.StringArrayFilter `json:"device,omitempty"`
}

func (rc *RefreshCache) FromModel(mod *models.RefreshCache) error {
	rc.RuleId = mod.RuleID
	rc.Publisher = mod.Publisher
	rc.Domain = mod.Domain
	rc.RefreshCache = mod.RefreshCache
	rc.RuleId = mod.RuleID

	if mod.Os.Valid {
		rc.OS = mod.Os.String
	}

	if mod.Country.Valid {
		rc.Country = mod.Country.String
	}

	if mod.Device.Valid {
		rc.Device = mod.Device.String
	}

	if mod.PlacementType.Valid {
		rc.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		rc.Browser = mod.Browser.String
	}

	return nil
}

func (cs *RefreshCacheSlice) FromModel(slice models.RefreshCacheSlice) error {
	for _, mod := range slice {
		c := RefreshCache{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (*RefreshCacheService) GetRefreshCache(ctx context.Context, ops *GetRefreshCacheOptions) (RefreshCacheSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.RefreshCacheColumns.Publisher).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.RefreshCaches(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve refresh cache")
	}

	res := make(RefreshCacheSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (filter *RefreshCacheFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.RefreshCacheColumns.Publisher))
	}

	if len(filter.Device) > 0 {
		mods = append(mods, filter.Device.AndIn(models.RefreshCacheColumns.Device))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.RefreshCacheColumns.Domain))
	}

	if len(filter.Country) > 0 {
		mods = append(mods, filter.Country.AndIn(models.RefreshCacheColumns.Country))
	}

	return mods
}

func (*RefreshCacheService) UpdateMetaData(ctx context.Context, data dto.RefreshCacheUpdateRequest) error {
	var err error

	go func() {
		err = SendRefreshCacheToRT(context.Background(), data)
	}()

	if err != nil {
		return fmt.Errorf("error in SendRefreshCacheToRT function")
	}

	return nil
}

func (lr *RefreshCache) GetFormula() string {
	p := lr.Publisher
	if p == "" {
		p = "*"
	}

	d := lr.Domain
	if d == "" {
		d = "*"
	}

	c := lr.Country
	if c == "" {
		c = "*"
	}

	os := lr.OS
	if os == "" {
		os = "*"
	}

	dt := lr.Device
	if dt == "" {
		dt = "*"
	}

	pt := lr.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := lr.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, c, os, dt, pt, b)

}

func (rc *RefreshCache) GetRuleID() string {
	if len(rc.RuleId) > 0 {
		return rc.RuleId
	} else {
		return bcguid.NewFrom(rc.GetFormula())
	}
}

func CreateRefreshCacheMetadata(modBC models.RefreshCacheSlice, finalRules []RefreshCacheRealtimeRecord) []RefreshCacheRealtimeRecord {
	if len(modBC) != 0 {
		refreshCaches := make(RefreshCacheSlice, 0)
		refreshCaches.FromModel(modBC)

		for _, lr := range refreshCaches {
			rule := RefreshCacheRealtimeRecord{
				Rule:         utils.GetFormulaRegex(lr.Country, lr.Domain, lr.Device, lr.PlacementType, lr.OS, lr.Browser, lr.Publisher),
				RefreshCache: lr.RefreshCache,
				RuleID:       lr.GetRuleID(),
			}
			finalRules = append(finalRules, rule)
		}
	}

	helpers.SortBy(finalRules, func(i, j RefreshCacheRealtimeRecord) bool {
		return strings.Count(i.Rule, "*") < strings.Count(j.Rule, "*")
	})

	return finalRules
}

func (rc *RefreshCache) ToModel() *models.RefreshCache {

	mod := models.RefreshCache{
		RuleID:       rc.GetRuleID(),
		RefreshCache: rc.RefreshCache,
		Publisher:    rc.Publisher,
		Domain:       rc.Domain,
	}

	if rc.Country != "" {
		mod.Country = null.StringFrom(rc.Country)
	} else {
		mod.Country = null.String{}
	}

	if rc.OS != "" {
		mod.Os = null.StringFrom(rc.OS)
	} else {
		mod.Os = null.String{}
	}

	if rc.Device != "" {
		mod.Device = null.StringFrom(rc.Device)
	} else {
		mod.Device = null.String{}
	}

	if rc.PlacementType != "" {
		mod.PlacementType = null.StringFrom(rc.PlacementType)
	} else {
		mod.PlacementType = null.String{}
	}

	if rc.Browser != "" {
		mod.Browser = null.StringFrom(rc.Browser)
	}

	return &mod

}

func (r *RefreshCacheService) UpdateRefreshCache(ctx context.Context, data *dto.RefreshCacheUpdateRequest) (bool, error) {
	var isInsert bool

	rc := RefreshCache{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		Device:        data.Device,
		RefreshCache:  data.RefreshCache,
		Browser:       data.Browser,
		OS:            data.OS,
		PlacementType: data.PlacementType,
	}

	mod := rc.ToModel()

	old, err, isInsert := rc.prepareHistory(ctx, mod, isInsert)

	err = mod.Upsert(
		ctx,
		bcdb.DB(),
		true,
		[]string{models.RefreshCacheColumns.RuleID},
		boil.Blacklist(models.RefreshCacheColumns.CreatedAt),
		boil.Infer(),
	)
	if err != nil {
		return false, err
	}

	r.historyModule.SaveAction(ctx, old, mod, nil)

	return isInsert, nil
}

func (rc *RefreshCache) prepareHistory(ctx context.Context, mod *models.RefreshCache, isInsert bool) (any, error, bool) {
	var old any

	oldMod, err := models.RefreshCaches(
		models.RefreshCacheWhere.RuleID.EQ(mod.RuleID),
	).One(ctx, bcdb.DB())

	if err != nil && err != sql.ErrNoRows {
		return err, nil, false
	}

	if oldMod == nil {
		isInsert = true
	} else {
		old = oldMod
	}
	return old, err, isInsert
}

func SendRefreshCacheToRT(c context.Context, updateRequest dto.RefreshCacheUpdateRequest) error {

	value, err := json.Marshal(updateRequest.RefreshCache)
	if err != nil {
		return eris.Wrap(err, "error marshaling record for refresh cache")
	}

	if updateRequest.Domain == "" {
		updateRequest.Domain = "*"
	}

	key := utils.GetMetadataObject(updateRequest)
	metadataKey := utils.CreateMetadataKey(key, utils.RefreshCacheMetaDataKeyPrefix)
	metadataValue := utils.CreateMetadataObject(updateRequest, metadataKey, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for refresh cache")
	}

	return nil
}
