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
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"time"
)

var softDeleteRefreshCacheQuery = `UPDATE refresh_cache SET active = false WHERE rule_id IN (%s);`

var getRefreshCacheQuery = `SELECT * FROM refresh_cache 
        WHERE (publisher, domain) IN (%s) AND active = true;`

const insertMetadataQuery = "INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) VALUES "

type RefreshCache struct {
	Publisher       string    `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain          string    `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country         string    `boil:"country" json:"country,omitempty" toml:"country" yaml:"country,omitempty"`
	Device          string    `boil:"device" json:"device,omitempty" toml:"device" yaml:"device,omitempty"`
	RefreshCache    int16     `boil:"refresh_cache" json:"refresh_cache" toml:"refresh_cache" yaml:"refresh_cache"`
	CreatedAt       time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt       null.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
	RuleID          string    `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	DemandPartnerID string    `boil:"demand_partner_id" json:"demand_partner_id" toml:"demand_partner_id" yaml:"demand_partner_id"`
	Browser         string    `boil:"browser" json:"browser,omitempty" toml:"browser" yaml:"browser,omitempty"`
	Os              string    `boil:"os" json:"os,omitempty" toml:"os" yaml:"os,omitempty"`
	PlacementType   string    `boil:"placement_type" json:"placement_type,omitempty" toml:"placement_type" yaml:"placement_type,omitempty"`
	Active          bool      `boil:"active" json:"active" toml:"active" yaml:"active"`
}

type RefreshCacheService struct {
	historyModule history.HistoryModule
}

func NewRefreshCacheService(historyModule history.HistoryModule) *RefreshCacheService {
	return &RefreshCacheService{
		historyModule: historyModule,
	}
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
	Active    filter.StringArrayFilter `json:"active,omitempty"`
}

func (r *RefreshCacheService) CreateRefreshCache(ctx context.Context, data *dto.RefreshCacheUpdateRequest) error {

	rc := dto.RefreshCache{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		Device:        data.Device,
		RefreshCache:  data.RefreshCache,
		Browser:       data.Browser,
		OS:            data.OS,
		PlacementType: data.PlacementType,
		Active:        true,
	}

	mod := rc.ToModel()

	err := mod.Upsert(
		ctx,
		bcdb.DB(),
		true,
		[]string{models.RefreshCacheColumns.RuleID},
		boil.Blacklist(models.RefreshCacheColumns.CreatedAt),
		boil.Infer(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert refresh cache table %s", err)
	}

	r.historyModule.SaveAction(ctx, nil, mod, nil)
	err = SendRefreshCacheToRT(context.Background(), *data)

	if err != nil {
		return fmt.Errorf("failed to update refresh cache metadata table %s", err)
	}

	return nil
}

func (r *RefreshCacheService) prepareHistory(ctx context.Context, mod *models.RefreshCache) (any, error) {

	oldMod, err := models.RefreshCaches(
		models.RefreshCacheWhere.RuleID.EQ(mod.RuleID),
	).One(ctx, bcdb.DB())

	if err != nil && err != sql.ErrNoRows {
		return err, nil
	}

	return oldMod, err
}

func LoadRefreshCacheByPublisherAndDomain(ctx context.Context, pubDom models.PublisherDomainSlice) (map[string][]models.RefreshCache, error) {
	refreshCacheMap := make(map[string][]models.RefreshCache)

	var refreshCache []models.RefreshCache
	query := createGetRefreshCacheQuery(pubDom)
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &refreshCache)
	if err != nil {
		return nil, err
	}

	for _, refreshCache := range refreshCache {
		key := refreshCache.Publisher + ":" + refreshCache.Domain.String
		refreshCacheMap[key] = append(refreshCacheMap[key], refreshCache)
	}

	return refreshCacheMap, err
}

func createGetRefreshCacheQuery(pubDom models.PublisherDomainSlice) string {
	tupleCondition := ""
	for i, mod := range pubDom {
		if i > 0 {
			tupleCondition += ","
		}
		tupleCondition += fmt.Sprintf("('%s','%s')", mod.PublisherID, mod.Domain)
	}

	return fmt.Sprintf(getRefreshCacheQuery, tupleCondition)
}

func (*RefreshCacheService) GetRefreshCache(ctx context.Context, ops *GetRefreshCacheOptions) (dto.RefreshCacheSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.RefreshCacheColumns.Publisher).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.RefreshCaches(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve refresh cache")
	}

	res := make(dto.RefreshCacheSlice, 0)
	err = res.FromModel(mods)

	if err != nil {
		return nil, fmt.Errorf("error creating model in refresh cache %s", err)

	}

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
	if len(filter.Active) > 0 {
		mods = append(mods, filter.Active.AndIn(models.RefreshCacheColumns.Active))
	}

	return mods
}

func UpdateRefreshCacheMetaData(ctx context.Context, data *dto.RefreshCacheUpdRequest) error {
	mod, err := models.RefreshCaches(models.RefreshCacheWhere.RuleID.EQ(data.RuleId)).One(ctx, bcdb.DB())

	domainValue := handleEmptyDomainValue(mod)
	res := dto.RefreshCacheUpdateRequest{
		Publisher:    mod.Publisher,
		Domain:       domainValue,
		RefreshCache: data.RefreshCache,
	}

	err = SendRefreshCacheToRT(context.Background(), res)
	if err != nil {
		return fmt.Errorf("error in SendRefreshCacheToRT function")
	}

	return nil
}

func (b *RefreshCacheService) UpdateRefreshCache(ctx context.Context, data *dto.RefreshCacheUpdRequest) error {
	mod, err := models.RefreshCaches(models.RefreshCacheWhere.RuleID.EQ(data.RuleId)).One(ctx, bcdb.DB())

	if err != nil {
		return fmt.Errorf("error while selecting from db %s", err)
	}
	mod.RefreshCache = data.RefreshCache
	mod.Active = true

	old, err := b.prepareHistory(ctx, mod)

	if err != nil {
		return fmt.Errorf("error while prepering history %s", err)
	}

	_, err = mod.Update(
		ctx,
		bcdb.DB(),
		boil.Infer(),
	)

	if err != nil {
		return fmt.Errorf("failed to update refresh cache table %s", err)
	}

	err = UpdateRefreshCacheMetaData(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to update refresh cache  metadata table %s", err)
	}

	b.historyModule.SaveAction(ctx, old, mod, nil)

	return nil
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

func createSoftDeleteQueryRefreshCache(refreshCache []string) string {
	var wrappedStrings []string
	for _, ruleId := range refreshCache {
		wrappedStrings = append(wrappedStrings, fmt.Sprintf(`'%s'`, ruleId))
	}

	return fmt.Sprintf(
		softDeleteRefreshCacheQuery,
		helpers.JoinStrings(wrappedStrings),
	)
}

func createDeleteQueryRefreshCache(ctx context.Context, refreshCache []string) error {
	mods, err := models.RefreshCaches(models.RefreshCacheWhere.RuleID.IN(refreshCache)).All(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to fetch refresh cache records: %w", err)
	}

	if len(mods) == 0 {
		return fmt.Errorf("no value found for these keys: %s", helpers.JoinStrings(refreshCache))
	}

	return DeleteFromMetadata(ctx, mods, err)
}

func DeleteFromMetadata(ctx context.Context, mods models.RefreshCacheSlice, err error) error {
	var (
		valueStrings []string
		valueArgs    []interface{}
	)

	multiplier := 7

	for i, data := range mods {

		domainValue := handleEmptyDomainValue(data)

		rc := dto.RefreshCacheUpdateRequest{
			Publisher:    data.Publisher,
			Domain:       domainValue,
			RefreshCache: constant.RefreshCacheDeleteValue,
		}

		metadataValue := CreateMetadataObjectRefreshCache(
			rc,
			generateMetadataKey(rc),
			[]byte(strconv.Itoa(int(rc.RefreshCache))),
		)

		offset := i * multiplier
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7))

		valueArgs = append(valueArgs,
			metadataValue.TransactionID,
			metadataValue.Key,
			metadataValue.Version,
			metadataValue.Value,
			metadataValue.CommitedInstances,
			constant.CurrentTime,
			constant.CurrentTime,
		)
	}

	query := insertMetadataQuery + strings.Join(valueStrings, ", ")

	_, err = bcdb.DB().ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to delete data from metadata: %w", err)
	}

	return nil
}

func handleEmptyDomainValue(data *models.RefreshCache) string {
	var domainValue string
	if data.Domain.Valid {
		domainValue = data.Domain.String
	} else {
		domainValue = "*"
	}
	return domainValue
}

func CreateMetadataObjectRefreshCache(res dto.RefreshCacheUpdateRequest, key string, b []byte) models.MetadataQueue {
	modMeta := models.MetadataQueue{
		TransactionID: bcguid.NewFromf(res.Publisher, res.Domain, time.Now()),
		Key:           key,
		Value:         b,
	}
	return modMeta
}

func generateMetadataKey(rc dto.RefreshCacheUpdateRequest) string {
	if rc.Domain == "" {
		rc.Domain = "*"
	}
	return fmt.Sprintf("%s:%s:%s", utils.RefreshCacheMetaDataKeyPrefix, rc.Publisher, rc.Domain)
}

func (rc *RefreshCacheService) DeleteRefreshCache(ctx context.Context, refreshCache []string) error {

	mods, err := models.RefreshCaches(models.RefreshCacheWhere.RuleID.IN(refreshCache)).All(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed getting refresh cache for soft deleting: %w", err)
	}

	oldMods := make([]any, 0, len(mods))
	newMods := make([]any, 0, len(mods))

	for i := range mods {
		oldMods = append(oldMods, mods[i])
		newMods = append(newMods, nil)
	}

	softDeleteQuery := createSoftDeleteQueryRefreshCache(refreshCache)

	_, err = queries.Raw(softDeleteQuery).Exec(bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed soft deleting refresh cache: %w", err)
	}

	err = createDeleteQueryRefreshCache(ctx, refreshCache)
	if err != nil {
		return fmt.Errorf("failed to delete from metadata table %s", err)
	}

	rc.historyModule.SaveAction(ctx, oldMods, newMods, nil)

	return nil
}
