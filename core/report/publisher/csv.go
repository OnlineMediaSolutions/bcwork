package publisher

import (
	"context"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

type PublisherReportOptions struct {
	WithHeaders     bool      `json:"with_headers"`
	WithTotals      bool      `json:"with_headers"`
	PublisherFilter []string  `json:"publisher_filter"`
	FromTime        time.Time `json:"from_time"`
	ToTime          time.Time `json:"to_time"`
}

func PublisherReportHourlyCSV(ctx context.Context, ops PublisherReportOptions) ([][]interface{}, error) {

	boil.DebugMode = true
	publisherMod := qm.And("TRUE")
	if len(ops.PublisherFilter) > 0 {
		publisherMod = models.PublisherHourlyWhere.PublisherID.IN(ops.PublisherFilter)
	}

	records, err := models.PublisherHourlies(
		models.PublisherHourlyWhere.Time.GTE(ops.FromTime),
		models.PublisherHourlyWhere.Time.LT(ops.ToTime),
		publisherMod,
		qm.OrderBy(models.PublisherHourlyColumns.Time+" DESC,"+models.PublisherHourlyColumns.SupplyTotal+" DESC")).All(ctx, bcdb.DB())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch records for hourly publisher report")
	}

	var totalBidRequests int64
	var totalBidResponses int64
	var totalImp int64
	var totalPubImp int64
	var totalSupRev float64
	var totalDemRev float64
	var totalRPM float64
	var totalCPM float64
	var totalDemRPM float64
	var totalProfit float64
	var totalResponsePriceSum float64
	var totalResponsePriceCount int64
	var totalMissedOpp int64
	var totalDemandPartnerFee float64

	headers := []interface{}{
		"time",
		"publisher_id",
		"domain",
		"os",
		"country",
		"device_type",
		"bid_requests",
		"bid_responses",
		"bid_rate",
		"avg_bid_price",
		"pub_imps",
		"sold_imps",
		"cpm",
		"cost",
		"dp_rpm",
		"rpm",
		"revenue",
		"dp_fee",
		"gp",
		"gp_pct",
		"ratio",
		"looping_ratio",
	}

	res := [][]interface{}{}
	if ops.WithHeaders {
		res = append(res, headers)
	}

	for _, rec := range records {
		bidRate := float64(0)
		if rec.BidRequests > 0 {
			bidRate = float64(rec.BidResponses) / float64(rec.BidRequests)
		}

		cpm := float64(0)
		if rec.BidPriceCount > 0 {
			cpm = rec.SupplyTotal / float64(rec.PublisherImpressions) * 1000
		}

		//supplyRPM := float64(0)
		//if rec.PublisherImpressions > 0 {
		//	supplyRPM = rec.SupplyTotal / float64(rec.PublisherImpressions) * float64(1000)
		//}

		demandRPM := float64(0)
		if rec.DemandImpressions > 0 {
			demandRPM = rec.DemandTotal / float64(rec.DemandImpressions) * float64(1000)
		}

		rpm := float64(0)
		if rec.PublisherImpressions > 0 {
			rpm = rec.DemandTotal / float64(rec.PublisherImpressions) * float64(1000)
		}

		//responseRate := float64(0)
		//if rec.BidRequests > 0 {
		//	responseRate = float64(rec.BidResponses) / float64(rec.BidRequests) * 100
		//}
		//
		//fillRate := float64(0)
		//if rec.BidRequests > 0 {
		//	fillRate = float64(rec.PublisherImpressions) / float64(rec.BidRequests) * 100
		//}

		bidPrice := float64(0)
		if rec.BidPriceCount > 0 {
			bidPrice = float64(rec.BidPriceTotal) / float64(rec.BidPriceCount)
		}

		gp := rec.DemandTotal - rec.SupplyTotal
		var gpp float64
		if rec.SupplyTotal > 0 {
			gpp = (rec.DemandTotal - rec.SupplyTotal) / rec.DemandTotal * 100
		}

		ratio := float64(0)
		loopRatio := float64(0)
		if rec.PublisherImpressions > 0 {
			ratio = float64(rec.DemandImpressions) / float64(rec.PublisherImpressions)
			loopRatio = (float64(rec.DemandImpressions) + float64(rec.MissedOpportunities)) / float64(rec.PublisherImpressions)

		}

		res = append(res, []interface{}{
			rec.Time,
			rec.PublisherID,
			rec.Domain,
			rec.Os,
			rec.Country,
			rec.DeviceType,
			rec.BidRequests,
			rec.BidResponses,
			bidRate,
			bidPrice,
			rec.PublisherImpressions,
			rec.DemandImpressions,
			cpm,
			rec.SupplyTotal,
			demandRPM,
			rpm,
			rec.DemandTotal,
			rec.DemandPartnerFee,
			gp,
			gpp,
			ratio,
			loopRatio,
		})

		totalBidRequests += rec.BidRequests
		totalBidResponses += rec.BidResponses
		totalImp += rec.DemandImpressions
		totalPubImp += rec.PublisherImpressions
		totalSupRev += rec.SupplyTotal
		totalDemRev += rec.SupplyTotal
		totalProfit += (rec.SupplyTotal - rec.SupplyTotal)
		totalResponsePriceSum += rec.BidPriceTotal
		totalResponsePriceCount += rec.BidPriceCount
		totalMissedOpp += rec.MissedOpportunities
		totalDemandPartnerFee += rec.DemandPartnerFee

	}

	if totalImp > 0 {
		totalDemRPM = totalDemRev / float64(totalImp) * 1000
	}

	if totalPubImp > 0 {
		totalCPM = totalSupRev / float64(totalPubImp) * 1000
	}

	totalBidRate := float64(0)
	if totalBidRequests > 0 {
		totalBidRate = float64(totalBidResponses) / float64(totalBidRequests)
	}

	totalBidPrice := float64(0)
	if totalResponsePriceCount > 0 {
		totalBidPrice = float64(totalResponsePriceSum) / float64(totalResponsePriceCount)
	}

	var gpp float64
	if totalDemRev > 0 {
		gpp = (totalDemRev - totalSupRev) / totalDemRev * 100
	}

	totalRatio := float64(0)
	totalLoopRatio := float64(0)
	if totalPubImp > 0 {
		totalRatio = float64(totalImp) / float64(totalPubImp)
		totalLoopRatio = (float64(totalImp) + float64(totalMissedOpp)) / float64(totalPubImp)

	}

	if ops.WithTotals {
		res = append(res, []interface{}{
			"total",
			"total",
			"total",
			"total",
			"total",
			"total",
			totalBidRequests,
			totalBidResponses,
			totalBidPrice,
			totalBidRate,
			totalPubImp,
			totalImp,
			totalCPM,
			0,
			totalDemRPM,
			totalRPM,
			totalDemRev,
			totalDemandPartnerFee,
			totalProfit,
			gpp,
			totalRatio,
			totalLoopRatio,
		})
	}

	return res, nil
}

