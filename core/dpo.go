package core

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"time"
)

type Dpo struct {
	DemandPartnerID   string     `json:"demand_partner_id"`
	IsInclude         bool       `json:"is_include"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DemandPartnerName string     `json:"demand_partner_name"`
	Active            bool       `json:"active"`
}

type DpoSlice []*Dpo

type DPOGetOptions struct {
	Filter     DPOGetFilter           `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type DPOGetFilter struct {
	DemandPartnerId   filter.StringArrayFilter `json:"demand_partner_id,omitempty"`
	DemandPartnerName filter.StringArrayFilter `json:"demand_partner_name,omitempty"`
	Active            filter.StringArrayFilter `json:"active,omitempty"`
}

func (filter *DPOGetFilter) QueryMod() qmods.QueryModsSlice {

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

func GetDpos(ctx context.Context, ops *DPOGetOptions) (DpoSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.DpoColumns.DemandPartnerID).AddArray(ops.Pagination.Do())

	qmods = qmods.Add(qm.Select("DISTINCT *"))

	mods, err := models.Dpos(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve Dpos")
	}

	res := make(DpoSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (dpo *Dpo) FromModel(mod *models.Dpo) {

	dpo.DemandPartnerID = mod.DemandPartnerID
	dpo.IsInclude = mod.IsInclude
	dpo.CreatedAt = mod.CreatedAt
	dpo.UpdatedAt = mod.UpdatedAt.Ptr()
	dpo.DemandPartnerName = mod.DemandPartnerName.String
	dpo.Active = mod.Active

}

func ValidateFactorValue(c *fiber.Ctx, factorStr string) (error, bool, float64) {
	if factorStr == "" {
		c.SendString("'Factor' (factor is mandatory (0-100)")
		return c.SendStatus(http.StatusBadRequest), true, 0
	}

	factor, err := strconv.ParseFloat(factorStr, 64)
	if err != nil {
		c.SendString("'Factor' should be numeric (0-100)")
		return c.SendStatus(http.StatusBadRequest), true, 0
	}
	if factor > 100 || factor < 0 {
		c.SendString("'Factor' should be numeric (0-100)")
		return c.SendStatus(http.StatusBadRequest), true, 0
	}

	return nil, false, factor
}

func (dpos *DpoSlice) FromModel(slice models.DpoSlice) {

	for _, mod := range slice {
		dpo := Dpo{}
		dpo.FromModel(mod)
		*dpos = append(*dpos, &dpo)
	}

}
