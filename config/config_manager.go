package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/spf13/viper"
)

const (
	DBEnvKey            = "dbenv"
	LogSeverityKey      = "logsev"
	CronExpressionKey   = "cron"
	BucketKey           = "bucket"
	PrefixKey           = "prefix"
	DaysBeforeKey       = "days_before"
	BaseURLKey          = "base_url"
	TestCasesPathKey    = "test_cases"
	ManagersMapPathKey  = "managers_map"
	QuestKey            = "quest"
	SkipInitRunKey      = "skip_init_run"
	ChunkSizeKey        = "chunk_size"
	SQLDebugKey         = "sql_debug"
	EmailToKey          = "email_to"
	EmailBCCKey         = "email_bcc"
	CreateAdsTxtLineKey = "create_ads_txt_lines"

	APIChunkSizeKey             = "api.chunkSize"
	CronWorkerAPIKeyKey         = "cron_worker_api_key" //nolint:gosec
	AWSWorkerAPIKeyKey          = "aws_worker_api_key"  //nolint:gosec
	LogSizeLimitKey             = "log_size_limit"
	SearchViewUpdateRateKey     = "search_view_update_rate"
	AdsTxtMetadataUpdateRateKey = "ads_txt_metadata_update_rate"
	// compass
	CompassModuleKey = "compassModule"
	CompassURLKey    = "compassURL"
	ReportingURLKey  = "reportingURL"
	SshServerKey     = "sshServer"
	SshTimeoutKey    = "sshTimeout"
	SshUserKey       = "sshUser"
	SshKey           = "sshKey"
	TokenApiKey      = "tokenApiKey"
	ReportNBBodyPath = "reportNBBodyPath"
)

type ConfigApi struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func FetchConfigValues(keys []string) (map[string]string, error) {
	endpoint := constant.ProductionApiUrl + constant.ConfigEndpoint

	requestBody := map[string]interface{}{
		"filter": map[string][]string{"key": keys},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error creating config request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.Header.Add(constant.HeaderOMSWorkerAPIKey, viper.GetString(CronWorkerAPIKeyKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error data config from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, %w", resp.StatusCode, err)
	}

	var response []ConfigApi
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing data from API: %w", err)
	}

	ConfigMap := make(map[string]string)
	for _, item := range response {
		ConfigMap[item.Key] = item.Value
	}

	return ConfigMap, nil
}
