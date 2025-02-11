package compass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/config"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"net/http"
	"time"
)

type Compass struct {
	compassURL        string
	reportingURL      string
	token             string
	tokenExpiration   time.Time
	tokenDuration     time.Duration
	client            *http.Client
	userName          string
	fetchTokenIp      string
	fetchTokenPostfix string
	apiHost           string
}

type CompassConfig struct {
	Login    string `json:"COMPASS_LOGIN"`
	Password string `json:"COMPASS_PASSWORD"`
}

type Data struct {
	Token string `json:"token"`
}

type Result struct {
	Data Data `json:"data"`
}

//Example usage
//compassClient := compass.NewCompass()
//For request compass-reporting
//reportData, err := compassClient.Request(/report-dashboard/report-new-bidder, "POST", requestData, true)
//For request compass
//reportData, err := compassClient.Request(/report-dashboard/report-new-bidder, "POST", requestData,false)

const timeout = 60

func NewCompass() *Compass {
	return &Compass{
		compassURL:        "http://10.166.10.36:8080",
		reportingURL:      "https://compass-reporting.deliverimp.com",
		apiHost:           "10.166.10.36:8080",
		userName:          "ec2-user",
		fetchTokenIp:      "52.3.132.64:22",
		fetchTokenPostfix: "/api/auth/token",
		client: &http.Client{
			Timeout: 100 * time.Second,
		},
	}
}

func (c *Compass) Request(url, method string, body []byte, isReportingRequest bool) ([]byte, error) {
	if c.token == "" || c.isTokenExpired() {
		if err := c.login(); err != nil {
			return nil, fmt.Errorf("login failed: %w", err)
		}
	}

	req, err := http.NewRequest(method, c.getURL(url, isReportingRequest), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.getHeaders() {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseBody, nil
}

func (c *Compass) login() error {
	compassCredentialsMap, err := config.FetchConfigValues([]string{"compass"})
	if err != nil {
		return fmt.Errorf("error fetching config values: %w", err)
	}

	creds := compassCredentialsMap["compass"]

	var compassConfig CompassConfig
	err = json.Unmarshal([]byte(creds), &compassConfig)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	data := map[string]string{
		"login":    compassConfig.Login,
		"password": compassConfig.Password,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal login data- %w", err)
	}

	return c.getCompassToken(body)
}

func (c *Compass) getHeaders() map[string]string {
	headers := make(map[string]string)
	headers["x-access-token"] = c.token
	headers["Content-Type"] = "application/json"
	return headers
}

func (c *Compass) getURL(path string, isReportingRequest bool) string {
	baseURL := c.compassURL
	if isReportingRequest {
		baseURL = c.reportingURL
	}
	return fmt.Sprintf("%s/api%s", baseURL, path)
}

func (c *Compass) getCompassToken(body []byte) error {
	client, err := createSSHClient(c.userName, c.fetchTokenIp)
	if err != nil {
		return fmt.Errorf("error creating client with SSH tunnel: %w", err)
	}

	conn, err := client.Dial("tcp", c.apiHost)
	if err != nil {
		return fmt.Errorf("error creating SSH tunnel: %w", err)
	}

	httpClientWithSSH := &http.Client{
		Transport: &http.Transport{Dial: func(_, _ string) (net.Conn, error) { return conn, nil }},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", c.fetchTokenIp, c.fetchTokenPostfix), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClientWithSSH.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status code: %d", resp.StatusCode)
	}

	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	if result.Data.Token == "" {
		return fmt.Errorf("login succeeded but returned an empty token")
	}

	c.token = result.Data.Token
	c.tokenDuration = 24 * time.Hour
	c.tokenExpiration = time.Now().Add(c.tokenDuration)
	return nil
}

func (c *Compass) isTokenExpired() bool {
	return time.Now().After(c.tokenExpiration)
}

func createSSHClient(user string, host string) (*ssh.Client, error) {
	key := viper.GetString(config.SshKey)
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	conf := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout * time.Second,
	}

	client, err := ssh.Dial("tcp", host, conf)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to SSH server: %v", err)
	}

	return client, nil
}

//func GetSSHKey() string {
//	key, err := ioutil.ReadFile("/Users/sonaisrayel/.ssh/amiram.pem")
//	if err != nil {
//		log.Fatalf("Unable to read SSH key: %v", err)
//	}
//	return string(key)
//}
