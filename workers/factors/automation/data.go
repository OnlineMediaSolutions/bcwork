package factors_autmation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/m6yf/bcwork/core/bulk"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/quest"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

var QuestImpressions = `SELECT to_date('%s','yyyy-MM-ddTHH:mm:ssZ') time,
       publisher publisher_id,
       domain,
       country,
       dtype device_type,
       sum(dbpr)/1000 revenue,
       sum(sbpr)/1000 cost,
       sum(dpfee)/1000 demand_partner_fee,
       count(1) sold_impressions,
       sum(CASE WHEN loop=false THEN 1 ELSE 0 END) publisher_impressions,
       sum(CASE WHEN uidsrc='iiq' THEN dbpr/1000 ELSE 0 END) * 0.045 data_fee
FROM impression
WHERE timestamp >= '%s'
  AND timestamp < '%s'
  AND publisher IS NOT NULL
  AND domain IS NOT NULL
  AND country IS NOT NULL
  AND dtype IS NOT NULL
  AND domain in('%s')
GROUP BY 1, 2, 3, 4, 5`

var QuestRequests = `SELECT to_date('%s','yyyy-MM-ddTHH:mm:ssZ') time,
  pubid publisher_id,
  domain,
  country,
  dtype device_type,
  sum(count) bid_requests
FROM
  request_placement
WHERE timestamp >= '%s'
  AND timestamp < '%s'
  AND dtype is not null
  AND country is not null
  AND pubid is not null
  AND domain is not null
  AND domain in('%s')
GROUP BY 1,2,3,4,5`

var inactiveKeysQuery = `SELECT *
FROM (SELECT publisher, domain, country, device, SUM(CASE WHEN new_factor >= %f THEN 1 ELSE 0 END) AS positive_cases
		FROM public.price_factor_log
		WHERE eval_time >= TO_TIMESTAMP('%s','YYYY-MM-DDTHH24:MI:SS')
		GROUP BY 1, 2, 3, 4) AS t
WHERE positive_cases < 1;`

func (worker *Worker) FetchData(ctx context.Context) (map[string]*FactorReport, map[string]*Factor, error) {
	var recordsMap map[string]*FactorReport
	var factors map[string]*Factor
	var err error

	worker.Domains, err = worker.FetchAutomationSetup(ctx)
	if err != nil {
		return nil, nil, err
	}

	worker.Fees, worker.ConsultantFees, err = worker.FetchFees(ctx)
	if err != nil {
		return nil, nil, err
	}

	recordsMap, err = worker.FetchFromQuest(ctx, worker.Start, worker.End)
	if err != nil {
		return nil, nil, err
	}

	factors, err = worker.FetchFactors(ctx)
	if err != nil {
		return nil, nil, err
	}

	worker.InactiveKeys, err = worker.FetchInactiveKeys(ctx)
	if err != nil {
		return nil, nil, err
	}

	return recordsMap, factors, nil
}

// Fetch performance data from quest
func (worker *Worker) FetchFromQuest(ctx context.Context, start time.Time, stop time.Time) (map[string]*FactorReport, error) {
	log.Debug().Msg("fetch records from QuestDB")
	var impressionsRecords []*FactorReport
	var bidRequestRecords []*FactorReport

	startString := start.Format("2006-01-02T15:04:05Z")
	stopString := stop.Format("2006-01-02T15:04:05Z")
	domains := worker.AutomationDomains()

	impressionsQuery := fmt.Sprintf(QuestImpressions, startString, startString, stopString, strings.Join(domains, "', '"))
	log.Info().Str("query", impressionsQuery).Msg("processImpressionsCounters")

	bidRequestQuery := fmt.Sprintf(QuestRequests, startString, startString, stopString, strings.Join(domains, "', '"))
	log.Info().Str("query", bidRequestQuery).Msg("processRequestCounters")

	var impressionsMap = make(map[string]*FactorReport)
	var bidRequestMap = make(map[string]*FactorReport)
	for _, instance := range worker.Quest {
		err := quest.InitDB(instance)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("Failed to initialize Quest instance: %s", instance))
		}

		err = queries.Raw(impressionsQuery).Bind(ctx, quest.DB(), &impressionsRecords)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("failed to query impressions from Quest instance: %s", instance))
		}

		err = queries.Raw(bidRequestQuery).Bind(ctx, quest.DB(), &bidRequestRecords)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("failed to query requests from Quest instance: %s", instance))
		}

		bidRequestMap = worker.GenerateBidRequestMap(bidRequestMap, bidRequestRecords)
		impressionsMap = worker.GenerateImpressionsMap(impressionsMap, impressionsRecords)

		impressionsRecords = nil
		bidRequestRecords = nil

		err = quest.CloseDB()
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("Failed to close Quest instance: %s", instance))
		}

	}
	return worker.MergeReports(bidRequestMap, impressionsMap)
}

