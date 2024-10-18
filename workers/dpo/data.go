package dpo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"io"
	"net/http"
)

var SelectQuery = `
	SELECT
		time as time,
		demand_partner_id as dp,
		domain as domain,
		os as os,
		CASE
			WHEN country IN ('us', 'ca', 'kr') THEN country
			ELSE 'other'
		END AS country,
		SUM(bid_requests) AS bid_request,
		SUM(revenue) AS revenue
	FROM
		nb_demand_hourly
	WHERE
		time >= '%s'
		AND time < '%s'
	GROUP BY
		time, demand_partner_id, domain, os,
		CASE
			WHEN country IN ('us', 'ca', 'kr') THEN country
			ELSE 'other'
		END;
        `

func (worker *Worker) FetchData(ctx context.Context) (map[string]*DpoReport, map[string]*PlacementReport, map[string]*DpReport, map[string]*DpoApi, error) {
	var recordsMap map[string]*DpoReport
	var placementMap map[string]*PlacementReport
	var dpMap map[string]*DpReport
	var dpoApiMap map[string]*DpoApi

	recordsMap, err := worker.FetchFromPostgres(ctx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	placementMap = GroupByPlacement(recordsMap)

	dpMap = GroupByDP(recordsMap)

	dpoApiMap, err = worker.FetchDpoApi()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return recordsMap, placementMap, dpMap, dpoApiMap, nil
}

// Fetch performance data from Postgres
func (worker *Worker) FetchFromPostgres(ctx context.Context) (map[string]*DpoReport, error) {
	log.Debug().Msg("fetch records from QuestDB")
	var reportRecords []*DpoReport
	reportMap := make(map[string]*DpoReport)

	startString := worker.Start.Format("2006-01-02 15:04:05")
	stopString := worker.End.Format("2006-01-02 15:04:05")

	startString = "2024-10-02 04:00:00"
	stopString = "2024-10-02 05:00:00"

	query := fmt.Sprintf(SelectQuery, startString, stopString)
	log.Info().Str("query", query).Msg("Fetching report")

	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &reportRecords)
	if err != nil {
		return nil, err
	}

	for _, record := range reportRecords {
		key := record.Key()
		reportMap[key] = &DpoReport{
			Time:       worker.End,
			EvalTime:   worker.Start,
			DP:         record.DP,
			Domain:     record.Domain,
			Publisher:  record.Publisher,
			Country:    record.Country,
			Os:         record.Os,
			Revenue:    record.Revenue,
			BidRequest: record.BidRequest,
			Erpm:       (record.Revenue / float64(record.BidRequest)) * 1000,
		}
	}

	return reportMap, nil
}

func GroupByPlacement(reports map[string]*DpoReport) map[string]*PlacementReport {
	placementMap := make(map[string]*PlacementReport)

	for _, report := range reports {
		key := report.PlacementKey()
		if placement, exists := placementMap[key]; exists {
			placement.Revenue += report.Revenue
		} else {
			placementMap[key] = &PlacementReport{
				Domain:    report.Domain,
				Os:        report.Os,
				Country:   report.Country,
				Publisher: report.Publisher,
				Revenue:   report.Revenue,
			}
		}
	}

	return placementMap
}

// GroupByDP groups the DpoReport data by DP (Demand Partner)
func GroupByDP(reports map[string]*DpoReport) map[string]*DpReport {
	dpMap := make(map[string]*DpReport)

	for _, report := range reports {
		if dp, exists := dpMap[report.DP]; exists {
			dp.Revenue += report.Revenue
			dp.BidRequest += report.BidRequest
		} else {
			dpMap[report.DP] = &DpReport{
				DP:         report.DP,
				Revenue:    report.Revenue,
				BidRequest: report.BidRequest,
			}
		}
	}

	return dpMap
}

func (worker *Worker) FetchDpoApi() (map[string]*DpoApi, error) {
	log.Debug().Msg("fetch records from Factors API")
	// Create the request body using a map
	requestBody := map[string]interface{}{
		"pagination": map[string]interface{}{
			"page":      0,
			"page_size": 100000,
		}}

	log.Debug().Msg(fmt.Sprintf("request body: %v", requestBody))

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error creating DPO request body: %v", requestBody))
		return nil, errors.Wrapf(err, "Error creating DPO request body")
	}

	resp, err := http.Post("http://localhost:8000/dpo/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching DPO from API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error Fetching DPO from API. Request failed with status code: %d", resp.StatusCode))
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	// Print the response body
	fmt.Println("Response body:")
	fmt.Println(string(bodyBytes))

	// Create a new reader with the body bytes for json.NewDecoder
	bodyReader := bytes.NewReader(bodyBytes)

	var factorsResponse []*DpoApi
	if err := json.NewDecoder(bodyReader).Decode(&factorsResponse); err != nil {
		return nil, errors.Wrapf(err, "Error parsing factors from API")
	}

	// Convert the response slice to a map
	factorsMap := make(map[string]*DpoApi)
	for _, item := range factorsResponse {
		factorsMap[item.Key()] = item
	}

	return factorsMap, nil
}

func (record *DpoChanges) UpdateFactor() error {
	requestBody := map[string]interface{}{
		"publisher":         record.Publisher,
		"demand_partner_id": record.DP,
		"domain":            record.Domain,
		"country":           record.Country,
		"os":                record.Os,
		"factor":            record.NewFactor,
	}

	log.Debug().Msg(fmt.Sprintf("request body: %s", requestBody))

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error creating factors request body: %s", requestBody))
		return errors.Wrapf(err, "Error creating factors request body")
	}

	resp, err := http.Post("http://localhost:8000/dpo/set", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Error updating DPO factor from API for key %s", record.Key()))
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
