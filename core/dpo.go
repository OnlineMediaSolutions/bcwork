package core

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
	"time"
)

var deleteQuery = `UPDATE dpo_rule
SET active = false
WHERE rule_id in (%s)`

type DemandPartnerOptimizationUpdateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DemandPartnerOptimizationUpdateRequest struct {
	DemandPartner string  `json:"demand_partner_id"`
	Publisher     string  `json:"publisher"`
	Domain        string  `json:"domain,omitempty"`
	Country       string  `json:"country,omitempty"`
	Browser       string  `json:"browser,omitempty"`
	OS            string  `json:"os,omitempty"`
	DeviceType    string  `json:"device_type,omitempty"`
	PlacementType string  `json:"placement_type,omitempty"`
	Factor        float64 `json:"factor"`
}

type Dpo struct {
	DemandPartnerID   string     `json:"demand_partner_id"  validate:"required"`
	IsInclude         bool       `json:"is_include"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DemandPartnerName string     `json:"demand_partner_name"`
	Active            bool       `json:"active"`
	Factor            float64    `json:"factor" validate:"required,factorDpo"`
	Country           string     `json:"country" validate:"required,country"`
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

func (dpos *DpoSlice) FromModel(slice models.DpoSlice) {

	for _, mod := range slice {
		dpo := Dpo{}
		dpo.FromModel(mod)
		*dpos = append(*dpos, &dpo)
	}

}

func CreateDeleteQuery(dpoRules []string) string {
	var wrappedStrings []string
	for _, ruleId := range dpoRules {
		wrappedStrings = append(wrappedStrings, fmt.Sprintf(`'%s'`, ruleId))
	}

	return fmt.Sprintf(deleteQuery, strings.Join(wrappedStrings, ","))
}
