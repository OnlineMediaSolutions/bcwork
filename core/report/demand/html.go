package demand

//
//import (
//	"bytes"
//	"context"
//	"github.com/friendsofgo/errors"
//	"github.com/m6yf/bcwork/bcdb"
//	"github.com/m6yf/bcwork/models"
//	"github.com/valyala/fasttemplate"
//	"github.com/volatiletech/sqlboiler/v4/boil"
//	"github.com/volatiletech/sqlboiler/v4/queries/qm"
//	"golang.org/x/text/message"
//	"time"
//)
//
//func DemandReportHourly(ctx context.Context, from time.Time, to time.Time, publisherIds []string, demandPartnerIds []string, full bool, domains []string) (string, error) {
//
//	boil.DebugMode = true
//	publisherMod := qm.And("TRUE")
//	if len(publisherIds) > 0 {
//		publisherMod = models.DemandHourlyWhere.PublisherID.IN(publisherIds)
//	}
//
//	demandPartnerMod := qm.And("TRUE")
//	if len(demandPartnerIds) > 0 {
//		demandPartnerMod = models.DemandHourlyWhere.DemandPartnerID.IN(demandPartnerIds)
//	}
//
//	domainMod := qm.And("TRUE")
//	if len(domains) > 0 {
//		domainMod = models.DemandHourlyWhere.Domain.IN(domains)
//	}
//
//	filterZeroMod := qm.And("TRUE")
//	if !full {
//		filterZeroMod = models.DemandHourlyWhere.BidResponse.GT(0)
//	}
//
//	records, err := models.DemandHourlies(
//		models.DemandHourlyWhere.Time.GTE(from),
//		models.DemandHourlyWhere.Time.LT(to),
//		publisherMod,
//		demandPartnerMod,
//		filterZeroMod,
//		domainMod,
//		qm.OrderBy(models.DemandHourlyColumns.Time+" DESC,"+models.DemandHourlyColumns.Impressions+" DESC")).All(ctx, bcdb.DB())
//	if err != nil {
//		return "", errors.Wrapf(err, "failed to fetch records")
//	}
//
//	b := bytes.Buffer{}
//	p := message.NewPrinter(message.MatchLanguage("en"))
//
//	var totalBidRequests int64
//	var totalBidResponses int64
//	var totalImp int64
//	var totalPubImp int64
//	var totalSupRev float64
//	var totalDemRev float64
//	var totalRPM float64
//	var totalSupRPM float64
//	var totalDemRPM float64
//	var totalProfit float64
//	var totalResponsePriceSum float64
//	var totalResponsePriceCount int64
//	var totalDemandPartnerFee float64
//
//	for _, rec := range records {
//		supplyRPM := float64(0)
//		if rec.SupplyImpressions > 0 {
//			supplyRPM = rec.SupplyTotal / float64(rec.SupplyImpressions) * float64(1000)
//		}
//
//		demandRPM := float64(0)
//		if rec.Impressions > 0 {
//			demandRPM = rec.DemandTotal / float64(rec.Impressions) * float64(1000)
//		}
//
//		rpm := float64(0)
//		if rec.SupplyImpressions > 0 {
//			rpm = rec.DemandTotal / float64(rec.SupplyImpressions) * float64(1000)
//		}
//
//		responseRate := float64(0)
//		if rec.BidRequests > 0 {
//			responseRate = float64(rec.BidResponses) / float64(rec.BidRequests) * 100
//		}
//
//		fillRate := float64(0)
//		if rec.BidRequests > 0 {
//			fillRate = float64(rec.Impressions) / float64(rec.BidRequests) * 100
//		}
//		bidPrice := float64(0)
//		if rec.DemandResponsePriceCount > 0 {
//			bidPrice = float64(rec.DemandResponsePriceSum) / float64(rec.DemandResponsePriceCount)
//		}
//
//		gp := rec.DemandTotal - rec.SupplyTotal
//		var gpp float64
//		if rec.DemandTotal > 0 {
//			gpp = (rec.DemandTotal - rec.SupplyTotal) / rec.DemandTotal * 100
//		}
//
//		loopRatio := float64(0)
//		if rec.SupplyImpressions > 0 {
//			loopRatio = float64(rec.Impressions) / float64(rec.SupplyImpressions)
//		}
//
//		b.WriteString(p.Sprintf(rowDemand, rec.Time.Format("2006-01-02 15:00"), rec.DemandPartnerID, rec.PublisherID, rec.Domain, rec.BidRequests, rec.BidResponses, responseRate, bidPrice, rec.Impressions, rec.SupplyImpressions, loopRatio, fillRate, rpm, demandRPM, supplyRPM, rec.DemandTotal, rec.SupplyTotal, rec.DemandPartnerFee, gp, gpp))
//		totalBidRequests += rec.BidRequests
//		totalBidResponses += rec.BidResponses
//		totalImp += rec.Impressions
//		totalPubImp += rec.SupplyImpressions
//		totalSupRev += rec.SupplyTotal
//		totalDemRev += rec.DemandTotal
//		totalResponsePriceSum += rec.DemandResponsePriceSum
//		totalResponsePriceCount += rec.DemandResponsePriceCount
//		totalDemandPartnerFee += rec.DemandPartnerFee
//		totalProfit += (rec.DemandTotal - rec.SupplyTotal - rec.DemandPartnerFee)
//
//	}
//
//	if totalPubImp > 0 {
//		totalSupRPM = totalSupRev / float64(totalPubImp) * 1000
//	}
//
//	if totalImp > 0 {
//		totalDemRPM = totalDemRev / float64(totalImp) * 1000
//		totalRPM = totalDemRev / float64(totalPubImp) * 1000
//	}
//
//	bidPrice := float64(0)
//	if totalResponsePriceCount > 0 {
//		bidPrice = float64(totalResponsePriceSum) / float64(totalResponsePriceCount)
//	}
//
//	gp := totalDemRev - totalSupRev - totalDemandPartnerFee
//	var gpp float64
//	if totalDemRev > 0 {
//		gpp = gp / totalDemRev * 100
//	}
//
//	totalLoopRatio := float64(0)
//	if totalSupRev > 0 {
//		totalLoopRatio = float64(totalImp) / float64(totalPubImp)
//	}
//
//	t := fasttemplate.New(htmlDemand, "{{", "}}")
//	s := t.ExecuteString(map[string]interface{}{
//		"period": "Hourly",
//		"data":   b.String(),
//		"totals": p.Sprintf(rowBoldDemand, "Total", "", "", "", totalBidRequests, totalBidResponses, float64(totalBidResponses)/float64(totalBidRequests)*100, bidPrice, totalImp, totalPubImp, totalLoopRatio, float64(totalImp)/float64(totalBidRequests)*100, totalRPM, totalDemRPM, totalSupRPM, totalDemRev, totalSupRev, totalDemandPartnerFee, gp, gpp),
//	})
//	return s, nil
//}
//
//func DemandReportDaily(ctx context.Context, from time.Time, to time.Time, publisherIds []string, demandPartnerIds []string, full bool, domains []string) (string, error) {
//
//	boil.DebugMode = true
//	publisherMod := qm.And("TRUE")
//	if len(publisherIds) > 0 {
//		publisherMod = models.DemandDailyWhere.PublisherID.IN(publisherIds)
//	}
//
//	demandPartnerMod := qm.And("TRUE")
//	if len(demandPartnerIds) > 0 {
//		demandPartnerMod = models.DemandDailyWhere.DemandPartnerID.IN(demandPartnerIds)
//	}
//
//	domainMod := qm.And("TRUE")
//	if len(domains) > 0 {
//		domainMod = models.DemandDailyWhere.Domain.IN(domains)
//	}
//
//	filterZeroMod := qm.And("TRUE")
//	if !full {
//		filterZeroMod = models.DemandDailyWhere.BidResponses.GT(0)
//	}
//
//	records, err := models.DemandDailies(
//		models.DemandDailyWhere.Time.GTE(from),
//		models.DemandDailyWhere.Time.LT(to),
//		publisherMod,
//		demandPartnerMod,
//		domainMod,
//		filterZeroMod,
//		qm.OrderBy(models.DemandDailyColumns.Time+" DESC,"+models.DemandDailyColumns.Impressions+" DESC")).All(ctx, bcdb.DB())
//	if err != nil {
//		return "", errors.Wrapf(err, "failed to fetch records")
//	}
//
//	b := bytes.Buffer{}
//	p := message.NewPrinter(message.MatchLanguage("en"))
//
//	var totalBidRequests int64
//	var totalBidResponses int64
//	var totalImp int64
//	var totalPubImp int64
//	var totalSupRev float64
//	var totalDemRev float64
//	var totalSupRPM float64
//	var totalRPM float64
//	var totalDemRPM float64
//	var totalProfit float64
//	var totalResponsePriceSum float64
//	var totalResponsePriceCount int64
//	var totalDemandPartnerFee float64
//
//	for _, rec := range records {
//		supplyRPM := float64(0)
//		if rec.SupplyImpressions > 0 {
//			supplyRPM = rec.SupplyTotal / float64(rec.SupplyImpressions) * float64(1000)
//		}
//
//		demandRPM := float64(0)
//		if rec.Impressions > 0 {
//			demandRPM = rec.DemandTotal / float64(rec.Impressions) * float64(1000)
//		}
//
//		rpm := float64(0)
//		if rec.SupplyImpressions > 0 {
//			rpm = rec.DemandTotal / float64(rec.SupplyImpressions) * float64(1000)
//		}
//
//		responseRate := float64(0)
//		if rec.BidRequests > 0 {
//			responseRate = float64(rec.BidResponses) / float64(rec.BidRequests) * 100
//		}
//
//		fillRate := float64(0)
//		if rec.BidRequests > 0 {
//			fillRate = float64(rec.Impressions) / float64(rec.BidRequests) * 100
//		}
//		bidPrice := float64(0)
//		if rec.DemandResponsePriceCount > 0 {
//			bidPrice = float64(rec.DemandResponsePriceSum) / float64(rec.DemandResponsePriceCount)
//		}
//
//		gp := rec.DemandTotal - rec.SupplyTotal - rec.DemandPartnerFee
//		var gpp float64
//		if rec.DemandTotal > 0 {
//			gpp = gp / rec.DemandTotal * 100
//		}
//
//		loopRatio := float64(0)
//		if rec.SupplyImpressions > 0 {
//			loopRatio = float64(rec.Impressions) / float64(rec.SupplyImpressions)
//
//		}
//
//		b.WriteString(p.Sprintf(rowDemand, rec.Time.Format("2006-01-02 15:00"), rec.DemandPartnerID, rec.PublisherID, rec.Domain, rec.BidRequests, rec.BidResponses, responseRate, bidPrice, rec.Impressions, rec.SupplyImpressions, loopRatio, fillRate, rpm, demandRPM, supplyRPM, rec.DemandTotal, rec.SupplyTotal, rec.DemandPartnerFee, gp, gpp))
//		totalBidRequests += rec.BidRequests
//		totalBidResponses += rec.BidResponses
//		totalImp += rec.Impressions
//		totalPubImp += rec.SupplyImpressions
//		totalSupRev += rec.SupplyTotal
//		totalDemRev += rec.DemandTotal
//		totalResponsePriceSum += rec.DemandResponsePriceSum
//		totalResponsePriceCount += rec.DemandResponsePriceCount
//		totalDemandPartnerFee += rec.DemandPartnerFee
//		totalProfit += (rec.DemandTotal - rec.SupplyTotal - rec.DemandPartnerFee)
//
//	}
//
//	if totalPubImp > 0 {
//		totalSupRPM = totalSupRev / float64(totalPubImp) * 1000
//		totalRPM = totalDemRev / float64(totalPubImp) * 1000
//
//	}
//
//	if totalImp > 0 {
//		totalDemRPM = totalDemRev / float64(totalImp) * 1000
//	}
//
//	bidPrice := float64(0)
//	if totalResponsePriceCount > 0 {
//		bidPrice = float64(totalResponsePriceSum) / float64(totalResponsePriceCount)
//	}
//
//	gp := totalDemRev - totalSupRev - totalDemandPartnerFee
//	var gpp float64
//	if totalDemRev > 0 {
//		gpp = gp / totalDemRev * 100
//	}
//
//	totalLoopRatio := float64(0)
//	if totalSupRev > 0 {
//		totalLoopRatio = float64(totalImp) / float64(totalPubImp)
//	}
//
//	t := fasttemplate.New(htmlDemand, "{{", "}}")
//	s := t.ExecuteString(map[string]interface{}{
//		"period": "Daily",
//		"data":   b.String(),
//		"totals": p.Sprintf(rowBoldDemand, "Total", "", "", "", totalBidRequests, totalBidResponses, float64(totalBidResponses)/float64(totalBidRequests)*100, bidPrice, totalImp, totalPubImp, totalLoopRatio, float64(totalImp)/float64(totalBidRequests)*100, totalRPM, totalDemRPM, totalSupRPM, totalDemRev, totalSupRev, totalDemandPartnerFee, gp, gpp),
//	})
//	return s, nil
//}
//
//var htmlDemand = `
//<html>
//<head>
//     <link href="https://unpkg.com/tailwindcss@^1.0/dist/tailwind.min.css" rel="stylesheet">
//</head>
//<body>
//<div class="md:flex justify-center md:items-center">
//   <div class="mt-1 flex md:mt-0 md:ml-4">
//    <img class="filter invert h-40 w-40" src="https://onlinemediasolutions.com/wp-content/themes/brightcom/assets/images/oms-logo.svg" alt="">
//  </div>
//<div class="min-w-0">
//    <h2 class="p-3 text-2xl font-bold leading-7 text-purple-600 sm:text-3xl sm:truncate">
//      {{period}} Demand Report
//    </h2>
//  </div>
//
//</div>
//
//
//<div class="flex flex-col">
//  <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
//    <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
//      <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
//        <table class="min-w-full divide-y divide-gray-200">
//          <thead class="bg-gray-50">
//            <tr>
//              <th scope="col" class="font-bold px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                Time
//              </th>
//              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                Demand Partner
//              </th>
//             <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                Publisher
//              </th>
//              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                Domain
//              </th>
//              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                Bid Requests
//               </th>
//               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                Bid Responses
//               </th>
//               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                 Responses Rate
//               </th>
//              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                Bid Price
//              </th>
//               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                  Impressions
//               </th>
//               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                  Pub Imps
//               </th>
//               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                  Ratio
//               </th>
//              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                 Fill Rate
//               </th>
//               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                RPM
//              </th>
//              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                  DP RPM
//              </th>
//               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                 Supply CPM
//              </th>
//               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                  Revenue
//              </th>
//               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                 Cost
//              </th>
//             <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                 DP Fee
//              </th>
//              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                 GP
//              </th>
//              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
//                 GP%
//              </th>
//            </tr>
//          </thead>
//          <tbody class="bg-white divide-y divide-gray-200">
//              {{data}}
//             {{totals}}
//          </tbody>
//        </table>
//      </div>
//    </div>
//  </div>
//</div>
//</body>
//</html>`
//
//var rowDemand = `<tr>
//                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %s
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %s
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %s
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %s
//                 </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %d
//                  </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %d
//                  </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %.2f%%
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %d
//                  </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %d
//                  </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %.2f%%
//                 </td>
//                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %.2f%%
//                 </td>
//            </tr>`
//
//var rowBoldDemand = `<tr class="font-bold">
//                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %s
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %s
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %s
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %s
//                 </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %d
//                  </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %d
//                  </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %.2f%%
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %d
//                   </td>
//                   <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %d
//                  </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %.2f%%
//                 </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     $%.2f
//                 </td>
//                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
//                     %.2f%%
//                 </td>
//            </tr>`
