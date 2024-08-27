package factors_autmation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/quest"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var QuestSelect = `SELECT to_date('%s','yyyy-MM-ddTHH:mm:ssZ') time,
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

var InactiveKeysQuery = `SELECT to_date('%s','yyyy-MM-ddTHH:mm:ssZ') time,
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

func (worker *Worker) FetchData(ctx context.Context) (map[string]*FactorReport, map[string]*Factor, error) {
	var recordsMap map[string]*FactorReport
	var factors map[string]*Factor
	var err error

	worker.Domains, err = FetchAutomationSetup()
	if err != nil {
		return nil, nil, err
	}

	recordsMap, err = worker.FetchFromQuest(ctx, worker.Start, worker.End)
	if err != nil {
		return nil, nil, err
	}

	factors, err = worker.FetchFactors()
	if err != nil {
		return nil, nil, err
	}

	return recordsMap, factors, nil
}

// Fetch performance data from quest
func (worker *Worker) FetchFromQuest(ctx context.Context, start time.Time, stop time.Time) (map[string]*FactorReport, error) {
	log.Debug().Msg("fetch records from QuestDB")
	var records []*FactorReport

	startString := start.Format("2006-01-02T15:04:05Z")
	stopString := stop.Format("2006-01-02T15:04:05Z")
	domains := worker.AutomationDomains()

	query := fmt.Sprintf(QuestSelect, startString, startString, stopString, strings.Join(domains, "', '"))
	log.Info().Str("query", query).Msg("processImpressionsCounters")

	var RecordsMap = make(map[string]*FactorReport)
	for _, instance := range worker.Quest {
		err := quest.InitDB(instance)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("Failed to initialize Quest instance: %s", instance))
		}

		err = queries.Raw(query).Bind(ctx, quest.DB(), &records)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("failed to query impressions from Quest instance: %s", instance))
		}

		// Check if the key exists on factors
		for _, record := range records {
			if !worker.CheckDomain(record) {
				continue
			}

			key := record.Key()
			item, exists := RecordsMap[key]
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
				mergedItem.CalculateGP()
				RecordsMap[key] = mergedItem
			} else {
				record.CalculateGP()
				RecordsMap[key] = record
			}

		}

		records = nil
		err = quest.CloseDB()
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("Failed to close Quest instance: %s", instance))
		}

	}

	return RecordsMap, nil
}

func (worker *Worker) FetchFactors() (map[string]*Factor, error) {
	log.Debug().Msg("fetch records from Factors API")
	// Create the request body using a map
	requestBody := map[string]interface{}{
		"filter": map[string][]string{"domain": worker.AutomationDomains()},
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

	resp, err := http.Post("http://localhost:8000/factor/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching factors from API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error Fetching factors from API. Request failed with status code: %d", resp.StatusCode))
	}

	var factorsResponse []*Factor
	if err := json.NewDecoder(resp.Body).Decode(&factorsResponse); err != nil {
		return nil, errors.Wrapf(err, "Error parsing factors from API")
	}

	// Convert the response slice to a map
	factorsMap := make(map[string]*Factor)
	for _, item := range factorsResponse {
		factorsMap[item.Key()] = item
	}

	return factorsMap, nil
}

func (record *FactorChanges) UpdateFactor() error {
	requestBody := map[string]interface{}{
		"publisher": record.Publisher,
		"domain":    record.Domain,
		"country":   record.Country,
		"device":    record.Device,
		"factor":    record.NewFactor,
	}

	log.Debug().Msg(fmt.Sprintf("request body: %s", requestBody))

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error creating factors request body: %s", requestBody))
		return errors.Wrapf(err, "Error creating factors request body")
	}

	resp, err := http.Post("http://localhost:8000/factor", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, "Error updating factors from API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	record.RespStatus = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("Error updating factor. Request failed with status code: %d. %s", resp.StatusCode, string(bodyBytes)))
	}
	return nil
}

func FetchAutomationSetup() (map[string]*DomainSetup, error) {
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

	resp, err := http.Post("http://localhost:8000/publisher/domain/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching automation setup from API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error Fetching automation setup from API. Request failed with status code: %d", resp.StatusCode))
	}

	var AutomationResponse []*AutomationApi
	if err := json.NewDecoder(resp.Body).Decode(&AutomationResponse); err != nil {
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

// Fetch inactive keys from postgres
func (worker *Worker) FetchInactiveKeys(ctx context.Context) ([]string, error) {
	log.Debug().Msg("fetch inactive keys from postgres")
	var records []*FactorReport
	var inactiveKeys []string

	startString := time.Now().UTC().Add(-time.Duration(worker.InactiveDays) * time.Hour).Format("2006-01-02T15:04:05Z")
	stopString := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	domains := worker.AutomationDomains()

	query := fmt.Sprintf(QuestSelect, startString, startString, stopString, strings.Join(domains, "', '"))
	log.Info().Str("query", query).Msg("processImpressionsCounters")
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
		return 0, errors.Wrapf(err, "failed to fetch gpp target value from API")
	}

	GppTargetFloat, err := strconv.ParseFloat(GppTargetString["factor-automation:gpp-target"], 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to convert GppTarget value to float")
	}

	GppTargetFloat = transformGppTarget(GppTargetFloat)

	return GppTargetFloat, nil
}
