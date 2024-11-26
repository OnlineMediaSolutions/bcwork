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

type BidCashingService struct {
	historyModule history.HistoryModule
}

func NewBidCashingService(historyModule history.HistoryModule) *BidCashingService {
	return &BidCashingService{
		historyModule: historyModule,
	}
}

type BidCashing struct {
	RuleId        string `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string `boil:"device" json:"device" toml:"device" yaml:"device"`
	BidCashing    int16  `boil:"bid_cashing" json:"bid_cashing,omitempty" toml:"bid_cashing" yaml:"bid_cashing,omitempty"`
	Browser       string `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
}

type BidCashingSlice []*BidCashing

type BidCashingRealtimeRecord struct {
	Rule       string `json:"rule"`
	BidCashing int16  `json:"bid_cashing"`
	RuleID     string `json:"rule_id"`
}

type GetBidCashingOptions struct {
	Filter     BidCashingFilter       `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type BidCashingFilter struct {
	Publisher filter.StringArrayFilter `json:"publisher,omitempty"`
	Domain    filter.StringArrayFilter `json:"domain,omitempty"`
	Country   filter.StringArrayFilter `json:"country,omitempty"`
	Device    filter.StringArrayFilter `json:"device,omitempty"`
}

func (bc *BidCashing) FromModel(mod *models.BidCashing) error {
	bc.RuleId = mod.RuleID
	bc.Publisher = mod.Publisher
	bc.Domain = mod.Domain
	bc.BidCashing = mod.BidCashing
	bc.RuleId = mod.RuleID

	if mod.Os.Valid {
		bc.OS = mod.Os.String
	}

	if mod.Country.Valid {
		bc.Country = mod.Country.String
	}

	if mod.Device.Valid {
		bc.Device = mod.Device.String
	}

	if mod.PlacementType.Valid {
		bc.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		bc.Browser = mod.Browser.String
	}

	return nil
}

func (cs *BidCashingSlice) FromModel(slice models.BidCashingSlice) error {
	for _, mod := range slice {
		c := BidCashing{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (bc *BidCashingService) GetBidCashing(ctx context.Context, ops *GetBidCashingOptions) (BidCashingSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.BidCashingColumns.Publisher).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.BidCashings(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve bid cashing")
	}

	res := make(BidCashingSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (filter *BidCashingFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.BidCashingColumns.Publisher))
	}

	if len(filter.Device) > 0 {
		mods = append(mods, filter.Device.AndIn(models.BidCashingColumns.Device))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.BidCashingColumns.Domain))
	}

	if len(filter.Country) > 0 {
		mods = append(mods, filter.Country.AndIn(models.BidCashingColumns.Country))
	}

	return mods
}

func (bc *BidCashingService) UpdateMetaData(ctx context.Context, data dto.BidCashingUpdateRequest) error {
	_, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to parse hash value for bidCashing: %w", err)
	}

	go func() {
		err = SendBidCashingToRT(context.Background(), data)
	}()

	if err != nil {
		return err
	}

	return nil
}

func BidCashingQuery(ctx context.Context, updateRequest dto.BidCashingUpdateRequest) (models.BidCashingSlice, error) {
	modBidCashing, err := models.BidCashings(
		models.BidCashingWhere.Domain.EQ(updateRequest.Domain),
		models.BidCashingWhere.Publisher.EQ(updateRequest.Publisher),
	).All(ctx, bcdb.DB())

	return modBidCashing, err
}

func (bc *BidCashing) GetFormula() string {
	p := bc.Publisher
	if p == "" {
		p = "*"
	}

	d := bc.Domain
	if d == "" {
		d = "*"
	}

	c := bc.Country
	if c == "" {
		c = "*"
	}

	os := bc.OS
	if os == "" {
		os = "*"
	}

	dt := bc.Device
	if dt == "" {
		dt = "*"
	}

	pt := bc.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := bc.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, c, os, dt, pt, b)

}

func (bc *BidCashing) GetRuleID() string {
	if len(bc.RuleId) > 0 {
		return bc.RuleId
	} else {
		return bcguid.NewFrom(bc.GetFormula())
	}
}

func CreateBidCashingMetadata(modBC models.BidCashingSlice, finalRules []BidCashingRealtimeRecord) []BidCashingRealtimeRecord {
	if len(modBC) != 0 {
		bidCashings := make(BidCashingSlice, 0)
		bidCashings.FromModel(modBC)

		for _, bc := range bidCashings {
			rule := BidCashingRealtimeRecord{
				Rule:       utils.GetFormulaRegex(bc.Country, bc.Domain, bc.Device, bc.PlacementType, bc.OS, bc.Browser, bc.Publisher),
				BidCashing: bc.BidCashing,
				RuleID:     bc.GetRuleID(),
			}
			finalRules = append(finalRules, rule)
		}
	}

	helpers.SortBy(finalRules, func(i, j BidCashingRealtimeRecord) bool {
		return strings.Count(i.Rule, "*") < strings.Count(j.Rule, "*")
	})

	return finalRules
}

func (bc *BidCashing) ToModel() *models.BidCashing {

	mod := models.BidCashing{
		RuleID:     bc.GetRuleID(),
		BidCashing: bc.BidCashing,
		Publisher:  bc.Publisher,
		Domain:     bc.Domain,
	}

	if bc.Country != "" {
		mod.Country = null.StringFrom(bc.Country)
	} else {
		mod.Country = null.String{}
	}

	if bc.OS != "" {
		mod.Os = null.StringFrom(bc.OS)
	} else {
		mod.Os = null.String{}
	}

	if bc.Device != "" {
		mod.Device = null.StringFrom(bc.Device)
	} else {
		mod.Device = null.String{}
	}

	if bc.PlacementType != "" {
		mod.PlacementType = null.StringFrom(bc.PlacementType)
	} else {
		mod.PlacementType = null.String{}
	}

	if bc.Browser != "" {
		mod.Browser = null.StringFrom(bc.Browser)
	}

	return &mod

}

func (b *BidCashingService) UpdateBidCashing(ctx context.Context, data *dto.BidCashingUpdateRequest) (bool, error) {
	var isInsert bool

	bc := BidCashing{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		Device:        data.Device,
		BidCashing:    data.BidCashing,
		Browser:       data.Browser,
		OS:            data.OS,
		PlacementType: data.PlacementType,
	}

	mod := bc.ToModel()

	var old any
	oldMod, err := models.BidCashings(
		models.BidCashingWhere.RuleID.EQ(mod.RuleID),
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
		[]string{models.BidCashingColumns.RuleID},
		boil.Blacklist(models.BidCashingColumns.CreatedAt),
		boil.Infer(),
	)
	if err != nil {
		return false, err
	}

	b.historyModule.SaveAction(ctx, old, mod, nil)

	return isInsert, nil
}

func SendBidCashingToRT(c context.Context, updateRequest dto.BidCashingUpdateRequest) error {
	modBidCashing, err := BidCashingQuery(c, updateRequest)

	if err != nil && err != sql.ErrNoRows {
		return eris.Wrapf(err, "Failed to fetch bid cashing for publisher %s", updateRequest.Publisher)
	}

	var finalRules []BidCashingRealtimeRecord

	finalRules = CreateBidCashingMetadata(modBidCashing, finalRules)

	finalOutput := struct {
		Rules []BidCashingRealtimeRecord `json:"rules"`
	}{Rules: finalRules}

	value, err := json.Marshal(finalOutput)
	if err != nil {
		return eris.Wrap(err, "failed to marshal bidCashingRT to JSON")
	}

	key := utils.GetMetadataObject(updateRequest)
	metadataKey := utils.CreateMetadataKey(key, utils.BidCashingMetaDataKeyPrefix)
	metadataValue := utils.CreateMetadataObject(updateRequest, metadataKey, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for bid cashing")
	}

	return nil
}
