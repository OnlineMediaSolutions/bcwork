package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateBidCachings(t *testing.T) {
	endpoint := "/test/bid_caching/set"

	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{

		{
			name:         "Invalid device",
			body:         `{"publisher":"1234", "device": "mm", "country": "US", "bid_caching": 1, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Device should be in the allowed list","status":"error"}`,
		},
		{
			name:         "Invalid country",
			body:         `{"publisher": "test", "device": "tablet", "country": "USA", "bid_caching": 1, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Country code must be 2 characters long and should be in the allowed list","status":"error"}`,
		},
		{
			name:         "Valid request",
			body:         `{"publisher": "test","bid_caching": 1, "domain": "example.com"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"success","message":"Bid Caching successfully created"}`,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", endpoint, strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := appTest.Test(req)
		if err != nil {
			t.Errorf("Test %s failed: %s", test.name, err)
			continue
		}

		if resp.StatusCode != test.expectedCode {
			t.Errorf("Test %s failed: expected status code %d, got %d", test.name, test.expectedCode, resp.StatusCode)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := strings.TrimSpace(string(bodyBytes))
		if bodyString != test.expectedBody {
			t.Errorf("Test %s failed: expected body %s, got %s", test.name, test.expectedBody, bodyString)
		}
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
			body:         `{"rule_id": "123456", "bid_caching": 8}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"success","message":"Bid Caching successfully updated"}`,
		},
		{
			name:         "Invalid Request",
			body:         `{"bid_caching": 8}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"RuleId is mandatory, validation failed","status":"error"}`,
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
			requestBody: `{"filter": {"active": true} }`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"rule_id":"123456","publisher":"21038","domain":"oms.com","country":"","device":"","bid_caching":10,"browser":"","os":"","placement_type":"","active":true}]`,
			},
		},
		{
			name:        "Bid caching with active filter=false",
			requestBody: `{"filter": {"active": false} }`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"rule_id":"1234567","publisher":"21000","domain":"brightcom.com","country":"","device":"","bid_caching":10,"browser":"","os":"","placement_type":"","active":false}]`,
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
