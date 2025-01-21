package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestBidCachingGetHandler(t *testing.T) {
	endpoint := "/test/bid_caching/get"

	type want struct {
		statusCode int
		response   string
	}

	tests := []struct {
		name        string
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "Bid caching with active filter=true",
			requestBody: `{"filter": {"active":true}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"rule_id":"123456","publisher":"21038","domain":"oms.com","country":"","device":"","bid_caching":10,"browser":"","os":"","placement_type":"","active":true,"control_percentage":0.5,"created_at":"2024-10-01T13:51:28.407Z","updated_at":"2024-10-01T13:51:28.407Z"}]`,
			},
		},
		{
			name:        "Bid caching with active filter=false",
			requestBody: `{"filter": {"active": false} }`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"rule_id":"1234567","publisher":"21000","domain":"brightcom.com","country":"","device":"","bid_caching":10,"browser":"","os":"","placement_type":"","active":false,"control_percentage":0,"created_at":"2024-10-01T13:51:28.407Z","updated_at":"2024-10-01T13:51:28.407Z"}]`,
			},
		},
		{
			name:        "Bid caching with all values",
			requestBody: `{}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"rule_id":"123456","publisher":"21038","domain":"oms.com","country":"","device":"","bid_caching":10,"browser":"","os":"","placement_type":"","active":true,"control_percentage":0.5,"created_at":"2024-10-01T13:51:28.407Z","updated_at":"2024-10-01T13:51:28.407Z"},{"rule_id":"1234567","publisher":"21000","domain":"brightcom.com","country":"","device":"","bid_caching":10,"browser":"","os":"","placement_type":"","active":false,"control_percentage":0,"created_at":"2024-10-01T13:51:28.407Z","updated_at":"2024-10-01T13:51:28.407Z"}]`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(fiber.MethodPost, baseURL+endpoint, strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			resp, err := http.DefaultClient.Do(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}

func TestBidCachingSetHandler(t *testing.T) {
	endpoint := "/test/bid_caching/set"

	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Invalid device",
			body:         `{"publisher":"1234", "device": "mm", "country": "us", "bid_caching": 1, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"error","message":"could not validate bid cache","errors":["Device should be in the allowed list"]}`,
		},
		{
			name:         "Invalid country",
			body:         `{"publisher": "test", "device": "tablet", "country": "USA", "bid_caching": 1, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"error","message":"could not validate bid cache","errors":["Country code must be 2 characters long and should be in the allowed list"]}`,
		},
		{
			name:         "Valid request",
			body:         `{"publisher": "test","bid_caching": 1, "domain": "example.com", "control_percentage":0.5}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"success","message":"Bid Caching successfully created"}`,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", endpoint, strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := appTest.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedCode, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, test.expectedBody, string(body))
	}
}

func TestBidCachingDeleteHandler(t *testing.T) {
	endpoint := "/test/bid_caching/delete"
	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid Request",
			body:         `["123456"]`,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"success","message":"Bid Caching successfully deleted"}`,
		},
		{
			name:         "Rule id does not exists",
			body:         `["444444"]`,
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":"error","message":"Failed to delete from Bid Caching table","error":"no bid caching records found for provided rule IDs"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", endpoint, strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := appTest.Test(req)
			if err != nil {
				t.Errorf("Test %q failed: %s", test.name, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Errorf("Test %q failed: expected status code %d, got %d", test.name, test.expectedCode, resp.StatusCode)
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Test %q failed: error reading response body: %s", test.name, err)
				return
			}
			bodyString := strings.TrimSpace(string(bodyBytes))
			if bodyString != test.expectedBody {
				t.Errorf("Test %q failed: expected body %q, got %q", test.name, test.expectedBody, bodyString)
			}
		})
	}
}

func TestBidCachingUpdateHandler(t *testing.T) {
	endpoint := "/test/bid_caching/update"
	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid Request",
			body:         `{"rule_id": "123456", "bid_caching": 8, "control_percentage":0.5}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"success","message":"Bid Caching successfully updated"}`,
		},
		{
			name:         "Invalid Request",
			body:         `{"bid_caching": 8}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"status":"error","message":"could not validate bid cache update","errors":["RuleId is mandatory, validation failed"]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", endpoint, strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := appTest.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, test.expectedBody, string(body))
		})
	}
}
