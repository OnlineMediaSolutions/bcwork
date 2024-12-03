package nbdemand

import (
	"context"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/quest"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"strings"
	"time"
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
	QuestSFO    *sqlx.DB
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

		w.QuestSFO, err = quest.Connect("sfoquest" + conf.GetStringValueWithDefault("quest", "2"))
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

	var records []*models.NBDemandHourly
	var err error
	if w.FromDB {
		log.Info().Msg("fetch demand records from BCDB")
		records, err = w.FetchFromBCDB(ctx)
		if err != nil {
			return eris.Wrapf(err, "failed to fetch counters data from postgres")
		}

	} else {
		log.Info().Msg("fetch demand records from QuestDB")
		records, err = w.FetchFromQuest(ctx)
		if err != nil {
			return eris.Wrapf(err, "failed to fetch counters data from quest")
		}
	}

	log.Info().Int("len", len(records)).Msg("data ready for DB")

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to open transaction")
	}
	soldImps := int64(0)

	values := make([]string, 0)
	//boil.DebugMode = true
	for _, rec := range records {
		values = append(values, fmt.Sprintf(`('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s','%s', '%s','%s', '%s', %d, %d, %f, %d, %f,%d,%f,%f,%d,%f)`, rec.Time.Format("2006-01-02 15")+":00:00", rec.DemandPartnerID, "other", rec.PublisherID, rec.Domain, rec.Os, rec.Country, rec.PlacementType, rec.DeviceType, rec.RequestType, rec.Size, rec.PaymentType, rec.Datacenter, rec.BidRequests, rec.BidResponses, rec.AvgBidPrice, rec.AuctionWins, rec.Auction, rec.SoldImpressions, rec.Revenue, rec.DPFee, rec.DataImpressions, rec.DataFee))
		soldImps += rec.SoldImpressions
	}

	q := fmt.Sprint(`INSERT INTO "nb_demand_hourly" ("time", "demand_partner_id", "demand_partner_placement_id", "publisher_id", "domain", "os", "country", "device_type", "placement_type","request_type","size","payment_type","datacenter","bid_requests", "bid_responses", "avg_bid_price", "auction_wins", "auction","sold_impressions","revenue","dp_fee","data_impressions","data_fee") VALUES `,
		strings.Join(values, ","),
		`ON CONFLICT ("time", "demand_partner_id", "demand_partner_placement_id", "publisher_id", "domain", "os", "country", "placement_type", "device_type","request_type","size","payment_type","datacenter") DO UPDATE SET "bid_requests" = EXCLUDED."bid_requests","bid_responses" = EXCLUDED."bid_responses","avg_bid_price" = EXCLUDED."avg_bid_price","dp_fee" = EXCLUDED."dp_fee","auction_wins" = EXCLUDED."auction_wins","auction" = EXCLUDED."auction","sold_impressions" = EXCLUDED."sold_impressions","revenue" = EXCLUDED."revenue","data_impressions" = EXCLUDED."data_impressions","data_fee" = EXCLUDED."data_fee"`)

	_, err = queries.Raw(q).Exec(tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "failed to update report nbsupply")
	}

	updatedAt := time.Now().UTC()
	updated := models.ReportUpdate{
		Report:   "new_bidder_demand_hourly",
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
	log.Info().Msg("bcdb saved")

	//err = copyToRawTable(ctx)
	//if err != nil {
	//	return errors.Wrapf(err, "failed to load csv file to dwh")
	//}

	log.Info().Int64("soldImps", soldImps).Msg("data saved to DB")

	compassRecords := ConvertToCompass(records)
	err = Send(compassRecords)
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

func (w *Worker) FetchFromQuest(ctx context.Context) ([]*models.NBDemandHourly, error) {
	var err error
	start := time.Now().UTC().Add(-1 * time.Duration(w.Hours) * time.Hour)
	stop := time.Now().UTC().Add(time.Duration(1) * time.Hour)
	if w.Start != "" {
		start, err = time.Parse("2006010215", w.Start)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse start hour")
		}

		stop = start.Add(time.Duration(1) * time.Hour)
	}

	data := make(map[string]*models.NBDemandHourly)
	log.Info().Str("start", start.Format("2006-01-02T15")+":00:00Z").Str("stop", stop.Format("2006-01-02T15")+":00:00Z").Msg("pulling data from questdb")

	err = w.processBidRequestCounters(ctx, start.Format("2006-01-02T15")+":00:00Z", stop.Format("2006-01-02T15")+":00:00Z", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query request counters")
	}

	err = w.processBidResponseCounters(ctx, start.Format("2006-01-02T15")+":00:00Z", stop.Format("2006-01-02T15")+":00:00Z", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query response counters")
	}

	err = w.processAuctionCounters(ctx, start.Format("2006-01-02T15")+":00:00Z", stop.Format("2006-01-02T15")+":00:00Z", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query auction counters")
	}

	err = w.processImpressionsCounters(ctx, start.Format("2006-01-02T15")+":00:00Z", stop.Format("2006-01-02T15")+":00:00Z", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query impressions counters")
	}

	var records []*models.NBDemandHourly
	for _, v := range data {
		if v.RequestType == "js" {
			v.PaymentType = "cpm"
		} else {
			v.PaymentType = "bid"
		}

		records = append(records, v)

	}

	return records, nil

}

