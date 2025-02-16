package nbsupply

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/quest"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

type Worker struct {
	Sleep       time.Duration `json:"sleep"`
	Hours       int           `json:"hours"`
	Start       string        `json:"start"`
	DatabaseEnv string        `json:"dbenv"`
	FromDB      bool          `json:"fromdb"`
	Debug       bool          `json:"debug"`
	DisableDWH  bool          `json:"dwhoff"`
	QuestNYC    *sqlx.DB
	QuestAMS    *sqlx.DB
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

	//err = bcdwh.InitDB(w.DatabaseEnv)
	//if err != nil {
	//	return errors.Wrapf(err, "failed to initalize DWH")
	//}

	w.FromDB = conf.GetBoolValueWithDefault("fromdb", false)
	if !w.FromDB {
		w.QuestNYC, err = quest.Connect("nycquest" + conf.GetStringValueWithDefault("quest", "2"))
		if err != nil {
			return errors.Wrapf(err, "failed to initalize DB")
		}

		w.QuestAMS, err = quest.Connect("amsquest" + conf.GetStringValueWithDefault("quest", "2"))
		if err != nil {
			return errors.Wrapf(err, "failed to initalize DB")
		}
	} else {
		w.Sleep = 0
	}

	w.Start, _ = conf.GetStringValue("start")
	if w.Start != "" {
		w.Sleep = 0
	}

	w.Debug = conf.GetBoolValueWithDefault("debug", false)
	if w.Debug {
		boil.DebugMode = true
	}

	w.DisableDWH = conf.GetBoolValueWithDefault("dwhoff", false)

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	var records []*models.NBSupplyHourly
	var err error
	if w.FromDB {
		log.Info().Msg("fetch supply records from BCDB")
		records, err = w.FetchFromBCDB(ctx)
		if err != nil {
			return err
		}
	} else {
		log.Info().Msg("fetch supply records from QuestDB")
		records, err = w.FetchFromQuest(ctx)
		if err != nil {
			return err
		}
	}
	log.Info().Int("len", len(records)).Msg("data ready for DB")

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to open transaction")
	}
	pubImps := int64(0)
	soldImps := int64(0)

	values := make([]string, 0)
	for _, rec := range records {
		values = append(values, fmt.Sprintf(`('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s','%s','%s', '%s', %d, %d, %d, %d, %f, %f, %f, %d, %f,%d, %f)`,
			rec.Time.Format("2006-01-02 15")+":00:00",
			rec.PublisherID,
			rec.Domain,
			rec.Os,
			rec.Country,
			rec.DeviceType,
			rec.PlacementType,
			rec.RequestType,
			rec.Size,
			rec.PaymentType,
			rec.Datacenter,
			rec.BidRequests,
			rec.BidResponses,
			rec.SoldImpressions,
			rec.PublisherImpressions,
			rec.Cost,
			rec.Revenue,
			rec.AvgBidPrice,
			rec.MissedOpportunities,
			rec.DemandPartnerFee,
			rec.PublisherImpressions,
			rec.DataFee))
		soldImps += rec.SoldImpressions
		pubImps += rec.PublisherImpressions
	}

	q := fmt.Sprint(`INSERT INTO "nb_supply_hourly" ("time", "publisher_id", "domain", "os", "country",  "device_type", "placement_type","request_type","size","payment_type","datacenter","bid_requests", "bid_responses", "sold_impressions","publisher_impressions", "cost", "revenue", "avg_bid_price", "missed_opportunities", "demand_partner_fee","data_impressions","data_fee") VALUES `,
		strings.Join(values, ","),
		`ON CONFLICT ("time", "publisher_id", "domain", "os", "country",  "device_type", "placement_type","request_type","size","payment_type","datacenter") DO UPDATE SET "bid_requests" = EXCLUDED."bid_requests","bid_responses" = EXCLUDED."bid_responses","sold_impressions" = EXCLUDED."sold_impressions","publisher_impressions" = EXCLUDED."publisher_impressions","cost" = EXCLUDED."cost","revenue" = EXCLUDED."revenue","avg_bid_price" = EXCLUDED."avg_bid_price","missed_opportunities" = EXCLUDED."missed_opportunities","demand_partner_fee" = EXCLUDED."demand_partner_fee","data_impressions" = EXCLUDED."data_impressions","data_fee" = EXCLUDED."data_fee"`)

	_, err = queries.Raw(q).Exec(tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "failed to update report nbdemand")
	}

	updatedAt := time.Now().UTC()
	updated := models.ReportUpdate{
		Report:   "new_bidder_supply_hourly",
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

	log.Info().Int64("pubImps", pubImps).Int64("soldImps", soldImps).Msg("data saved to DB")

	compassRecords := ConvertToCompass(ctx, records)
	err = Send(ctx, compassRecords)
	if err != nil {
		return errors.Wrapf(err, "failed to send data to compass")
	}

	log.Info().Msg("Json update successfully posted to compass")

	if !w.DisableDWH {
		// New DWH loader
		log.Info().Msg("loading dwh")

		err = createCSV(records)
		if err != nil {
			return errors.Wrapf(err, "failed to create csv file for dwh")
		}
		log.Info().Msg("csv created")

		err = scpCSV()
		if err != nil {
			return errors.Wrapf(err, "failed to send csv file to dwh")
		}
		log.Info().Msg("csv copied to dwh")
		err = createTempTable()
		if err != nil {
			return errors.Wrapf(err, "failed to create temp table")
		}
		log.Info().Msg("temp table create on dwh")
		err = createTempTable()
		if err != nil {
			return errors.Wrapf(err, "failed to create temp table")
		}

		err = loadToRawTable(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to load temp table")
		}
		log.Info().Msg("csv loaded to temp table  on dwh")

		err = updateFromRawToHourly()
		if err != nil {
			return errors.Wrapf(err, "failed to update from temp table to hourly")
		}
		log.Info().Msg("hourly table updated")

		err = dropTempTable()
		if err != nil {
			return errors.Wrapf(err, "failed to drop  temp table")
		}
		log.Info().Msg("temp table dropped")
	}

	log.Info().Msg("DONE")

	return nil
}

