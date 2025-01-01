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
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

//
// This worker fetch once a week all domains without factor automation and sends an email
//

const fetchTokenIp = "52.3.132.64:22"
const fetchTokenPostfix = "/api/auth/token"
const loginUser = "compass-service"
const compassReportingUrl = "https://compass-reporting.deliverimp.com/api/report-dashboard/report-new-bidder"
const managerEmail = "Maayan Bar"
const defaultPath = "/etc/oms/reportNBBody.json"
const sshKeyPath = "/etc/oms/bcwork_ny01.pem"

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
	Cron        string `json:"cron"`
	httpClient  httpclient.Doer
	skipInitRun bool
	pass        string
	jsonPath    string
	sshPath     string
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {

	worker.DatabaseEnv = conf.GetStringValueWithDefault(config.DBEnvKey, "local")
	worker.jsonPath = conf.GetStringValueWithDefault("jsonPath", defaultPath)
	worker.sshPath = conf.GetStringValueWithDefault("sshPath", sshKeyPath)
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
	log.Info().Msg("fetched domains successfully")

	token, err := worker.fetchToken(ctx)
	if err != nil {
		return err
	}
	log.Info().Msg("fetched token " + token)

	requestJson, err := worker.createJsonForBody(domains)
	if err != nil {
		return err
	}
	log.Info().Msg("created body")

	results, err := fetchPubImps(token, requestJson)
	log.Info().Msg("fetched pubimp successfully")

	if err != nil || len(results.Data.Result) == 0 {
		if err != nil {
			return fmt.Errorf("error while fetching pub imps: %w", err)
		} else {
			return fmt.Errorf("no results found")
		}
	}

	domainsPerAccountManager, managerList := generateDomainLists(results.Data.Result)
	log.Info().Msg("generated domains per account managers")

	emails, err := fetchEmailAddresses(ctx)
	log.Info().Msg("fetched emails successfully")
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

	log.Info().Msg("login user: " + loginUser)
	log.Info().Msg("token: " + worker.pass)
	log.Info().Msg("token URL: " + "http://" + fetchTokenIp + fetchTokenPostfix)

	body, err := json.Marshal(map[string]interface{}{
		"login":    loginUser,
		"password": worker.pass,
	})
	if err != nil {
		return "", eris.Wrapf(err, "failed to marshall body request")
	}

	client, err := worker.createSSHClient("ec2-user", fetchTokenIp)
	if err != nil {
		return "", eris.Wrap(err, "error creating client with SSH tunnel")
	}
	// Create a tunnel (optional: only needed if accessing API through SSH)
	apiHost := "10.166.10.36:8080"
	conn, err := client.Dial("tcp", apiHost)
	if err != nil {
		fmt.Printf("Error creating SSH tunnel: %v\n", err)
		return "", eris.Wrap(err, "error creating SSH tunnel")
	}

	httpclientWithSSH := &http.Client{
		Transport: &http.Transport{Dial: func(_, _ string) (net.Conn, error) { return conn, nil }},
	}

	response, err := httpclientWithSSH.Post("http://"+fetchTokenIp+fetchTokenPostfix, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error making API call: %v\n", err)
		return "", eris.Wrap(err, "error making API call")
	}

	var token Token
	body, err = ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(body, &token); err != nil {
		return "", eris.Wrapf(err, "failed to marshall response body")
	}
	defer response.Body.Close()

	return token.Data.Token, nil
}

func (worker *Worker) createSSHClient(user string, host string) (*ssh.Client, error) {
	// Load the private key
	key, err := ioutil.ReadFile(worker.sshPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	// Parse the private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	// Configure the SSH client
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to SSH server: %v", err)
	}

	return client, nil
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
