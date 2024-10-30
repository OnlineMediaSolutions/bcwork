package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/volatiletech/null/v8"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Floor struct {
	RuleId        string  `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	PublisherName string  `boil:"publisher_name" json:"publisher_name" toml:"publisher_name" yaml:"publisher_name"`
	Domain        string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Country       string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Floor         float64 `boil:"floor" json:"floor" toml:"floor" yaml:"floor"`
	Browser       string  `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string  `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string  `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
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
}

type FloorSlice []*Floor

type FloorRealtimeRecord struct {
	Rule   string  `json:"rule"`
	Floor  float64 `json:"floor"`
	RuleID string  `json:"rule_id"`
}

func (floor *Floor) FromModel(mod *models.Floor) error {

	floor.Publisher = mod.Publisher
	floor.Domain = mod.Domain
	floor.Floor = mod.Floor

	if mod.Os.Valid {
		floor.OS = mod.Os.String
	}

	if mod.Country.Valid {
		floor.Country = mod.Country.String
	}

	if mod.Device.Valid {
		floor.Device = mod.Device.String
	}

	if mod.PlacementType.Valid {
		floor.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		floor.Browser = mod.Browser.String
	}

	return nil
}

func (floor *Floor) GetRuleID() string {
	if len(floor.RuleId) > 0 {
		return floor.RuleId
	} else {
		return bcguid.NewFrom(floor.GetFormula())
	}
}

func (floor *Floor) GetFormula() string {
	p := floor.Publisher
	if p == "" {
		p = "*"
	}

	d := floor.Domain
	if d == "" {
		d = "*"
	}

	c := floor.Country
	if c == "" {
		c = "*"
	}

	os := floor.OS
	if os == "" {
		os = "*"
	}

	dt := floor.Device
	if dt == "" {
		dt = "*"
	}

	pt := floor.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := floor.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, c, os, dt, pt, b)

}

func (cs *FloorSlice) FromModel(slice models.FloorSlice) error {

	for _, mod := range slice {
		c := Floor{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func GetFloors(ctx context.Context, ops GetFloorOptions) (FloorSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.FloorColumns.Publisher).AddArray(ops.Pagination.Do())

	qmods = qmods.Add(qm.Select("DISTINCT *"))
	qmods = qmods.Add(qm.Load(models.FloorRels.FloorPublisher))

	mods, err := models.Floors(qmods...).All(ctx, bcdb.DB())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrap(err, "Failed to retrieve floors")
	}

	res := make(FloorSlice, 0)
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

	return mods
}

func UpdateFloorMetaData(c *fiber.Ctx, data constant.FloorUpdateRequest) error {
	_, err := json.Marshal(data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to parse hash value for floor", err)
	}

	err = SendFloorToRT(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}

func (floor *Floor) ToModel() *models.Floor {

	mod := models.Floor{
		RuleID:    floor.GetRuleID(),
		Floor:     floor.Floor,
		Publisher: floor.Publisher,
		Domain:    floor.Domain,
	}

	if floor.Country != "" {
		mod.Country = null.StringFrom(floor.Country)
	} else {
		mod.Country = null.String{}
	}

	if floor.OS != "" {
		mod.Os = null.StringFrom(floor.OS)
	} else {
		mod.Os = null.String{}
	}

	if floor.Device != "" {
		mod.Device = null.StringFrom(floor.Device)
	} else {
		mod.Device = null.String{}
	}

	if floor.PlacementType != "" {
		mod.PlacementType = null.StringFrom(floor.PlacementType)
	} else {
		mod.PlacementType = null.String{}
	}

	if floor.Browser != "" {
		mod.Browser = null.StringFrom(floor.Browser)
	} else {
		mod.Browser = null.String{}
	}

	return &mod
}

func UpdateFloors(c *fiber.Ctx, data constant.FloorUpdateRequest) (bool, error) {
	isInsert := false

	exists, err := models.Floors(
		models.FloorWhere.Publisher.EQ(data.Publisher),
		models.FloorWhere.Domain.EQ(data.Domain),
	).Exists(c.Context(), bcdb.DB())

	if err != nil {
		return false, err
	}

	if !exists {
		isInsert = true
	}

	floor := Floor{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		Device:        data.Device,
		Floor:         data.Floor,
		Browser:       data.Browser,
		OS:            data.OS,
		PlacementType: data.PlacementType,
	}

	modConf := floor.ToModel()

	//modConf := models.Floor{
	//	Publisher:     data.Publisher,
	//	Domain:        data.Domain,
	//	Device:        null.StringFrom(data.Device),
	//	Floor:         data.Floor,
	//	Country:        null.StringFrom(data.Country),
	//	Browser:        null.StringFrom(data.Browser),
	//	Os:             null.StringFrom(data.OS),
	//	PlacementType:  null.StringFrom(data.PlacementType),
	//	RuleID:        floor.GetRuleID(),
	//}

	err = modConf.Upsert(
		c.Context(),
		bcdb.DB(),
		true,
		[]string{models.FloorColumns.RuleID},
		boil.Blacklist(models.FloorColumns.CreatedAt),
		boil.Infer(),
	)
	if err != nil {
		return false, err
	}

	return isInsert, nil
}

func SendFloorToRT(c context.Context, updateRequest constant.FloorUpdateRequest) error {
	const PREFIX string = "price:floor:v2"
	modFloor, err := FloorQuery(c, updateRequest)

	if err != nil && err != sql.ErrNoRows {
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
	metadataKey := utils.CreateMetadataKey(key, PREFIX)
	metadataValue := utils.CreateMetadataObject(updateRequest, metadataKey, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for floor")
	}

	return nil
}

func FloorQuery(c context.Context, updateRequest constant.FloorUpdateRequest) (models.FloorSlice, error) {
	modFloor, err := models.Floors(
		models.FloorWhere.Domain.EQ(updateRequest.Domain),
		models.FloorWhere.Publisher.EQ(updateRequest.Publisher),
	).All(c, bcdb.DB())

	return modFloor, err
}

func CreateFloorMetadata(modFloor models.FloorSlice, finalRules []FloorRealtimeRecord) []FloorRealtimeRecord {
	if len(modFloor) != 0 {
		floors := make(FloorSlice, 0)
		floors.FromModel(modFloor)

		for _, floor := range floors {
			rule := FloorRealtimeRecord{
				Rule:   utils.GetFormulaRegex(floor.Country, floor.Domain, floor.Device, floor.PlacementType, floor.OS, floor.Browser, floor.Publisher),
				Floor:  floor.Floor,
				RuleID: floor.GetRuleID(),
			}
			finalRules = append(finalRules, rule)
		}
	}
	return finalRules
}
