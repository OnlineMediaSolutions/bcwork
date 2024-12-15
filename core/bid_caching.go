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
	"time"
)

type BidCachingRT struct {
	Rules []BidCachingRealtimeRecord `json:"rules"`
}

var getBidCacheQuery = `SELECT * FROM bid_caching 
        WHERE (publisher, domain) IN (%s) AND active = true`

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
	Publisher       string      `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain          string      `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Country         null.String `boil:"country" json:"country,omitempty" toml:"country" yaml:"country,omitempty"`
	Device          null.String `boil:"device" json:"device,omitempty" toml:"device" yaml:"device,omitempty"`
	BidCaching      int16       `boil:"bid_caching" json:"bid_caching" toml:"bid_caching" yaml:"bid_caching"`
	CreatedAt       time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt       null.Time   `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	RuleID          string      `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	DemandPartnerID string      `boil:"demand_partner_id" json:"demand_partner_id" toml:"demand_partner_id" yaml:"demand_partner_id"`
	Browser         null.String `boil:"browser" json:"browser,omitempty" toml:"browser" yaml:"browser,omitempty"`
	Os              null.String `boil:"os" json:"os,omitempty" toml:"os" yaml:"os,omitempty"`
	PlacementType   null.String `boil:"placement_type" json:"placement_type,omitempty" toml:"placement_type" yaml:"placement_type,omitempty"`
	Active          bool        `boil:"active" json:"active" toml:"active" yaml:"active"`
}

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

func (bc *BidCachingService) GetBidCaching(ctx context.Context, ops *GetBidCachingOptions) (dto.BidCachingSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.BidCachingColumns.Publisher).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.BidCachings(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve bid caching")
	}

	res := make(dto.BidCachingSlice, 0)
	err = res.FromModel(mods)
	if err != nil {
		return nil, fmt.Errorf("error creating model in bid caching %s", err)
	}

	return res, nil
}

func LoadBidCacheByPublisherAndDomain(ctx context.Context, pubDom models.PublisherDomainSlice) (map[string][]models.BidCaching, error) {
	bidCacheMap := make(map[string][]models.BidCaching)

	var bidCache []models.BidCaching
	query := createGetBidCacheQuery(pubDom)
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &bidCache)
	if err != nil {
		return nil, err
	}

	for _, bidCache := range bidCache {
		key := bidCache.Publisher + ":" + bidCache.Domain.String
		bidCacheMap[key] = append(bidCacheMap[key], bidCache)
	}

	return bidCacheMap, err
}

func createGetBidCacheQuery(pubDom models.PublisherDomainSlice) string {
	tupleCondition := ""
	for i, mod := range pubDom {
		if i > 0 {
			tupleCondition += ","
		}
		tupleCondition += fmt.Sprintf("('%s','%s')", mod.PublisherID, mod.Domain)
	}

	return fmt.Sprintf(getBidCacheQuery, tupleCondition)
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

	err := SendBidCachingToRT(context.Background(), data)

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

func CreateBidCachingMetadata(modBC models.BidCachingSlice, finalRules []BidCachingRealtimeRecord) []BidCachingRealtimeRecord {
	if len(modBC) != 0 {
		bidCachings := make(dto.BidCachingSlice, 0)
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
		finalRules = []BidCachingRealtimeRecord{}
	}

	helpers.SortBy(finalRules, func(i, j BidCachingRealtimeRecord) bool {
		return strings.Count(i.Rule, "*") < strings.Count(j.Rule, "*")
	})

	return finalRules
}

func (b *BidCachingService) CreateBidCaching(ctx context.Context, data *dto.BidCachingUpdateRequest) error {

	bc := dto.BidCaching{
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

	err := mod.Upsert(
		ctx,
		bcdb.DB(),
		true,
		[]string{models.BidCachingColumns.RuleID},
		boil.Blacklist(models.BidCachingColumns.CreatedAt),
		boil.Infer(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert bid caching to bid_cache table: %s", err)
	}

	b.historyModule.SaveAction(ctx, nil, mod, nil)

	err = UpdateBidCachingMetaData(*data)
	if err != nil {
		return fmt.Errorf("failed to create metadata for bid caching %s", err)
	}
	return nil
}

func (b *BidCachingService) UpdateBidCaching(ctx context.Context, data *dto.BidCachingUpdateRequest) error {
	mod, err := models.BidCachings(models.BidCachingWhere.RuleID.EQ(data.RuleId)).One(ctx, bcdb.DB())

	if err != nil {
		return fmt.Errorf("failed to fetch data from bid caching table %s", err)
	}

	mod.BidCaching = data.BidCaching
	mod.Active = true

	old, err := b.prepareHistory(ctx, mod)

	if err != nil {
		return fmt.Errorf("error in creating history record in update id caching  %s", err)
	}

	_, err = mod.Update(
		ctx,
		bcdb.DB(),
		boil.Infer(),
	)

	if err != nil {
		return fmt.Errorf("failed to update bid caching table %s", err)
	}

	err = UpdateBidCachingMetaData(*data)
	if err != nil {
		return fmt.Errorf("failed to update metadata table %s", err)
	}

	b.historyModule.SaveAction(ctx, old, mod, nil)

	return nil
}

func (b *BidCachingService) DeleteBidCaching(ctx context.Context, bidCaching []string) error {

	mods, err := models.BidCachings(models.BidCachingWhere.RuleID.IN(bidCaching)).All(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed getting bid caching for soft deleting: %w", err)
	}

	if len(mods) == 0 {
		return fmt.Errorf("no bid caching records found for provided rule IDs")
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

	err = DeleteBidCachingFromRT(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete  value from metadata table for bid caching %s", err)
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
		return nil, err
	}

	return oldMod, err
}

func SendBidCachingToRT(ctx context.Context, updateRequest dto.BidCachingUpdateRequest) error {
	modBidCaching, err := BidCachingQuery(ctx)

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

	err = metadataValue.Insert(ctx, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for bid caching")
	}

	return nil
}

func DeleteBidCachingFromRT(c context.Context) error {
	modBidCaching, err := BidCachingQuery(c)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to fetch bid cachings for delete %s", err)
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

	metadataValue := CreateMetadataObjectBidCachingDelete(utils.BidCachingMetaDataKeyPrefix, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for bid caching")
	}

	return nil
}

func CreateMetadataObjectBidCachingDelete(key string, b []byte) models.MetadataQueue {
	modMeta := models.MetadataQueue{
		TransactionID: bcguid.NewFromf(time.Now()),
		Key:           key,
		Value:         b,
	}
	return modMeta
}