func (w *Worker) GetSleep() int {
	return int(w.Sleep.Seconds())
}

type NBSupplyHourly struct {
	Time                 time.Time `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string    `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain               string    `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Os                   string    `boil:"os" json:"os" toml:"os" yaml:"os"`
	Country              string    `boil:"country" json:"country" toml:"country" yaml:"country"`
	DeviceType           string    `boil:"device_type" json:"device_type" toml:"device_type" yaml:"device_type"`
	PlacementType        string    `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
	RequestType          string    `boil:"request_type" json:"request_type" toml:"request_type" yaml:"request_type"`
	Size                 string    `boil:"size" json:"size" toml:"size" yaml:"size"`
	PaymentType          string    `boil:"payment_type" json:"payment_type" toml:"payment_type" yaml:"payment_type"`
	Loop                 bool      `boil:"loop" json:"loop" toml:"loop" yaml:"loop"`
	Datacenter           string    `boil:"datacenter" json:"datacenter" toml:"datacenter" yaml:"datacenter"`
	BidRequests          float64   `boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	BidResponses         float64   `boil:"bid_responses" json:"bid_repsponses" toml:"bid_repsponses" yaml:"bid_repsponses"`
	SoldImpressions      float64   `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	PublisherImpressions float64   `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
	Cost                 float64   `boil:"cost" json:"cost" toml:"cost" yaml:"cost"`
	Revenue              float64   `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	AvgBidPrice          float64   `boil:"avg_bid_price" json:"avg_bid_price" toml:"avg_bid_price" yaml:"avg_bid_price"`
	MissedOpportunities  float64   `boil:"missed_opportunities" json:"missed_opportunities" toml:"missed_opportunities" yaml:"missed_opportunities"`
	DemandPartnerFee     float64   `boil:"demand_partner_fee" json:"demand_partner_fee" toml:"demand_partner_fee" yaml:"demand_partner_fee"`
	DataImpressions      float64   `boil:"data_impressions" json:"data_impressions" toml:"data_impressions" yaml:"data_impressions"`
	DataFee              float64   `boil:"data_fee" json:"data_fee" toml:"data_fee" yaml:"data_fee"`
}

