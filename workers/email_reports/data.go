package email_reports

import (
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/quest"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"golang.org/x/net/context"
	"os"
	"time"
)

type RealTimeReport struct {
	Time                 string  `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string  `boil:"pubid" json:"pubid" toml:"pubid" yaml:"pubid"`
	Domain               string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	BidRequests          float64 `boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	Device               string  `boil:"dtype" json:"dtype" toml:"dtype" yaml:"dtype"`
	Country              string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Revenue              string  `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	Cost                 string  `boil:"cost" json:"cost" toml:"cost" yaml:"cost"`
	DemandPartnerFee     string  `boil:"demand_partner_fee" json:"demand_partner_fee" toml:"demand_partner_fee" yaml:"demand_partner_fee"`
	SoldImpressions      string  `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	PublisherImpressions string  `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
}

var QuestRequests = `SELECT to_date('%s','yyyy-MM-ddTHH:mm:ssZ') time,
  pubid,
  domain,
  country,
  dtype,
  sum(count) bid_requests
FROM
  request_placement
WHERE timestamp >= '%s'
  AND timestamp < '%s'
  AND dtype is not null
  AND country is not null
  AND pubid is not null
  AND domain is not null
GROUP BY 1,2,3,4,5`

var QuestImpressions = `SELECT to_date('%s','yyyy-MM-ddTHH:mm:ssZ') time,
       publisher pubid,
       domain,
       country,
       dtype,
       sum(dbpr)/1000 revenue,
       sum(sbpr)/1000 cost,
       sum(dpfee)/1000 demand_partner_fee,
       count(1) sold_impressions,
       sum(CASE WHEN loop=false THEN 1 ELSE 0 END) publisher_impressions
FROM impression
WHERE timestamp >= '%s'
  AND timestamp < '%s'
  AND publisher IS NOT NULL
  AND domain IS NOT NULL
  AND country IS NOT NULL
  AND dtype IS NOT NULL
GROUP BY 1, 2, 3, 4, 5`

func (worker *Worker) FetchFromQuest(ctx context.Context, start time.Time, end time.Time) (map[string]*RealTimeReport, error) {
	var realTimeRecordsSlice []*RealTimeReport
	var impressionsRecords []*RealTimeReport
	var bidRequestRecords []*RealTimeReport

	startString := start.Format("2006-01-02T15:04:05Z")
	endString := end.Format("2006-01-02T15:04:05Z")

	bidRequestsQuery := fmt.Sprintf(QuestRequests, startString, startString, endString)
	log.Info().Str("query", bidRequestsQuery).Msg("processBidRequestsCounters")

	impressionsQuery := fmt.Sprintf(QuestImpressions, startString, startString, endString)
	log.Info().Str("query", impressionsQuery).Msg("processImpressionsCounters")

	impressionsMap := make(map[string]*RealTimeReport)
	bidRequestMap := make(map[string]*RealTimeReport)

	for _, instance := range worker.Quest {
		if err := quest.InitDB(instance); err != nil {
			return nil, errors.Wrapf(err, "failed to initialize Quest instance: %s", instance)
		}

		// Fetch impressions
		if err := queries.Raw(impressionsQuery).Bind(ctx, quest.DB(), &impressionsRecords); err != nil {
			return nil, errors.Wrapf(err, "failed to query impressions from Quest instance: %s", instance)
		}

		// Fetch bid requests
		if err := queries.Raw(bidRequestsQuery).Bind(ctx, quest.DB(), &bidRequestRecords); err != nil {
			return nil, errors.Wrapf(err, "failed to query bid requests from Quest instance: %s", instance)
		}

		// Generate maps from fetched records
		bidRequestMap = worker.GenerateBidRequestMap(bidRequestMap, bidRequestRecords)
		impressionsMap = worker.GenerateImpressionsMap(impressionsMap, impressionsRecords)

		// Clear records for the next iteration
		impressionsRecords = nil
		bidRequestRecords = nil

		// Save results to file
		filename := fmt.Sprintf("real_time_records_%s.json", instance)
		if err := saveResultsToFile(realTimeRecordsSlice, filename); err != nil {
			return nil, errors.Wrapf(err, "failed to save results to file for instance: %s", instance)
		}
	}

	// Return merged reports
	return worker.MergeReports(bidRequestMap, impressionsMap)
}

// todo delete before push
func saveResultsToFile(data []*RealTimeReport, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Failed to marshal data to JSON")
	}

	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to create file: %s", filename)
	}
	defer file.Close()

	// Write JSON data to the file
	if _, err := file.Write(jsonData); err != nil {
		return errors.Wrapf(err, "Failed to write data to file: %s", filename)
	}

	fmt.Printf("Data for instance saved to %s\n", filename)
	return nil
}

func (worker *Worker) MergeReports(bidRequestMap map[string]*RealTimeReport, impressionsMap map[string]*RealTimeReport) (map[string]*RealTimeReport, error) {
	reportMap := make(map[string]*RealTimeReport)
	var err error
	for _, record := range impressionsMap {
		key := record.Key()
		requestsItem, exists := bidRequestMap[key]
		if exists {
			mergedRecord := &RealTimeReport{
				Time:                 record.Time,
				PublisherID:          record.PublisherID,
				Domain:               record.Domain,
				Country:              record.Country,
				Device:               record.Device,
				Revenue:              record.Revenue,
				Cost:                 record.Cost,
				DemandPartnerFee:     record.DemandPartnerFee,
				SoldImpressions:      record.SoldImpressions,
				PublisherImpressions: record.PublisherImpressions,
				BidRequests:          requestsItem.BidRequests,
			}
			reportMap[key] = mergedRecord
		} else {
			reportMap[key] = record
		}
	}

	return reportMap, err
}

func (record *RealTimeReport) Key() string {
	return fmt.Sprintf("%s - %s - %s - %s", record.PublisherID, record.Domain, record.Device, record.Country)
}

func (worker *Worker) GenerateImpressionsMap(impressionsMap map[string]*RealTimeReport, impressionsRecords []*RealTimeReport) map[string]*RealTimeReport {
	for _, record := range impressionsRecords {

		key := record.Key()
		item, exists := impressionsMap[key]
		if exists {
			mergedItem := &RealTimeReport{
				Time:                 record.Time,
				PublisherID:          record.PublisherID,
				Domain:               record.Domain,
				Country:              record.Country,
				Device:               record.Device,
				Revenue:              record.Revenue + item.Revenue,
				Cost:                 record.Cost + item.Cost,
				DemandPartnerFee:     record.DemandPartnerFee + item.DemandPartnerFee,
				SoldImpressions:      record.SoldImpressions + item.SoldImpressions,
				PublisherImpressions: record.PublisherImpressions + item.PublisherImpressions,
			}
			impressionsMap[key] = mergedItem
		} else {
			impressionsMap[key] = record
		}

	}
	return impressionsMap
}

func (worker *Worker) GenerateBidRequestMap(bidRequestMap map[string]*RealTimeReport, bidRequestRecords []*RealTimeReport) map[string]*RealTimeReport {
	for _, record := range bidRequestRecords {
		key := record.Key()
		item, exists := bidRequestMap[key]
		if exists {
			mergedItem := &RealTimeReport{
				Time:        record.Time,
				PublisherID: record.PublisherID,
				Domain:      record.Domain,
				Country:     record.Country,
				Device:      record.Device,
				BidRequests: item.BidRequests + record.BidRequests,
			}
			bidRequestMap[key] = mergedItem
		} else {
			bidRequestMap[key] = record
		}
	}
	return bidRequestMap
}
