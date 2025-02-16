package nbsupply

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdwh"
	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

const csvFileName = "/tmp/nbsupply.csv"

func createCSV(records []*models.NBSupplyHourly) error {
	file, err := os.Create(csvFileName)
	if err != nil {
		return errors.Wrapf(err, "failed to open csv file nbsupply")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, rec := range records {
		if rec.PublisherID == "" {
			rec.PublisherID = "0"
		}
		row := []string{
			rec.Time.Format("2006-01-02 15") + ":00:00",
			rec.PublisherID,
			rec.Domain,
			rec.Os,
			rec.Country,
			rec.DeviceType,
			rec.PlacementType,
			rec.RequestType,
			rec.Size,
			rec.PaymentType,
			strconv.FormatInt(rec.BidRequests, 10),
			strconv.FormatInt(rec.BidResponses, 10),
			strconv.FormatInt(rec.SoldImpressions, 10),
			strconv.FormatInt(rec.PublisherImpressions, 10),
			strconv.FormatFloat(rec.Cost, 'f', 8, 64),
			strconv.FormatFloat(rec.Revenue, 'f', 8, 64),
			strconv.FormatFloat(rec.AvgBidPrice, 'f', 8, 64),
			strconv.FormatInt(rec.MissedOpportunities, 10),
			strconv.FormatFloat(rec.DemandPartnerFee, 'f', 8, 64),
			strconv.FormatInt(rec.PublisherImpressions, 10),
			strconv.FormatFloat(rec.DataFee, 'f', 8, 64),
		}

		if err := writer.Write(row); err != nil {
			return errors.Wrapf(err, "error writing record to csv file")
		}
	}
	writer.Flush()

	return nil
}

func scpCSV() error {
	user := "root"
	host := "bcdwh-nyc-01"
	source := csvFileName
	destination := csvFileName

	cmd := exec.Command("scp", source, user+"@"+host+":"+destination)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(out)

		return errors.Wrapf(err, "failed to scp csv file")
	}

	return nil
}

func loadToRawTable(ctx context.Context) error {
	_, err := queries.Raw(`copy supply_raw ("time", "publisher_id", "domain", "os", "country", "device_type", "placement_type", "request_type","size","payment_type", "bid_requests", "bid_responses","sold_impressions","publisher_impressions","cost","revenue","avg_bid_price","missed_opportunities","demand_partner_fee","data_impressions","data_fee") from '`+csvFileName+`' with (delimiter ',', format csv)`).ExecContext(ctx, bcdwh.DB())
	if err != nil {
		return errors.Wrapf(err, "failed load csv to raw table")
	}

	return nil
}

func updateFromRawToHourly() error {
	_, err := queries.Raw(`INSERT INTO supply_hourly 
	SELECT * FROM supply_raw
	ON CONFLICT ("time", "publisher_id", "domain", "os", "country", "device_type", "placement_type","request_type","size","payment_type")
	DO UPDATE SET 
        publisher_impressions=EXCLUDED.publisher_impressions,
        sold_impressions=EXCLUDED.sold_impressions,
		bid_requests=EXCLUDED.bid_requests,
        bid_responses=EXCLUDED.bid_responses,
        avg_bid_price=EXCLUDED.avg_bid_price,
        missed_opportunities=EXCLUDED.missed_opportunities,
        cost=EXCLUDED.cost,
		revenue=EXCLUDED.revenue,
		demand_partner_fee=EXCLUDED.demand_partner_fee,
		data_fee=EXCLUDED.data_fee,
		data_impressions=EXCLUDED.data_impressions`).Exec(bcdwh.DB())

	if err != nil {
		return errors.Wrapf(err, "failed to update from supply raw to hourly")
	}

	return nil
}

func createTempTable() error {
	_, err := queries.Raw(`
create table if not exists supply_raw
(
    time                  timestamp                                       not null,
    publisher_id          varchar(36)                                     not null,
    domain                varchar(256)     default '-'::character varying not null,
    os                    varchar(64)      default '-'::character varying not null,
    country               varchar(64)      default '-'::character varying not null,
    device_type           varchar(64)      default '-'::character varying not null,
    placement_type        varchar(16)      default '-'::character varying not null,
    request_type           varchar(16)      not null default '-',
    size           varchar(16)      not null default '-',
    payment_type           varchar(16)      not null default '-',
    bid_requests          bigint           default 0                      not null,
    bid_responses         bigint           default 0                      not null,
    sold_impressions      bigint           default 0                      not null,
    publisher_impressions bigint           default 0                      not null,
    cost                  double precision default 0                      not null,
    revenue               double precision default 0                      not null,
    avg_bid_price         double precision default 0                      not null,
    missed_opportunities  bigint           default 0                      not null,
    demand_partner_fee    double precision default 0                      not null,
    data_impressions      bigint           default 0                      not null,
    data_fee              double precision default 0                      not null
)`).Exec(bcdwh.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to create temp table")
	}

	return nil
}

func dropTempTable() error {
	_, err := queries.Raw(`drop table supply_raw`).Exec(bcdwh.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to drop temp table")
	}

	return nil
}
