package nbdemand

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

const csvFileName = "/tmp/nbdemand.csv"

func createCSV(records []*models.NBDemandHourly) error {
	file, err := os.Create(csvFileName)
	if err != nil {
		return errors.Wrapf(err, "failed to open csv file nbdemand")
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
			rec.DemandPartnerID,
			"other",
			rec.PublisherID,
			rec.Domain,
			rec.Os,
			rec.Country,
			rec.PlacementType,
			rec.DeviceType,
			rec.RequestType,
			rec.Size,
			rec.PaymentType,
			strconv.FormatInt(rec.BidRequests, 10),
			strconv.FormatInt(rec.BidResponses, 10),
			strconv.FormatFloat(rec.AvgBidPrice, 'f', 8, 64),
			strconv.FormatInt(rec.AuctionWins, 10),
			strconv.FormatFloat(rec.Auction, 'f', 8, 64),
			strconv.FormatInt(rec.SoldImpressions, 10),
			strconv.FormatFloat(rec.Revenue, 'f', 8, 64),
			strconv.FormatFloat(rec.DPFee, 'f', 8, 64),
			strconv.FormatInt(rec.DataImpressions, 10),
			strconv.FormatFloat(rec.DataFee, 'f', 8, 64)}
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
	_, err := queries.Raw(`copy demand_raw 
     ("time", "demand_partner_id", "demand_partner_placement_id", "publisher_id", "domain", "os", "country", "device_type", "placement_type", "request_type","size","payment_type", "bid_requests", "bid_responses", "avg_bid_price", "auction_wins", "auction","sold_impressions","revenue","dp_fee","data_impressions","data_fee") from '`+csvFileName+`' with (delimiter ',', format csv)`).ExecContext(ctx, bcdwh.DB())
	if err != nil {
		return errors.Wrapf(err, "failed load csv to raw table")
	}

	return nil
}

func updateFromRawToHourly() error {
	_, err := queries.Raw(`INSERT INTO demand_hourly 
	SELECT * FROM demand_raw
	ON CONFLICT ("time", "demand_partner_id", "demand_partner_placement_id", "publisher_id", "domain", "os", "country", "device_type", "placement_type","request_type","size","payment_type")
	DO UPDATE SET 
		sold_impressions=EXCLUDED.sold_impressions,
		bid_requests=EXCLUDED.bid_requests,
        bid_responses=EXCLUDED.bid_responses,
        avg_bid_price=EXCLUDED.avg_bid_price,
        auction_wins=EXCLUDED.auction_wins,
        auction=EXCLUDED.auction,
		revenue=EXCLUDED.revenue,
		dp_fee=EXCLUDED.dp_fee,
		data_fee=EXCLUDED.data_fee,
		data_impressions=EXCLUDED.data_impressions`).Exec(bcdwh.DB())

	if err != nil {
		return errors.Wrapf(err, "failed to update from demand raw to hourly")
	}

	return nil
}

func createTempTable() error {
	_, err := queries.Raw(`create table if not exists demand_raw 
(
    time                     timestamp        not null,
    demand_partner_id        varchar(36)      not null,
    demand_partner_placement_id        varchar(36)      not null,
    publisher_id             varchar(36)      not null,
    domain                   varchar(256)     not null default '-',
    os                       varchar(64)      not null default '-',
    country                  varchar(64)      not null default '-',
    device_type              varchar(64)      not null default '-',
    placement_type           varchar(16)      not null default '-',
    request_type           varchar(16)      not null default '-',
    size           varchar(16)      not null default '-',
    payment_type           varchar(16)      not null default '-',
    bid_requests             int8             not null default 0,
    bid_responses           int8             not null default 0,
    avg_bid_price            double precision not null default 0,
    dp_fee                   double precision not null default 0,
    auction_wins             int8             not null default 0,
    auction                  double precision not null default 0,
    sold_impressions         int8             not null default 0,
    revenue                  double precision not null default 0,
    data_impressions         int8 not null default 0,
    data_fee                  double precision not null default 0
)`).Exec(bcdwh.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to create temp table")
	}

	return nil
}

func dropTempTable() error {
	_, err := queries.Raw(`drop table demand_raw`).Exec(bcdwh.DB())

	if err != nil {
		return errors.Wrapf(err, "failed to drop temp table")
	}

	return nil
}
