package logs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"os"
	"path/filepath"
	"time"
)

type Worker struct {
	Sleep       time.Duration `json:"sleep"`
	Hours       int           `json:"hours"`
	Days        int           `json:"days"`
	DatabaseEnv string        `json:"dbenv"`
	Start       time.Time     `json:"start"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	var err error
	w.Sleep, _ = conf.GetDurationValueWithDefault("sleep", 0)
	w.Hours, err = conf.GetIntValueWithDefault("hours", 2)
	if err != nil {
		log.Warn().Err(err).Msg("failed to fetch hours config value (will user default)")
	}
	w.Start, err = conf.GetDateHourValueWithDefault("start", time.Now().Add(-1*time.Hour).UTC())
	if err != nil {
		return errors.Wrapf(err, "failed to read start value")
	}

	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "prod")
	err = bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	w.Days, err = conf.GetIntValueWithDefault("days", 2)
	QueryUpdateDaily = fmt.Sprintf(QueryUpdateDaily, w.Days)

	return nil
}

func (w *Worker) Do(ctx context.Context) error {

	fmt.Println("Impression Log Report Do")

	files, err := w.listFiles()
	if err != nil {
		return errors.WithStack(err)
	}
	log.Info().Interface("files", files).Msg("processing files")

	res := make(map[string]*models.ImpressionLogHourly)
	for _, f := range files {
		err := ProcessFile(f, res)
		if err != nil {
			return errors.Wrapf(err, "failed to process file (filename:%s)", f)
		}
		log.Info().Str("filename", f).Int("records", len(res)).Msg("file processed")
	}

	log.Info().Msg("data ready for DB")
	values := funk.Values(res).([]*models.ImpressionLogHourly)
	for i, v := range values {
		err := v.Upsert(ctx, bcdb.DB(), true, []string{
			models.ImpressionLogHourlyColumns.Time,
			models.ImpressionLogHourlyColumns.PublisherID,
			models.ImpressionLogHourlyColumns.DemandPartnerID,
			models.ImpressionLogHourlyColumns.Domain,
			models.ImpressionLogHourlyColumns.DeviceType,
			models.ImpressionLogHourlyColumns.Os,
			models.ImpressionLogHourlyColumns.Country,
			models.ImpressionLogHourlyColumns.Size,
			models.ImpressionLogHourlyColumns.HadFollowup,
			models.ImpressionLogHourlyColumns.IsFirst,
		}, boil.Infer(), boil.Infer())
		if err != nil {
			return errors.Wrapf(err, "failed to upsert log hourly record")
		}
		if i%1000 == 0 {
			log.Info().Msgf("DB updated (%d of %d)", i, len(values))
		}
	}

	_, err = queries.Raw(QueryUpdateDaily).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to update daily impressions log tables")
	}

	log.Info().Msg("data saved to DB")
	return nil
}

func (w *Worker) GetSleep() int {
	return int(w.Sleep.Seconds())
}

func ProcessFile(filename string, res map[string]*models.ImpressionLogHourly) error {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "error opening file(filename:%s)", filename)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	l := 0
	// Read each line and print it to the console
	for scanner.Scan() {
		l++
		rec := Record{}
		err = json.Unmarshal(scanner.Bytes(), &rec)
		if err != nil {
			return errors.Wrapf(err, "failed to read line and convert it to record (filename:%s,line:%d)", filename, l)
		}
		key := rec.getKey()
		mod := res[key]
		if mod == nil {
			mod = &models.ImpressionLogHourly{
				Time:            time.Date(rec.Time.Year(), rec.Time.Month(), rec.Time.Day(), rec.Time.Hour(), 0, 0, 0, time.UTC),
				DemandPartnerID: rec.Get("dpid"),
				PublisherID:     rec.Get("pubid"),
				Domain:          rec.Get("domain"),
				Os:              rec.Get("os"),
				DeviceType:      rec.Get("dtype"),
				Country:         rec.Get("country"),
				Size:            rec.Get("size"),
				HadFollowup:     rec.Get("hadfup") == "1",
				IsFirst:         rec.Get("loop") == "0",
			}
			res[key] = mod
		}

		if rec.Get("loop") == "0" {
			mod.PubImpressions++
			cpm, err := rec.GetFloat64("sbpr")
			if err != nil {
				return errors.Wrapf(err, "failed to convert sbpr(filename:%s,line:%d)", filename, l)
			}
			mod.Cost += (cpm / 1000)
		}
		mod.SoldImpressions++
		drpm, err := rec.GetFloat64("dbpr")
		if err != nil {
			return errors.Wrapf(err, "failed to convert dbpr(filename:%s,line:%d)", filename, l)
		}
		mod.Revenue += (drpm / 1000)

		dpfee, err := rec.GetFloat64("dpfee")
		if err != nil {
			dpfee = 0
		}
		mod.DemandPartnerFees += (dpfee / 1000)
	}

	// Check for any errors that may have occurred during scanning
	if err := scanner.Err(); err != nil {
		return errors.Wrapf(err, "error reading file(filename:%s)", filename)
	}

	return nil
}

func (w *Worker) listFiles() ([]string, error) {

	var list []string
	for h := 0; h < w.Hours; h++ {
		format := w.Start.Add(time.Duration(h) * time.Hour).Format("2006010215")
		match := fmt.Sprintf("/var/log/bcrt/*%s*", format)
		files, err := filepath.Glob(match)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list files")
		}
		list = append(list, files...)
	}

	return list, nil
}

var QueryUpdateDaily = `INSERT INTO impression_log_daily
SELECT date_trunc('day',"time") "time",
       publisher_id,
       demand_partner_id,
       domain,
       os,
       country,
       device_type,
       size,
       is_first,
       had_followup,
       SUM(sold_impressions),
       SUM(pub_impressions),
       SUM(cost),
       SUM(revenue),
       SUM(demand_partner_fees) 
FROM impression_log_hourly
WHERE date_trunc('day',"time") > (now())- interval '%d day'
GROUP BY (date_trunc('day',"time"),publisher_id,demand_partner_id,domain,os,country,device_type,size,is_first,had_followup)
ON CONFLICT ("time",publisher_id,demand_partner_id,domain,os,country,device_type,size,is_first,had_followup)
DO UPDATE SET sold_impressions=EXCLUDED.sold_impressions,
              pub_impressions=EXCLUDED.pub_impressions,
              cost=EXCLUDED.cost, 
              revenue=EXCLUDED.revenue,
              demand_partner_fees=EXCLUDED.demand_partner_fees`
