package revenue

import (
	"bytes"
	"context"
	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/models"
	"github.com/valyala/fasttemplate"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"golang.org/x/text/message"
	"time"
)

func DailyHtmlReport(ctx context.Context, month string, db *sqlx.DB) (string, error) {

	records := make(models.RevenueDailySlice, 0)
	sql := `SELECT time,publisher_impressions,sold_impressions,dp_bid_requests,revenue,missed_opportunities,cost,demand_partner_fees,data_fee FROM revenue_daily where time>=$1 and time<$2 group by time order by time asc`

	from, err := time.Parse("20060102", month+"01")
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse 'from' date")
	}
	to := from.AddDate(0, 1, 0)
	err = queries.Raw(sql, from, to).Bind(ctx, db, &records)
	if err != nil {
		return "", errors.Wrapf(err, "failed to fetch summary records")
	}

	b := bytes.Buffer{}
	p := message.NewPrinter(message.MatchLanguage("en"))

	var totalImp int64
	var totalPubImp int64
	var totalSupRev float64
	var totalDemRev float64
	var totalRPM float64
	var totalSupRPM float64
	var totalDemRPM float64
	var totalDPBidRequest int64
	var totalDPFillRate float64
	var totalMissedOpp int64
	var totalDemandPartnerFee float64
	var totalDataFee float64

	for _, rec := range records {
		supplyRPM := float64(0)
		if rec.PublisherImpressions > 0 {
			supplyRPM = rec.Cost / float64(rec.PublisherImpressions) * float64(1000)
		}

		dpFillRate := float64(0)
		demandRPM := float64(0)
		if rec.SoldImpressions > 0 {
			demandRPM = rec.Revenue / float64(rec.SoldImpressions) * float64(1000)
			dpFillRate = float64(rec.SoldImpressions) / float64(rec.DPBidRequests) * float64(100)

		}

		rpm := float64(0)
		if rec.PublisherImpressions > 0 {
			rpm = rec.Revenue / float64(rec.PublisherImpressions) * float64(1000)
		}

		gp := rec.Revenue - rec.Cost - rec.DemandPartnerFees - rec.DataFee
		var gpp float64
		if rec.Revenue > 0 {
			gpp = (rec.Revenue - rec.Cost - rec.DemandPartnerFees - rec.DataFee) / rec.Revenue * 100
		}

		ratio := float64(0)
		loopRatio := float64(0)
		if rec.PublisherImpressions > 0 {
			ratio = float64(rec.SoldImpressions) / float64(rec.PublisherImpressions)
			loopRatio = (float64(rec.SoldImpressions) + float64(rec.MissedOpportunities)) / float64(rec.PublisherImpressions)

		}

		b.WriteString(p.Sprintf(rowDemand, rec.Time.Format("2006-01-02"), rec.PublisherImpressions, loopRatio, ratio, rec.SoldImpressions, dpFillRate, rec.Cost, supplyRPM, rec.Revenue, rpm, demandRPM, rec.DemandPartnerFees, rec.DataFee, gp, gpp))
		totalImp += rec.SoldImpressions
		totalPubImp += rec.PublisherImpressions
		totalSupRev += rec.Cost
		totalDemRev += rec.Revenue
		totalDataFee += rec.DataFee
		totalMissedOpp += rec.MissedOpportunities
		totalDemandPartnerFee += rec.DemandPartnerFees
		totalDPBidRequest += rec.DPBidRequests

	}

	if totalImp > 0 {
		totalDemRPM = totalDemRev / float64(totalImp) * 1000
	}

	if totalPubImp > 0 {
		totalSupRPM = totalSupRev / float64(totalPubImp) * 1000
	}

	if totalImp > 0 {
		totalRPM = totalDemRev / float64(totalPubImp) * 1000
		totalDPFillRate = float64(totalImp) / float64(totalDPBidRequest) * 100

	}

	gp := totalDemRev - totalSupRev - totalDemandPartnerFee - totalDataFee
	var gpp float64
	if totalDemRev > 0 {
		gpp = (totalDemRev - totalSupRev - totalDemandPartnerFee - totalDataFee) / totalDemRev * 100
	}

	totalRatio := float64(0)
	totalLoopRatio := float64(0)
	if totalPubImp > 0 {
		totalRatio = float64(totalImp) / float64(totalPubImp)
		totalLoopRatio = (float64(totalImp) + float64(totalMissedOpp)) / float64(totalPubImp)

	}

	t := fasttemplate.New(htmlDemandRevenueReport, "{{", "}}")
	s := t.ExecuteString(map[string]interface{}{
		"period": "Daily",
		"data":   b.String(),
		"totals": p.Sprintf(rowBoldDemand, "Total", totalPubImp, totalLoopRatio, totalRatio, totalImp, totalDPFillRate, totalSupRev, totalSupRPM, totalDemRev, totalRPM, totalDemRPM, totalDemandPartnerFee, totalDataFee, gp, gpp),
	})
	return s, nil
}

