package rest

import (
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateRefreshCache(t *testing.T) {
	endpoint := "/test/refresh_cache/set"
	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Missing refresh cache",
			body:         `{"publisher": "21038", "domain": "bubu8.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Refresh cache value not allowed, it should be \u003c= 500 and \u003e= 1","status":"error"}`,
		},

		{
			name:         "Wrong value refresh cache",
			body:         `{"publisher": "21038","refresh_cache":1800, "domain": "bubu8.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Refresh cache value not allowed, it should be \u003c= 500 and \u003e= 1","status":"error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", endpoint, strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := appTest.Test(req)
			if err != nil {
				t.Errorf("Test %s failed: %s", test.name, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Errorf("Test %s failed: expected status code %d, got %d", test.name, test.expectedCode, resp.StatusCode)
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Test %s failed: error reading response body: %s", test.name, err)
				return
			}
			bodyString := strings.TrimSpace(string(bodyBytes))
			if bodyString != test.expectedBody {
				t.Errorf("Test %s failed: expected body %s, got %s", test.name, test.expectedBody, bodyString)
			}
		})
	}
}

func TestSetRefreshCache(t *testing.T) {
	endpoint := "/test/refresh_cache/set"
	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{

		{
			name:         "Correct values in refresh cache",
			body:         `{"publisher": "21038","refresh_cache":180, "domain": "test.com"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"success","message":"Refresh cache successfully created"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", endpoint, strings.NewReader(test.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := appTest.Test(req)
			if err != nil {
				t.Errorf("Test %s failed: %s", test.name, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Errorf("Test %s failed: expected status code %d, got %d", test.name, test.expectedCode, resp.StatusCode)
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Test %s failed: error reading response body: %s", test.name, err)
				return
			}
			bodyString := strings.TrimSpace(string(bodyBytes))
			if bodyString != test.expectedBody {
				t.Errorf("Test %s failed: expected body %s, got %s", test.name, test.expectedBody, bodyString)
			}
		})
	}
}

func TestRefreshCacheUpdateHandler(t *testing.T) {
	endpoint := "/test/refresh_cache/update"
	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid Request",
			body:         `{"rule_id": "123456", "refresh_cache": 8}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"success","message":"Refresh Cache successfully updated"}`,
		},
		{
			name:         "Invalid Request",
			body:         `{"refresh_cache": 8}`,
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

func TestRefreshCacheDeleteHandler(t *testing.T) {
	endpoint := "/test/refresh_cache/delete"
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
			expectedBody: `{"status":"success","message":"Refresh cache successfully deleted"}`,
		},
		{
			name:         "Rule id does not exists",
			body:         `["444444"]`,
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"status":"error","message":"Failed to delete from  Refresh cache table","error":"failed to delete from metadata table no value found for these keys: 444444"}`,
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

func Test_LR_ToModel(t *testing.T) {
	t.Parallel()

	type args struct {
		refreshCache *dto.RefreshCache
	}

	tests := []struct {
		name     string
		args     args
		expected *models.RefreshCache
	}{
		{
			name: "All fields populated",
			args: args{
				refreshCache: &dto.RefreshCache{
					RuleId:       "50afedac-d41a-53b0-a922-2c64c6e80623",
					Publisher:    "Publisher1",
					Domain:       "example.com",
					RefreshCache: 1,
				},
			},
			expected: &models.RefreshCache{
				RuleID:       "50afedac-d41a-53b0-a922-2c64c6e80623",
				Publisher:    "Publisher1",
				Domain:       null.StringFrom("example.com"),
				RefreshCache: 1,
			},
		},
		{
			name: "Domain value empty",
			args: args{
				refreshCache: &dto.RefreshCache{
					RuleId:       "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
					Publisher:    "Publisher2",
					Domain:       "",
					RefreshCache: 1,
				},
			},
			expected: &models.RefreshCache{
				RuleID:       "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
				Publisher:    "Publisher2",
				Domain:       null.String{String: "", Valid: false},
				RefreshCache: 1,
			},
		},
		{
			name: "All fields empty",
			args: args{
				refreshCache: &dto.RefreshCache{
					RuleId:       "966affd7-d087-57a2-baff-55b926f4c32d",
					Publisher:    "",
					Domain:       "",
					RefreshCache: 1,
				},
			},
			expected: &models.RefreshCache{
				RuleID:       "966affd7-d087-57a2-baff-55b926f4c32d",
				Publisher:    "",
				Domain:       null.String{String: "", Valid: false},
				RefreshCache: 1,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mod := tt.args.refreshCache.ToModel()
			assert.Equal(t, tt.expected, mod)
		})
	}
}
