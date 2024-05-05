package revenue

import (
	"context"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdwh"
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
	DisableDWH  bool          `json:"dwhoff"`
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

	err = bcdwh.InitDB(w.DatabaseEnv)
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DWH")
	}
	if conf.GetBoolValueWithDefault("debug", false) {
		log.Info().Msg("debug mode: on")
		boil.DebugMode = true
	}

	w.DisableDWH = conf.GetBoolValueWithDefault("dwhoff", false)

	return nil
}

func (w *Worker) Do(ctx context.Context) error {

	log.Info().Int("days", w.Days).Int("hours", w.Hours).Msg("Revenue Report Do")
	now := time.Now()
	query := fmt.Sprintf(hourlyUpdate, now.Add(time.Duration(w.Hours)*-1*time.Hour).Format("2006-01-02T15")+":00:00")
	_, err := queries.Raw(query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update hourly revenue tables")
	}

	query = fmt.Sprintf(hourlyUpdateIiq, now.Add(time.Duration(w.Hours)*-1*time.Hour).Format("2006-01-02T15")+":00:00")
	_, err = queries.Raw(query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update hourly revenue iiq tables")
	}

	query = fmt.Sprintf(hourlyUpdateDpBidRequest, now.Add(time.Duration(w.Hours)*-1*time.Hour).Format("2006-01-02T15")+":00:00")
	_, err = queries.Raw(query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update hourly revenue iiq tables")
	}

	query = fmt.Sprintf(dailyUpdate, now.AddDate(0, 0, -1*w.Days).Format("2006-01-02"))
	_, err = queries.Raw(query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update daily(yesterday) revenue tables")
	}

	// Update new DWH
	log.Info().Msg("updating new DWH")
	if !w.DisableDWH {
		query = fmt.Sprintf(dwhHourlyUpdate, now.Add(time.Duration(w.Hours)*-1*time.Hour).Format("2006-01-02T15")+":00:00")
		_, err = queries.Raw(query).ExecContext(ctx, bcdwh.DB())
		if err != nil {
			return errors.Wrapf(err, "failed to update dwh hourly revenue tables")
		}

		query = fmt.Sprintf(hourlyUpdateIiq, now.Add(time.Duration(w.Hours)*-1*time.Hour).Format("2006-01-02T15")+":00:00")
		_, err = queries.Raw(query).ExecContext(ctx, bcdwh.DB())
		if err != nil {
			return errors.Wrapf(err, "failed to update dwh hourly revenue iiq tables")
		}

		query = fmt.Sprintf(dwhHourlyUpdateDpBidRequest, now.Add(time.Duration(w.Hours)*-1*time.Hour).Format("2006-01-02T15")+":00:00")
		_, err = queries.Raw(query).ExecContext(ctx, bcdwh.DB())
		if err != nil {
			return errors.Wrapf(err, "failed to update dwh hourly revenue from bid requests tables")
		}

		query = fmt.Sprintf(dailyUpdate, now.AddDate(0, 0, -1*w.Days).Format("2006-01-02"))
		_, err = queries.Raw(query).ExecContext(ctx, bcdwh.DB())
		if err != nil {
			return errors.Wrapf(err, "failed to update dwh daily(yesterday) revenue tables")
		}
	}

	log.Info().Msg("Done")
	return nil
}

var hourlyUpdate = `INSERT INTO revenue_hourly (time, publisher_impressions, sold_impressions, cost, revenue, demand_partner_fees, missed_opportunities)
SELECT time,sum(publisher_impressions),sum(sold_impressions),sum(cost),sum(revenue),sum(demand_partner_fee),sum(missed_opportunities) 
FROM nb_supply_hourly WHERE time >='%s' GROUP BY time
ON CONFLICT (time)
DO UPDATE SET publisher_impressions=EXCLUDED.publisher_impressions,
              sold_impressions=EXCLUDED.sold_impressions,
              cost=EXCLUDED.cost,
              revenue=EXCLUDED.revenue,
              demand_partner_fees=EXCLUDED.demand_partner_fees,
              missed_opportunities=EXCLUDED.missed_opportunities`

var hourlyUpdateIiq = `INSERT INTO revenue_hourly (time, data_fee)
SELECT time,sum(revenue)*0.045 data_fee FROM iiq_hourly WHERE time >='%s' GROUP BY time
ON CONFLICT (time)
DO UPDATE SET data_fee=EXCLUDED.data_fee`

var hourlyUpdateDpBidRequest = `INSERT INTO revenue_hourly (time, dp_bid_requests)
SELECT time,sum(bid_requests) dp_bid_requests FROM nb_demand_hourly WHERE time >='%s' GROUP BY time
ON CONFLICT (time)
DO UPDATE SET dp_bid_requests=EXCLUDED.dp_bid_requests`

var dailyUpdate = `INSERT INTO revenue_daily (time, publisher_impressions, sold_impressions, dp_bid_requests,cost, revenue, demand_partner_fees,data_fee, missed_opportunities)
SELECT date_trunc('day',"time") "time",sum(publisher_impressions),sum(sold_impressions),sum(dp_bid_requests),sum(cost),sum(revenue),sum(demand_partner_fees),sum(data_fee),sum(missed_opportunities) FROM revenue_hourly WHERE date_trunc('day',"time") >='%s' GROUP BY date_trunc('day',"time")
ON CONFLICT (time)
DO UPDATE SET publisher_impressions=EXCLUDED.publisher_impressions,
              sold_impressions=EXCLUDED.sold_impressions,
              dp_bid_requests=EXCLUDED.dp_bid_requests,
              cost=EXCLUDED.cost,
              revenue=EXCLUDED.revenue,
              demand_partner_fees=EXCLUDED.demand_partner_fees,
              data_fee=EXCLUDED.data_fee,
              missed_opportunities=EXCLUDED.missed_opportunities`

func (w *Worker) GetSleep() int {
	return int(w.Sleep.Seconds())
}
