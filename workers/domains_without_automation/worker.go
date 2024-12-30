package domains_without_automation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//
// This worker fetch once a week all domains without factor automation and sends an email
//

const fetchTokenUrl = "http://10.166.10.36:8080/api/auth/token"
const loginUuser = "compass-service"
const compassReportingUrl = "https://compass-reporting.deliverimp.com/api/report-dashboard/report-new-bidder"
const jsonUrl = "workers/domains_without_automation/reportNBBody.json"
const managerEmail = "Maayan Bar"

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
	Cron        string `json:"cron"`
	httpClient  httpclient.Doer
	skipInitRun bool
}

type TransactionIds struct {
	TransactionId string `json:"transaction_id"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {

	worker.DatabaseEnv = conf.GetStringValueWithDefault(config.DBEnvKey, "local")
	worker.Cron, _ = conf.GetStringValue("cron")
	worker.httpClient = httpclient.New(true)
	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")

	err := bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "failed to initalize DB")
	}
	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	if worker.skipInitRun {
		log.Info().Msg("Skipping work as per the skip_init_run flag.")
		worker.skipInitRun = false
		return nil
	}

	log.Info().Msg("Starting domains_without_automation worker")

	domains, err := fetchDomainsWithoutAutomation(ctx)
	if err != nil {
		return err
	}

	token, err := fetchToken(worker.httpClient, ctx)
	if err != nil {
		return err
	}

	requestJson, err := createJsonForBody(domains)
	if err != nil {
		return err
	}

	results, err := fetchPubImps(token, requestJson)
	if err != nil || len(results.Data.Result) == 0 {
		return fmt.Errorf("error while fetching pub imps: %v", err)
	}

	domainsPerAccountManager, managerList := generateDomainLists(results.Data.Result)
	emails, err := fetchEmailAddresses(ctx)
	if err != nil {
		return fmt.Errorf("error fetching emails for account managers: %v", err)
	}

	SendEmails(emails, domainsPerAccountManager, managerList)

	log.Info().Msg("domains_without_automation worker Finished")
	return nil
}

func fetchEmailAddresses(ctx context.Context) (map[string]string, error) {
	var domainData = make(map[string]string)

	users, err := models.Users().All(ctx, bcdb.DB())
	if err != nil {
		return nil, eris.Wrap(err, "failed to fetch domains without Automation from DB")
	}

	for _, user := range users {
		domainData[user.FirstName+" "+user.LastName] = user.Email
	}
	return domainData, nil
}

func generateDomainLists(results []Result) (finalList map[string][]Result, managerList []Result) {

	var domainsPerAccountManager = make(map[string][]Result)
	var domainsForManager = []Result{}

	for _, result := range results {
		if result.PubImps > constant.NewBidderAutomationThreshold {
			domainsForManager = append(domainsForManager, result)
			domainsPerAccountManager[result.AccountManager] = append(domainsPerAccountManager[result.AccountManager], result)
		}
	}

	return domainsPerAccountManager, domainsForManager
}

func createJsonForBody(domains []string) ([]byte, error) {
	jsonFile, err := os.Open(jsonUrl)
	if err != nil {
		return nil, eris.Wrapf(err, "error opening file %s", jsonUrl)
	}
	defer jsonFile.Close()

	fileBytes, err := ioutil.ReadAll(jsonFile)
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

	// Modify Dates
	if nestedData, ok := data["data"].(map[string]interface{}); ok {
		filters := nestedData["date"].(map[string]interface{})
		if dates, ok := filters["range"].([]interface{}); ok {
			filters["range"] = append(dates, sevenDaysAgo, today)
		}
	}
	// Modify filters
	if nestedData, ok := data["data"].(map[string]interface{}); ok {
		filters := nestedData["filters"].(map[string]interface{})
		filters["Domain"] = domains
	}

	modifiedJSON, err := json.Marshal(data)
	if err != nil {
		return nil, eris.Wrapf(err, "Error marshalling JSON")
	}

	return modifiedJSON, nil
}

func fetchPubImps(token string, requestJson []byte) (*ReportResults, error) {

	request, err := http.NewRequest(http.MethodPost, compassReportingUrl, bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, eris.Wrapf(err, "error creating request for compass reporting endpoint")
	}
	request.Header.Set("x-access-token", token)
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	response, err := http.DefaultClient.Do(request)
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to report new bidder failed with status: %d", response.StatusCode)
	}

	var reportResults ReportResults
	if err := json.Unmarshal(body, &reportResults); err != nil {
		return nil, eris.Wrapf(err, "failed to marshall response body")
	}

	defer response.Body.Close()
	return &reportResults, nil
}

func fetchToken(client httpclient.Doer, ctx context.Context) (string, error) {

	pass := viper.GetString(config.TokenApiKey)

	body, err := json.Marshal(map[string]interface{}{
		"login":    loginUuser,
		"password": pass,
	})
	if err != nil {
		return "", eris.Wrapf(err, "failed to marshall body request")
	}

	response, statusCode, err := client.Do(ctx, http.MethodPost, fetchTokenUrl, bytes.NewBuffer(body))
	if statusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status: %d", statusCode)
	}

	var token Token
	if err := json.Unmarshal(response, &token); err != nil {
		return "", eris.Wrapf(err, "failed to marshall response body")
	}

	return token.Data.Token, nil
}

func fetchDomainsWithoutAutomation(ctx context.Context) ([]string, error) {

	var domains []string
	results, err := models.PublisherDomains(models.PublisherDomainWhere.Automation.EQ(false),
		qm.Select(models.PublisherDomainColumns.Domain)).All(ctx, bcdb.DB())

	if err != nil {
		return nil, eris.Wrapf(err, "failed to fetch domains without Automation from DB")
	}

	for _, result := range results {
		domains = append(domains, result.Domain)
	}
	return domains, nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}
