package compass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Compass struct {
	CompassURL   string
	ReportingURL string
	Token        string
	Client       *http.Client
}

func NewCompass() *Compass {
	return &Compass{
		CompassURL:   "http://10.166.10.36:8080",
		ReportingURL: "https://compass-reporting.deliverimp.com",
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Compass) Login() error {

	data := map[string]string{
		"login":    "compass-service",
		"password": "HdkwLFpvkfAmfQMNEEv9WqudVZRt8",
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal login data- %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/auth/token", c.CompassURL), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create login request- %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed- %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode login response: %w", err)
	}

	if result.Data.Token == "" {
		return fmt.Errorf("login succeeded but returned an empty token")
	}

	c.Token = result.Data.Token
	return nil
}

func (c *Compass) getHeaders(auth bool) map[string]string {
	headers := make(map[string]string)
	if auth && c.Token != "" {
		headers["x-access-token"] = c.Token
	}
	headers["Content-Type"] = "application/json"
	return headers
}

func (c *Compass) GetURL(path string, isReportingRequest bool) string {
	baseURL := c.CompassURL
	if isReportingRequest {
		baseURL = c.ReportingURL
	}
	return fmt.Sprintf("%s/api%s", baseURL, path)
}

func (c *Compass) Request(url, method string, data interface{}, auth, isReportingRequest bool) (map[string]interface{}, error) {
	if auth && c.Token == "" {
		if err := c.Login(); err != nil {
			return nil, fmt.Errorf("login failed: %w", err)
		}
	}

	var body []byte
	var err error
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
	}

	req, err := http.NewRequest(method, c.GetURL(url, isReportingRequest), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.getHeaders(auth) {
		req.Header.Set(key, value)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %d", err)
	}

	return result, nil
}
