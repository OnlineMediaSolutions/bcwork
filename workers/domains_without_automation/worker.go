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
	"github.com/m6yf/bcwork/modules"
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
	"strings"
	"text/template"
	"time"
)

//
// This worker fetch once a week all domains without factor automation and sends an email
//

const fetchTokenUrl = "http://10.166.10.36:8080/api/auth/token"
const loginUuser = "compass-service"
const compassReportingUrl = "https://compass-reporting.deliverimp.com/api/report-dashboard/report-new-bidder"
const jsonUrl = "workers/domains_without_automation/reportNBBody.json"

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
	domainsPerAccountManager := calculatePubImpsPerDomain(results.Data.Result)

	for _, accountManager := range domainsPerAccountManager {
		body, _ := createEmailBody(accountManager)
		sendEmail(body)
	}

	log.Info().Msg("domains_without_automation worker Finished")

	return nil
}

func createEmailBody(accountManagerData []Result) (string, error) {
	const tpl = `
<html>
    <head>
        <title>Domains without automation</title>
        <style>
            table { width: 100%; border-collapse: collapse; }
            th, td { border: 1px solid black; padding: 8px; text-align: left; }
            th { background-color: #f2f2f2; }
            .no-changes { color: red; font-weight: bold; }
        </style>
    </head>
    <body>
        <h3>{{Domains without automation}}</h3>
                <table>
                    <tr>
                        <th>Publisher</th>
                        <th>Domain</th>
                        <th>Account Manager</th>
                        <th>Pub Imps</th>
                        <th>Looping Ratio</th>
                        <th>Cost</th>
                        <th>CPM</th>
                        <th>Revenue</th>
                        <th>RPM</th>
                        <th>DP RPM</th>
                        <th>GP</th>
                        <th>GP%</th>
                    </tr>
                </table>
    </body>
</html>
`
	data := struct {
		Body            string
		CompetitorsData []Result
	}{
		Body:            "data per account manager",
		CompetitorsData: accountManagerData,
	}

	t, err := template.New("emailTemplate").Parse(tpl)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err := t.Execute(&tplBuffer, data); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}

func sendEmail(body string) {

	emailCredsMap, _ := config.FetchConfigValues([]string{"real_time_report"})

	var emailProperties EmailProperties
	if err := json.Unmarshal([]byte(emailCredsMap["real_time_report"]), &emailProperties); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal email credentials")
	}

	emailReq := modules.EmailRequest{
		To:      strings.Split(emailProperties.TO, ","),
		Bcc:     strings.Split(emailProperties.BCC, ","),
		Subject: "Domains without automation",
		Body:    body,
		IsHTML:  true,
	}

	modules.SendEmail(emailReq)

}

func calculatePubImpsPerDomain(results []Result) (finalList map[string][]Result) {

	var domainData = make(map[string][]Result)
	for _, result := range results {
		if result.PubImps > constant.NewBidderAutomationThreshold {
			domainData[result.AccountManager] = append(domainData[result.AccountManager], result)
		}
	}

	return domainData
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
