package real_time_report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"net/http"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/quest"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/m6yf/bcwork/utils/helpers"

	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"golang.org/x/net/context"
)

type RealTimeReport struct {
	Time                 string  `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string  `boil:"pubid" json:"pubid" toml:"pubid" yaml:"pubid"`
	Publisher            string  `boil:"pub" json:"pub" toml:"pub" yaml:"pub"`
	Domain               string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	BidRequests          float64 `boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	BidResponses         float64 `boil:"bid_response" json:"bid_response" toml:"bid_response" yaml:"bid_response"`
	Device               string  `boil:"dtype" json:"dtype" toml:"dtype" yaml:"dtype"`
	Country              string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Revenue              float64 `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	Cost                 float64 `boil:"cost" json:"cost" toml:"cost" yaml:"cost"`
	SoldImpressions      float64 `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	PublisherImpressions float64 `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
	PubFillRate          float64 `boil:"fill_rate" json:"fill_rate" toml:"fill_rate" yaml:"fill_rate"`
	CPM                  float64 `boil:"cpm" json:"cpm" toml:"cpm" yaml:"cpm"`
	RPM                  float64 `boil:"rpm" json:"rpm" toml:"rpm" yaml:"rpm"`
	DpRPM                float64 `boil:"dp_rpm" json:"dp_rpm" toml:"dp_rpm" yaml:"dp_rpm"`
	GP                   float64 `boil:"gp" json:"gp" toml:"gp" yaml:"gp"`
	GPP                  float64 `boil:"gpp" json:"gpp" toml:"gpp" yaml:"gpp"`
	ConsultantFee        float64 `boil:"consultant_fee" json:"consultant_fee" toml:"consultant_fee" yaml:"consultant_fee"`
	TamFee               float64 `boil:"tam_fee" json:"tam_fee" toml:"tam_fee" yaml:"tam_fee"`
	TechFee              float64 `boil:"tech_fee" json:"tech_fee" toml:"tech_fee" yaml:"tech_fee"`
	DemandPartnerFee     float64 `boil:"demand_partner_fee" json:"demand_partner_fee" toml:"demand_partner_fee" yaml:"demand_partner_fee"`
	DataFee              float64 `boil:"data_fee" json:"data_fee" toml:"data_fee" yaml:"data_fee"`
}

type DBRealTimeReport struct {
	Time                 string  `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string  `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Publisher            string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain               string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	BidRequests          float64 `boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	BidResponses         float64 `boil:"bid_responses" json:"bid_responses" toml:"bid_responses" yaml:"bid_responses"`
	Device               string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Country              string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Revenue              float64 `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	Cost                 float64 `boil:"cost" json:"cost" toml:"cost" yaml:"cost"`
	SoldImpressions      float64 `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	PublisherImpressions float64 `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
	PubFillRate          float64 `boil:"fill_rate" json:"fill_rate" toml:"fill_rate" yaml:"fill_rate"`
	CPM                  float64 `boil:"cpm" json:"cpm" toml:"cpm" yaml:"cpm"`
	RPM                  float64 `boil:"rpm" json:"rpm" toml:"rpm" yaml:"rpm"`
	DpRPM                float64 `boil:"dp_rpm" json:"dp_rpm" toml:"dp_rpm" yaml:"dp_rpm"`
	GP                   float64 `boil:"gp" json:"gp" toml:"gp" yaml:"gp"`
	GPP                  float64 `boil:"gpp" json:"gpp" toml:"gpp" yaml:"gpp"`
	ConsultantFee        float64 `boil:"consultant_fee" json:"consultant_fee" toml:"consultant_fee" yaml:"consultant_fee"`
	TamFee               float64 `boil:"tam_fee" json:"tam_fee" toml:"tam_fee" yaml:"tam_fee"`
	TechFee              float64 `boil:"tech_fee" json:"tech_fee" toml:"tech_fee" yaml:"tech_fee"`
	DemandPartnerFee     float64 `boil:"demand_partner_fee" json:"demand_partner_fee" toml:"demand_partner_fee" yaml:"demand_partner_fee"`
	DataFee              float64 `boil:"data_fee" json:"data_fee" toml:"data_fee" yaml:"data_fee"`
}

func (worker *Worker) FetchAndMergeQuestReports(ctx context.Context) (map[string]*RealTimeReport, error) {
	var impressionsRecords []*RealTimeReport
	var bidRequestRecords []*RealTimeReport
	var bidResponseRecords []*RealTimeReport
	var err error

	worker.Fees, worker.ConsultantFees, err = worker.FetchFees(ctx)
	if err != nil {
		return nil, err
	}

	startString := worker.Start.Format(constant.PostgresTimestamp)
	endString := worker.End.Format(constant.PostgresTimestamp)

	bidRequestsQuery := fmt.Sprintf(QuestRequests, startString, endString)
	log.Info().Str("query", bidRequestsQuery).Msg("processBidRequestsCounters")

	impressionsQuery := fmt.Sprintf(QuestImpressions, startString, endString)
	log.Info().Str("query", impressionsQuery).Msg("processImpressionsCounters")

	bidResponseQuery := fmt.Sprintf(QuestBidResponse, startString, endString)
	log.Info().Str("query", bidResponseQuery).Msg("processBidResponseCounters")

	impressionsMap := make(map[string]*RealTimeReport)
	bidRequestMap := make(map[string]*RealTimeReport)
	bidResponseMap := make(map[string]*RealTimeReport)

	for _, instance := range worker.Quest {
		if err := quest.InitDB(instance); err != nil {
			return nil, fmt.Errorf("failed to initialize Quest instance: %s", instance)
		}

		if err := queries.Raw(impressionsQuery).Bind(ctx, quest.DB(), &impressionsRecords); err != nil {
			return nil, fmt.Errorf("failed to query impressions from Quest instance: %s", instance)
		}

		if err := queries.Raw(bidRequestsQuery).Bind(ctx, quest.DB(), &bidRequestRecords); err != nil {
			return nil, fmt.Errorf("failed to query bid requests from Quest instance: %s", instance)
		}

		if err := queries.Raw(bidResponseQuery).Bind(ctx, quest.DB(), &bidResponseRecords); err != nil {
			return nil, fmt.Errorf("failed to query bid response from Quest instance: %s", instance)
		}

		bidRequestMap = worker.GenerateBidRequestMap(bidRequestMap, bidRequestRecords)
		bidResponseMap = worker.GenerateBidResponseMap(bidResponseMap, bidResponseRecords)
		impressionsMap = worker.GenerateImpressionsMap(impressionsMap, impressionsRecords)

		impressionsRecords = nil
		bidRequestRecords = nil
	}

	worker.Publishers, _ = FetchPublishers(context.Background(), worker)
	return worker.MergeReports(bidRequestMap, impressionsMap, bidResponseMap)
}

func (worker *Worker) MergeReports(bidRequestMap map[string]*RealTimeReport, impressionsMap map[string]*RealTimeReport, bidResponseMap map[string]*RealTimeReport) (map[string]*RealTimeReport, error) {
	reportMap := make(map[string]*RealTimeReport)

	for _, record := range impressionsMap {
		key := record.Key()
		requestsItem, existsRequest := bidRequestMap[key]
		responseItem, existsResponse := bidResponseMap[key]
		publisherName, _ := worker.Publishers[record.PublisherID]

		mergedRecord := &RealTimeReport{
			Time:                 record.Time,
			PublisherID:          record.PublisherID,
			Publisher:            publisherName,
			Domain:               record.Domain,
			Country:              record.Country,
			Device:               record.Device,
			Revenue:              record.Revenue,
			Cost:                 record.Cost,
			SoldImpressions:      record.SoldImpressions,
			PublisherImpressions: record.PublisherImpressions,
			DataFee:              record.DataFee,
		}

		if existsRequest {
			mergedRecord.BidRequests = requestsItem.BidRequests
		} else {
			mergedRecord.BidRequests = 0
		}
		if existsResponse {
			mergedRecord.BidResponses = responseItem.BidResponses
		} else {
			mergedRecord.BidResponses = 0
		}

		mergedRecord.PubFillRate = constant.PubFillRate(int64(mergedRecord.PublisherImpressions), int64(mergedRecord.BidRequests))
		mergedRecord.CPM = constant.CPM(mergedRecord.Cost, mergedRecord.PublisherImpressions)
		mergedRecord.RPM = constant.RPM(mergedRecord.Revenue, mergedRecord.PublisherImpressions)
		mergedRecord.DpRPM = constant.DpRPM(mergedRecord.Revenue, mergedRecord.SoldImpressions)

		mergedRecord.CalculateGP(worker.Fees, worker.ConsultantFees)

		reportMap[key] = mergedRecord
	}

	return reportMap, nil
}

func (record *RealTimeReport) Key() string {
	return fmt.Sprintf("%s - %s - %s - %s - %s", record.PublisherID, record.Domain, record.Device, record.Country, record.Time)
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
				SoldImpressions:      record.SoldImpressions + item.SoldImpressions,
				PublisherImpressions: record.PublisherImpressions + item.PublisherImpressions,
				DataFee:              record.DataFee,
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

func (worker *Worker) GenerateBidResponseMap(bidResponseMap map[string]*RealTimeReport, bidResponseRecords []*RealTimeReport) map[string]*RealTimeReport {
	for _, record := range bidResponseRecords {
		key := record.Key()
		item, exists := bidResponseMap[key]
		if exists {
			mergedItem := &RealTimeReport{
				Time:         record.Time,
				PublisherID:  record.PublisherID,
				Domain:       record.Domain,
				Country:      record.Country,
				Device:       record.Device,
				BidResponses: item.BidResponses + record.BidResponses,
			}
			bidResponseMap[key] = mergedItem
		} else {
			bidResponseMap[key] = record
		}
	}
	return bidResponseMap
}

func (rec *RealTimeReport) CalculateGP(fees map[string]float64, consultantFees map[string]float64) {
	rec.TamFee = helpers.RoundFloat((fees["tam_fee"] * rec.Cost))
	rec.TechFee = helpers.RoundFloat(fees["tech_fee"] * rec.BidRequests / constant.ConversionToMillion)
	rec.ConsultantFee = 0.0
	value, exists := consultantFees[rec.PublisherID]
	if exists {
		rec.ConsultantFee = rec.Cost * value
	}

	rec.GP = helpers.RoundFloat(rec.Revenue - rec.Cost - rec.DemandPartnerFee - rec.DataFee - rec.TamFee - rec.TechFee - rec.ConsultantFee)
	rec.GPP = 0
	if rec.Revenue != 0 {
		rec.GPP = helpers.RoundFloat(rec.GP / rec.Revenue)
	}

}

func FetchPublishers(ctx context.Context, worker *Worker) (map[string]string, error) {

	requestBody := map[string]interface{}{
		"filter": map[string]interface{}{},
	}

	body, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}

	publisherData, statusCode, err := worker.HttpClient.Do(ctx, http.MethodPost, constant.ProductionApiUrl+"/publisher/get", bytes.NewBuffer(body))

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", statusCode)
	}

	var publishers []dto.Publisher
	if err := json.Unmarshal(publisherData, &publishers); err != nil {
		return nil, fmt.Errorf("error parsing publisher data  from API")
	}

	publisherMap := make(map[string]string)
	for _, publisher := range publishers {
		publisherMap[publisher.PublisherID] = publisher.Name
	}

	return publisherMap, nil
}

func (worker *Worker) FetchFees(ctx context.Context) (map[string]float64, map[string]float64, error) {
	log.Info().Msg("fetch global fees for Full Publisher Requests")

	requestBody := map[string]interface{}{}
	jsonData, err := json.Marshal(requestBody)

	if err != nil {
		return nil, nil, fmt.Errorf("error creating fees request body for Full Publisher Requests")
	}

	data, statusCode, err := worker.HttpClient.Do(ctx, http.MethodPost, constant.ProductionApiUrl+constant.GlobalFactorEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, fmt.Errorf("error Fetching fees from API")
	}

	if statusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("error fetching fees from API. Request failed with status code: %d", statusCode)
	}

	var FeesResponse []*dto.GlobalFactor
	if err := json.Unmarshal(data, &FeesResponse); err != nil {
		return nil, nil, fmt.Errorf("error parsing fees from API")
	}

	// Collect fee rates
	fees := make(map[string]float64)
	consultantFees := make(map[string]float64)

	for _, item := range FeesResponse {
		if item.Key == constant.GlobalFactorConsultantFeeType && item.PublisherID != "" {
			consultantFees[item.PublisherID] = item.Value
		} else if item.Key == constant.GlobalFactorTAMFeeType {
			fees[item.Key] = 0
		} else if item.Key == constant.GlobalFactorTechFeeType {
			fees[item.Key] = item.Value
		}
	}

	return fees, consultantFees, nil
}

func (worker *Worker) FetchRealTimeData(ctx context.Context) ([]*DBRealTimeReport, error) {
	var records []*DBRealTimeReport

	startString := worker.Start.Format(constant.PostgresTimestamp)
	endString := worker.End.Format(constant.PostgresTimestamp)

	realTimeReportQuery := fmt.Sprintf(RealTimeReportQuery, startString, endString)
	log.Info().Str("query", realTimeReportQuery).Msg("processBidRequestsCounters")

	if err := queries.Raw(realTimeReportQuery).Bind(ctx, bcdb.DB(), &records); err != nil {
		return nil, fmt.Errorf("failed to query real_time_request from postgres: %w", err)
	}

	return records, nil
}

func (worker *Worker) RemoveOldDataFromDB(ctx context.Context) error {

	realTimeReportDeleteQuery := fmt.Sprintf(DeleteQuery, worker.Start)
	log.Info().Msg("delete query: " + realTimeReportDeleteQuery)

	_, err := queries.Raw(realTimeReportDeleteQuery).Exec(bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to remove old data from real_time_report table: %w", err)
	}
	return nil
}
