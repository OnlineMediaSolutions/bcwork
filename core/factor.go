package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

type FactorUpdateRequest struct {
	Publisher string  `json:"publisher"`
	Domain    string  `json:"domain"`
	Device    string  `json:"device"`
	Factor    float64 `json:"factor"`
	Country   string  `json:"country"`
}

type FactorRealtimeRecord struct {
	Rule     string  `json:"rule"`
	Factor   float64 `json:"factor"`
	FactorID string  `json:"factor_id"`
}

type FactorSlice []*FactorUpdateRequest

func (fs *FactorSlice) FromModel(slice models.FactorSlice) error {
	for _, mod := range slice {
		factor := &FactorUpdateRequest{
			Publisher: mod.Publisher,
			Domain:    mod.Domain,
			Device:    mod.Device,
			Factor:    mod.Factor,
			Country:   mod.Country,
		}
		*fs = append(*fs, factor)
	}
	return nil
}

type Factor struct {
	Publisher string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain    string  `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country   string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device    string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Factor    float64 `boil:"factor" json:"factor,omitempty" toml:"factor" yaml:"factor,omitempty"`
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

func (factor *FactorUpdateRequest) FromModel(mod *models.Factor) error {
	factor.Publisher = mod.Publisher
	factor.Domain = mod.Domain
	factor.Device = mod.Device
	factor.Factor = mod.Factor
	factor.Country = mod.Country
	return nil
}

func (factor *FactorUpdateRequest) ToRtRule() *FactorRealtimeRecord {
	return &FactorRealtimeRecord{
		Rule:     utils.GetFormulaRegex(factor.Country, factor.Domain, factor.Device),
		Factor:   factor.Factor,
		FactorID: factor.Publisher,
	}
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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to parse hash value for factor metadata")
	}

	err = SendFactorToRT(context.Background(), *data)
	if err != nil {
		return err
	}
	return nil
}

func UpdateFactor(c *fiber.Ctx, data *FactorUpdateRequest) error {

	modConf := models.Factor{
		Publisher: data.Publisher,
		Domain:    data.Domain,
		Device:    data.Device,
		Factor:    data.Factor,
		Country:   data.Country,
	}

	return modConf.Upsert(c.Context(), bcdb.DB(), true, []string{models.FactorColumns.Publisher, models.FactorColumns.Domain, models.FactorColumns.Device, models.FactorColumns.Country}, boil.Infer(), boil.Infer())
}

func (fs *FactorSlice) FromModelFactor(modFactors *models.Factor) {
	*fs = append(*fs, &FactorUpdateRequest{
		Publisher: modFactors.Publisher,
		Domain:    modFactors.Domain,
		Device:    modFactors.Device,
		Factor:    modFactors.Factor,
		Country:   modFactors.Country,
	})
}

func SendFactorToRT(c context.Context, updateRequest FactorUpdateRequest) error {

	modFactor, err := factorQuery(c, updateRequest)

	if err != nil && err != sql.ErrNoRows {
		return eris.Wrapf(err, "failed to fetch factors")
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

	key := getMetadataKey(updateRequest)
	metadataKey := utils.CreateMetadataKey(key, "price:factor:v2")
	metadataValue := CreateMetadataValue(updateRequest, metadataKey, value)

	err = metadataValue.Insert(c, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert metadata record")
	}

	return nil
}

func getMetadataKey(updateRequest FactorUpdateRequest) utils.MetadataKey {
	key := utils.MetadataKey{
		Publisher: updateRequest.Publisher,
		Domain:    updateRequest.Domain,
		Device:    updateRequest.Device,
		Country:   updateRequest.Country,
	}
	return key
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
				Rule:     utils.GetFormulaRegex(factor.Country, factor.Domain, factor.Device),
				Factor:   factor.Factor,
				FactorID: factor.Publisher,
			}
			finalRules = append(finalRules, rule)
		}
	}

	newRule := FactorRealtimeRecord{
		Rule:     utils.GetFormulaRegex(updateRequest.Country, updateRequest.Domain, updateRequest.Device),
		Factor:   updateRequest.Factor,
		FactorID: updateRequest.Publisher,
	}
	finalRules = append(finalRules, newRule)
	return finalRules
}

func CreateMetadataValue(updateRequest FactorUpdateRequest, key string, b []byte) models.MetadataQueue {
	modMeta := models.MetadataQueue{
		TransactionID: bcguid.NewFromf(updateRequest.Publisher, updateRequest.Domain, time.Now()),
		Key:           key,
		Value:         b,
	}
	return modMeta
}
