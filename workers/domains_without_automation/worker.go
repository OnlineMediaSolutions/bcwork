package domains_without_automation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
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
const loginUser = "compass-service"
const compassReportingUrl = "https://compass-reporting.deliverimp.com/api/report-dashboard/report-new-bidder"
const managerEmail = "Maayan Bar"
const defaultPath = "/etc/oms/reportNBBody.json"

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
	Cron        string `json:"cron"`
	httpClient  httpclient.Doer
	skipInitRun bool
	pass        string
	jsonPath    string
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {

	worker.DatabaseEnv = conf.GetStringValueWithDefault(config.DBEnvKey, "local")
	worker.jsonPath = conf.GetStringValueWithDefault("jsonPath", defaultPath)
	worker.Cron, _ = conf.GetStringValue("cron")
	worker.httpClient = httpclient.New(true)
	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")
	worker.pass = viper.GetString(config.TokenApiKey)

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

	token, err := worker.fetchToken(ctx)
	if err != nil {
		return err
	}

	requestJson, err := worker.createJsonForBody(domains)
	if err != nil {
		return err
	}

	results, err := fetchPubImps(token, requestJson)
	if err != nil || len(results.Data.Result) == 0 {
		if err != nil {
			return fmt.Errorf("error while fetching pub imps: %w", err)
		} else {
			return fmt.Errorf("no results found")
		}
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

func (worker *Worker) createJsonForBody(domains []string) ([]byte, error) {

	jsonFile, err := os.Open(worker.jsonPath)
	if err != nil {
		return nil, eris.Wrapf(err, "error opening file %s", worker.jsonPath)
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
		err := errors.New(fmt.Sprintf("Error Fetching pubImps with status code: %d", response.StatusCode))
		return nil, eris.Wrapf(err, "request to report new bidder failed with status")
	}

	var reportResults ReportResults
	if err := json.Unmarshal(body, &reportResults); err != nil {
		return nil, eris.Wrapf(err, "failed to marshall response body")
	}

	defer response.Body.Close()
	return &reportResults, nil
}

func (worker *Worker) fetchToken(ctx context.Context) (string, error) {

	body, err := json.Marshal(map[string]interface{}{
		"login":    loginUser,
		"password": worker.pass,
	})
	if err != nil {
		return "", eris.Wrapf(err, "failed to marshall body request")
	}

	response, statusCode, err := worker.httpClient.Do(ctx, http.MethodPost, fetchTokenUrl, bytes.NewBuffer(body))
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
