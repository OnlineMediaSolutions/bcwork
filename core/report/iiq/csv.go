package iiq

import (
	"context"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

type IiqReportOptions struct {
	WithHeaders  bool      `json:"with_headers"`
	WithTotals   bool      `json:"with_headers"`
	DemandFilter []string  `json:"demand_filter"`
	FromTime     time.Time `json:"from_time"`
	ToTime       time.Time `json:"to_time"`
}

func IiqReportHourlyCSV(ctx context.Context, ops IiqReportOptions) ([][]interface{}, error) {

	boil.DebugMode = true
	demandMod := qm.And("TRUE")
	if len(ops.DemandFilter) > 0 {
		demandMod = models.IiqTestingWhere.DemandPartnerID.IN(ops.DemandFilter)
	}

	records, err := models.IiqTestings(
		models.IiqTestingWhere.Time.GTE(ops.FromTime),
		models.IiqTestingWhere.Time.LT(ops.ToTime),
		models.IiqTestingWhere.IiqRequests.GT(0),
		demandMod,
		qm.OrderBy(models.IiqTestingColumns.Time+" DESC,"+models.IiqTestingColumns.IiqImpressions+" DESC")).All(ctx, bcdb.DB())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch records for hourly iiq report")
	}

	var totalIiqRequests int64
	var totalNonIiqRequests int64
	var totalIiqImp int64
	var totalNonIiqImp int64
	var totalIiqRate float64
	var totalNonIiqRate float64
	var totalIiqLift float64

	headers := []interface{}{
		"time",
		"demand_partner_id",
		"iiq_requests",
		"non_iiq_requests",
		"iiq_impressions",
		"non_iiq_impressions",
		"iiq_fill_rate",
		"non_iiq_fill_rate",
		"iiq_lift",
	}

	res := [][]interface{}{}
	if ops.WithHeaders {
		res = append(res, headers)
	}

	for _, rec := range records {

		iiqFillRate := float64(0)
		if rec.IiqRequests > 0 {
			iiqFillRate = float64(rec.IiqImpressions) / float64(rec.IiqRequests)
		}

		nonNonIiqFillRate := float64(0)
		if rec.NonIiqRequests > 0 {
			nonNonIiqFillRate = float64(rec.NonIiqImpressions) / float64(rec.NonIiqRequests)
		}

		iiqLift := float64(0)
		if iiqFillRate > 0 {
			iiqLift = iiqFillRate / nonNonIiqFillRate
		}

		res = append(res, []interface{}{
			rec.Time,
			rec.DemandPartnerID,
			rec.IiqRequests,
			rec.NonIiqRequests,
			rec.IiqImpressions,
			rec.NonIiqImpressions,
			iiqFillRate,
			nonNonIiqFillRate,
			iiqLift,
		})

		totalIiqRequests += rec.IiqRequests
		totalNonIiqRequests += rec.NonIiqRequests
		totalIiqImp += rec.IiqImpressions
		totalNonIiqImp += rec.NonIiqImpressions

	}

	if totalIiqRequests > 0 {
		totalIiqRate = float64(totalIiqImp) / float64(totalIiqRequests)
	}

	if totalNonIiqRequests > 0 {
		totalNonIiqRate = float64(totalNonIiqImp) / float64(totalNonIiqRequests)
	}

	totalIiqLift = float64(0)
	if totalNonIiqRate > 0 {
		totalIiqLift = float64(totalIiqRate) / float64(totalNonIiqRate)
	}

	if ops.WithTotals {
		res = append(res, []interface{}{
			"total",
			"total",
			totalIiqRequests,
			totalNonIiqRequests,
			totalIiqImp,
			totalNonIiqImp,
			totalIiqRate,
			totalNonIiqRate,
			totalIiqLift,
		})
	}

	return res, nil
}