func (worker *Worker) FetchFactors(ctx context.Context) (map[string]*Factor, error) {
	log.Debug().Msg("fetch records from Factors API")
	// Create the request body using a map
	requestBody := map[string]interface{}{
		"filter": map[string][]string{"domain": worker.AutomationDomains(),
			"browser":        make([]string, 0),
			"os":             make([]string, 0),
			"placement_type": make([]string, 0),
		},
		"pagination": map[string]interface{}{
			"page":      0,
			"page_size": 10000,
		}}

	log.Debug().Msg(fmt.Sprintf("request body: %s", requestBody))

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error creating factors request body: %s", requestBody))
		return nil, errors.Wrapf(err, "Error creating factors request body")
	}

	data, statusCode, err := worker.HttpClient.Do(ctx, http.MethodPost, "http://localhost:8000/factor/get", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching factors from API")
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error Fetching factors from API. Request failed with status code: %d", statusCode))
	}

	var factorsResponse []*Factor
	if err := json.Unmarshal(data, &factorsResponse); err != nil {
		return nil, errors.Wrapf(err, "Error parsing factors from API")
	}

	// Convert the response slice to a map
	factorsMap := make(map[string]*Factor)
	for _, item := range factorsResponse {
		factorsMap[item.Key()] = item
	}

	return factorsMap, nil
}

func ToBulkRequest(newRules map[string]*FactorChanges) []bulk.FactorUpdateRequest {
	var body []bulk.FactorUpdateRequest

	for _, record := range newRules {
		tempBody := bulk.FactorUpdateRequest{
			Publisher: record.Publisher,
			Domain:    record.Domain,
			Device:    record.Device,
			Country:   record.Country,
			Factor:    record.NewFactor,
		}
		body = append(body, tempBody)
	}

	return body
}

func UpdateResponseStatus(newRules map[string]*FactorChanges, respStatus int) map[string]*FactorChanges {
	for _, record := range newRules {
		record.UpdateResponseStatus(respStatus)
	}
	return newRules
}

func UpsertLogs(ctx context.Context, newRules map[string]*FactorChanges) error {
	stringErrors := make([]string, 0)

	for _, record := range newRules {
		logJSON, err := json.Marshal(record) //Create log json to log it
		if err != nil {
			message := fmt.Sprintf("Error marshalling log for key:%v entry: %v", record.Key(), err)
			stringErrors = append(stringErrors, message)
			log.Error().Msg(message)
		}
		log.Info().Msg(fmt.Sprintf("%s", logJSON))

		mod, err := record.ToModel()
		if err != nil {
			message := fmt.Sprintf("failed to convert to model for key:%v. error: %v", record.Key(), err)
			stringErrors = append(stringErrors, message)
			log.Error().Msg(message)
		}

		err = mod.Upsert(ctx, bcdb.DB(), true, Columns, boil.Infer(), boil.Infer())
		if err != nil {
			message := fmt.Sprintf("failed to push log to postgres for key %s. Err: %s", record.Key(), err)
			stringErrors = append(stringErrors, message)
			log.Error().Err(err).Msg(message)
		}
	}
	if len(stringErrors) > 0 {
		return errors.New(strings.Join(stringErrors, "\n"))
	}
	return nil
}

