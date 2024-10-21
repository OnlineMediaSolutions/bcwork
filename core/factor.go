package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

type Factor struct {
	RuleId        string  `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string  `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Factor        float64 `boil:"factor" json:"factor,omitempty" toml:"factor" yaml:"factor,omitempty"`
	Browser       string  `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string  `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string  `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
}

type FactorSlice []*Factor

type FactorRealtimeRecord struct {
	Rule   string  `json:"rule"`
	Factor float64 `json:"factor"`
	RuleID string  `json:"rule_id"`
}

type GetFactorOptions struct {
	Filter     FactorFilter           `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type FactorFilter struct {
	Publisher filter.StringArrayFilter `json:"publisher,omitempty"`
	Domain    filter.StringArrayFilter `json:"domain,omitempty"`
	Country   filter.StringArrayFilter `json:"country,omitempty"`
	Device    filter.StringArrayFilter `json:"device,omitempty"`
}

func (factor *Factor) FromModel(mod *models.Factor) error {

	factor.Publisher = mod.Publisher
	factor.Domain = mod.Domain
	factor.Country = mod.Country
	factor.Device = mod.Device
	factor.Factor = mod.Factor

	return nil
}

func (cs *FactorSlice) FromModel(slice models.FactorSlice) error {

	for _, mod := range slice {
		c := Factor{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func GetFactors(ctx context.Context, ops *GetFactorOptions) (FactorSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.FactorColumns.Publisher).AddArray(ops.Pagination.Do())

	qmods = qmods.Add(qm.Select("DISTINCT *"))

	mods, err := models.Factors(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve factors")
	}

	res := make(FactorSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (filter *FactorFilter) QueryMod() qmods.QueryModsSlice {

	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.FactorColumns.Publisher))
	}

	if len(filter.Device) > 0 {
		mods = append(mods, filter.Device.AndIn(models.FactorColumns.Device))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.FactorColumns.Domain))
	}

	if len(filter.Country) > 0 {
		mods = append(mods, filter.Country.AndIn(models.FactorColumns.Country))
	}

	return mods
}

func UpdateMetaData(c *fiber.Ctx, data constant.FactorUpdateRequest) error {
	_, err := json.Marshal(data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to parse hash value for factor", err)
	}

	go func() {
		err = SendFactorToRT(context.Background(), data)
	}()

	if err != nil {
		return err
	}
	return nil
}

func FactorQuery(c context.Context, updateRequest constant.FactorUpdateRequest) (models.FactorSlice, error) {
	modFactor, err := models.Factors(
		models.FactorWhere.Domain.EQ(updateRequest.Domain),
		models.FactorWhere.Publisher.EQ(updateRequest.Publisher),
	).All(c, bcdb.DB())

	return modFactor, err
}

func (factor *Factor) GetFormula() string {
	p := factor.Publisher
	if p == "" {
		p = "*"
	}

	d := factor.Domain
	if d == "" {
		d = "*"
	}

	c := factor.Country
	if c == "" {
		c = "*"
	}

	os := factor.OS
	if os == "" {
		os = "*"
	}

	dt := factor.Device
	if dt == "" {
		dt = "*"
	}

	pt := factor.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := factor.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, c, os, dt, pt, b)

}

func (factor *Factor) GetRuleID() string {
	if len(factor.RuleId) > 0 {
		return factor.RuleId
	} else {
		return bcguid.NewFrom(factor.GetFormula())
	}
}

func CreateFactorMetadata(modFactor models.FactorSlice, finalRules []FactorRealtimeRecord) []FactorRealtimeRecord {
	if len(modFactor) != 0 {
		factors := make(FactorSlice, 0)
		factors.FromModel(modFactor)

		for _, factor := range factors {
			rule := FactorRealtimeRecord{
				Rule:   utils.GetFormulaRegex(factor.Country, factor.Domain, factor.Device, factor.PlacementType, factor.OS, factor.Browser, factor.Publisher),
				Factor: factor.Factor,
				RuleID: factor.GetRuleID(),
			}
			finalRules = append(finalRules, rule)
		}
	}
	return finalRules
}

func SendFactorToRT(c context.Context, updateRequest constant.FactorUpdateRequest) error {
	modFactor, err := FactorQuery(c, updateRequest)

	if err != nil && err != sql.ErrNoRows {
		return eris.Wrapf(err, "Failed to fetch factors for publisher %s", updateRequest.Publisher)
	}

	var finalRules []FactorRealtimeRecord

	finalRules = CreateFactorMetadata(modFactor, finalRules)

	finalOutput := struct {
		Rules []FactorRealtimeRecord `json:"rules"`
	}{Rules: finalRules}

	value, err := json.Marshal(finalOutput)
	if err != nil {
		return eris.Wrap(err, "failed to marshal factorRT to JSON")
	}

	key := utils.GetMetadataObject(updateRequest)
	metadataKey := utils.CreateMetadataKey(key, utils.FactorMetaDataKeyPrefix)
	metadataValue := utils.CreateMetadataObject(updateRequest, metadataKey, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record for factor")
	}

	return nil
}

func UpdateFactor(c *fiber.Ctx, data *constant.FactorUpdateRequest) (bool, error) {
	isInsert := false

	exists, err := models.Factors(
		models.FactorWhere.Publisher.EQ(data.Publisher),
		models.FactorWhere.Domain.EQ(data.Domain),
		models.FactorWhere.Device.EQ(data.Device),
		models.FactorWhere.Country.EQ(data.Country),
		models.FactorWhere.Os.EQ(data.OS),
		models.FactorWhere.Browser.EQ(data.Browser),
		models.FactorWhere.PlacementType.EQ(data.PlacementType),
	).Exists(c.Context(), bcdb.DB())

	if err != nil {
		return false, err
	}

	if !exists {
		isInsert = true
	}

	factor := Factor{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		Device:        data.Device,
		Factor:        data.Factor,
		Browser:       data.Browser,
		OS:            data.OS,
		PlacementType: data.PlacementType,
	}

	modConf := models.Factor{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Device:        data.Device,
		Factor:        data.Factor,
		Country:       data.Country,
		Os:            data.OS,
		Browser:       data.Browser,
		PlacementType: data.PlacementType,
		RuleID:        factor.GetRuleID(),
	}

	err = modConf.Upsert(
		c.Context(),
		bcdb.DB(),
		true,
		[]string{models.FactorColumns.RuleID},
		boil.Blacklist(models.FactorColumns.CreatedAt),
		boil.Infer(),
	)

	if err != nil {
		return false, err
	}

	return isInsert, nil
}
