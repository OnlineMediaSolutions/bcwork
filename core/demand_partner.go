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
	"github.com/m6yf/bcwork/utils/helpers"
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
	Active            *filter.BoolFilter       `json:"active,omitempty"`
	Automation        *filter.BoolFilter       `json:"automation,omitempty"`
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

	if filter.Active != nil {
		mods = append(mods, filter.Active.Where(models.DpoColumns.Active))
	}

	if filter.Automation != nil {
		mods = append(mods, filter.Automation.Where(models.DpoColumns.Automation))
	}

	return mods
}

func (d *DemandPartnerService) GetDemandPartners(ctx context.Context, ops *DemandPartnerGetOptions) ([]*dto.DemandPartner, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.DpoColumns.DemandPartnerID).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *")).
		Add(
			qm.Load(models.DpoRels.Manager),
			qm.Load(models.DpoRels.SeatOwner),
			qm.Load(
				qm.Rels(
					models.DpoRels.DemandPartnerDemandPartnerConnections,
					models.DemandPartnerConnectionRels.DPConnectionDemandPartnerChildren,
				),
			),
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
		return fmt.Errorf("demand partner with name [%v] already exists", data.DemandPartnerName)
	}

	demandPartnerID := strings.ReplaceAll(strings.ToLower(data.DemandPartnerName), " ", "")

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	mod := data.ToModel(demandPartnerID)
	err = mod.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to insert demand partner: %w", err)
	}

	_, err = processDemandPartnerConnections(ctx, tx, demandPartnerID, data.Connections)
	if err != nil {
		return fmt.Errorf("failed to process demand partner connections: %w", err)
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

	isConnectionsChanged, err := processDemandPartnerConnections(ctx, tx, mod.DemandPartnerID, data.Connections)
	if err != nil {
		return fmt.Errorf("failed to process demand partner connections: %w", err)
	}

	newMod := data.ToModel(mod.DemandPartnerID)
	newMod.UpdatedAt = null.TimeFrom(time.Now().UTC())

	columns, err := getModelsColumnsToUpdate(
		mod, newMod,
		[]string{
			models.DpoColumns.DemandPartnerID,
			models.DpoColumns.CreatedAt,
		},
	)
	if err != nil {
		return fmt.Errorf("error getting demand partner columns for update: %w", err)
	}

	// if updating only updated_at
	if len(columns) == 1 && !isConnectionsChanged {
		return errors.New("there are no new values to update demand partner")
	}

	_, err = newMod.Update(ctx, tx, boil.Whitelist(columns...))
	if err != nil {
		return fmt.Errorf("failed to update demand partner: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to make commit for updating of demand partner: %w", err)
	}

	return nil
}

func processDemandPartnerConnections(
	ctx context.Context,
	tx *sql.Tx,
	demandPartnerID string,
	connections []*dto.DemandPartnerConnection,
) (bool, error) {
	var (
		isChanged         bool
		isChildrenChanged bool
	)

	modConnections, err := models.DemandPartnerConnections(
		models.DemandPartnerConnectionWhere.DemandPartnerID.EQ(demandPartnerID),
	).All(ctx, tx)
	if err != nil {
		return false, fmt.Errorf("failed to get current connections of demand partner: %w", err)
	}

	modConnectionsMap := make(map[string]*models.DemandPartnerConnection, len(modConnections))
	for _, modConnection := range modConnections {
		modConnectionsMap[modConnection.PublisherAccount] = modConnection
	}

	for _, connection := range connections {
		mod := connection.ToModel(demandPartnerID)
		mod.UpdatedAt = null.TimeFrom(time.Now().UTC())

		oldMod, ok := modConnectionsMap[mod.PublisherAccount]
		if !ok {
			isChanged = true

			err := mod.Insert(ctx, tx, boil.Blacklist(models.DemandPartnerChildColumns.UpdatedAt))
			if err != nil {
				return false, fmt.Errorf("failed to insert demand partner connection: %w", err)
			}
		} else {
			mod.ID = oldMod.ID

			columns, err := getModelsColumnsToUpdate(
				oldMod, mod,
				[]string{
					models.DemandPartnerConnectionColumns.ID,
					models.DemandPartnerConnectionColumns.CreatedAt,
					models.DemandPartnerConnectionColumns.DemandPartnerID,
					models.DemandPartnerConnectionColumns.PublisherAccount,
				},
			)
			if err != nil {
				return false, fmt.Errorf("error getting demand partner connection columns for update: %w", err)
			}

			// if updating not only updated_at
			if len(columns) > 1 {
				isChanged = true

				_, err := mod.Update(ctx, tx, boil.Whitelist(columns...))
				if err != nil {
					return false, fmt.Errorf("failed to update demand partner connection: %w", err)
				}
			}
		}

		isConnectionChildrenChanged, err := processDemandPartnerChildren(ctx, tx, mod.ID, connection.Children)
		if err != nil {
			return false, fmt.Errorf("failed to process demand partner children: %w", err)
		}

		isChildrenChanged = isChildrenChanged || isConnectionChildrenChanged

		delete(modConnectionsMap, mod.PublisherAccount)
	}

	// deleting demand partner connections which weren't been in request
	for _, modConnection := range modConnectionsMap {
		// if connection was deleted, delete all its ads.txt lines
		_, err := models.AdsTXTS(models.AdsTXTWhere.DemandPartnerConnectionID.EQ(null.IntFrom(modConnection.ID))).DeleteAll(ctx, tx)
		if err != nil {
			return false, fmt.Errorf("failed to delete demand partner connection ads txt lines: %w", err)
		}

		// if connection was deleted, delete all its children
		_, err = processDemandPartnerChildren(ctx, tx, modConnection.ID, nil)
		if err != nil {
			return false, fmt.Errorf("failed to delete demand partner connection children: %w", err)
		}

		_, err = modConnection.Delete(ctx, tx)
		if err != nil {
			return false, fmt.Errorf("failed to delete demand partner connection: %w", err)
		}
	}

	return isChanged || isChildrenChanged, nil
}

