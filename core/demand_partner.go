package core

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"errors"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Important note! *models.Dpo = demand partners

type DemandPartnerService struct {
	historyModule history.HistoryModule
}

func NewDemandPartnerService(historyModule history.HistoryModule) *DemandPartnerService {
	return &DemandPartnerService{
		historyModule: historyModule,
	}
}

type DemandPartnerGetOptions struct {
	Filter     DemandPartnerGetFilter `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type DemandPartnerGetFilter struct {
	DemandPartnerId   filter.StringArrayFilter `json:"demand_partner_id,omitempty"`
	DemandPartnerName filter.StringArrayFilter `json:"demand_partner_name,omitempty"`
	Active            filter.StringArrayFilter `json:"active,omitempty"`
}

func (filter *DemandPartnerGetFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	if len(filter.DemandPartnerId) > 0 {
		mods = append(mods, filter.DemandPartnerId.AndIn(models.DpoColumns.DemandPartnerID))
	}

	if len(filter.DemandPartnerName) > 0 {
		mods = append(mods, filter.DemandPartnerName.AndIn(models.DpoColumns.DemandPartnerName))
	}

	if len(filter.Active) > 0 {
		mods = append(mods, filter.Active.AndIn(models.DpoColumns.Active))
	}

	return mods
}

func (d *DemandPartnerService) GetDemandPartners(ctx context.Context, ops *DemandPartnerGetOptions) ([]*dto.DemandPartner, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.DpoColumns.DemandPartnerID).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *")).
		Add(
			qm.Load(models.DpoRels.DPParentDemandPartnerChildren),
			qm.Load(models.DpoRels.DemandPartnerDemandPartnerConnections),
		)

	mods, err := models.Dpos(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to retrieve demand partners: %w", err)
	}

	dps := make([]*dto.DemandPartner, 0, len(mods))
	for _, mod := range mods {
		dp := &dto.DemandPartner{}
		dp.FromModel(mod)
		dps = append(dps, dp)
	}

	return dps, nil
}

func (d *DemandPartnerService) CreateDemandPartner(ctx context.Context, data *dto.DemandPartner) error {
	isExists, err := models.Dpos(models.DpoWhere.DemandPartnerName.EQ(data.DemandPartnerName)).
		Exists(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to check existance of demand partner: %w", err)
	}

	if isExists {
		return errors.New("demand partner with such parameters already exists")
	}

	demandPartnerID := strings.ReplaceAll(strings.ToLower(data.DemandPartnerName), " ", "")

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = processDemandPartnerChildren(ctx, tx, demandPartnerID, data.Children)
	if err != nil {
		return fmt.Errorf("failed to process demand partner children: %w", err)
	}

	err = processDemandPartnerConnections(ctx, tx, demandPartnerID, data.Connections)
	if err != nil {
		return fmt.Errorf("failed to process demand partner connections: %w", err)
	}

	mod := data.ToModel(demandPartnerID)
	err = mod.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to insert demand partner: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to make commit for creating of demand partner: %w", err)
	}

	return nil
}

func (d *DemandPartnerService) UpdateDemandPartner(ctx context.Context, data *dto.DemandPartner) error {
	mod, err := models.Dpos(models.DpoWhere.DemandPartnerID.EQ(data.DemandPartnerID)).
		One(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to get demand partner with id [%v] to update: %w", data.DemandPartnerID, err)
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = processDemandPartnerChildren(ctx, tx, mod.DemandPartnerID, data.Children)
	if err != nil {
		return fmt.Errorf("failed to process demand partner children: %w", err)
	}

	err = processDemandPartnerConnections(ctx, tx, mod.DemandPartnerID, data.Connections)
	if err != nil {
		return fmt.Errorf("failed to process demand partner connections: %w", err)
	}

	newMod := data.ToModel(mod.DemandPartnerID)
	newMod.UpdatedAt = null.TimeFrom(time.Now())

	columns, err := getDemandPartnerColumnsToUpdate(newMod, mod)
	if err != nil {
		return fmt.Errorf("error getting columns for update: %w", err)
	}

	// if updating only updated_at
	if len(columns) == 1 {
		return errors.New("there are no new values to update demand partner")
	}

	_, err = newMod.Update(ctx, tx, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to update demand partner: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to make commit for updating of demand partner: %w", err)
	}

	return nil
}

func processDemandPartnerChildren(
	ctx context.Context,
	tx *sql.Tx,
	demandPartnerID string,
	children []*dto.DemandPartnerChild,
) error {
	for _, child := range children {
		isExists, err := models.DemandPartnerChildren(
			models.DemandPartnerChildWhere.DPParentID.EQ(demandPartnerID),
			models.DemandPartnerChildWhere.DPChildName.EQ(child.DPChildName),
		).Exists(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to check existance of child of demand partner: %w", err)
		}

		mod := child.ToModel(demandPartnerID)

		if !isExists {
			err := mod.Insert(ctx, tx, boil.Infer())
			if err != nil {
				return fmt.Errorf("failed to insert demand partner child: %w", err)
			}
		}

		// TODO: process updating children
	}

	return nil
}

func processDemandPartnerConnections(
	ctx context.Context,
	tx *sql.Tx,
	demandPartnerID string,
	connections []*dto.DemandPartnerConnection,
) error {
	for _, connection := range connections {
		isExists, err := models.DemandPartnerConnections(
			models.DemandPartnerConnectionWhere.DemandPartnerID.EQ(demandPartnerID),
			models.DemandPartnerConnectionWhere.PublisherAccount.EQ(connection.PublisherAccount),
		).Exists(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to check existance of connection of demand partner: %w", err)
		}

		mod := connection.ToModel(demandPartnerID)

		if !isExists {
			err := mod.Insert(ctx, tx, boil.Infer())
			if err != nil {
				return fmt.Errorf("failed to insert demand partner connection: %w", err)
			}
		}

		// TODO: process updating connections
	}

	return nil
}

func getDemandPartnerColumnsToUpdate(newData, oldData *models.Dpo) ([]string, error) {
	const boilTagName = "boil"

	blacklistColumns := []string{
		models.DpoColumns.DemandPartnerID,
		models.DpoColumns.CreatedAt,
	}
	columns := make([]string, 0, 12)

	oldValueReflection := reflect.ValueOf(oldData).Elem()
	newValueReflection := reflect.ValueOf(newData).Elem()

	for i := 0; i < oldValueReflection.NumField(); i++ {
		field := oldValueReflection.Type().Field(i)
		property := strings.Split(field.Tag.Get(boilTagName), ",")[0]
		oldFieldValue := oldValueReflection.Field(i)
		newFieldValue := newValueReflection.Field(i)

		if !reflect.DeepEqual(oldFieldValue.Interface(), newFieldValue.Interface()) &&
			!slices.Contains(blacklistColumns, property) {
			columns = append(columns, property)
		}
	}

	return columns, nil
}
