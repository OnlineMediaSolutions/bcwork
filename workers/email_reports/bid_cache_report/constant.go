package bid_cache_report

const MINIMUM_IMPRESSIONS = 5000
const THRESHOLD = 1.1

const questQuery = `SELECT DATE_TRUNC('day', to_timezone(timestamp, 'America/New_York')) time,
       publisher publisher_id,
       domain,
       tgrp,
       sum(dbpr)/1000 revenue,
       sum(sbpr)/1000 cost,
       sum(dpfee)/1000 demand_partner_fee,
       count(1) sold_impressions,
       sum(CASE WHEN loop=false THEN 1 ELSE 0 END) publisher_impressions,
       sum(CASE WHEN uidsrc='iiq' THEN dbpr/1000 ELSE 0 END) * 0.045 data_fee
FROM impression
WHERE publisher IS NOT NULL
AND DATE_TRUNC('day', to_timezone(timestamp, 'America/New_York')) = '%s'
  AND domain IS NOT NULL
  AND country IS NOT NULL
  AND dtype IS NOT NULL
  AND %s

GROUP BY 1, 2, 3, 4`
