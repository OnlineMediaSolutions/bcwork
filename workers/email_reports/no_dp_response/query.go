package nodpresponse

var (
	questQuery = `
		SELECT
			DATE_TRUNC('day', to_timezone(timestamp, 'America/New_York')) as time,
			pubid,
			domain,
			dpid,
			sum(count) as %s
		FROM %s
		WHERE
			to_timezone(timestamp, 'America/New_York') >= '%s'
			AND to_timezone(timestamp, 'America/New_York') < '%s'
		GROUP BY 1,2,3,4;
	`

	postgresQuery = `
		WITH res AS (
			SELECT
				demand_partner_id AS dpid,
				publisher_id AS pubid,
				"domain",
				sum(bid_requests) AS bid_requests
			FROM no_dp_response_report AS x 
			WHERE
				"time" >= $1 AND "time" < $2
			GROUP BY demand_partner_id, publisher_id, "domain"
			HAVING count(demand_partner_id) = $3
		)
		SELECT
			r.dpid,
			p."name" AS publisher_name,
			r.pubid,
			r."domain",
			r.bid_requests
		FROM res AS r
		JOIN publisher AS p ON p.publisher_id = r.pubid
		ORDER BY r.bid_requests DESC;
	`
)
