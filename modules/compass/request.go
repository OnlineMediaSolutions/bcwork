package compass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/config"
	"io"
	"net/http"
	"time"
)

type Compass struct {
	compassURL      string
	reportingURL    string
	token           string
	tokenExpiration time.Time
	tokenDuration   time.Duration
	client          *http.Client
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

func NewCompass() *Compass {
	return &Compass{
		compassURL:   "http://10.166.10.36:8080",
		reportingURL: "https://compass-reporting.deliverimp.com",
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/auth/token", c.compassURL), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create token request- %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed- %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token failed with status code: %d", resp.StatusCode)
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
