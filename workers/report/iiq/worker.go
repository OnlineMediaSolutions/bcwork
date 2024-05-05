package iiq

import (
	"context"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdwh"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/quest"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"strings"
	"time"
)

type Worker struct {
	Sleep       time.Duration `json:"sleep"`
	Hours       int           `json:"hours"`
	Days        int           `json:"days"`
	Start       string        `json:"start"`
	DatabaseEnv string        `json:"dbenv"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	var err error
	w.Sleep, _ = conf.GetDurationValueWithDefault("sleep", time.Duration(5*time.Minute))
	w.Hours, err = conf.GetIntValueWithDefault("hours", 1)
	if err != nil {
		log.Warn().Err(err).Msg("failed to fetch hours config value (will user default)")
	}
	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local_prod")
	err = bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	err = bcdwh.InitDB(w.DatabaseEnv)
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DWH DB")
	}

	err = quest.InitDB("quest2")
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	w.Start, _ = conf.GetStringValue("start")
	if w.Start != "" {
		w.Sleep = 0
	}

	w.Days, err = conf.GetIntValueWithDefault("days", 1)
	if err != nil {
		log.Warn().Err(err).Msg("failed to fetch hours config value (will user default)")
	}
	//boil.DebugMode = true

	return nil

}

func (w *Worker) Do(ctx context.Context) error {

	var err error
	log.Info().Msg("New bidder supply report go")
	now := time.Now()
	start := now.UTC().Add(-1 * time.Duration(w.Hours) * time.Hour)
	stop := now.UTC().Add(time.Duration(1) * time.Hour)
	if w.Start != "" {
		start, err = time.Parse("2006010215", w.Start)
		if err != nil {
			return errors.Wrapf(err, "failed to parse start hour")
		}

		stop = start.Add(time.Duration(1) * time.Hour)
	}

	log.Info().Str("start", start.Format("2006-01-02T15")+":00:00Z").Str("stop", stop.Format("2006-01-02T15")+":00:00Z").Msg("pulling data from questdb")
	data, err := processIiqHourly(ctx, start.Format("2006-01-02T15")+":00:00Z", stop.Format("2006-01-02T15")+":00:00Z")
	if err != nil {
		return errors.Wrapf(err, "failed to query iiq counters")
	}

	log.Info().Int("len", len(data)).Msg("data ready for DB")

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to open transaction")
	}

	values := make([]string, 0)
	for _, rec := range data {
		values = append(values, fmt.Sprintf(`('%s', '%s',  %d, %d, %d, %f)`,
			rec.Time.Format("2006-01-02 15")+":00:00",
			rec.Dpid,
			rec.Request,
			rec.Response,
			rec.Impression,
			rec.Revenue))

	}

	q := fmt.Sprint(`INSERT INTO "iiq_hourly" ("time", "dpid","request", "response", "impression","revenue") VALUES `,
		strings.Join(values, ","),
		`ON CONFLICT ("time", "dpid") DO UPDATE SET "request" = EXCLUDED."request","response" = EXCLUDED."response","impression" = EXCLUDED."impression","revenue" = EXCLUDED."revenue"`)

	_, err = queries.Raw(q).Exec(tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "failed to update report iiq")
	}

	query := fmt.Sprintf(dailyUpdate, now.AddDate(0, 0, -1*w.Days).Format("2006-01-02"))
	_, err = queries.Raw(query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update daily(yesterday) revenue tables")
	}

	updatedAt := time.Now().UTC()
	updated := models.ReportUpdate{
		Report:   "iiq",
		UpdateAt: updatedAt,
	}
	err = updated.Upsert(ctx, tx, true, []string{models.ReportUpdateColumns.Report}, boil.Infer(), boil.Infer())
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "failed to update report update timestamp")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "failed to commit transaction")
	}
	log.Info().Msg("data saved to DB")

	// Update DWH
	log.Info().Msg("saving to DWH")

	_, err = queries.Raw(q).Exec(bcdwh.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update dwh iiq hourly table")
	}

	query = fmt.Sprintf(dailyUpdate, now.AddDate(0, 0, -1*w.Days).Format("2006-01-02"))
	_, err = queries.Raw(query).ExecContext(ctx, bcdwh.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update dwh iiq hourly table")
	}
	log.Info().Msg("data saved to DWH")

	return nil
}

func (w *Worker) GetSleep() int {
	return int(w.Sleep.Seconds())
}

type IiqHourly struct {
	Time       time.Time    `boil:"time" json:"time" toml:"time" yaml:"time"`
	Dpid       string       `boil:"dpid" json:"dpid" toml:"dpid" yaml:"dpid"`
	Request    null.Float64 `boil:"request" json:"request" toml:"request" yaml:"request"`
	Response   null.Float64 `boil:"response" json:"response" toml:"response" yaml:"response"`
	Impression null.Float64 `boil:"impression" json:"impression" toml:"impression" yaml:"impression"`
	Revenue    null.Float64 `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
}

func processIiqHourly(ctx context.Context, start string, stop string) (models.IiqHourlySlice, error) {

	log.Info().Msg("processIiqHourly")

	var records []*IiqHourly

	q := fmt.Sprintf(`SELECT date_trunc('hour',timestamp) as time,
                                    dpid,
                                    sum(request) request,
                                    sum(response) response,
                                    sum(impression) impression,
                                    sum(revenue) revenue
                              FROM iiq WHERE date_trunc('hour',timestamp)>='%s' AND date_trunc('hour',timestamp)<='%s'
                              GROUP BY date_trunc('hour',timestamp),dpid`, start, stop)
	err := queries.Raw(q).Bind(ctx, quest.DB(), &records)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query iiq from questdb")
	}

	var res models.IiqHourlySlice
	for _, r := range records {
		res = append(res, &models.IiqHourly{
			Time:       r.Time,
			Dpid:       r.Dpid,
			Request:    int64(r.Request.Float64),
			Response:   int64(r.Response.Float64),
			Impression: int64(r.Impression.Float64),
			Revenue:    r.Revenue.Float64,
		})
	}
	return res, nil
}

var dailyUpdate = `INSERT INTO iiq_daily (time, dpid, request, response,impression, revenue)
SELECT date_trunc('day',"time") "time",dpid,sum(request),sum(response),sum(impression),sum(revenue) FROM iiq_hourly WHERE date_trunc('day',"time") >='%s' GROUP BY date_trunc('day',"time"),dpid
ON CONFLICT (time,dpid)
DO UPDATE SET request=EXCLUDED.request,
             response=EXCLUDED.response,
             impression=EXCLUDED.impression,
              revenue=EXCLUDED.revenue`
