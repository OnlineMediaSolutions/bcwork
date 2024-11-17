package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"net/http"
)

func FetchFees(ctx context.Context) (map[string]float64, map[string]float64, error) {
	HttpClient := httpclient.New(true)

	log.Debug().Msg("fetch global fees for real time report")

	requestBody := map[string]interface{}{}
	jsonData, err := json.Marshal(requestBody)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error creating fees request body for real time report: %s", requestBody))
		return nil, nil, errors.Wrapf(err, "Error creating fees request body for real time report")
	}

	data, statusCode, err := HttpClient.Do(ctx, http.MethodPost, "http://localhost:8000/global/factor/get", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Error Fetching fees from API")
	}

	if statusCode != http.StatusOK {
		return nil, nil, errors.New(fmt.Sprintf("Error Fetching fees from API. Request failed with status code: %d", statusCode))
	}

	var FeesResponse []*constant.GlobalFactor
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
			fees[item.Key] = 0
		} else if item.Key == "tech_fee" {
			fees[item.Key] = item.Value
		}
	}

	return fees, consultantFees, nil
}
