package revenue

var dwhHourlyUpdate = `INSERT INTO revenue_hourly (time, publisher_impressions, sold_impressions, cost, revenue, demand_partner_fees, missed_opportunities)
SELECT time,sum(publisher_impressions),sum(sold_impressions),sum(cost),sum(revenue),sum(demand_partner_fee),sum(missed_opportunities) FROM supply_hourly WHERE time >='%s' GROUP BY time
ON CONFLICT (time)
DO UPDATE SET publisher_impressions=EXCLUDED.publisher_impressions,
              sold_impressions=EXCLUDED.sold_impressions,
              cost=EXCLUDED.cost,
              revenue=EXCLUDED.revenue,
              demand_partner_fees=EXCLUDED.demand_partner_fees,
              missed_opportunities=EXCLUDED.missed_opportunities`

var dwhHourlyUpdateDpBidRequest = `INSERT INTO revenue_hourly (time, dp_bid_requests)
SELECT time,sum(bid_requests) dp_bid_requests FROM demand_hourly WHERE time >='%s' GROUP BY time
ON CONFLICT (time)
DO UPDATE SET dp_bid_requests=EXCLUDED.dp_bid_requests`
