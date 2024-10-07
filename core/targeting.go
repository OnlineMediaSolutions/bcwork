package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TargetingRealtimeRecord struct {
	RuleID     string  `json:"rule_id"`
	Rule       string  `json:"rule"`
	PriceModel string  `json:"price_model"`
	Value      float64 `json:"value"`
	DailyCap   int     `json:"daily_cap"`
}

type TargetingOptions struct {
	Filter     TargetingFilter        `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type TargetingFilter struct {
	Publisher     filter.StringArrayFilter   `json:"publisher,omitempty"`
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

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.TargetingColumns.Publisher))
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

func GetTargetings(ctx context.Context, ops *TargetingOptions) ([]*constant.Targeting, error) {
	qmods := ops.Filter.queryMod().
		Order(ops.Order, nil, models.DpoRuleColumns.RuleID).
		AddArray(ops.Pagination.Do())

	mods, err := models.Targetings(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve targetings")
	}

	targetings := make([]*constant.Targeting, 0, len(mods))
	for _, mod := range mods {
		targeting := new(constant.Targeting)
		err := targeting.FromModel(mod)
		if err != nil {
			return nil, eris.Wrap(err, "failed to map data from model")
		}

		targetings = append(targetings, targeting)
	}

	return targetings, nil
}

func CreateTargeting(ctx context.Context, data *constant.Targeting) error {
	data.PrepareData()

	duplicate, err := getTargetingByProps(ctx, data)
	if err != nil {
		return eris.Wrap(err, "failed to check for duplicates")
	}

	if duplicate != nil {
		duplicateString := fmt.Sprintf(
			"country=%v,device_type=%v,browser=%v,os=%v,kv=%v",
			duplicate.Country, duplicate.DeviceType, duplicate.Browser, duplicate.Os, string(duplicate.KV.JSON),
		)
		return fmt.Errorf(
			"could not create targeting: there is same targeting with such parameters [%v]", duplicateString)
	}

	mod, err := data.ToModel()
	if err != nil {
		return eris.Wrap(err, "failed to map data to model")
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return eris.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	err = mod.Insert(ctx, tx, boil.Infer())
	if err != nil && err != sql.ErrNoRows {
		return eris.Wrap(err, "failed to upsert targeting")
	}

	err = updateTargetingMetaData(ctx, data, tx)
	if err != nil {
		return eris.Wrapf(err, "failed to update targeting metadata")
	}

	err = tx.Commit()
	if err != nil {
		return eris.Wrapf(err, "failed to commit targeting and metadata")
	}

	return nil
}

func UpdateTargeting(ctx context.Context, data *constant.Targeting) error {
	data.PrepareData()

	mod, err := getTargetingByProps(ctx, data)
	if err != nil {
		return eris.Wrap(err, "failed to get targeting to update")
	}

	if mod == nil {
		return errors.New("no targeting found to update")
	}

	columns := getColumnsToUpdate(data, mod)
	if len(columns) == 0 {
		return errors.New("there are no new values to update targeting")
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return eris.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	_, err = mod.Update(ctx, tx, boil.Whitelist(columns...))
	if err != nil {
		return eris.Wrap(err, "failed to update targeting")
	}

	err = updateTargetingMetaData(ctx, data, tx)
	if err != nil {
		return eris.Wrapf(err, "failed to update targeting metadata")
	}

	err = tx.Commit()
	if err != nil {
		return eris.Wrapf(err, "failed to commit targeting updates and metadata")
	}

	return nil
}
func getTargetingsByHash(ctx context.Context, hash string, exec boil.ContextExecutor) (models.TargetingSlice, error) {
	mods, err := models.Targetings(
		models.TargetingWhere.Hash.EQ(hash),
		models.TargetingWhere.Status.NEQ(constant.TargetingStatusArchived),
		qm.OrderBy(models.TargetingColumns.Country),
		qm.OrderBy(models.TargetingColumns.DeviceType),
		qm.OrderBy(models.TargetingColumns.Browser),
		qm.OrderBy(models.TargetingColumns.Os),
		qm.OrderBy(models.TargetingColumns.KV),
	).All(ctx, exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return mods, nil
}

func getTargetingByProps(ctx context.Context, data *constant.Targeting) (*models.Targeting, error) {
	var qmods qmods.QueryModsSlice
	qmods = qmods.Add(
		models.TargetingWhere.Hash.EQ(data.Hash),
		models.TargetingWhere.Status.NEQ(constant.TargetingStatusArchived),
		qm.Where(models.TargetingColumns.Country+" && ARRAY['"+strings.Join(data.Country, "','")+"']"),
		qm.Where(models.TargetingColumns.DeviceType+" && ARRAY['"+strings.Join(data.DeviceType, "','")+"']"),
		qm.Where(models.TargetingColumns.Browser+" && ARRAY['"+strings.Join(data.Browser, "','")+"']"),
		qm.Where(models.TargetingColumns.Os+" && ARRAY['"+strings.Join(data.OS, "','")+"']"),
	).Add(
		getKeyValueWhereQueries(data.KV)...,
	)

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

	var i int = 3 // depends on getTargetingByProps
	for key, value := range kv {
		mods = append(mods, qm.Where(
			fmt.Sprintf("%v ->> '%v' = $%v", models.TargetingColumns.KV, key, i),
			value,
		))
		i++
	}

	return mods
}

func updateTargetingMetaData(ctx context.Context, data *constant.Targeting, exec boil.ContextExecutor) error {
	mods, err := getTargetingsByHash(ctx, data.Hash, exec)
	if err != nil && err != sql.ErrNoRows {
		return eris.Wrap(err, "failed to get targetings for metadata update")
	}

	modMeta, err := createTargetingMetaData(mods, data.Publisher, data.Domain)
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
		rule, err := constant.GetTargetingRegExp(mod)
		if err != nil {
			return nil, err
		}

		records = append(records, TargetingRealtimeRecord{
			RuleID:     mod.RuleID,
			Rule:       rule,
			PriceModel: mod.PriceModel,
			Value:      getTargetingValue(mod),
			DailyCap:   mod.DailyCap.Int,
		})
	}

	b, err := json.Marshal(records)
	if err != nil {
		return nil, err
	}

	return &models.MetadataQueue{
		TransactionID: bcguid.NewFromf(time.Now()),
		Key:           constant.GetTargetingKey(publisher, domain),
		Value:         b,
	}, nil
}

func getTargetingValue(mod *models.Targeting) float64 {
	switch mod.Status {
	case constant.TargetingStatusActive:
		return mod.Value
	default:
		// if rule is untargeted return 0 for CPM and Rev Share
		return 0
	}
}

// getColumnsToUpdate update only cost model, value or/and status
func getColumnsToUpdate(newData *constant.Targeting, currentData *models.Targeting) []string {
	columns := make([]string, 0, 4)
	if newData.PriceModel != "" && newData.PriceModel != currentData.PriceModel {
		currentData.PriceModel = newData.PriceModel
		columns = append(columns, models.TargetingColumns.PriceModel)
	}

	if newData.Value != 0 && newData.Value != currentData.Value {
		currentData.Value = newData.Value
		columns = append(columns, models.TargetingColumns.Value)
	}

	if newData.Status != "" && newData.Status != currentData.Status {
		currentData.Status = newData.Status
		columns = append(columns, models.TargetingColumns.Status)
	}

	if newData.DailyCap != 0 && newData.DailyCap != currentData.DailyCap.Int {
		currentData.DailyCap = null.IntFrom(newData.DailyCap)
		columns = append(columns, models.TargetingColumns.DailyCap)
	}

	return columns
}