func (w *Worker) FetchFromBCDB(ctx context.Context) ([]*models.NBSupplyHourly, error) {
	var err error
	log.Info().Msg("New bidder supply report go")
	start := time.Now().UTC().Add(-1 * time.Duration(w.Hours) * time.Hour)
	stop := time.Now().UTC().Add(time.Duration(1) * time.Hour)
	if w.Start != "" {
		start, err = time.Parse("2006010215", w.Start)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse start hour")
		}

		stop = start.Add(time.Duration(1) * time.Hour)
	}

	data, err := models.NBSupplyHourlies(
		models.NBSupplyHourlyWhere.Time.GTE(start),
		models.NBSupplyHourlyWhere.Time.LT(stop)).All(ctx, bcdb.DB())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch from bcdb")
	}

	var records []*models.NBSupplyHourly
	for _, v := range data {
		if v.PublisherID != "" && (v.BidResponses > 0 || v.PublisherImpressions > 0 || v.SoldImpressions > 0) {
			records = append(records, v)
		}
	}

	return records, nil
}

func (w *Worker) FetchFromQuest(ctx context.Context) ([]*models.NBSupplyHourly, error) {
	var err error
	log.Info().Msg("New bidder supply report go")
	start := time.Now().UTC().Add(-1 * time.Duration(w.Hours) * time.Hour)
	stop := time.Now().UTC().Add(time.Duration(1) * time.Hour)
	if w.Start != "" {
		start, err = time.Parse("2006010215", w.Start)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse start hour")
		}

		stop = start.Add(time.Duration(1) * time.Hour)
	}

	data := make(map[string]*models.NBSupplyHourly)
	log.Info().Str("start", start.Format("2006-01-02T15")+":00:00Z").Str("stop", stop.Format("2006-01-02T15")+":00:00Z").Msg("pulling data from questdb")
	err = w.processBidRequestCounters(ctx, start.Format("2006-01-02T15")+":00:00Z", stop.Format("2006-01-02T15")+":00:00Z", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query request counters")
	}

	err = w.processBidResponseCounters(ctx, start.Format("2006-01-02T15")+":00:00Z", stop.Format("2006-01-02T15")+":00:00Z", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query response counters")
	}

	err = w.processMopsCounters(ctx, start.Format("2006-01-02T15")+":00:00Z", stop.Format("2006-01-02T15")+":00:00Z", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query mops counters")
	}

	err = w.processImpressionsCounters(ctx, start.Format("2006-01-02T15")+":00:00Z", stop.Format("2006-01-02T15")+":00:00Z", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query impressions counters")
	}

	var records []*models.NBSupplyHourly

	for _, v := range data {
		if v.PublisherID != "" && (v.BidResponses > 0 || v.PublisherImpressions > 0 || v.SoldImpressions > 0) {
			if v.RequestType == "js" {
				v.PaymentType = "cpm"
			} else {
				v.PaymentType = "bid"
			}
			records = append(records, v)
		}
	}

	return records, nil
}
func (rec *NBSupplyHourly) Key() string {
	return fmt.Sprint(rec.Time.Format("2006-01-02-15"), rec.PublisherID, rec.Domain, rec.Os, rec.Country, rec.DeviceType, rec.PlacementType, rec.RequestType, rec.Size, rec.PaymentType, rec.Datacenter)
}

