package factors_autmation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/quest"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"io"
	"net/http"
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
GROUP BY 1, 2, 3, 4, 5`

// Main function to fetch factors & quest data
func (w *Worker) FetchData(ctx context.Context) (map[string]*FactorReport, map[string]*Factor, error) {
	var RecordsMap map[string]*FactorReport
	var factors map[string]*Factor
	var err error

	log.Info().Msg("fetch records from QuestDB")
	RecordsMap, err = w.FetchFromQuest(ctx, w.Start, w.End)
	if err != nil {
		return nil, nil, err
	}

	log.Info().Msg("fetch records from Factors API")
	factors, err = w.FetchFactors()
	if err != nil {
		return nil, nil, err
	}

	return RecordsMap, factors, nil
}

// Fetch performance data from quest
func (w *Worker) FetchFromQuest(ctx context.Context, start time.Time, stop time.Time) (map[string]*FactorReport, error) {
	var records []*FactorReport

	startString := start.Format("2006-01-02T15:04:05Z")
	stopString := stop.Format("2006-01-02T15:04:05Z")

	query := fmt.Sprintf(QuestSelect, startString, startString, stopString)
	log.Info().Str("query", query).Msg("processImpressionsCounters")

	var RecordsMap = make(map[string]*FactorReport)
	for _, instance := range w.Quest {
		err := quest.InitDB(instance)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("Failed to initialize Quest instance: %s", instance))
		}

		err = queries.Raw(query).Bind(ctx, quest.DB(), &records)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to query impressions from questdb")
		}

		// Check if the key exists on factors
		for _, record := range records {
			if !w.CheckDomain(record.Domain) {
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

// Fetch the factors from the Factors API
func (w *Worker) FetchFactors() (map[string]*Factor, error) {
	// Create the request body using a map
	requestBody := map[string]interface{}{
		"pagination": map[string]interface{}{
			"page":      0,
			"page_size": 10000,
		}}

	if w.FilterExists {
		requestBody["filter"] = map[string]interface{}{"domain": w.Domains}
	}
	fmt.Println(requestBody)

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating factors request body")
	}

	// Perform the HTTP request
	resp, err := http.Post("http://localhost:8000/factor/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching factors from API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, fmt.Sprintf("Error Fetching factors from API. Request failed with status code: %d", resp.StatusCode))
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

// Update a factor via the API
func (record *FactorChanges) updateFactor() error {
	requestBody := map[string]interface{}{
		"publisher": record.Publisher,
		"domain":    record.Domain,
		"country":   record.Country,
		"device":    record.Device,
		"factor":    record.NewFactor,
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrapf(err, "Error creating factors request body")
	}

	// Perform the HTTP request
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

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, fmt.Sprintf("Error Fetching factors from API. Request failed with status code: %d", resp.StatusCode))
	}
	return nil
}