func (w *Worker) FetchFromBCDB(ctx context.Context) ([]*models.NBDemandHourly, error) {
	var err error
	start := time.Now().UTC().Add(-1 * time.Duration(w.Hours) * time.Hour)
	stop := time.Now().UTC().Add(time.Duration(1) * time.Hour)
	if w.Start != "" {
		start, err = time.Parse("2006010215", w.Start)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse start hour")
		}

		stop = start.Add(time.Duration(1) * time.Hour)
	}

	data, err := models.NBDemandHourlies(
		models.NBDemandHourlyWhere.Time.GTE(start),
		models.NBDemandHourlyWhere.Time.LT(stop),
		models.NBDemandHourlyWhere.PublisherID.NEQ(""),
	).All(ctx, bcdb.DB())

	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch from bcdb")
	}

	var records []*models.NBDemandHourly
	for _, v := range data {
		records = append(records, v)
	}

	return records, nil
}

func (w *Worker) GetSleep() int {
	return int(w.Sleep.Seconds())
}

type NBDemandHourly struct {
	Time                     time.Time `csv:"time" boil:"time" json:"time" toml:"time" yaml:"time"`
	DemandPartnerID          string    `csv:"demand_partner_id" boil:"demand_partner_id" json:"demand_partner_id" toml:"demand_partner_id" yaml:"demand_partner_id"`
	DemandPartnerPlacementID string    `csv:"-" boil:"demand_partner_placement_id" json:"demand_partner_placement_id" toml:"demand_partner_placement_id" yaml:"demand_partner_placement_id"`
	PublisherID              string    `csv:"publisher_id" boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain                   string    `csv:"domain" boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Os                       string    `csv:"os" boil:"os" json:"os" toml:"os" yaml:"os"`
	Country                  string    `csv:"country" boil:"country" json:"country" toml:"country" yaml:"country"`
	DeviceType               string    `csv:"device_type" boil:"device_type" json:"device_type" toml:"device_type" yaml:"device_type"`
	PlacementType            string    `csv:"placement_type"  boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
	RequestType              string    `boil:"request_type" json:"request_type" toml:"request_type" yaml:"request_type"`
	Size                     string    `boil:"size" json:"size" toml:"size" yaml:"size"`
	PaymentType              string    `boil:"payment_type" json:"payment_type" toml:"payment_type" yaml:"payment_type"`
	Datacenter               string    `boil:"datacenter" json:"datacenter" toml:"datacenter" yaml:"datacenter"`
	BidRequests              float64   `csv:"bid_requests" boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	BidResponses             float64   `csv:"bid_responses" boil:"bid_responses" json:"bid_responses" toml:"bid_responses" yaml:"bid_responses"`
	AvgBidPrice              float64   `csv:"avg_bid_price" boil:"avg_bid_price" json:"avg_bid_price" toml:"avg_bid_price" yaml:"avg_bid_price"`
	DPFee                    float64   `csv:"dp_fee" boil:"dp_fee" json:"dp_fee" toml:"dp_fee" yaml:"dp_fee"`
	AuctionWins              float64   `csv:"auction_wins" boil:"auction_wins" json:"auction_wins" toml:"auction_wins" yaml:"auction_wins"`
	Auction                  float64   `csv:"auction" boil:"auction" json:"auction" toml:"auction" yaml:"auction"`
	SoldImpressions          float64   `csv:"sold_impressions" boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	Revenue                  float64   `csv:"revenue" boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	DataImpressions          float64   `csv:"data_impressions" boil:"data_impressions" json:"data_impressions" toml:"data_impressions" yaml:"data_impressions"`
	DataFee                  float64   `csv:"data_fee" boil:"data_fee" json:"data_fee" toml:"data_fee" yaml:"data_fee"`
}

