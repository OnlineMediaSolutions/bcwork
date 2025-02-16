package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/m6yf/bcwork/dto"
	"github.com/rs/zerolog/log"

	"sort"
	"strings"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/utils"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type FloorService struct {
	historyModule history.HistoryModule
}

func NewFloorService(historyModule history.HistoryModule) *FloorService {
	return &FloorService{
		historyModule: historyModule,
	}
}

type GetFloorOptions struct {
	Filter     FloorFilter            `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type FloorFilter struct {
	Publisher filter.StringArrayFilter `json:"publisher"`
	Domain    filter.StringArrayFilter `json:"domain"`
	Country   filter.StringArrayFilter `json:"country"`
	Device    filter.StringArrayFilter `json:"device"`
	Active    *filter.BoolFilter       `json:"active"`
}

type FloorRealtimeRecord struct {
	Rule   string  `json:"rule"`
	Floor  float64 `json:"floor"`
	RuleID string  `json:"rule_id"`
}

func (f *FloorService) GetFloors(ctx context.Context, ops GetFloorOptions) (dto.FloorSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.FloorColumns.Publisher).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *")).
		Add(qm.Load(models.FloorRels.FloorPublisher))

	mods, err := models.Floors(qmods...).All(ctx, bcdb.DB())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrap(err, "Failed to retrieve floors")
	}

	res := make(dto.FloorSlice, 0)
	err = res.FromModel(mods)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (filter *FloorFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.FloorColumns.Publisher))
	}

	if len(filter.Device) > 0 {
		mods = append(mods, filter.Device.AndIn(models.FloorColumns.Device))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.FloorColumns.Domain))
	}

	if len(filter.Country) > 0 {
		mods = append(mods, filter.Country.AndIn(models.FloorColumns.Country))
	}

	if filter.Active != nil {
		mods = append(mods, filter.Active.Where(models.FloorColumns.Active))
	}

	return mods
}

func (f *FloorService) UpdateFloorMetaData(ctx context.Context, data dto.FloorUpdateRequest) error {
	_, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to parse hash value for floor: %w", err)
	}

	err = sendFloorToRT(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (f *FloorService) UpdateFloors(ctx context.Context, data dto.FloorUpdateRequest) (bool, error) {
	var isInsert bool

	floor := dto.Floor{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		Device:        data.Device,
		Floor:         data.Floor,
		Browser:       data.Browser,
		OS:            data.OS,
		PlacementType: data.PlacementType,
		RuleId:        data.RuleId,
		Active:        data.Active,
	}

	if len(floor.RuleId) == 0 {
		floor.RuleId = floor.GetRuleID()
	}
	mod := floor.ToModel()

	var old any
	oldMod, err := models.Floors(
		models.FloorWhere.RuleID.EQ(mod.RuleID),
	).One(ctx, bcdb.DB())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
		[]string{models.FloorColumns.RuleID},
		boil.Blacklist(models.FloorColumns.CreatedAt),
		boil.Infer(),
	)
	if err != nil {
		return false, err
	}

	f.historyModule.SaveAction(ctx, old, mod, nil)

	return isInsert, nil
}

func sendFloorToRT(ctx context.Context, updateRequest dto.FloorUpdateRequest) error {
	modFloor, err := models.Floors(
		models.FloorWhere.Domain.EQ(updateRequest.Domain),
		models.FloorWhere.Publisher.EQ(updateRequest.Publisher),
	).All(ctx, bcdb.DB())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return eris.Wrapf(err, "Failed to fetch floors for publisher %s", updateRequest.Publisher)
	}

	var finalRules []FloorRealtimeRecord

	finalRules = CreateFloorMetadata(modFloor, finalRules)

	finalOutput := struct {
		Rules []FloorRealtimeRecord `json:"rules"`
	}{Rules: finalRules}

	value, err := json.Marshal(finalOutput)
	if err != nil {
		return eris.Wrap(err, "failed to marshal floorRT to JSON")
	}

	key := utils.GetMetadataObject(updateRequest)
	metadataKey := utils.CreateMetadataKey(key, utils.FloorMetaDataKeyPrefix)
	metadataValue := utils.CreateMetadataObject(updateRequest, metadataKey, value)

	err = metadataValue.Insert(ctx, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for floor")
	}

	return nil
}

func CreateFloorMetadata(modFloor models.FloorSlice, finalRules []FloorRealtimeRecord) []FloorRealtimeRecord {
	if len(modFloor) != 0 {
		floors := make(dto.FloorSlice, 0)
		err := floors.FromModel(modFloor)
		if err != nil {
			log.Error().Err(err).Msg("failed to map floors")
		}

		for _, floor := range floors {
			rule := FloorRealtimeRecord{
				Rule:   utils.GetFormulaRegex(floor.Country, floor.Domain, floor.Device, floor.PlacementType, floor.OS, floor.Browser, floor.Publisher),
				Floor:  floor.Floor,
				RuleID: floor.RuleId,
			}
			if len(rule.RuleID) == 0 {
				rule.RuleID = floor.GetRuleID()
			}

			finalRules = append(finalRules, rule)
		}
	}
	sortFloorRules(finalRules)

	return finalRules
}

func sortFloorRules(floors []FloorRealtimeRecord) {
	sort.Slice(floors, func(i, j int) bool {
		return strings.Count(floors[i].Rule, "*") < strings.Count(floors[j].Rule, "*")
	})
}