func PublisherReportDailyCSV(ctx context.Context, ops PublisherReportOptions) ([][]interface{}, error) {

	boil.DebugMode = true
	publisherMod := qm.And("TRUE")
	if len(ops.PublisherFilter) > 0 {
		publisherMod = models.PublisherHourlyWhere.PublisherID.IN(ops.PublisherFilter)
	}

	records, err := models.PublisherDailies(
		models.PublisherDailyWhere.Time.GTE(ops.FromTime),
		models.PublisherDailyWhere.Time.LT(ops.ToTime),
		publisherMod,
		qm.OrderBy(models.PublisherDailyColumns.Time+" DESC,"+models.PublisherDailyColumns.SupplyTotal+" DESC")).All(ctx, bcdb.DB())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch records for hourly publisher report")
	}

	var totalBidRequests int64
	var totalBidResponses int64
	var totalImp int64
	var totalPubImp int64
	var totalSupRev float64
	var totalDemRev float64
	var totalRPM float64
	var totalCPM float64
	var totalDemRPM float64
	var totalProfit float64
	var totalResponsePriceSum float64
	var totalResponsePriceCount int64
	var totalMissedOpp int64
	var totalDemandPartnerFee float64

	headers := []interface{}{
		"time",
		"publisher_id",
		"domain",
		"os",
		"country",
		"device_type",
		"bid_requests",
		"bid_responses",
		"bid_rate",
		"avg_bid_price",
		"pub_imps",
		"sold_imps",
		"cpm",
		"cost",
		"dp_rpm",
		"rpm",
		"revenue",
		"dp_fee",
		"gp",
		"gp_pct",
		"ratio",
		"looping_ratio",
	}

	res := [][]interface{}{}
	if ops.WithHeaders {
		res = append(res, headers)
	}

	for _, rec := range records {
		bidRate := float64(0)
		if rec.BidRequests > 0 {
			bidRate = float64(rec.BidResponses) / float64(rec.BidRequests)
		}

		cpm := float64(0)
		if rec.BidPriceCount > 0 {
			cpm = rec.SupplyTotal / float64(rec.PublisherImpressions) * 1000
		}

		//supplyRPM := float64(0)
		//if rec.PublisherImpressions > 0 {
		//	supplyRPM = rec.SupplyTotal / float64(rec.PublisherImpressions) * float64(1000)
		//}

		demandRPM := float64(0)
		if rec.DemandImpressions > 0 {
			demandRPM = rec.DemandTotal / float64(rec.DemandImpressions) * float64(1000)
		}

		rpm := float64(0)
		if rec.PublisherImpressions > 0 {
			rpm = rec.DemandTotal / float64(rec.PublisherImpressions) * float64(1000)
		}

		//responseRate := float64(0)
		//if rec.BidRequests > 0 {
		//	responseRate = float64(rec.BidResponses) / float64(rec.BidRequests) * 100
		//}
		//
		//fillRate := float64(0)
		//if rec.BidRequests > 0 {
		//	fillRate = float64(rec.PublisherImpressions) / float64(rec.BidRequests) * 100
		//}

		bidPrice := float64(0)
		if rec.BidPriceCount > 0 {
			bidPrice = float64(rec.BidPriceTotal) / float64(rec.BidPriceCount)
		}

		gp := rec.DemandTotal - rec.SupplyTotal
		var gpp float64
		if rec.SupplyTotal > 0 {
			gpp = (rec.DemandTotal - rec.SupplyTotal) / rec.DemandTotal * 100
		}

		ratio := float64(0)
		loopRatio := float64(0)
		if rec.PublisherImpressions > 0 {
			ratio = float64(rec.DemandImpressions) / float64(rec.PublisherImpressions)
			loopRatio = (float64(rec.DemandImpressions) + float64(rec.MissedOpportunities)) / float64(rec.PublisherImpressions)

		}

		res = append(res, []interface{}{
			rec.Time,
			rec.PublisherID,
			rec.Domain,
			rec.Os,
			rec.Country,
			rec.DeviceType,
			rec.BidRequests,
			rec.BidResponses,
			bidRate,
			bidPrice,
			rec.PublisherImpressions,
			rec.DemandImpressions,
			cpm,
			rec.SupplyTotal,
			demandRPM,
			rpm,
			rec.DemandTotal,
			rec.DemandPartnerFee,
			gp,
			gpp,
			ratio,
			loopRatio,
		})

		totalBidRequests += rec.BidRequests
		totalBidResponses += rec.BidResponses
		totalImp += rec.DemandImpressions
		totalPubImp += rec.PublisherImpressions
		totalSupRev += rec.SupplyTotal
		totalDemRev += rec.SupplyTotal
		totalProfit += (rec.SupplyTotal - rec.SupplyTotal)
		totalResponsePriceSum += rec.BidPriceTotal
		totalResponsePriceCount += rec.BidPriceCount
		totalMissedOpp += rec.MissedOpportunities
		totalDemandPartnerFee += rec.DemandPartnerFee

	}

	if totalImp > 0 {
		totalDemRPM = totalDemRev / float64(totalImp) * 1000
	}

	if totalPubImp > 0 {
		totalCPM = totalSupRev / float64(totalPubImp) * 1000
	}

	if totalPubImp > 0 {
		totalRPM = totalDemRev / float64(totalPubImp) * 1000
	}

	totalBidRate := float64(0)
	if totalBidRequests > 0 {
		totalBidRate = float64(totalBidResponses) / float64(totalBidRequests)
	}

	totalBidPrice := float64(0)
	if totalResponsePriceCount > 0 {
		totalBidPrice = float64(totalResponsePriceSum) / float64(totalResponsePriceCount)
	}

	var gpp float64
	if totalDemRev > 0 {
		gpp = (totalDemRev - totalSupRev) / totalDemRev * 100
	}

	totalRatio := float64(0)
	totalLoopRatio := float64(0)
	if totalPubImp > 0 {
		totalRatio = float64(totalImp) / float64(totalPubImp)
		totalLoopRatio = (float64(totalImp) + float64(totalMissedOpp)) / float64(totalPubImp)

	}

	if ops.WithTotals {
		res = append(res, []interface{}{
			"total",
			"total",
			"total",
			"total",
			"total",
			"total",
			totalBidRequests,
			totalBidResponses,
			totalBidPrice,
			totalBidRate,
			totalPubImp,
			totalImp,
			totalCPM,
			totalSupRev,
			totalDemRPM,
			totalRPM,
			totalDemRev,
			totalDemandPartnerFee,
			totalProfit,
			gpp,
			totalRatio,
			totalLoopRatio,
		})
	}

	return res, nil
}