func (rec *NBDemandHourly) Key() string {
	return fmt.Sprint(rec.Time.Format("2006-01-02-15"), rec.PublisherID, rec.Domain, rec.Os, rec.Country, rec.DemandPartnerID, rec.DeviceType, rec.PlacementType, rec.RequestType, rec.Size, rec.PaymentType, rec.Datacenter)
}

func (r *NBDemandHourly) ToModel() *models.NBDemandHourly {
	return &models.NBDemandHourly{
		Time:            r.Time,
		DemandPartnerID: r.DemandPartnerID,
		PublisherID:     r.PublisherID,
		Domain:          r.Domain,
		Os:              r.Os,
		Country:         r.Country,
		DeviceType:      r.DeviceType,
		PlacementType:   r.PlacementType,
		RequestType:     r.RequestType,
		Size:            r.Size,
		PaymentType:     r.PaymentType,
		Datacenter:      r.Datacenter,
	}
}

func (w *Worker) processBidRequestCounters(ctx context.Context, start string, stop string, data map[string]*models.NBDemandHourly) error {

	log.Info().Msg("processBidRequestCounters")

	var recordsNYC []*NBDemandHourly
	var recordsAMS []*NBDemandHourly
	var recordsSFO []*NBDemandHourly

	q := `SELECT date_trunc('hour',timestamp) as time,dtype device_type,os,country,dpid demand_partner_id,ptype placement_type,pubid publisher_id,domain,size,reqtyp request_type,'%s' datacenter,
                             sum(count) bid_requests
                             FROM demand_request_placement WHERE date_trunc('hour',timestamp)>='%s' AND date_trunc('hour',timestamp)<='%s' and size!= '' and size is not null`

	//log.Info().Str("q", q).Msg("processBidRequestCounters")
	err := queries.Raw(fmt.Sprintf(q, "nyc", start, stop)).Bind(ctx, w.QuestNYC, &recordsNYC)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb NYC1")
	}
	log.Info().Str("q", fmt.Sprintf(q, "nyc", start, stop)).Int("records", len(recordsNYC)).Msg("processImpressionsCounters NYC1")

	err = queries.Raw(fmt.Sprintf(q, "ams", start, stop)).Bind(ctx, w.QuestAMS, &recordsAMS)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb AMS3")
	}
	log.Info().Str("q", fmt.Sprintf(q, "ams", start, stop)).Int("records", len(recordsAMS)).Msg("processImpressionsCounters AMS3")

	err = queries.Raw(fmt.Sprintf(q, "sfo", start, stop)).Bind(ctx, w.QuestSFO, &recordsSFO)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb SFO3")
	}
	log.Info().Str("q", fmt.Sprintf(q, "sfo", start, stop)).Int("records", len(recordsAMS)).Msg("processImpressionsCounters SFO3")

	records := append(recordsNYC, recordsAMS...)
	records = append(records, recordsSFO...)

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

