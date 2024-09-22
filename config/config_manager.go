package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/friendsofgo/errors"
)

const (
	DBEnvKey          = "dbenv"
	LogSeverityKey    = "logsev"
	CronExpressionKey = "cron"
)

type ConfigApi struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func FetchConfigValues(keys []string) (map[string]string, error) {
	endpoint := "http://localhost:8000/config/get"

	requestBody := map[string]interface{}{
		"filter": map[string][]string{"key": keys},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating config request body")
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching factors from API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, fmt.Sprintf("Error Fetching factors from API. Request failed with status code: %d", resp.StatusCode))
	}

	var response []ConfigApi
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrapf(err, "Error parsing factors from API")
	}

	// Convert the response slice to a map
	ConfigMap := make(map[string]string)
	for _, item := range response {
		ConfigMap[item.Key] = item.Value
	}
	return ConfigMap, nil
}