func (worker *Worker) UpdateFactors(ctx context.Context, newRules map[string]*FactorChanges) (error, map[string]*FactorChanges) {
	bulkBody := ToBulkRequest(newRules)
	err := worker.bulkService.BulkInsertFactors(ctx, bulkBody)
	if err != nil {
		log.Error().Err(err).Msg("Error updating DPO factor from API. Err:")
		newRules = UpdateResponseStatus(newRules, 500)
		return err, newRules
	}

	newRules = UpdateResponseStatus(newRules, 200)

	return nil, newRules
}

func (worker *Worker) FetchAutomationSetup(ctx context.Context) (map[string]*DomainSetup, error) {
	log.Debug().Msg("fetch automation domains setup")

	requestBody := map[string]interface{}{
		"filter": map[string][]string{
			"automation": {"true"},
		},
	}

	log.Debug().Msg(fmt.Sprintf("request body: %s", requestBody))

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error creating factors request body: %s", requestBody))
		return nil, errors.Wrapf(err, "Error creating automation setup request body")
	}

	data, statusCode, err := worker.HttpClient.Do(ctx, http.MethodPost, "http://localhost:8000/publisher/domain/get", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching automation setup from API")
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error Fetching automation setup from API. Request failed with status code: %d", statusCode))
	}

	var AutomationResponse []*AutomationApi
	if err := json.Unmarshal(data, &AutomationResponse); err != nil {
		return nil, errors.Wrapf(err, "Error parsing automation setup from API")
	}

	// Append the domains to the list of active domains & get Targets
	domainsMap := make(map[string]*DomainSetup)
	for _, item := range AutomationResponse {
		gppTarget := transformGppTarget(item.GppTarget)

		domainsMap[item.Key()] = &DomainSetup{
			Domain:    item.Domain,
			GppTarget: gppTarget,
		}

	}

	return domainsMap, nil
}

func (worker *Worker) FetchInactiveKeys(ctx context.Context) ([]string, error) {
	log.Log().Msg("fetch inactive keys from postgres")
	var records []*Factor
	var inactiveKeys []string

	startString := time.Now().UTC().Truncate(time.Hour).Add(-time.Duration(worker.InactiveDaysThreshold) * 24 * time.Hour).Format("2006-01-02T15:04:05Z")

	query := fmt.Sprintf(inactiveKeysQuery, worker.InactiveFactorThreshold, startString)
	log.Log().Str("InactiveKeysQuery", query)
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &records)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fetch inactive keys from postgres")
	}

	for _, record := range records {
		inactiveKeys = append(inactiveKeys, record.Key())
	}
	return inactiveKeys, nil
}

// We need GPP Target in percentages
func transformGppTarget(gppTarget float64) float64 {
	if gppTarget != 0 {
		gppTarget /= 100
	}
	return gppTarget
}

func (worker *Worker) FetchGppTargetDefault() (float64, error) {
	GppTargetString, err := config.FetchConfigValues([]string{"factor-automation:gpp-target"})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to fetch system default gpp target value from API")
	}

	GppTargetFloat, err := strconv.ParseFloat(GppTargetString["factor-automation:gpp-target"], 64)
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("failed to convert GppTarget value to float. Gpp Target: %s", GppTargetString))
	}

	GppTargetFloat = transformGppTarget(GppTargetFloat)

	return GppTargetFloat, nil
}

