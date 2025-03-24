package core

import (
	"context"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DpAPIService struct{}

type GetDPApiOptions struct {
	Filter        DPApiFilter            `json:"filter"`
	Pagination    *pagination.Pagination `json:"pagination"`
	Order         order.Sort             `json:"order"`
	StartDate     string                 `json:"startDate"`
	EndDate       string                 `json:"endDate"`
	DateStamp     string                 `json:"date_stamp"`
	DemandPartner string                 `json:"demand_partner"`
}

type DPApiFilter struct {
	DemandPartner filter.StringArrayFilter `json:"demand_partner,omitempty"`
	Domain        filter.StringArrayFilter `json:"domain,omitempty"`
}

func (d *DpAPIService) GetDpApiReport(ctx context.Context, ops *GetDPApiOptions) (dto.DpApiSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.DPAPIReportColumns.DemandPartner).
		AddArray(ops.Pagination.Do()).
		Add(
			qm.Select("demand_partner, TO_CHAR(date_stamp, 'YYYY-MM-DD') AS date, SUM(revenue) AS platformRevenue"),
			qm.Where("date_stamp >= ? AND date_stamp <= ?", ops.StartDate, ops.EndDate),
			qm.GroupBy("TO_CHAR(date_stamp, 'YYYY-MM-DD'), demand_partner"),
		)

	mods, err := models.DPAPIReports(qmods...).All(ctx, bcdb.DB())

	res := make(dto.DpApiSlice, 0)
	err = res.FromModel(mods)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *DpAPIService) GetDemandPartners(ctx context.Context) (dto.DpApiSlice, error) {
	query := `SELECT DISTINCT demand_partner,
              MAX(day) AS last_full_day,
              counter
				FROM (
					SELECT demand_partner,
						   TO_CHAR(date_stamp, 'YYYY-MM-DD') AS day,
						   COUNT(DISTINCT date_stamp) AS counter
					FROM dp_api_report
					GROUP BY demand_partner, day
				) AS t1
				WHERE counter = 24 OR counter = 1
				group by demand_partner,counter`

	var records dto.DpApiSlice
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &records)
	if err != nil {
		return nil, err
	}

	return records, nil
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