func (r *NBSupplyHourly) ToModel() *models.NBSupplyHourly {
	return &models.NBSupplyHourly{
		Time:          r.Time,
		PublisherID:   r.PublisherID,
		Domain:        r.Domain,
		Os:            r.Os,
		Country:       r.Country,
		DeviceType:    r.DeviceType,
		PlacementType: r.PlacementType,
		RequestType:   r.RequestType,
		Size:          r.Size,
		PaymentType:   r.PaymentType,
		Datacenter:    r.Datacenter,
	}
}

func (w *Worker) processBidRequestCounters(ctx context.Context, start string, stop string, data map[string]*models.NBSupplyHourly) error {
	log.Info().Msg("processBidRequestCounters")

	var recordsNYC []*NBSupplyHourly
	var recordsAMS []*NBSupplyHourly

	q := `SELECT date_trunc('hour',timestamp) as time,
                                        dtype device_type,
os,country,ptype placement_type,pubid publisher_id,domain,size,reqtyp request_type,'%s' as datacenter,
                              sum(count) bid_requests
                              FROM request_placement WHERE date_trunc('hour',timestamp)>='%s' AND date_trunc('hour',timestamp)<='%s' AND dtype is not null
                              GROUP BY date_trunc('hour',timestamp),dtype,os,country,ptype,pubid,domain,size,request_type`
	//log.Info().Str("q", q).Msg("processBidRequestCounters")
	err := queries.Raw(fmt.Sprintf(q, "nyc", start, stop)).Bind(ctx, w.QuestNYC, &recordsNYC)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb nyc")
	}
	err = queries.Raw(fmt.Sprintf(q, "ams", start, stop)).Bind(ctx, w.QuestAMS, &recordsAMS)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb ams")
	}
	records := append(recordsNYC, recordsAMS...)

	for _, r := range records {
		key := r.Key()
		mod, ok := data[key]
		if !ok {
			mod = r.ToModel()
			data[key] = mod
		}

		mod.BidRequests = int64(r.BidRequests)
	}

	return nil
}

func (w *Worker) processBidResponseCounters(ctx context.Context, start string, stop string, data map[string]*models.NBSupplyHourly) error {
	log.Info().Msg("processBidResponseCounters")
	var recordsNYC []*NBSupplyHourly
	var recordsAMS []*NBSupplyHourly

	q := `SELECT date_trunc('hour',timestamp) as time,
                              dtype device_type,
                     os,country,ptype placement_type,pubid publisher_id,domain,size,reqtyp request_type,'%s' as datacenter,
                              sum(count) bid_responses,sum("sum") avg_bid_price
                              FROM bid_price WHERE date_trunc('hour',timestamp)>='%s' AND date_trunc('hour',timestamp)<='%s' AND dtype is not null
                              GROUP BY date_trunc('hour',timestamp),dtype,os,country,ptype,pubid,domain,size,request_type`

	//log.Info().Str("q", q).Msg("processBidResponseCounters")
	err := queries.Raw(fmt.Sprintf(q, "nyc", start, stop)).Bind(ctx, w.QuestNYC, &recordsNYC)
	if err != nil {
		return errors.Wrapf(err, "failed to query processBidResponseCounters from questdb nyc")
	}
	err = queries.Raw(fmt.Sprintf(q, "ams", start, stop)).Bind(ctx, w.QuestAMS, &recordsAMS)
	if err != nil {
		return errors.Wrapf(err, "failed to query processBidResponseCounters from questdb ams")
	}
	records := append(recordsNYC, recordsAMS...)

	for _, r := range records {
		key := r.Key()
		mod, ok := data[key]
		if !ok {
			mod = r.ToModel()
			data[key] = mod
		}
		mod.BidResponses = int64(r.BidResponses)
		mod.AvgBidPrice += r.AvgBidPrice
	}

	return nil
}