func (worker *Worker) FetchFees(ctx context.Context) (map[string]float64, map[string]float64, error) {
	log.Debug().Msg("fetch global fees")

	requestBody := map[string]interface{}{}

	log.Debug().Msg(fmt.Sprintf("request body: %s", requestBody))

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error creating fees request body: %s", requestBody))
		return nil, nil, errors.Wrapf(err, "Error creating fees request body")
	}

	data, statusCode, err := worker.HttpClient.Do(ctx, http.MethodPost, "http://localhost:8000/global/factor/get", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Error Fetching fees from API")
	}

	if statusCode != http.StatusOK {
		return nil, nil, errors.New(fmt.Sprintf("Error Fetching fees from API. Request failed with status code: %d", statusCode))
	}

	var FeesResponse []*core.GlobalFactor
	if err := json.Unmarshal(data, &FeesResponse); err != nil {
		return nil, nil, errors.Wrapf(err, "Error parsing fees from API")
	}

	// Collect fee rates
	fees := make(map[string]float64)
	consultantFees := make(map[string]float64)
	for _, item := range FeesResponse {
		if item.Key == "consultant_fee" && item.PublisherID != "" {
			consultantFees[item.PublisherID] = item.Value
		} else if item.Key == "tam_fee" {
			//fees[item.Key] = item.Value For now Zeroing Tam Fee
			fees[item.Key] = 0
		} else if item.Key == "tech_fee" {
			fees[item.Key] = item.Value
		}
	}

	return fees, consultantFees, nil
}

func (worker *Worker) GenerateBidRequestMap(bidRequestMap map[string]*FactorReport, bidRequestRecords []*FactorReport) map[string]*FactorReport {
	for _, record := range bidRequestRecords {
		if !worker.CheckDomain(record) {
			continue
		}
		key := record.Key()
		item, exists := bidRequestMap[key]
		if exists {
			mergedItem := &FactorReport{
				Time:        record.Time,
				PublisherID: record.PublisherID,
				Domain:      record.Domain,
				Country:     record.Country,
				DeviceType:  record.DeviceType,
				BidRequests: item.BidRequests + record.BidRequests,
			}
			bidRequestMap[key] = mergedItem
		} else {
			bidRequestMap[key] = record
		}
	}
	return bidRequestMap
}

func (worker *Worker) GenerateImpressionsMap(impressionsMap map[string]*FactorReport, impressionsRecords []*FactorReport) map[string]*FactorReport {
	for _, record := range impressionsRecords {
		if !worker.CheckDomain(record) {
			continue
		}

		key := record.Key()
		item, exists := impressionsMap[key]
		if exists {
			mergedItem := &FactorReport{
				Time:                 record.Time,
				PublisherID:          record.PublisherID,
				Domain:               record.Domain,
				Country:              record.Country,
				DeviceType:           record.DeviceType,
				Revenue:              record.Revenue + item.Revenue,
				Cost:                 record.Cost + item.Cost,
				DemandPartnerFee:     record.DemandPartnerFee + item.DemandPartnerFee,
				SoldImpressions:      record.SoldImpressions + item.SoldImpressions,
				PublisherImpressions: record.PublisherImpressions + item.PublisherImpressions,
				DataFee:              record.DataFee + item.DataFee,
			}
			impressionsMap[key] = mergedItem
		} else {
			impressionsMap[key] = record
		}

	}
	return impressionsMap
}

func (worker *Worker) MergeReports(bidRequestMap map[string]*FactorReport, impressionsMap map[string]*FactorReport) (map[string]*FactorReport, error) {
	reportMap := make(map[string]*FactorReport)
	var err error
	for _, record := range impressionsMap {
		key := record.Key()
		requestsItem, exists := bidRequestMap[key]
		if exists {
			mergedRecord := &FactorReport{
				Time:                 record.Time,
				PublisherID:          record.PublisherID,
				Domain:               record.Domain,
				Country:              record.Country,
				DeviceType:           record.DeviceType,
				Revenue:              record.Revenue,
				Cost:                 record.Cost,
				DemandPartnerFee:     record.DemandPartnerFee,
				SoldImpressions:      record.SoldImpressions,
				PublisherImpressions: record.PublisherImpressions,
				BidRequests:          requestsItem.BidRequests,
				DataFee:              record.DataFee,
			}
			mergedRecord.CalculateGP(worker.Fees, worker.ConsultantFees)
			reportMap[key] = mergedRecord
		} else {
			record.CalculateGP(worker.Fees, worker.ConsultantFees)
			reportMap[key] = record
		}
	}
	return reportMap, err
}
