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
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type FactorUpdateRequest struct {
	RuleId        string  `json:"rule_id"`
	Publisher     string  `json:"publisher"`
	Domain        string  `json:"domain"`
	Device        string  `json:"device"`
	Factor        float64 `json:"factor"`
	Country       string  `json:"country"`
	Browser       string  `json:"browser"`
	OS            string  `json:"os"`
	PlacementType string  `json:"placement_type"`
}

type Factor struct {
	RuleID        string  `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
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

type FactorRealtimeRecord struct {
	Rule   string  `json:"rule"`
	Factor float64 `json:"factor"`
	RuleID string  `json:"rule_id"`
}

func (f FactorUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f FactorUpdateRequest) GetDomain() string        { return f.Domain }
func (f FactorUpdateRequest) GetDevice() string        { return f.Device }
func (f FactorUpdateRequest) GetCountry() string       { return f.Country }
func (f FactorUpdateRequest) GetBrowser() string       { return f.Browser }
func (f FactorUpdateRequest) GetOS() string            { return f.OS }
func (f FactorUpdateRequest) GetPlacementType() string { return f.PlacementType }

func (factor *Factor) FromModel(mod *models.Factor) error {

	factor.Publisher = mod.Publisher
	factor.Domain = mod.Domain
	factor.Country = mod.Country
	factor.Device = mod.Device
	factor.Factor = mod.Factor
	factor.RuleID = mod.RuleID
	factor.PlacementType = helpers.GetStringOrEmpty(mod.PlacementType)
	factor.OS = helpers.GetStringOrEmpty(mod.Os)
	factor.Browser = helpers.GetStringOrEmpty(mod.Browser)

	return nil
}

func (factor *Factor) GetRuleID() string {
	return bcguid.NewFrom(factor.GetFormula())
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

func UpdateMetaData(c *fiber.Ctx, data *FactorUpdateRequest) error {
	_, err := json.Marshal(data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to parse hash value for factor metadata", err)
	}

	err = SendFactorToRT(context.Background(), *data)
	if err != nil {
		return err
	}
	return nil
}

func SendFactorToRT(c context.Context, updateRequest FactorUpdateRequest) error {

	const PREFIX string = "price:factor:v2"
	modFactor, err := factorQuery(c, updateRequest)

	if err != nil && err != sql.ErrNoRows {
		return eris.Wrapf(err, "failed to fetch  for publisher %s", updateRequest.Publisher)
	}

	var finalRules []FactorRealtimeRecord

	finalRules = createFactorMetadata(modFactor, finalRules, updateRequest)

	finalOutput := struct {
		Rules []FactorRealtimeRecord `json:"rules"`
	}{Rules: finalRules}

	value, err := json.Marshal(finalOutput)
	if err != nil {
		return eris.Wrap(err, "failed to marshal factorRT to JSON")
	}

	key := utils.GetMetadataObject(updateRequest)
	metadataKey := utils.CreateMetadataKey(key, PREFIX)
	metadataValue := utils.CreateMetadataObject(updateRequest, metadataKey, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record")
	}

	return nil
}

func factorQuery(c context.Context, updateRequest FactorUpdateRequest) (models.FactorSlice, error) {
	modFactor, err := models.Factors(
		models.FactorWhere.Country.EQ(updateRequest.Country),
		models.FactorWhere.Domain.EQ(updateRequest.Domain),
		models.FactorWhere.Device.EQ(updateRequest.Device),
		models.FactorWhere.Publisher.EQ(updateRequest.Publisher),
	).All(c, bcdb.DB())
	return modFactor, err
}

func createFactorMetadata(modFactor models.FactorSlice, finalRules []FactorRealtimeRecord, updateRequest FactorUpdateRequest) []FactorRealtimeRecord {
	if len(modFactor) != 0 {
		factors := make(FactorSlice, 0)
		factors.FromModel(modFactor)

		for _, factor := range factors {
			rule := FactorRealtimeRecord{
				Rule:   utils.GetFormulaRegex(factor.Country, factor.Domain, factor.Device, factor.PlacementType, factor.OS, factor.Browser, factor.Publisher, false),
				Factor: factor.Factor,
				RuleID: factor.RuleID,
			}
			finalRules = append(finalRules, rule)

		}
	}

	newFactor := Factor{
		Publisher:     updateRequest.Publisher,
		Domain:        updateRequest.Domain,
		Country:       updateRequest.Country,
		Device:        updateRequest.Device,
		Factor:        updateRequest.Factor,
		Browser:       updateRequest.Browser,
		OS:            updateRequest.OS,
		PlacementType: updateRequest.PlacementType,
	}

	newRule := FactorRealtimeRecord{
		Rule: utils.GetFormulaRegex(
			updateRequest.Country,
			updateRequest.Domain,
			updateRequest.Device,
			updateRequest.PlacementType,
			updateRequest.OS,
			updateRequest.Browser,
			updateRequest.Publisher,
			false),
		Factor: updateRequest.Factor,
		RuleID: newFactor.GetRuleID(),
	}
	finalRules = append(finalRules, newRule)
	return finalRules
}

func UpdateFactor(c *fiber.Ctx, data *FactorUpdateRequest) (bool, error) {

	isInsert := false

	exists, err := models.Factors(
		models.FactorWhere.Publisher.EQ(data.Publisher),
		models.FactorWhere.Domain.EQ(data.Domain),
		models.FactorWhere.Device.EQ(data.Device),
		models.FactorWhere.Country.EQ(data.Country),
	).Exists(c.Context(), bcdb.DB())

	if err != nil {
		return false, err
	}

	if !exists {
		isInsert = true
	}

	if data.RuleId == "" {
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
		data.RuleId = factor.GetRuleID()
	}

	modConf := models.Factor{
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Device:        data.Device,
		Factor:        data.Factor,
		Country:       data.Country,
		PlacementType: null.StringFrom(data.PlacementType),
		Browser:       null.StringFrom(data.Browser),
		Os:            null.StringFrom(data.OS),
		RuleID:        data.RuleId,
	}

	fmt.Println("Upserting Factor with Rule ID:", data.RuleId)
	err = modConf.Upsert(c.Context(), bcdb.DB(), true,
		[]string{
			models.FactorColumns.Publisher,
			models.FactorColumns.Domain,
			models.FactorColumns.Device,
			models.FactorColumns.Country},
		boil.Infer(), boil.Infer())
	return isInsert, err
}
