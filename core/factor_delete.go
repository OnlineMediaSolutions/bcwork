package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	sort "sort"
	"strings"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type FactorService struct {
	historyModule history.HistoryModule
}

func NewFactorService(historyModule history.HistoryModule) *FactorService {
	return &FactorService{
		historyModule: historyModule,
	}
}

type Factor struct {
	RuleId        string  `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string  `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Factor        float64 `boil:"factor" json:"factor,omitempty" toml:"factor" yaml:"factor,omitempty"`
	Browser       string  `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string  `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string  `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
	Active        bool    `boil:"active" json:"active" toml:"active" yaml:"active"`
}

type FactorSlice []*Factor

type FactorRealtimeRecord struct {
	Rule   string  `json:"rule"`
	Factor float64 `json:"factor"`
	RuleID string  `json:"rule_id"`
}

type GetFactorOptions struct {
	Filter     FactorFilter           `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type FactorFilter struct {
	Publisher filter.StringArrayFilter `json:"publisher,omitempty"`
	Domain    filter.StringArrayFilter `json:"domain,omitempty"`
	Country   filter.StringArrayFilter `json:"country,omitempty"`
	Device    filter.StringArrayFilter `json:"device,omitempty"`
}

func (factor *Factor) FromModel(mod *models.Factor) error {
	factor.RuleId = mod.RuleID
	factor.Publisher = mod.Publisher
	factor.Domain = mod.Domain
	factor.Factor = mod.Factor
	factor.RuleId = mod.RuleID
	factor.Active = mod.Active

	if mod.Os.Valid {
		factor.OS = mod.Os.String
	}

	if mod.Country.Valid {
		factor.Country = mod.Country.String
	}

	if mod.Device.Valid {
		factor.Device = mod.Device.String
	}

	if mod.PlacementType.Valid {
		factor.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		factor.Browser = mod.Browser.String
	}

	return nil
}

func (cs *FactorSlice) FromModel(slice models.FactorSlice) error {
	for _, mod := range slice {
		c := Factor{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (f *FactorService) GetFactors(ctx context.Context, ops *GetFactorOptions) (FactorSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.FactorColumns.Publisher).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.Factors(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve factors")
	}

	res := make(FactorSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (filter *FactorFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.FactorColumns.Publisher))
	}

	if len(filter.Device) > 0 {
		mods = append(mods, filter.Device.AndIn(models.FactorColumns.Device))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.FactorColumns.Domain))
	}

	if len(filter.Country) > 0 {
		mods = append(mods, filter.Country.AndIn(models.FactorColumns.Country))
	}

	return mods
}

func (f *FactorService) UpdateMetaData(ctx context.Context, data constant.FactorUpdateRequest) error {
	_, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to parse hash value for factor: %w", err)
	}

	go func() {
		err = SendFactorToRT(context.Background(), data)
	}()

	if err != nil {
		return err
	}

	return nil
}

func FactorQuery(ctx context.Context, updateRequest constant.FactorUpdateRequest) (models.FactorSlice, error) {
	modFactor, err := models.Factors(
		models.FactorWhere.Domain.EQ(updateRequest.Domain),
		models.FactorWhere.Publisher.EQ(updateRequest.Publisher),
		models.FactorWhere.Active.EQ(true),
	).All(ctx, bcdb.DB())

	return modFactor, err
}

func (factor *Factor) GetFormula() string {
	p := factor.Publisher
	if p == "" {
		p = "*"
	}

	d := factor.Domain
	if d == "" {
		d = "*"
	}

	c := factor.Country
	if c == "" {
		c = "*"
	}

	os := factor.OS
	if os == "" {
		os = "*"
	}

	dt := factor.Device
	if dt == "" {
		dt = "*"
	}

	pt := factor.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := factor.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, c, os, dt, pt, b)

}

func (factor *Factor) GetRuleID() string {
	if len(factor.RuleId) > 0 {
		return factor.RuleId
	} else {
		return bcguid.NewFrom(factor.GetFormula())
	}
}