func (w *Worker) processBidResponseCounters(ctx context.Context, start string, stop string, data map[string]*models.NBDemandHourly) error {

	log.Info().Msg("processBidResponseCounters")
	var recordsNYC []*NBDemandHourly
	var recordsAMS []*NBDemandHourly
	var recordsSFO []*NBDemandHourly

	q := `SELECT date_trunc('hour',timestamp) as time,dtype device_type,os,country,dpid demand_partner_id,ptype placement_type,pubid publisher_id,domain,size,reqtyp request_type,'%s' datacenter,
                             sum(count) bid_responses,sum("sum") avg_bid_price
                             FROM demand_response_placement WHERE date_trunc('hour',timestamp)>='%s' AND date_trunc('hour',timestamp)<='%s' and size!= '' and size is not null`

	//log.Info().Str("q", q).Msg("processBidResponseCounters")
	err := queries.Raw(fmt.Sprintf(q, "nyc", start, stop)).Bind(ctx, w.QuestNYC, &recordsNYC)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb NYC1")
	}
	log.Info().Str("q", fmt.Sprintf(q, "nyc", start, stop)).Int("records", len(recordsNYC)).Msg("processImpressionsCounters NYC1")

	err = queries.Raw(fmt.Sprintf(q, "ams", start, stop)).Bind(ctx, w.QuestAMS, &recordsAMS)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb AMS3")
	}
	log.Info().Str("q", fmt.Sprintf(q, "ams", start, stop)).Int("records", len(recordsAMS)).Msg("processImpressionsCounters AMS3")

	err = queries.Raw(fmt.Sprintf(q, "sfo", start, stop)).Bind(ctx, w.QuestSFO, &recordsSFO)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb SFO3")
	}
	log.Info().Str("q", fmt.Sprintf(q, "sfo", start, stop)).Int("records", len(recordsSFO)).Msg("processImpressionsCounters SFO3")

	records := append(recordsNYC, recordsAMS...)
	records = append(records, recordsSFO...)

	for _, r := range records {
		key := r.Key()
		mod, ok := data[key]
		if !ok {
			mod = r.ToModel()
			data[key] = mod
		}
		mod.BidResponses = int64(r.BidResponses)
		mod.Auction = float64(r.BidResponses)

		if r.BidResponses > 0 {
			var dpfee float64
			if mod.DemandPartnerID == "pubmatic-pbs" || mod.DemandPartnerID == "pubmatic-sahar" {
				dpfee = 0.2
			}
			mod.AvgBidPrice = (r.AvgBidPrice * (1 - dpfee)) / r.BidResponses
		}
	}

	return nil
}

