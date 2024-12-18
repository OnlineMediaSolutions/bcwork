package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TargetingService struct {
	historyModule history.HistoryModule
}

func NewTargetingService(historyModule history.HistoryModule) *TargetingService {
	return &TargetingService{
		historyModule: historyModule,
	}
}

type ExportTagsRequest struct {
	IDs     []int `json:"ids"`
	AddGDPR bool  `json:"add_gdpr"`
}

type TargetingRealtimeRecord struct {
	RuleID     string  `json:"rule_id"`
	Rule       string  `json:"rule"`
	PriceModel string  `json:"price_model"`
	Value      float64 `json:"value"`
	DailyCap   *int    `json:"daily_cap"`
}

type TargetingOptions struct {
	Filter     TargetingFilter        `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type TargetingFilter struct {
	PublisherID   filter.StringArrayFilter   `json:"publisher_id,omitempty"`
	Domain        filter.StringArrayFilter   `json:"domain,omitempty"`
	UnitSize      filter.StringArrayFilter   `json:"unit_size,omitempty"`
	PlacementType filter.StringArrayFilter   `json:"placement_type,omitempty"`
	Country       filter.String2DArrayFilter `json:"country,omitempty"`
	DeviceType    filter.String2DArrayFilter `json:"device_type,omitempty"`
	Browser       filter.String2DArrayFilter `json:"browser,omitempty"`
	Os            filter.String2DArrayFilter `json:"os,omitempty"`
	Status        filter.StringArrayFilter   `json:"status,omitempty"`
	PriceModel    filter.StringArrayFilter   `json:"price_model,omitempty"`
}

func (filter *TargetingFilter) queryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.TargetingColumns.PublisherID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.TargetingColumns.Domain))
	}

	if len(filter.UnitSize) > 0 {
		mods = append(mods, filter.UnitSize.AndIn(models.TargetingColumns.UnitSize))
	}

	if len(filter.PlacementType) > 0 {
		mods = append(mods, filter.PlacementType.AndIn(models.TargetingColumns.PlacementType))
	}

	if len(filter.Country) > 0 {
		mods = append(mods, filter.Country.AndIn(models.TargetingColumns.Country))
	}

	if len(filter.DeviceType) > 0 {
		mods = append(mods, filter.DeviceType.AndIn(models.TargetingColumns.DeviceType))
	}

	if len(filter.Browser) > 0 {
		mods = append(mods, filter.Browser.AndIn(models.TargetingColumns.Browser))
	}

	if len(filter.Os) > 0 {
		mods = append(mods, filter.Os.AndIn(models.TargetingColumns.Os))
	}

	if len(filter.Status) > 0 {
		mods = append(mods, filter.Status.AndIn(models.TargetingColumns.Status))
	}

	if len(filter.PriceModel) > 0 {
		mods = append(mods, filter.PriceModel.AndIn(models.TargetingColumns.PriceModel))
	}

	return mods
}

func (t *TargetingService) GetTargetings(ctx context.Context, ops *TargetingOptions) ([]*dto.Targeting, error) {
	qmods := ops.Filter.queryMod().
		Order(ops.Order, nil, models.TargetingColumns.ID).
		AddArray(ops.Pagination.Do())

	mods, err := models.Targetings(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve targetings")
	}

	targetings := make([]*dto.Targeting, 0, len(mods))
	for _, mod := range mods {
		targeting := new(dto.Targeting)
		err := targeting.FromModel(mod)
		if err != nil {
			return nil, eris.Wrap(err, "failed to map data from model")
		}

		targetings = append(targetings, targeting)
	}

	return targetings, nil
}

func (t *TargetingService) CreateTargeting(ctx context.Context, data *dto.Targeting) (*dto.Targeting, error) {
	data.PrepareData()

	duplicate, err := checkForDuplicate(ctx, data)
	if err != nil {
		return duplicate, eris.Wrap(err, "checking for duplicates")
	}

	mod, err := data.ToModel()
	if err != nil {
		return nil, eris.Wrap(err, "failed to map data to model")
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, eris.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	err = mod.Insert(ctx, tx, boil.Infer())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to upsert targeting")
	}

	err = updateTargetingMetaData(ctx, data, tx)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to update targeting metadata")
	}

	err = tx.Commit()
	if err != nil {
		return nil, eris.Wrapf(err, "failed to commit targeting and metadata")
	}

	t.historyModule.SaveAction(ctx, nil, mod, nil)

	return nil, nil
}

func (t *TargetingService) UpdateTargeting(ctx context.Context, data *dto.Targeting) (*dto.Targeting, error) {
	data.PrepareData()

	mod, err := models.Targetings(models.TargetingWhere.ID.EQ(data.ID)).One(ctx, bcdb.DB())
	if err != nil {
		return nil, eris.Wrap(err, fmt.Sprintf("failed to get targeting with id [%v] to update", data.ID))
	}

	oldMod := *mod

	duplicate, err := checkForDuplicate(ctx, data)
	if err != nil {
		return duplicate, eris.Wrap(err, "checking for duplicates")
	}

	columns, err := getColumnsToUpdate(data, mod)
	if err != nil {
		return nil, eris.Wrap(err, "error getting columns for update")
	}
	// if updating only updated_at, rule_id
	if len(columns) == 2 {
		return nil, errors.New("there are no new values to update targeting")
	}

	ruleID, err := dto.CalculateTargetingRuleID(mod)
	if err != nil {
		return nil, eris.Wrap(err, "error getting calculating targeting rule id")
	}

	mod.RuleID = ruleID

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, eris.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	_, err = mod.Update(ctx, tx, boil.Whitelist(columns...))
	if err != nil {
		return nil, eris.Wrap(err, "failed to update targeting")
	}

	err = updateTargetingMetaData(ctx, data, tx)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to update targeting metadata")
	}

	err = tx.Commit()
	if err != nil {
		return nil, eris.Wrapf(err, "failed to commit targeting updates and metadata")
	}

	t.historyModule.SaveAction(ctx, &oldMod, mod, nil)

	return nil, nil
}

func (t *TargetingService) ExportTags(ctx context.Context, data *ExportTagsRequest) ([]dto.Tags, error) {
	mods, err := models.Targetings(
		models.TargetingWhere.ID.IN(data.IDs),
		qm.Load(models.TargetingRels.Publisher),
	).All(ctx, bcdb.DB())
	if err != nil {
		return nil, eris.Wrap(err, fmt.Sprintf("failed to get targetings with ids %v to export tags", data.IDs))
	}

	tags := make([]dto.Tags, 0, len(mods))
	for _, mod := range mods {
		tag, err := dto.GetJSTagString(mod, data.AddGDPR)
		if err != nil {
			return nil, eris.Wrap(err, fmt.Sprintf("failed to get js tag for id [%v]", mod.ID))
		}
		tags = append(tags, dto.Tags{ID: mod.ID, Tag: tag})
	}

	return tags, nil
}

// getTargetingsByData Get targetings for publisher, domain and unit size
func getTargetingsByData(ctx context.Context, data *dto.Targeting, exec boil.ContextExecutor) (models.TargetingSlice, error) {
	mods, err := models.Targetings(
		models.TargetingWhere.PublisherID.EQ(data.PublisherID),
		models.TargetingWhere.Domain.EQ(data.Domain),
		models.TargetingWhere.Status.NEQ(dto.TargetingStatusArchived),
		qm.OrderBy(models.TargetingColumns.UnitSize),
		qm.OrderBy(models.TargetingColumns.Country),
		qm.OrderBy(models.TargetingColumns.DeviceType),
		qm.OrderBy(models.TargetingColumns.Os),
		qm.OrderBy(models.TargetingColumns.Browser),
		qm.OrderBy(models.TargetingColumns.PlacementType),
		qm.OrderBy(models.TargetingColumns.KV),
	).All(ctx, exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return mods, nil
}

func getTargetingByProps(ctx context.Context, data *dto.Targeting) (*models.Targeting, error) {
	var qmods qmods.QueryModsSlice
	qmods = qmods.Add(
		models.TargetingWhere.PublisherID.EQ(data.PublisherID),
		models.TargetingWhere.Domain.EQ(data.Domain),
		models.TargetingWhere.UnitSize.EQ(data.UnitSize),
		models.TargetingWhere.PlacementType.EQ(null.StringFrom(data.PlacementType)),
	).
		Add(getMultipleValuesFieldsWhereQueries(data.Country, data.DeviceType, data.Browser, data.OS)...).
		Add(getKeyValueWhereQueries(data.KV)...)

	mod, err := models.Targetings(qmods...).One(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return mod, nil
}

func getKeyValueWhereQueries(kv map[string]string) qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0, len(kv)+1)
	if len(kv) == 0 {
		mods = append(mods, models.TargetingWhere.KV.IsNull())
		return mods
	}

	var i int = 5 // depends on getTargetingByProps
	for key, value := range kv {
		mods = append(mods, qm.Where(
			fmt.Sprintf("%v ->> '%v' = $%v", models.TargetingColumns.KV, key, i),
			value,
		))
		i++
	}

	return mods
}

func getMultipleValuesFieldsWhereQueries(country, deviceType, browser, os []string) qmods.QueryModsSlice {
	return qmods.QueryModsSlice{
		getMultipleValuesFieldsWhereQuery(country, models.TargetingColumns.Country),
		getMultipleValuesFieldsWhereQuery(deviceType, models.TargetingColumns.DeviceType),
		getMultipleValuesFieldsWhereQuery(browser, models.TargetingColumns.Browser),
		getMultipleValuesFieldsWhereQuery(os, models.TargetingColumns.Os),
	}
}

func getMultipleValuesFieldsWhereQuery(array []string, columnName string) qm.QueryMod {
	if len(array) > 0 {
		return qm.Where(columnName + " && ARRAY['" + strings.Join(array, "','") + "']")
	}
	return qm.Where(columnName + " IS NULL")
}

func checkForDuplicate(ctx context.Context, data *dto.Targeting) (*dto.Targeting, error) {
	duplicate, err := getTargetingByProps(ctx, data)
	if err != nil {
		return nil, eris.Wrap(err, "failed to get duplicate")
	}

	if isDuplicate(duplicate, data) {
		targeting := new(dto.Targeting)
		err := targeting.FromModel(duplicate)
		if err != nil {
			return nil, eris.Wrap(err, "failed to map duplicate data from model")
		}

		return targeting, fmt.Errorf("%w: there is targeting with such parameters", dto.ErrFoundDuplicate)
	}

	return nil, nil
}

func isDuplicate(mod *models.Targeting, data *dto.Targeting) bool {
	if mod == nil {
		return false
	}

	// id == 0 if creating new targeting
	if data.ID == 0 {
		return true
	}

	return data.ID != mod.ID
}

func updateTargetingMetaData(ctx context.Context, data *dto.Targeting, exec boil.ContextExecutor) error {
	mods, err := getTargetingsByData(ctx, data, exec)
	if err != nil && err != sql.ErrNoRows {
		return eris.Wrap(err, "failed to get targetings for metadata update")
	}

	modMeta, err := createTargetingMetaData(mods, data.PublisherID, data.Domain)
	if err != nil && err != sql.ErrNoRows {
		return eris.Wrap(err, "failed to create targeting metadata")
	}

	err = modMeta.Insert(ctx, exec, boil.Infer())
	if err != nil {
		return eris.Wrapf(err, "failed to insert metadata record")
	}

	return nil
}

func createTargetingMetaData(mods models.TargetingSlice, publisher, domain string) (*models.MetadataQueue, error) {
	records := make([]TargetingRealtimeRecord, 0, len(mods))

	for _, mod := range mods {
		rule, err := dto.GetTargetingRegExp(mod)
		if err != nil {
			return nil, err
		}

		records = append(records, TargetingRealtimeRecord{
			RuleID: mod.RuleID,
			Rule:   rule,
			PriceModel: strings.ToLower(
				strings.ReplaceAll(
					mod.PriceModel,
					" ",
					"",
				),
			),
			Value: getTargetingValue(mod),
			DailyCap: func() *int {
				if mod.DailyCap.Valid {
					return &mod.DailyCap.Int
				}
				return nil
			}(),
		})
	}

	b, err := json.Marshal(records)
	if err != nil {
		return nil, err
	}

	return &models.MetadataQueue{
		TransactionID: bcguid.NewFromf(time.Now()),
		Key:           dto.GetTargetingKey(publisher, domain),
		Value:         b,
	}, nil
}

func getTargetingValue(mod *models.Targeting) float64 {
	switch mod.Status {
	case dto.TargetingStatusActive:
		return mod.Value
	default:
		// if rule is untargeted return 0 for CPM and Rev Share
		return 0
	}
}

// getColumnsToUpdate update only multiple value field (country, device type, os, browser, kv),
// placement type, price model, value, daily cap or/and status
func getColumnsToUpdate(newData *dto.Targeting, currentData *models.Targeting) ([]string, error) {
	columns := make([]string, 0, 13)
	columns = append(columns, models.TargetingColumns.RuleID)
	columns = append(columns, models.TargetingColumns.UpdatedAt)

	if !slices.Equal(newData.Country, currentData.Country) {
		currentData.Country = newData.Country
		columns = append(columns, models.TargetingColumns.Country)
	}

	if !slices.Equal(newData.DeviceType, currentData.DeviceType) {
		currentData.DeviceType = newData.DeviceType
		columns = append(columns, models.TargetingColumns.DeviceType)
	}

	if !slices.Equal(newData.OS, currentData.Os) {
		currentData.Os = newData.OS
		columns = append(columns, models.TargetingColumns.Os)
	}

	if !slices.Equal(newData.Browser, currentData.Browser) {
		currentData.Browser = newData.Browser
		columns = append(columns, models.TargetingColumns.Browser)
	}

	if newData.PlacementType != currentData.PlacementType.String {
		currentData.PlacementType = null.StringFrom(newData.PlacementType)
		columns = append(columns, models.TargetingColumns.PlacementType)
	}

	var currentKV map[string]string
	err := json.Unmarshal(currentData.KV.JSON, &currentKV)
	if err != nil {
		if currentData.KV.Valid {
			return nil, err
		}
	}

	if !reflect.DeepEqual(currentKV, newData.KV) {
		modKV, err := dto.GetModelKV(newData.KV)
		if err != nil {
			return nil, err
		}

		currentData.KV = modKV
		columns = append(columns, models.TargetingColumns.KV)
	}

	if newData.PriceModel != currentData.PriceModel {
		currentData.PriceModel = newData.PriceModel
		columns = append(columns, models.TargetingColumns.PriceModel)
	}

	if newData.Value != currentData.Value {
		currentData.Value = newData.Value
		columns = append(columns, models.TargetingColumns.Value)
	}

	if newData.Status != currentData.Status {
		currentData.Status = newData.Status
		columns = append(columns, models.TargetingColumns.Status)
	}

	newDailyCap := null.IntFromPtr(newData.DailyCap)
	if newDailyCap != currentData.DailyCap {
		currentData.DailyCap = newDailyCap
		columns = append(columns, models.TargetingColumns.DailyCap)
	}

	return columns, nil
}
