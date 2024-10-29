package dpo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"net/http"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

var SelectQuery = `
	SELECT
		time as time,
		publisher_id as publisher, 
		demand_partner_id as dp,
		domain as domain,
		os as os,
		CASE
			WHEN country IN ('us', 'ca', 'kr', 'gb') THEN country
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
		1,2,3,4,5,6;
        `

func (worker *Worker) FetchData(ctx context.Context) DpoData {
	var data DpoData
	var err error

	data.DpoReport, err = worker.FetchFromPostgres(ctx)
	if err != nil {
		return DpoData{Error: err}
	}

	data.PlacementReport = GroupByPlacement(data.DpoReport)

	data.DpReport = GroupByDP(data.DpoReport)

	data.DpoApi, err = worker.FetchDpoApi(ctx)
	if err != nil {
		return DpoData{Error: err}
	}

	return data
}

// Fetch performance data from Postgres
func (worker *Worker) FetchFromPostgres(ctx context.Context) (map[string]*DpoReport, error) {
	log.Debug().Msg("fetch records from QuestDB")
	var reportRecords []*DpoReport
	reportMap := make(map[string]*DpoReport)

	startString := worker.Start.Format("2006-01-02 15:04:05")
	stopString := worker.End.Format("2006-01-02 15:04:05")

	query := fmt.Sprintf(SelectQuery, startString, stopString)
	log.Info().Str("query", query).Msg("Fetching report")

	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &reportRecords)
	if err != nil {
		return nil, errors.Wrapf(err, "error fetching records from postgres")
	}

	for _, record := range reportRecords {
		DpApiName := ""
		_, exists := worker.Demands[record.DP]
		if exists {
			DpApiName = worker.Demands[record.DP].ApiName
		}
		key := record.Key()
		reportMap[key] = &DpoReport{
			Time:       worker.End,
			EvalTime:   worker.Start,
			DpApiName:  DpApiName,
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

func (worker *Worker) FetchDpoApi(ctx context.Context) (map[string]*DpoApi, error) {
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

	data, statusCode, err := worker.httpClient.Do(ctx, http.MethodPost, constant.ProductionApiUrl+constant.DpoGetEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching DPO from API")
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error Fetching DPO from API. Request failed with status code: %d", statusCode))
	}

	var factorsResponse []*DpoApi
	if err := json.Unmarshal(data, &factorsResponse); err != nil {
		return nil, errors.Wrapf(err, "Error parsing factors from API")
	}

	// Convert the response slice to a map
	factorsMap := make(map[string]*DpoApi)
	for _, item := range factorsResponse {
		factorsMap[item.Key()] = item
	}

	return factorsMap, nil
}

func (worker *Worker) UpdateFactors(ctx context.Context, requestBody []map[string]interface{}) (error, int) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error creating factors request body: %s", requestBody))
		return errors.Wrapf(err, "Error creating factors request body"), 0
	}

	data, statusCode, err := worker.httpClient.Do(ctx, http.MethodPost, constant.ProductionApiUrl+constant.DpoSetEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, "Error updating DPO factor from API for key"), 0
	}

	if statusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Error updating factor. Request failed with status code: %d. %s", statusCode, string(data))), statusCode
	}

	return nil, statusCode
}

func UpsertLogs(ctx context.Context, newRules map[string]*DpoChanges, respStatus int) error {
	stringErrors := make([]string, 0)

	for _, record := range newRules {
		logJSON, err := json.Marshal(record) //Create log json to log it
		if err != nil {
			message := fmt.Sprintf("Error marshalling log for key:%v entry: %v", record.Key(), err)
			stringErrors = append(stringErrors, message)
			log.Error().Msg(message)

		}
		log.Info().Msg(fmt.Sprintf("%s", logJSON))

		mod, err := record.ToModel(respStatus)
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
	return errors.New(strings.Join(stringErrors, "\n"))
}
