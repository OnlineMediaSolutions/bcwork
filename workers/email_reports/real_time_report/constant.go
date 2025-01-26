package real_time_report

var ColumnNames = []string{
	"Time",
	"PublisherID",
	"Publisher",
	"Domain",
	"Device",
	"Country",
	"Fill Rate",
	"Bid Requests",
	"Publisher Impressions",
	"Sold Impressions",
	"Cost",
	"Revenue",
	"CPM",
	"RPM",
	"DP RPM",
	"GP",
	"GPP",
}

const RealTimeReportQuery = `
SELECT time,
  publisher,
  publisher_id,
  domain,
  country,
  device,
  bid_requests,
  revenue,
  cost,
    sold_impressions,
    publisher_impressions,
    pub_fill_rate,
    cpm,
    rpm,
    dp_rpm,
    gp,
    gpp,
    consultant_fee,
    tam_fee,
    tech_fee,
    demand_partner_fee,
    data_fee
FROM
  real_time_report
WHERE time >= '%s'
  AND time < '%s'
  AND device is not null
  AND country is not null
  AND publisher_id is not null
  AND domain is not null`

const QuestRequests = `
 SELECT DATE_TRUNC('day', to_timezone(timestamp, 'America/New_York')) time,
  pubid,
  domain,
  country,
  dtype,
  sum(count) bid_requests
FROM
  request_placement
WHERE to_timezone(timestamp, 'America/New_York') >= '%s'
  AND to_timezone(timestamp, 'America/New_York') < '%s'
  AND dtype is not null
  AND country is not null
  AND pubid is not null
  AND domain is not null
GROUP BY 1,2,3,4,5`

const QuestImpressions = `
     SELECT DATE_TRUNC('day', to_timezone(timestamp, 'America/New_York')) time,
       publisher pubid,
       domain,
       country,
       dtype,
       sum(dbpr)/1000 revenue,
       sum(sbpr)/1000 cost,
       count(1) sold_impressions,
       sum(CASE WHEN loop=false THEN 1 ELSE 0 END) publisher_impressions,      
       sum(dpfee)/1000 demand_partner_fee,
       sum(CASE WHEN uidsrc='iiq' THEN dbpr/1000 ELSE 0 END) * 0.045 data_fee
FROM impression
WHERE to_timezone(timestamp, 'America/New_York') >= '%s'
  AND to_timezone(timestamp, 'America/New_York') < '%s'
  AND publisher IS NOT NULL
  AND domain IS NOT NULL
  AND country IS NOT NULL
  AND dtype IS NOT NULL
GROUP BY 1, 2, 3, 4, 5`

const DeleteQuery = `
DELETE FROM real_time_report
WHERE time < '%s'
`
