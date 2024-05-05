package demand

import (
	"context"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"time"
)

type Worker struct {
	Sleep       time.Duration `json:"sleep"`
	Hours       int           `json:"hours"`
	Days        int           `json:"days"`
	DatabaseEnv string        `json:"dbenv"`
	Start       *time.Time    `json:"start"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	var err error
	w.Sleep, _ = conf.GetDurationValueWithDefault("sleep", 0)
	w.Hours, err = conf.GetIntValueWithDefault("hours", 2)
	if err != nil {
		log.Warn().Err(err).Msg("failed to fetch hours config value (will user default)")
	}

	w.Days, err = conf.GetIntValueWithDefault("days", 1)
	if err != nil {
		log.Warn().Err(err).Msg("failed to fetch hours config value (will user default)")
	}
	//w.Start, err = conf.GetDateHourValueWithDefault("start", time.Now().Add(-1*time.Hour).UTC())
	//if err != nil {
	//	return errors.Wrapf(err, "failed to read start value")
	//}

	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local_prod")
	err = bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	if conf.GetBoolValueWithDefault("debug", false) {
		log.Info().Msg("debug mode: on")
		boil.DebugMode = true
	}
	return nil
}

func (w *Worker) Do(ctx context.Context) error {

	log.Info().Int("days", w.Days).Int("hours", w.Hours).Msg("Demand Report Do")
	now := time.Now()
	query := fmt.Sprintf(hourlyUpdate, now.Add(time.Duration(w.Hours)*-1*time.Hour).Format("2006-01-02T15")+":00:00")
	_, err := queries.Raw(query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update hourly demand tables")
	}

	query = fmt.Sprintf(dailyUpdate, now.AddDate(0, 0, -1*w.Days).Format("2006-01-02"))
	_, err = queries.Raw(query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update daily(yesterday) demand tables")
	}

	log.Info().Msg("Done")
	return nil
}

var hourlyUpdate = `INSERT INTO demand_report_hourly (time, demand_partner_id,publisher_id,domain,bid_requests,bid_responses,avg_bid_price, sold_impressions,publisher_impressions, cost, revenue, demand_partner_fee)
SELECT time, demand_partner_id,publisher_id,domain,sum(bid_requests),sum(bid_responses),avg(avg_bid_price),sum(sold_impressions),sum(publisher_impressions),sum(cost),sum(revenue),sum(demand_partner_fee) FROM nb_demand_hourly WHERE time >='%s' GROUP BY time
ON CONFLICT (time,demand_partner_id,publisher_id,domain)
DO UPDATE SET bid_requests=EXCLUDED.bid_requests,
              bid_responses=EXCLUDED.bid_responses,
              avg_bid_price=EXCLUDED.avg_bid_price,
              publisher_impressions=EXCLUDED.publisher_impressions,
              sold_impressions=EXCLUDED.sold_impressions,
              cost=EXCLUDED.cost,
              revenue=EXCLUDED.revenue,
              demand_partner_fee=EXCLUDED.demand_partner_fee`

var dailyUpdate = `INSERT INTO demand_report_daiy (time, demand_partner_id,publisher_id,domain,bid_requests,bid_responses,avg_bid_price, sold_impressions,publisher_impressions, cost, revenue, demand_partner_fee)
SELECT date_trunc('day',"time") "time", demand_partner_id,publisher_id,domain,sum(bid_requests),sum(bid_responses),avg(avg_bid_price),sum(sold_impressions),sum(publisher_impressions),sum(cost),sum(revenue),sum(demand_partner_fee) FROM nb_demand_hourly WHERE date_trunc('day',"time") >='%s' GROUP BY date_trunc('day',"time")
ON CONFLICT (time,demand_partner_id,publisher_id,domain)
DO UPDATE SET bid_requests=EXCLUDED.bid_requests,
              bid_responses=EXCLUDED.bid_responses,
              avg_bid_price=EXCLUDED.avg_bid_price,
              publisher_impressions=EXCLUDED.publisher_impressions,
              sold_impressions=EXCLUDED.sold_impressions,
              cost=EXCLUDED.cost,
              revenue=EXCLUDED.revenue,
              demand_partner_fee=EXCLUDED.demand_partner_fee`

func (w *Worker) GetSleep() int {
	return int(w.Sleep.Seconds())
}