func CreateFactorMetadata(modFactor models.FactorSlice, finalRules []FactorRealtimeRecord) []FactorRealtimeRecord {
	if len(modFactor) != 0 {
		factors := make(FactorSlice, 0)
		factors.FromModel(modFactor)

		for _, factor := range factors {
			rule := FactorRealtimeRecord{
				Rule:   utils.GetFormulaRegex(factor.Country, factor.Domain, factor.Device, factor.PlacementType, factor.OS, factor.Browser, factor.Publisher),
				Factor: factor.Factor,
				RuleID: factor.GetRuleID(),
			}
			finalRules = append(finalRules, rule)
		}
	}

	sortRules(finalRules)
	return finalRules
}

func sortRules(factors []FactorRealtimeRecord) {
	sort.Slice(factors, func(i, j int) bool {
		return strings.Count(factors[i].Rule, "*") < strings.Count(factors[j].Rule, "*")
	})
}

func SendFactorToRT(c context.Context, updateRequest constant.FactorUpdateRequest) error {
	modFactor, err := FactorQuery(c, updateRequest)

	if err != nil && err != sql.ErrNoRows {
		return eris.Wrapf(err, "Failed to fetch factors for publisher %s", updateRequest.Publisher)
	}

	var finalRules []FactorRealtimeRecord

	finalRules = CreateFactorMetadata(modFactor, finalRules)

	finalOutput := struct {
		Rules []FactorRealtimeRecord `json:"rules"`
	}{Rules: finalRules}

	value, err := json.Marshal(finalOutput)
	if err != nil {
		return eris.Wrap(err, "failed to marshal factorRT to JSON")
	}

	key := utils.GetMetadataObject(updateRequest)
	metadataKey := utils.CreateMetadataKey(key, utils.FactorMetaDataKeyPrefix)
	metadataValue := utils.CreateMetadataObject(updateRequest, metadataKey, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for factor")
	}

	return nil
}

func (factor *Factor) ToModel() *models.Factor {

	mod := models.Factor{
		RuleID:    factor.GetRuleID(),
		Factor:    factor.Factor,
		Publisher: factor.Publisher,
		Domain:    factor.Domain,
	}

	if factor.Country != "" {
		mod.Country = null.StringFrom(factor.Country)
	} else {
		mod.Country = null.String{}
	}

	if factor.OS != "" {
		mod.Os = null.StringFrom(factor.OS)
	} else {
		mod.Os = null.String{}
	}

	if factor.Device != "" {
		mod.Device = null.StringFrom(factor.Device)
	} else {
		mod.Device = null.String{}
	}

	if factor.PlacementType != "" {
		mod.PlacementType = null.StringFrom(factor.PlacementType)
	} else {
		mod.PlacementType = null.String{}
	}

	if factor.Browser != "" {
		mod.Browser = null.StringFrom(factor.Browser)
	}

	return &mod

}

func (f *FactorService) UpdateFactor(ctx context.Context, data *constant.FactorUpdateRequest) (bool, error) {
	var isInsert bool

	factor := Factor{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		Device:        data.Device,
		Factor:        data.Factor,
		Browser:       data.Browser,
		OS:            data.OS,
		PlacementType: data.PlacementType,
	}

	mod := factor.ToModel()

	var old any
	oldMod, err := models.Factors(
		models.FactorWhere.RuleID.EQ(mod.RuleID),
	).One(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if oldMod == nil {
		isInsert = true
	} else {
		old = oldMod
	}

	err = mod.Upsert(
		ctx,
		bcdb.DB(),
		true,
		[]string{models.FactorColumns.RuleID},
		boil.Blacklist(models.FactorColumns.CreatedAt),
		boil.Infer(),
	)
	if err != nil {
		return false, err
	}

	f.historyModule.SaveAction(ctx, old, mod, nil)

	return isInsert, nil
}
