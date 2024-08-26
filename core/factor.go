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
	"strconv"
	"time"
)

type FactorUpdateRequest struct {
	Publisher string  `json:"publisher"`
	Domain    string  `json:"domain"`
	Device    string  `json:"device"`
	Factor    float64 `json:"factor"`
	Country   string  `json:"country"`
}

type Factor struct {
	Publisher string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain    string  `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country   string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device    string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Factor    float64 `boil:"factor" json:"factor,omitempty" toml:"factor" yaml:"factor,omitempty"`
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

func UpdateMetaData(c *fiber.Ctx, data *FactorUpdateRequest) error {
	_, err := json.Marshal(data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to parse hash value for factor metadata", err)
	}

	metadataKey := utils.MetadataKey{
		Publisher: data.Publisher,
		Domain:    data.Domain,
		Device:    data.Device,
		Country:   data.Country,
	}

	key := utils.CreateMetadataKey(metadataKey, "price:factor")

	factor := strconv.FormatFloat(data.Factor, 'f', 2, 64)
	mod := models.MetadataQueue{
		Key:           key,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         []byte(factor),
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to insert metadata update to queue", err)
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
