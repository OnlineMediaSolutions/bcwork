package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/spf13/viper"
)

const (
	DBEnvKey          = "dbenv"
	LogSeverityKey    = "logsev"
	CronExpressionKey = "cron"
	BucketKey         = "bucket"
	PrefixKey         = "prefix"
	DaysBeforeKey     = "days_before"
	BaseURLKey        = "base_url"
	TestCasesPathKey  = "test_cases"

	APIChunkSizeKey     = "api.chunkSize"
	CronWorkerAPIKeyKey = "cron_worker_api_key"
	AWSWorkerAPIKeyKey  = "aws_worker_api_key"
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

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating request")
	}
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.Header.Add(constant.HeaderOMSWorkerAPIKey, viper.GetString(CronWorkerAPIKeyKey))

	resp, err := http.DefaultClient.Do(req)
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