func processDemandPartnerChildren(
	ctx context.Context,
	tx *sql.Tx,
	connectionID int,
	children []*dto.DemandPartnerChild,
) (bool, error) {
	var isChanged bool

	modChildren, err := models.DemandPartnerChildren(
		models.DemandPartnerChildWhere.DPConnectionID.EQ(connectionID),
	).All(ctx, tx)
	if err != nil {
		return false, fmt.Errorf("failed to get current children of demand partner: %w", err)
	}

	modChildrenMap := make(map[string]*models.DemandPartnerChild, len(modChildren))
	for _, modChild := range modChildren {
		modChildrenMap[modChild.DPChildName] = modChild
	}

	for _, child := range children {
		mod := child.ToModel(connectionID)
		mod.UpdatedAt = null.TimeFrom(time.Now().UTC())

		oldMod, ok := modChildrenMap[mod.DPChildName]
		if !ok {
			isChanged = true

			err := mod.Insert(ctx, tx, boil.Blacklist(models.DemandPartnerChildColumns.UpdatedAt))
			if err != nil {
				return false, fmt.Errorf("failed to insert demand partner child: %w", err)
			}
		} else {
			columns, err := getModelsColumnsToUpdate(
				oldMod, mod,
				[]string{
					models.DemandPartnerChildColumns.ID,
					models.DemandPartnerChildColumns.CreatedAt,
					models.DemandPartnerChildColumns.DPConnectionID,
					models.DemandPartnerChildColumns.DPChildName,
				},
			)
			if err != nil {
				return false, fmt.Errorf("error getting demand partner child columns for update: %w", err)
			}

			// if updating not only updated_at
			if len(columns) > 1 {
				isChanged = true
				mod.ID = oldMod.ID

				_, err = mod.Update(ctx, tx, boil.Whitelist(columns...))
				if err != nil {
					return false, fmt.Errorf("failed to update demand partner child: %w", err)
				}
			}
		}

		delete(modChildrenMap, mod.DPChildName)
	}

	// delete demand partner children which weren't been in request
	for _, modChild := range modChildrenMap {
		// if demand partner child was deleted, delete all its ads.txt lines
		_, err := models.AdsTXTS(models.AdsTXTWhere.DemandPartnerChildID.EQ(null.IntFrom(modChild.ID))).DeleteAll(ctx, tx)
		if err != nil {
			return false, fmt.Errorf("failed to delete demand partner child ads txt lines: %w", err)
		}

		_, err = modChild.Delete(ctx, tx)
		if err != nil {
			return false, fmt.Errorf("failed to delete demand partner child: %w", err)
		}
	}

	return isChanged, nil
}

func getModelsColumnsToUpdate(oldData, newData any, blacklistColumns []string) ([]string, error) {
	const boilTagName = "boil"

	oldValueReflection, err := helpers.GetStructReflectValue(oldData)
	if err != nil {
		return nil, fmt.Errorf("cannot get reflection of old data: %w", err)
	}

	newValueReflection, err := helpers.GetStructReflectValue(newData)
	if err != nil {
		return nil, fmt.Errorf("cannot get reflection of new data: %w", err)
	}

	if oldValueReflection.Type().Name() != newValueReflection.Type().Name() {
		return nil, fmt.Errorf(
			"provided different structs: old [%v], new [%v]",
			oldValueReflection.Type().Name(), newValueReflection.Type().Name(),
		)
	}

	columns := make([]string, 0, oldValueReflection.NumField())

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