func (w *Worker) processMopsCounters(ctx context.Context, start string, stop string, data map[string]*models.NBSupplyHourly) error {
	log.Info().Msg("processMopsCounters")

	var recordsNYC []*NBSupplyHourly
	var recordsAMS []*NBSupplyHourly

	q := `SELECT date_trunc('hour',timestamp) as time,                              
          dtype device_type,
         os,country,ptype placement_type,pubid publisher_id,domain,size,reqtyp request_type,'%s' as datacenter,
                              sum(count) missed_opportunities
                              FROM mop WHERE date_trunc('hour',timestamp)>='%s' AND date_trunc('hour',timestamp)<='%s' AND dtype is not null
                              GROUP BY date_trunc('hour',timestamp),dtype,os,country,ptype,pubid,domain,size,request_type`
	//log.Info().Str("q", q).Msg("processMopsCounters")
	err := queries.Raw(fmt.Sprintf(q, "nyc", start, stop)).Bind(ctx, w.QuestNYC, &recordsNYC)
	if err != nil {
		return errors.Wrapf(err, "failed to query processMopsCounters from questdb nyc")
	}
	err = queries.Raw(fmt.Sprintf(q, "ams", start, stop)).Bind(ctx, w.QuestAMS, &recordsAMS)
	if err != nil {
		return errors.Wrapf(err, "failed to query processMopsCounters from questdb ams")
	}
	records := append(recordsNYC, recordsAMS...)

	for _, r := range records {
		key := r.Key()
		mod, ok := data[key]
		if !ok {
			mod = r.ToModel()
			data[key] = mod
		}

		mod.MissedOpportunities += int64(r.MissedOpportunities)
	}

	return nil
}

func (w *Worker) processImpressionsCounters(ctx context.Context, start string, stop string, data map[string]*models.NBSupplyHourly) error {
	log.Info().Msg("processImpressionsCounters")

	var recordsNYC []*NBSupplyHourly
	var recordsAMS []*NBSupplyHourly

	q := `  SELECT date_trunc('hour',timestamp) as time,
                                            dtype device_type,
os,country,ptype placement_type,publisher publisher_id,domain,size,reqtyp request_type,'%s' as datacenter,
              sum(dbpr)/1000 revenue,sum(sbpr)/1000 cost  ,sum(dpfee)/1000 demand_partner_fee ,count(1) sold_impressions,sum(case when loop=false then 1 else 0 end) publisher_impressions,sum(case when uidsrc='iiq' then 1 else 0 end) data_impressions,sum(case when uidsrc='iiq' then dbpr/1000 else 0 end) data_fee
                              FROM impression WHERE date_trunc('hour',timestamp)>='%s' AND date_trunc('hour',timestamp)<='%s' AND dtype is not null
                              GROUP BY date_trunc('hour',timestamp),dtype,os,country,ptype,publisher,domain,size,request_type`
	//log.Info().Str("q", q).Msg("processImpressionsCounters")
	err := queries.Raw(fmt.Sprintf(q, "nyc", start, stop)).Bind(ctx, w.QuestNYC, &recordsNYC)
	if err != nil {
		return errors.Wrapf(err, "failed to query processImpressionsCounters from questdb nyc")
	}
	err = queries.Raw(fmt.Sprintf(q, "ams", start, stop)).Bind(ctx, w.QuestAMS, &recordsAMS)
	if err != nil {
		return errors.Wrapf(err, "failed to query processImpressionsCounters from questdb ams")
	}
	records := append(recordsNYC, recordsAMS...)

	for _, r := range records {
		key := r.Key()
		mod, ok := data[key]
		if !ok {
			mod = r.ToModel()
			data[key] = mod
		}

		mod.Cost += r.Cost
		mod.Revenue += r.Revenue
		mod.DemandPartnerFee += r.DemandPartnerFee
		mod.PublisherImpressions += int64(r.PublisherImpressions)
		mod.SoldImpressions += int64(r.SoldImpressions)
		mod.DataImpressions += int64(r.DataImpressions)
		mod.DataFee += r.DataFee

		//log.Info().Interface("rec", mod).Msgf("rec")
	}

	return nil
}
