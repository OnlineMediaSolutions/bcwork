package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const fetchTokenIp = "52.3.132.64:22"
const apiHost = "10.166.10.36:8080"
const fetchTokenPostfix = "/api/auth/token"
const loginUser = "compass-service"
const compassReportingUrl = "https://compass-reporting.deliverimp.com/api/report-dashboard/report-new-bidder"
const userName = "ec2-user"
const timeout = 15

func (p *PublisherService) GetPubImpsPerPublisherDomain(ctx context.Context, ops *GetPublisherDetailsOptions) (map[string]map[string]dto.ActivityStatus, error) {

	token, err := fetchToken(ctx)
	if err != nil {
		return nil, err
	}

	requestJson, err := createJsonForBody(ops.Filter.Domain)
	if err != nil {
		return nil, err
	}

	results, err := fetchPubImps(token, requestJson)
	if err != nil {
		return nil, err
	}
	buildMap := buildResultMap(results.Data.Result)
	if err != nil {
		return nil, err
	}

	return buildMap, nil
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

func fetchToken(ctx context.Context) (string, error) {

	body, err := json.Marshal(map[string]interface{}{
		"login":    loginUser,
		"password": viper.GetString(config.TokenApiKey),
	})
	if err != nil {
		return "", eris.Wrapf(err, "failed to marshall body request")
	}

	client, err := createSSHClient(userName, fetchTokenIp)
	if err != nil {
		return "", eris.Wrap(err, "error creating client with SSH tunnel")
	}

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

	var token dto.Token
	body, err = ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(body, &token); err != nil {
		return "", eris.Wrapf(err, "failed to marshall response body")
	}
	defer response.Body.Close()

	return token.Data.Token, nil
}

func createSSHClient(user string, host string) (*ssh.Client, error) {

	key := viper.GetString(config.SshKey)
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout * time.Second,
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to SSH server: %v", err)
	}

	return client, nil
}

func createJsonForBody(pubDomains filter.StringArrayFilter) ([]byte, error) {

	jsonFile, err := os.Open("files/reportNBBody.json")
	if err != nil {
		return nil, eris.Wrapf(err, "error opening file %s", "/files/reportNBBody.json")
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

func fetchPubImps(token string, requestJson []byte) (*dto.ReportResults, error) {

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

	var reportResults dto.ReportResults
	if err := json.Unmarshal(body, &reportResults); err != nil {
		return nil, eris.Wrapf(err, "failed to marshall response body")
	}

	defer response.Body.Close()
	return &reportResults, nil
}