func HourlyHtmlReport(ctx context.Context, day string, db *sqlx.DB) (string, error) {

	records := make(models.RevenueDailySlice, 0)
	sql := `SELECT time,publisher_impressions,sold_impressions,dp_bid_requests,revenue,missed_opportunities,cost,demand_partner_fees,data_fee
             FROM revenue_hourly where time>=$1 and time<$2 group by time order by time asc`

	from, err := time.Parse("2006010215", day+"00")
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse 'from' date")
	}
	to := from.AddDate(0, 0, 1)
	err = queries.Raw(sql, from, to).Bind(ctx, db, &records)
	if err != nil {
		return "", errors.Wrapf(err, "failed to fetch summary records")
	}

	b := bytes.Buffer{}
	p := message.NewPrinter(message.MatchLanguage("en"))

	var totalImp int64
	var totalPubImp int64
	var totalSupRev float64
	var totalDemRev float64
	var totalRPM float64
	var totalSupRPM float64
	var totalDemRPM float64
	var totalDPBidRequest int64
	var totalDPFillRate float64
	var totalMissedOpp int64
	var totalDemandPartnerFee float64
	var totalDataFee float64

	for _, rec := range records {
		supplyRPM := float64(0)
		if rec.PublisherImpressions > 0 {
			supplyRPM = rec.Cost / float64(rec.PublisherImpressions) * float64(1000)
		}

		dpFillRate := float64(0)
		demandRPM := float64(0)
		if rec.SoldImpressions > 0 {
			demandRPM = rec.Revenue / float64(rec.SoldImpressions) * float64(1000)
			dpFillRate = float64(rec.SoldImpressions) / float64(rec.DPBidRequests) * float64(100)
		}

		rpm := float64(0)
		if rec.PublisherImpressions > 0 {
			rpm = rec.Revenue / float64(rec.PublisherImpressions) * float64(1000)
		}

		gp := rec.Revenue - rec.Cost - rec.DemandPartnerFees - rec.DataFee
		var gpp float64
		if rec.Revenue > 0 {
			gpp = (rec.Revenue - rec.Cost - rec.DemandPartnerFees - rec.DataFee) / rec.Revenue * 100
		}

		ratio := float64(0)
		loopRatio := float64(0)
		if rec.PublisherImpressions > 0 {
			ratio = float64(rec.SoldImpressions) / float64(rec.PublisherImpressions)
			loopRatio = (float64(rec.SoldImpressions) + float64(rec.MissedOpportunities)) / float64(rec.PublisherImpressions)

		}

		b.WriteString(p.Sprintf(rowDemand, rec.Time.Format("2006-01-02 15")+":00", rec.PublisherImpressions, loopRatio, ratio, rec.SoldImpressions, dpFillRate, rec.Cost, supplyRPM, rec.Revenue, rpm, demandRPM, rec.DemandPartnerFees, rec.DataFee, gp, gpp))
		totalImp += rec.SoldImpressions
		totalPubImp += rec.PublisherImpressions
		totalSupRev += rec.Cost
		totalDemRev += rec.Revenue
		totalMissedOpp += rec.MissedOpportunities
		totalDemandPartnerFee += rec.DemandPartnerFees
		totalDataFee += rec.DataFee
		totalDPBidRequest += rec.DPBidRequests

	}

	if totalImp > 0 {
		totalDemRPM = totalDemRev / float64(totalImp) * 1000
	}

	if totalPubImp > 0 {
		totalSupRPM = totalSupRev / float64(totalPubImp) * 1000
	}

	if totalImp > 0 {
		totalRPM = totalDemRev / float64(totalPubImp) * 1000
		totalDPFillRate = float64(totalImp) / float64(totalDPBidRequest) * 100
	}

	gp := totalDemRev - totalSupRev - totalDemandPartnerFee - totalDataFee
	var gpp float64
	if totalDemRev > 0 {
		gpp = (totalDemRev - totalSupRev - totalDemandPartnerFee - totalDataFee) / totalDemRev * 100
	}

	totalRatio := float64(0)
	totalLoopRatio := float64(0)
	if totalPubImp > 0 {
		totalRatio = float64(totalImp) / float64(totalPubImp)
		totalLoopRatio = (float64(totalImp) + float64(totalMissedOpp)) / float64(totalPubImp)

	}

	t := fasttemplate.New(htmlDemandRevenueReport, "{{", "}}")
	s := t.ExecuteString(map[string]interface{}{
		"period": "Hourly",
		"data":   b.String(),
		"totals": p.Sprintf(rowBoldDemand, "Total", totalPubImp, totalLoopRatio, totalRatio, totalImp, totalDPFillRate, totalSupRev, totalSupRPM, totalDemRev, totalRPM, totalDemRPM, totalDemandPartnerFee, totalDataFee, gp, gpp),
	})
	return s, nil
}

var htmlDemandRevenueReport = `
<html>
<head>
     <link href="https://unpkg.com/tailwindcss@^1.0/dist/tailwind.min.css" rel="stylesheet">
</head>
<body>
<div class="md:flex justify-center md:items-center">
   <div class="mt-1 flex md:mt-0 md:ml-4">
    <img class="filter invert h-40 w-40" src="https://onlinemediasolutions.com/wp-content/themes/brightcom/assets/images/oms-logo.svg" alt="">
  </div>
<div class="min-w-0">
    <h2 class="p-3 text-2xl font-bold leading-7 text-purple-600 sm:text-3xl sm:truncate">
      {{period}} Revenue Report 
    </h2>
  </div>
 
</div>


<div class="flex flex-col">
  <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
    <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
      <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="font-bold px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                Time
              </th>
                 <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  Pub Imps
               </th>
               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  Loop Ratio
               </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  Ratio
               </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  Sold Imps
               </th>
               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  DP Fillrate
               </th>
               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                 Cost
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  CPM
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  Revenue
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                RPM
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  DP RPM
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  DP Fee
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  Data Fee
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                 GP
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                 GP%
              </th> 
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
              {{data}}
             {{totals}}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>
</body>
</html>`

var rowDemand = `<tr>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %s
                 </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %d
                  </td>
                   <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %d
                  </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %.2f%%
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %.2f%%
                 </td>
                         
            </tr>`

var rowBoldDemand = `<tr class="font-bold">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %s
                 </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %d
                   </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %.2f
                 </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %d
                  </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    %.2f%%
                 </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    $%.2f
                 </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     $%.2f
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %.2f%%
                 </td>
            </tr>`
