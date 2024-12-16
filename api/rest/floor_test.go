package rest

import (
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"

	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/stretchr/testify/assert"

	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateFloors(t *testing.T) {
	endpoint := "/test/floor"

	tests := []struct {
		name     string
		body     string
		expected int
	}{
		{
			name:     "Valid request",
			body:     `{"publisher": "test", "device": "tablet", "country": "us", "floor": 1.0, "domain": "example.com"}`,
			expected: http.StatusOK,
		},
		{
			name:     "Missing publisher",
			body:     `{"device": "test", "country": "US", "floor": 1.0, "domain": "example.com"}`,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid device",
			body:     `{"publisher": "test", "country": "US","device":"test", "floor": 1.0, "domain": "example.com"}`,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid country",
			body:     `{"publisher": "test", "device": "test", "country": "USA", "floor": 1.0, "domain": "example.com"}`,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Missing floor",
			body:     `{"publisher": "test", "device": "test", "country": "US", "domain": "example.com"}`,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid JSON",
			body:     `{"publisher": "test" "device": "test", "country": "US", "floor": 1.0, "domain": "example.com"`,
			expected: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(fiber.MethodPost, endpoint, strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := appTest.Test(req)
		if err != nil {
			t.Errorf("Test %s failed: %s", test.name, err)
			continue
		}
		if resp.StatusCode != test.expected {
			t.Errorf("Test %s failed: expected status code %d, got %d", test.name, test.expected, resp.StatusCode)
		}
	}
}

func TestFloorPostHandler(t *testing.T) {
	endpoint := "/test/floor"
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:           "error parsing body",
			body:           `{"id":"}`,
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"message":"Invalid request body. Please ensure it's a valid JSON.","status":"error"}`,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(fiber.MethodPost, endpoint, bytes.NewBufferString(tt.body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := appTest.Test(req)
		if err != nil {
			t.Errorf("Test %s: %v", tt.name, err)
			continue
		}

		if resp.StatusCode != tt.expectedStatus {
			t.Errorf("Test [%v]: expected status code [%v], got [%v]", tt.name, tt.expectedStatus, resp.StatusCode)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Test %s: %v", tt.name, err)
			continue
		}

		if string(respBody) != tt.expectedJSON {
			t.Errorf("Test [%v]: expected JSON response [%v], got [%v]", tt.name, tt.expectedJSON, string(respBody))
		}
	}
}

func TestCreateFloorMetadataGeneration(t *testing.T) {

	tests := []struct {
		name         string
		modFloor     models.FloorSlice
		finalRules   []core.FloorRealtimeRecord
		expectedJSON string
	}{
		{
			name: "Country empty",
			modFloor: models.FloorSlice{
				{
					RuleID:    "177a2473-225c-4882-9969-93000fa75fe5",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Country:   null.StringFrom(""),
					Device:    null.StringFrom("mobile"),
					Floor:     0.11,
				},
			},
			finalRules:   []core.FloorRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=.*__os=.*__dt=mobile__pt=.*__b=.*)", "floor": 0.11, "rule_id": "177a2473-225c-4882-9969-93000fa75fe5"}]}`,
		},
		{
			name: "Same ruleId different input floor",
			modFloor: models.FloorSlice{
				{
					RuleID:    "a0d406cd-bf98-50ab-9ff2-1b314b27da65",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Country:   null.StringFrom("us"),
					Device:    null.StringFrom("mobile"),
					Floor:     0.14,
				},
			},
			finalRules:   []core.FloorRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=us__os=.*__dt=mobile__pt=.*__b=.*)", "floor": 0.14, "rule_id": "a0d406cd-bf98-50ab-9ff2-1b314b27da65"}]}`,
		},
		{
			name: "Floor sort mechanism puts most specific rule on top",
			modFloor: models.FloorSlice{
				{
					RuleID:    "a0d406cd-bf98-50ab-9ff2-1b314b27da65",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Device:    null.StringFrom("mobile"),
					Floor:     0.13,
				},
				{
					RuleID:    "a0d406cd-bf98-50ab-9ff2-1b314b27da65",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Country:   null.StringFrom("il"),
					Device:    null.StringFrom("mobile"),
					Floor:     0.14,
				},
				{
					RuleID:    "a0d406cd-bf98-50ab-9ff2-1b314b27da65",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Country:   null.StringFrom("us"),
					Device:    null.StringFrom("mobile"),
					Floor:     0.11,
				},
			},
			finalRules:   []core.FloorRealtimeRecord{},
			expectedJSON: `{"rules":[{"rule":"(p=20814__d=stream-together.org__c=il__os=.*__dt=mobile__pt=.*__b=.*)","floor":0.14,"rule_id":"a0d406cd-bf98-50ab-9ff2-1b314b27da65"},{"rule":"(p=20814__d=stream-together.org__c=us__os=.*__dt=mobile__pt=.*__b=.*)","floor":0.11,"rule_id":"a0d406cd-bf98-50ab-9ff2-1b314b27da65"},{"rule":"(p=20814__d=stream-together.org__c=.*__os=.*__dt=mobile__pt=.*__b=.*)","floor":0.13,"rule_id":"a0d406cd-bf98-50ab-9ff2-1b314b27da65"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.CreateFloorMetadata(tt.modFloor, tt.finalRules)

			resultJSON, err := json.Marshal(map[string]interface{}{"rules": result})
			if err != nil {
				t.Fatalf("Failed to marshal result to JSON: %v", err)
			}

			assert.JSONEq(t, tt.expectedJSON, string(resultJSON))
		})
	}
}

func TestFloorHistory(t *testing.T) {
	t.Parallel()

	endpoint := "/floor"
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
		history    dto.History
	}

	tests := []struct {
		name               string
		requestBody        string
		historyRequestBody string
		want               want
		wantErr            bool
	}{
		{
			name:               "validRequest_Created",
			requestBody:        `{"publisher":"333","domain":"3.com","country":"af","device":"tablet","os":"windowsphone","browser":"opera","placement_type":"rectangle","floor":0.02}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Floor"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Floor",
					Item:         "333_3.com_af_tablet_windowsphone_opera_rectangle",
					Changes: []dto.Changes{
						{Property: "floor", OldValue: nil, NewValue: float64(0.02)},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `{"publisher":"333","domain":"3.com","country":"af","device":"tablet","os":"windowsphone","browser":"opera","placement_type":"rectangle","floor":0.02}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Floor"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Floor",
					Item:         "333_3.com_af_tablet_windowsphone_opera_rectangle",
					Changes: []dto.Changes{
						{Property: "floor", OldValue: nil, NewValue: float64(0.02)},
					},
				},
			},
		},
		{
			name:               "validRequest_Updated",
			requestBody:        `{"publisher":"333","domain":"3.com","country":"af","device":"tablet","os":"windowsphone","browser":"opera","placement_type":"rectangle","floor":0.05}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Floor"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Floor",
					Item:         "333_3.com_af_tablet_windowsphone_opera_rectangle",
					Changes: []dto.Changes{
						{
							Property: "floor",
							OldValue: float64(0.02),
							NewValue: float64(0.05),
						},
					},
				},
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
			req.Header.Set(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))

			_, err = http.DefaultClient.Do(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			time.Sleep(250 * time.Millisecond)

			historyReq, err := http.NewRequest(fiber.MethodPost, baseURL+historyEndpoint, strings.NewReader(tt.historyRequestBody))
			if err != nil {
				t.Fatal(err)
			}
			historyReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			historyReq.Header.Set(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))

			historyResp, err := http.DefaultClient.Do(historyReq)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, historyResp.StatusCode)

			body, err := io.ReadAll(historyResp.Body)
			assert.NoError(t, err)
			defer historyResp.Body.Close()

			var got []dto.History
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)

			for i := range got {
				got[i].ID = 0
				got[i].Date = time.Time{}
				for j := range got[i].Changes {
					got[i].Changes[j].ID = ""
				}
			}

			assert.Contains(t, got, tt.want.history)
		})
	}
}
