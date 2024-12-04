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
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/net/context"
	"strings"
)

type BidCachingRT struct {
	Rules []BidCachingRealtimeRecord `json:"rules"`
}

var deleteBidCachingQuery = `UPDATE bid_caching SET active = false WHERE rule_id IN (%s);`

type BidCachingService struct {
	historyModule history.HistoryModule
}

func NewBidCachingService(historyModule history.HistoryModule) *BidCachingService {
	return &BidCachingService{
		historyModule: historyModule,
	}
}

type BidCaching struct {
	RuleId        string `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string `boil:"device" json:"device" toml:"device" yaml:"device"`
	BidCaching    int16  `boil:"bid_caching" json:"bid_caching,omitempty" toml:"bid_caching" yaml:"bid_caching,omitempty"`
	Browser       string `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
	Active        string `boil:"actvie" json:"actvie" toml:"actvie" yaml:"actvie"`
}

type BidCachingSlice []*BidCaching

type BidCachingRealtimeRecord struct {
	Rule       string `json:"rule"`
	BidCaching int16  `json:"bid_caching"`
	RuleID     string `json:"rule_id"`
}

type GetBidCachingOptions struct {
	Filter     BidCachingFilter       `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type BidCachingFilter struct {
	Publisher filter.StringArrayFilter `json:"publisher,omitempty"`
	Domain    filter.StringArrayFilter `json:"domain,omitempty"`
	Country   filter.StringArrayFilter `json:"country,omitempty"`
	Device    filter.StringArrayFilter `json:"device,omitempty"`
	Active    filter.StringArrayFilter `json:"active,omitempty"`
}

func (bc *BidCaching) FromModel(mod *models.BidCaching) error {
	bc.RuleId = mod.RuleID
	bc.Publisher = mod.Publisher
	bc.Domain = mod.Domain
	bc.BidCaching = mod.BidCaching
	bc.Active = fmt.Sprintf("%t", mod.Active)

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

func (cs *BidCachingSlice) FromModel(slice models.BidCachingSlice) error {
	for _, mod := range slice {
		c := BidCaching{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (bc *BidCachingService) GetBidCaching(ctx context.Context, ops *GetBidCachingOptions) (BidCachingSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.BidCachingColumns.Publisher).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.BidCachings(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve bid caching")
	}

	res := make(BidCachingSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (filter *BidCachingFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.BidCachingColumns.Publisher))
	}

	if len(filter.Device) > 0 {
		mods = append(mods, filter.Device.AndIn(models.BidCachingColumns.Device))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.BidCachingColumns.Domain))
	}

	if len(filter.Country) > 0 {
		mods = append(mods, filter.Country.AndIn(models.BidCachingColumns.Country))
	}

	if len(filter.Active) > 0 {
		mods = append(mods, filter.Active.AndIn(models.BidCachingColumns.Active))
	}

	return mods
}

func UpdateBidCachingMetaData(data dto.BidCachingUpdateRequest) error {
	var err error

	go func() {
		err = SendBidCachingToRT(context.Background(), data)
	}()

	if err != nil {
		return err
	}

	return nil
}

func BidCachingQuery(ctx context.Context) (models.BidCachingSlice, error) {
	modBidCaching, err := models.BidCachings(
		qm.Where(models.BidCachingColumns.Active),
	).All(ctx, bcdb.DB())

	return modBidCaching, err
}

func (bc *BidCaching) GetFormula() string {
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

func (bc *BidCaching) GetRuleID() string {
	if len(bc.RuleId) > 0 {
		return bc.RuleId
	} else {
		return bcguid.NewFrom(bc.GetFormula())
	}
}

func CreateBidCachingMetadata(modBC models.BidCachingSlice, finalRules []BidCachingRealtimeRecord) []BidCachingRealtimeRecord {
	if len(modBC) != 0 {
		bidCachings := make(BidCachingSlice, 0)
		bidCachings.FromModel(modBC)

		for _, bc := range bidCachings {
			rule := BidCachingRealtimeRecord{
				Rule:       utils.GetFormulaRegex(bc.Country, bc.Domain, bc.Device, bc.PlacementType, bc.OS, bc.Browser, bc.Publisher),
				BidCaching: bc.BidCaching,
				RuleID:     bc.GetRuleID(),
			}
			finalRules = append(finalRules, rule)
		}
	} else {
		finalRules = append(finalRules, BidCachingRealtimeRecord{})
	}

	helpers.SortBy(finalRules, func(i, j BidCachingRealtimeRecord) bool {
		return strings.Count(i.Rule, "*") < strings.Count(j.Rule, "*")
	})

	return finalRules
}

func (bc *BidCaching) ToModel() *models.BidCaching {

	mod := models.BidCaching{
		RuleID:     bc.GetRuleID(),
		BidCaching: bc.BidCaching,
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

func (b *BidCachingService) CreateBidCaching(ctx context.Context, data *dto.BidCachingUpdateRequest) error {

	bc := BidCaching{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		Device:        data.Device,
		BidCaching:    data.BidCaching,
		Browser:       data.Browser,
		OS:            data.OS,
		PlacementType: data.PlacementType,
	}

	mod := bc.ToModel()

	old, err := b.prepareHistory(ctx, mod)

	err = mod.Insert(
		ctx,
		bcdb.DB(),
		boil.Infer(),
	)

	if err != nil {
		return err
	}

	b.historyModule.SaveAction(ctx, old, mod, nil)

	return nil
}

func (b *BidCachingService) UpdateBidCaching(ctx context.Context, data *dto.BidCachingUpdateRequest) error {
	mod, err := models.BidCachings(models.BidCachingWhere.RuleID.EQ(data.RuleId)).One(ctx, bcdb.DB())

	mod.BidCaching = data.BidCaching
	old, err := b.prepareHistory(ctx, mod)

	err = mod.Upsert(
		ctx,
		bcdb.DB(),
		true,
		[]string{models.BidCachingColumns.RuleID},
		boil.Blacklist(models.BidCachingColumns.CreatedAt),
		boil.Infer(),
	)

	if err != nil {
		return err
	}

	b.historyModule.SaveAction(ctx, old, mod, nil)

	return nil
}

func (b *BidCachingService) DeleteBidCaching(ctx context.Context, bidCaching []string) error {

	mods, err := models.BidCachings(models.BidCachingWhere.RuleID.IN(bidCaching)).All(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed getting bid caching for soft deleting: %w", err)
	}

	oldMods := make([]any, 0, len(mods))
	newMods := make([]any, 0, len(mods))

	for i := range mods {
		oldMods = append(oldMods, mods[i])
		newMods = append(newMods, nil)
	}

	deleteQuery := createSoftDeleteQueryBidCaching(bidCaching)

	_, err = queries.Raw(deleteQuery).Exec(bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed soft deleting bid caching: %w", err)
	}

	b.historyModule.SaveAction(ctx, oldMods, newMods, nil)

	return nil
}

func createSoftDeleteQueryBidCaching(bidCaching []string) string {
	var wrappedStrings []string
	for _, ruleId := range bidCaching {
		wrappedStrings = append(wrappedStrings, fmt.Sprintf(`'%s'`, ruleId))
	}

	return fmt.Sprintf(
		deleteBidCachingQuery,
		strings.Join(wrappedStrings, ","),
	)
}

func (b *BidCachingService) prepareHistory(ctx context.Context, mod *models.BidCaching) (any, error) {

	oldMod, err := models.BidCachings(
		models.BidCachingWhere.RuleID.EQ(mod.RuleID),
	).One(ctx, bcdb.DB())

	if err != nil && err != sql.ErrNoRows {
		return err, nil
	}

	return oldMod, err
}

func SendBidCachingToRT(c context.Context, updateRequest dto.BidCachingUpdateRequest) error {
	modBidCaching, err := BidCachingQuery(c)

	if err != nil && err != sql.ErrNoRows {
		return eris.Wrapf(err, "Failed to fetch bid caching for publisher %s", updateRequest.Publisher)
	}

	var finalRules []BidCachingRealtimeRecord

	finalRules = CreateBidCachingMetadata(modBidCaching, finalRules)

	finalOutput := struct {
		Rules []BidCachingRealtimeRecord `json:"rules"`
	}{Rules: finalRules}

	value, err := json.Marshal(finalOutput)
	if err != nil {
		return eris.Wrap(err, "failed to marshal bidCachingRT to JSON")
	}

	metadataValue := utils.CreateMetadataObject(updateRequest, utils.BidCachingMetaDataKeyPrefix, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for bid caching")
	}

	return nil
}
