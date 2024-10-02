package core

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

type MetadataUpdateRequest struct {
	Key     string      `json:"key"`
	Version string      `json:"version"`
	Data    interface{} `json:"data"`
}

func InsertDataToMetaData(c *fiber.Ctx, data MetadataUpdateRequest, value []byte, now time.Time) error {
	mod := models.MetadataQueue{
		TransactionID: bcguid.NewFromf(data.Key, now),
		Key:           data.Key,
		Value:         value,
		CreatedAt:     now,
	}

	err := mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to insert metadata update to queue")
		return fmt.Errorf("failed to insert metadata update to queue", err)
	}

	return nil
}

type AdsTxtSlice []*AdsTxt

func GetAdsTxt(ctx *fasthttp.RequestCtx, ops *GetAdsTxtOptions) (AdsTxtSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.AdsTXTColumns.PublisherID).AddArray(ops.Pagination.Do())

	mods, err := models.AdsTXTS(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve AdsTxt")
	}

	res := make(AdsTxtSlice, 0)
	res.FromModel(mods)

	return res, nil
}

type GetAdsTxtOptions struct {
	Filter     AdsTxtFilter           `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type AdsTxtFilter struct {
	Publisher filter.StringArrayFilter `json:"publisher_id,omitempty"`
	Domain    filter.StringArrayFilter `json:"domain,omitempty"`
	Demand    filter.StringArrayFilter `json:"demand,omitempty"`
	Active    filter.StringArrayFilter `json:"active,omitempty"`
}

type AdsTxt struct {
	PublisherID string     `boil:"publisher_id" json:"publisher_id,omitempty" toml:"publisher_id" yaml:"publisher_id"`
	Domain      *string    `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Demand      *string    `boil:"demand" json:"demand,omitempty" toml:"demand" yaml:"demand"`
	Active      *bool      `boil:"active" json:"active,omitempty" toml:"active" yaml:"active"`
	CreatedAt   *time.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

func (filter *AdsTxtFilter) QueryMod() qmods.QueryModsSlice {

	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.AdsTXTColumns.PublisherID))
	}

	if len(filter.Demand) > 0 {
		mods = append(mods, filter.Demand.AndIn(models.AdsTXTColumns.DemandPartnerName))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.AdsTXTColumns.Domain))
	}

	if len(filter.Active) > 0 {
		mods = append(mods, filter.Active.AndIn(models.AdsTXTColumns.Active))
	}

	return mods
}

func (cs *AdsTxtSlice) FromModel(slice models.AdsTXTSlice) error {

	for _, mod := range slice {
		c := AdsTxt{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (pixalate *AdsTxt) FromModel(mod *models.AdsTXT) error {

	pixalate.PublisherID = mod.PublisherID
	pixalate.CreatedAt = &mod.CreatedAt
	pixalate.UpdatedAt = mod.UpdatedAt.Ptr()
	pixalate.Domain = &mod.Domain
	pixalate.Demand = &mod.DemandPartnerName
	pixalate.Active = &mod.Active

	return nil
}