func (w *Worker) processAuctionCounters(ctx context.Context, start string, stop string, data map[string]*models.NBDemandHourly) error {
	log.Info().Msg("processAuctionCounters")

	var recordsNYC []*NBDemandHourly
	var recordsAMS []*NBDemandHourly
	var recordsSFO []*NBDemandHourly

	q := `SELECT date_trunc('hour',timestamp) as time,dtype device_type,os,country,ptype placement_type,pubid publisher_id,dpid demand_partner_id,domain,size,reqtyp request_type,'%s' datacenter,sum(count) auction_wins
                             FROM demand_wins_placement WHERE date_trunc('hour',timestamp)>='%s' AND date_trunc('hour',timestamp)<='%s' and size!= '' and size is not null`
	//log.Info().Str("q", q).Msg("processMopsCounters")
	err := queries.Raw(fmt.Sprintf(q, "nyc", start, stop)).Bind(ctx, w.QuestNYC, &recordsNYC)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb NYC1")
	}
	log.Info().Str("q", fmt.Sprintf(q, "nyc", start, stop)).Int("records", len(recordsNYC)).Msg("processImpressionsCounters NYC1")

	err = queries.Raw(fmt.Sprintf(q, "ams", start, stop)).Bind(ctx, w.QuestAMS, &recordsAMS)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb AMS3")
	}
	log.Info().Str("q", fmt.Sprintf(q, "ams", start, stop)).Int("records", len(recordsAMS)).Msg("processImpressionsCounters AMS3")

	err = queries.Raw(fmt.Sprintf(q, "sfo", start, stop)).Bind(ctx, w.QuestSFO, &recordsSFO)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb SFO3")
	}
	log.Info().Str("q", fmt.Sprintf(q, "sfo", start, stop)).Int("records", len(recordsSFO)).Msg("processImpressionsCounters SFO3")

	records := append(recordsNYC, recordsAMS...)
	records = append(records, recordsSFO...)
	for _, r := range records {
		key := r.Key()
		mod, ok := data[key]
		if !ok {
			mod = r.ToModel()
			data[key] = mod
		}

		mod.AuctionWins += int64(r.AuctionWins)

	}

	return nil
}

func (w *Worker) processImpressionsCounters(ctx context.Context, start string, stop string, data map[string]*models.NBDemandHourly) error {
	log.Info().Msg("processImpressionsCounters")

	var recordsNYC []*NBDemandHourly
	var recordsAMS []*NBDemandHourly
	var recordsSFO []*NBDemandHourly

	q := `  SELECT date_trunc('hour',timestamp) as time,
             dtype device_type,os,country,ptype placement_type,publisher publisher_id,domain,dpid demand_partner_id,size,reqtyp request_type,'%s' datacenter,
             sum(dbpr)/1000 revenue,sum(dpfee)/1000 dp_fee ,count(1) sold_impressions,sum(case when uidsrc='iiq' then 1 else 0 end) data_impressions,sum(case when uidsrc='iiq' then dbpr/1000 else 0 end) data_fee
                             FROM impression WHERE date_trunc('hour',timestamp)>='%s' AND date_trunc('hour',timestamp)<='%s'`
	//log.Info().Str("q", q).Msg("processImpressionsCounters")
	err := queries.Raw(fmt.Sprintf(q, "nyc", start, stop)).Bind(ctx, w.QuestNYC, &recordsNYC)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb NYC1")
	}
	log.Info().Str("q", fmt.Sprintf(q, "nyc", start, stop)).Int("records", len(recordsNYC)).Msg("processImpressionsCounters NYC1")

	err = queries.Raw(fmt.Sprintf(q, "ams", start, stop)).Bind(ctx, w.QuestAMS, &recordsAMS)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb AMS3")
	}
	log.Info().Str("q", fmt.Sprintf(q, "ams", start, stop)).Int("records", len(recordsAMS)).Msg("processImpressionsCounters AMS3")

	err = queries.Raw(fmt.Sprintf(q, "sfo", start, stop)).Bind(ctx, w.QuestSFO, &recordsSFO)
	if err != nil {
		return errors.Wrapf(err, "failed to query impressions from questdb SFO3")
	}
	log.Info().Str("q", fmt.Sprintf(q, "sfo", start, stop)).Int("records", len(recordsSFO)).Msg("processImpressionsCounters SFO3")

	records := append(recordsNYC, recordsAMS...)
	records = append(records, recordsSFO...)
	for _, r := range records {
		key := r.Key()
		mod, ok := data[key]
		if !ok {
			mod = r.ToModel()
			data[key] = mod
		}

		mod.Revenue += r.Revenue
		mod.DPFee += r.DPFee
		mod.SoldImpressions += int64(r.SoldImpressions)
		mod.DataImpressions += int64(r.DataImpressions)
		mod.DataFee += r.DataFee

		//log.Info().Interface("rec", mod).Msgf("rec")

	}

	return nil
}
