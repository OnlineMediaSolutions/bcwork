package core

import (
	"context"
	"database/sql"
	"errors"
	"github.com/m6yf/bcwork/bcdb/qmods"

	"github.com/m6yf/bcwork/dto"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DpAPIService struct {
	historyModule history.HistoryModule
}

func NewDPAPIService(historyModule history.HistoryModule) *DpAPIService {
	return &DpAPIService{
		historyModule: historyModule,
	}
}

type GetDPApiOptions struct {
	Filter     DPApiFilter            `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type DPApiFilter struct {
	DemandPartner filter.StringArrayFilter `json:"demand_partner,omitempty"`
	Domain        filter.StringArrayFilter `json:"domain,omitempty"`
	SoldImps      filter.StringArrayFilter `json:"sold_imps,omitempty"`
	Revenue       filter.StringArrayFilter `json:"revenue,omitempty"`
}

func (d *DpAPIService) GetDpApiReport(ctx context.Context, ops *GetDPApiOptions) (dto.DpApiSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.DPAPIReportColumns.DemandPartner).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.DPAPIReports(qmods...).All(ctx, bcdb.DB())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrap(err, "failed to retrieve factors")
	}

	res := make(dto.DpApiSlice, 0)
	err = res.FromModel(mods)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (filter *DPApiFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.DemandPartner) > 0 {
		mods = append(mods, filter.DemandPartner.AndIn(models.DPAPIReportColumns.DemandPartner))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.DPAPIReportColumns.Domain))
	}

	return mods
}
