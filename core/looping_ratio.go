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

type LoopingRatioService struct {
	historyModule history.HistoryModule
}

func NewLoopingRatioService(historyModule history.HistoryModule) *LoopingRatioService {
	return &LoopingRatioService{
		historyModule: historyModule,
	}
}

type LoopingRatio struct {
	RuleId        string `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string `boil:"device" json:"device" toml:"device" yaml:"device"`
	LoopingRatio  int16  `boil:"looping_ratio" json:"looping_ratio,omitempty" toml:"looping_ratio" yaml:"looping_ratio,omitempty"`
	Browser       string `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
}

type LoopingRatioSlice []*LoopingRatio

type LoopingRatioRealtimeRecord struct {
	Rule         string `json:"rule"`
	LoopingRatio int16  `json:"looping_ratio"`
	RuleID       string `json:"rule_id"`
}

type GetLoopingRatioOptions struct {
	Filter     BidCashingFilter       `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type LoopingRatioFilter struct {
	Publisher filter.StringArrayFilter `json:"publisher,omitempty"`
	Domain    filter.StringArrayFilter `json:"domain,omitempty"`
	Country   filter.StringArrayFilter `json:"country,omitempty"`
	Device    filter.StringArrayFilter `json:"device,omitempty"`
}

func (lr *LoopingRatio) FromModel(mod *models.LoopingRatio) error {
	lr.RuleId = mod.RuleID
	lr.Publisher = mod.Publisher
	lr.Domain = mod.Domain
	lr.LoopingRatio = mod.LoopingRatio
	lr.RuleId = mod.RuleID

	if mod.Os.Valid {
		lr.OS = mod.Os.String
	}

	if mod.Country.Valid {
		lr.Country = mod.Country.String
	}

	if mod.Device.Valid {
		lr.Device = mod.Device.String
	}

	if mod.PlacementType.Valid {
		lr.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		lr.Browser = mod.Browser.String
	}

	return nil
}

func (cs *LoopingRatioSlice) FromModel(slice models.LoopingRatioSlice) error {
	for _, mod := range slice {
		c := LoopingRatio{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (lr *LoopingRatioService) GetLoopingRatio(ctx context.Context, ops *GetLoopingRatioOptions) (LoopingRatioSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.LoopingRatioColumns.Publisher).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.LoopingRatios(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve bid cashing")
	}

	res := make(LoopingRatioSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (filter *LoopingRatioFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.LoopingRatioColumns.Publisher))
	}

	if len(filter.Device) > 0 {
		mods = append(mods, filter.Device.AndIn(models.LoopingRatioColumns.Device))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.LoopingRatioColumns.Domain))
	}

	if len(filter.Country) > 0 {
		mods = append(mods, filter.Country.AndIn(models.LoopingRatioColumns.Country))
	}

	return mods
}

func (bc *LoopingRatioService) UpdateMetaData(ctx context.Context, data dto.LoopingRatioUpdateRequest) error {
	_, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to parse hash value for looping ratio: %w", err)
	}

	go func() {
		err = SendLoopingRationToRT(context.Background(), data)
	}()

	if err != nil {
		return err
	}

	return nil
}

func LoopingRatioQuery(ctx context.Context, updateRequest dto.LoopingRatioUpdateRequest) (models.LoopingRatioSlice, error) {
	modLoopingRatio, err := models.LoopingRatios(
		models.LoopingRatioWhere.Domain.EQ(updateRequest.Domain),
		models.LoopingRatioWhere.Publisher.EQ(updateRequest.Publisher),
	).All(ctx, bcdb.DB())

	return modLoopingRatio, err
}

func (lr *LoopingRatio) GetFormula() string {
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

func (lr *LoopingRatio) GetRuleID() string {
	if len(lr.RuleId) > 0 {
		return lr.RuleId
	} else {
		return bcguid.NewFrom(lr.GetFormula())
	}
}

func CreateLoopingRatioMetadata(modBC models.LoopingRatioSlice, finalRules []LoopingRatioRealtimeRecord) []LoopingRatioRealtimeRecord {
	if len(modBC) != 0 {
		loopingRatios := make(LoopingRatioSlice, 0)
		loopingRatios.FromModel(modBC)

		for _, lr := range loopingRatios {
			rule := LoopingRatioRealtimeRecord{
				Rule:         utils.GetFormulaRegex(lr.Country, lr.Domain, lr.Device, lr.PlacementType, lr.OS, lr.Browser, lr.Publisher),
				LoopingRatio: lr.LoopingRatio,
				RuleID:       lr.GetRuleID(),
			}
			finalRules = append(finalRules, rule)
		}
	}
	SortRules(finalRules)
	return finalRules
}

func (lr *LoopingRatio) ToModel() *models.LoopingRatio {

	mod := models.LoopingRatio{
		RuleID:       lr.GetRuleID(),
		LoopingRatio: lr.LoopingRatio,
		Publisher:    lr.Publisher,
		Domain:       lr.Domain,
	}

	if lr.Country != "" {
		mod.Country = null.StringFrom(lr.Country)
	} else {
		mod.Country = null.String{}
	}

	if lr.OS != "" {
		mod.Os = null.StringFrom(lr.OS)
	} else {
		mod.Os = null.String{}
	}

	if lr.Device != "" {
		mod.Device = null.StringFrom(lr.Device)
	} else {
		mod.Device = null.String{}
	}

	if lr.PlacementType != "" {
		mod.PlacementType = null.StringFrom(lr.PlacementType)
	} else {
		mod.PlacementType = null.String{}
	}

	if lr.Browser != "" {
		mod.Browser = null.StringFrom(lr.Browser)
	}

	return &mod

}

func (l *LoopingRatioService) UpdateLoopingRatio(ctx context.Context, data *dto.LoopingRatioUpdateRequest) (bool, error) {
	var isInsert bool

	lr := LoopingRatio{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		Device:        data.Device,
		LoopingRatio:  data.LoopingRatio,
		Browser:       data.Browser,
		OS:            data.OS,
		PlacementType: data.PlacementType,
	}

	mod := lr.ToModel()

	var old any
	oldMod, err := models.LoopingRatios(
		models.LoopingRatioWhere.RuleID.EQ(mod.RuleID),
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
		[]string{models.LoopingRatioColumns.RuleID},
		boil.Blacklist(models.LoopingRatioColumns.CreatedAt),
		boil.Infer(),
	)
	if err != nil {
		return false, err
	}

	l.historyModule.SaveAction(ctx, old, mod, nil)

	return isInsert, nil
}

func SendLoopingRationToRT(c context.Context, updateRequest dto.LoopingRatioUpdateRequest) error {
	modLoopingRatio, err := LoopingRatioQuery(c, updateRequest)

	if err != nil && err != sql.ErrNoRows {
		return eris.Wrapf(err, "Failed to fetch looping ratio for publisher %s", updateRequest.Publisher)
	}

	var finalRules []LoopingRatioRealtimeRecord

	finalRules = CreateLoopingRatioMetadata(modLoopingRatio, finalRules)

	finalOutput := struct {
		Rules []LoopingRatioRealtimeRecord `json:"rules"`
	}{Rules: finalRules}

	value, err := json.Marshal(finalOutput)
	if err != nil {
		return eris.Wrap(err, "failed to marshal loopingRatioRT to JSON")
	}

	key := utils.GetMetadataObject(updateRequest)
	metadataKey := utils.CreateMetadataKey(key, utils.LoopingRatioMetaDataKeyPrefix)
	metadataValue := utils.CreateMetadataObject(updateRequest, metadataKey, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for looping ratio")
	}

	return nil
}

func SortRules(lr []LoopingRatioRealtimeRecord) {
	helpers.SortBy(lr, func(i, j LoopingRatioRealtimeRecord) bool {
		return strings.Count(i.Rule, "*") < strings.Count(j.Rule, "*")
	})
}
