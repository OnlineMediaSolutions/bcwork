package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
)

func (p *PublisherService) GetPubImpsPerPublisherDomain(ctx context.Context, ops *GetPublisherDetailsOptions) (map[string]map[string]dto.ActivityStatus, error) {
	requestJson, err := createJsonForBody(ops.Filter.Domain)
	if err != nil {
		return nil, err
	}

	data, err := p.compassModule.Request("/report-dashboard/report-new-bidder", http.MethodPost, requestJson, true)
	if err != nil {
		return nil, err
	}

	var results dto.ReportResults
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, eris.Wrapf(err, "failed to unmarshall response data")
	}

	return buildResultMap(results.Data.Result), nil
}

func buildResultMap(results []dto.Result) map[string]map[string]dto.ActivityStatus {
	returnMap := make(map[string]map[string]dto.ActivityStatus)
	for _, result := range results {
		if len(returnMap[result.Domain]) == 0 {
			returnMap[result.Domain] = make(map[string]dto.ActivityStatus)
		}
		if result.PubImps >= dto.ActivePubs {
			returnMap[result.Domain][strconv.Itoa(result.PublisherId)] = dto.ActivityStatus(2)
		} else if result.PubImps >= dto.LowPubs && result.PubImps < dto.ActivePubs {
			returnMap[result.Domain][strconv.Itoa(result.PublisherId)] = dto.ActivityStatus(1)
		} else {
			returnMap[result.Domain][strconv.Itoa(result.PublisherId)] = dto.ActivityStatus(0)
		}
	}

	return returnMap
}

func createJsonForBody(pubDomains filter.StringArrayFilter) ([]byte, error) {
	jsonFile, err := os.Open(viper.GetString(config.ReportNBBodyPath))
	if err != nil {
		return nil, eris.Wrapf(err, "error opening file %s", "reportNBBody.json")
	}
	defer jsonFile.Close()

	fileBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, eris.Wrapf(err, "Error reading JSON file")
	}

	var data map[string]interface{}
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return nil, eris.Wrapf(err, "Error unmarshalling JSON")
	}

	today := time.Now().UTC().Truncate(24 * time.Hour).Format("2006-01-02 15:04:05")
	sevenDaysAgo := time.Now().UTC().AddDate(0, 0, -7).Truncate(24 * time.Hour).Format("2006-01-02 15:04:05")

	var domains []string
	for _, pubDomain := range pubDomains {
		domains = append(domains, pubDomain)
	}

	// Modify Dates
	if nestedData, ok := data["data"].(map[string]interface{}); ok {
		if filters, ok := nestedData["date"].(map[string]interface{}); ok {
			if dates, ok := filters["range"].([]interface{}); ok {
				filters["range"] = append(dates, sevenDaysAgo, today)
			}
		}
	}
	// Modify filters
	if nestedData, ok := data["data"].(map[string]interface{}); ok {
		if filters, ok := nestedData["filters"].(map[string]interface{}); ok {
			filters["Domain"] = domains
		}
	}

	modifiedJSON, err := json.Marshal(data)
	if err != nil {
		return nil, eris.Wrapf(err, "Error marshalling JSON")
	}

	return modifiedJSON, nil
}
