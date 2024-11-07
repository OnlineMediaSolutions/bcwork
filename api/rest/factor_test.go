package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/m6yf/bcwork/validations"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func TestValidateFactors(t *testing.T) {
	app := fiber.New()
	app.Post("/factor", validations.ValidateFactor, func(c *fiber.Ctx) error {
		return c.SendString("Factor created successfully")
	})

	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{

		{
			name:         "Missing publisher",
			body:         `{"device": "test", "country": "US", "factor": 1.0, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Publisher is mandatory, validation failed","status":"error"}`,
		},
		{
			name:         "Invalid device",
			body:         `{"publisher":"1234", "device": "test", "country": "US", "factor": 1.0, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Device should be in the allowed list","status":"error"}`,
		},
		{
			name:         "Invalid country",
			body:         `{"publisher": "test", "device": "tablet", "country": "USA", "factor": 1.0, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Country code must be 2 characters long and should be in the allowed list","status":"error"}`,
		},
		{
			name:         "Missing factor",
			body:         `{"publisher": "test", "device": "tablet", "country": "us", "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Factor is mandatory, validation failed","status":"error"}`,
		},
		{
			name:         "Invalid JSON",
			body:         `{"publisher": "test" "device": "test", "country": "US", "factor": 1.0, "domain": "example.com"`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Invalid request body for factor. Please ensure it's a valid JSON.","status":"error"}`,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/factor", strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
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

func TestCreateFactorMetadataGeneration(t *testing.T) {
	tests := []struct {
		name         string
		modFactor    models.FactorSlice
		finalRules   []core.FactorRealtimeRecord
		expectedJSON string
	}{
		{
			name: "Country with null value",
			modFactor: models.FactorSlice{
				{
					RuleID:    "",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Device:    null.StringFrom("mobile"),
					Factor:    0.11,
				},
			},
			finalRules:   []core.FactorRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=.*__os=.*__dt=mobile__pt=.*__b=.*)", "factor": 0.11, "rule_id": "cb45cb97-5ca2-503d-9008-317dbbe26d10"}]}`,
		},
		{
			name: "Device with null value",
			modFactor: models.FactorSlice{
				{
					RuleID:    "",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Country:   null.StringFrom("us"),
					Factor:    0.11,
				},
			},
			finalRules:   []core.FactorRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=us__os=.*__dt=.*__pt=.*__b=.*)", "factor": 0.11, "rule_id": "ad18394a-ee20-58c2-bb9b-dd459550a9f7"}]}`,
		},
		{
			name: "Same ruleId different input factor",
			modFactor: models.FactorSlice{
				{
					RuleID:    "",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Country:   null.StringFrom("us"),
					Device:    null.StringFrom("mobile"),
					Factor:    0.14,
				},
			},
			finalRules:   []core.FactorRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=us__os=.*__dt=mobile__pt=.*__b=.*)", "factor": 0.14, "rule_id": "a0d406cd-bf98-50ab-9ff2-1b314b27da65"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.CreateFactorMetadata(tt.modFactor, tt.finalRules)

			resultJSON, err := json.Marshal(map[string]interface{}{"rules": result})
			if err != nil {
				t.Fatalf("Failed to marshal result to JSON: %v", err)
			}

			assert.JSONEq(t, tt.expectedJSON, string(resultJSON))
		})
	}
}

func Test_Factor_ToModel_(t *testing.T) {
	tests := []struct {
		name     string
		factor   *core.Factor
		expected *models.Factor
	}{
		{
			name: "All fields populated",
			factor: &core.Factor{
				RuleId:        "50afedac-d41a-53b0-a922-2c64c6e80623",
				Publisher:     "Publisher1",
				Domain:        "example.com",
				Factor:        1,
				OS:            "Windows",
				Country:       "US",
				Device:        "Desktop",
				PlacementType: "Banner",
				Browser:       "Chrome",
			},
			expected: &models.Factor{
				RuleID:        "50afedac-d41a-53b0-a922-2c64c6e80623",
				Publisher:     "Publisher1",
				Domain:        "example.com",
				Factor:        1,
				Country:       null.StringFrom("US"),
				Os:            null.StringFrom("Windows"),
				Device:        null.StringFrom("Desktop"),
				PlacementType: null.StringFrom("Banner"),
				Browser:       null.StringFrom("Chrome"),
			},
		},
		{
			name: "Some fields empty",
			factor: &core.Factor{
				RuleId:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
				Publisher:     "Publisher2",
				Domain:        "example.org",
				Factor:        1,
				OS:            "",
				Country:       "CA",
				Device:        "",
				PlacementType: "Sidebar",
				Browser:       "",
			},
			expected: &models.Factor{
				RuleID:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
				Publisher:     "Publisher2",
				Domain:        "example.org",
				Factor:        1,
				Country:       null.StringFrom("CA"),
				Os:            null.String{},
				Device:        null.String{},
				PlacementType: null.StringFrom("Sidebar"),
				Browser:       null.String{},
			},
		},
		{
			name: "All fields empty",
			factor: &core.Factor{
				RuleId:        "966affd7-d087-57a2-baff-55b926f4c32d",
				Publisher:     "",
				Domain:        "",
				Factor:        1,
				OS:            "",
				Country:       "",
				Device:        "",
				PlacementType: "",
				Browser:       "",
			},
			expected: &models.Factor{
				RuleID:        "966affd7-d087-57a2-baff-55b926f4c32d",
				Publisher:     "",
				Domain:        "",
				Factor:        1,
				Country:       null.String{},
				Os:            null.String{},
				Device:        null.String{},
				PlacementType: null.String{},
				Browser:       null.String{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := tt.factor.ToModel()
			assert.Equal(t, tt.expected, mod)
		})
	}
}

func Test_Factor_ToModel(t *testing.T) {
	t.Parallel()

	type args struct {
		factor *core.Factor
	}

	tests := []struct {
		name     string
		args     args
		expected *models.Factor
	}{
		{
			name: "All fields populated",
			args: args{
				factor: &core.Factor{
					RuleId:        "50afedac-d41a-53b0-a922-2c64c6e80623",
					Publisher:     "Publisher1",
					Domain:        "example.com",
					Factor:        1,
					OS:            "Windows",
					Country:       "US",
					Device:        "Desktop",
					PlacementType: "Banner",
					Browser:       "Chrome",
				},
			},
			expected: &models.Factor{
				RuleID:        "50afedac-d41a-53b0-a922-2c64c6e80623",
				Publisher:     "Publisher1",
				Domain:        "example.com",
				Factor:        1,
				Country:       null.StringFrom("US"),
				Os:            null.StringFrom("Windows"),
				Device:        null.StringFrom("Desktop"),
				PlacementType: null.StringFrom("Banner"),
				Browser:       null.StringFrom("Chrome"),
			},
		},
		{
			name: "Some fields empty",
			args: args{
				factor: &core.Factor{
					RuleId:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
					Publisher:     "Publisher2",
					Domain:        "example.org",
					Factor:        1,
					OS:            "",
					Country:       "CA",
					Device:        "",
					PlacementType: "Sidebar",
					Browser:       "",
				},
			},
			expected: &models.Factor{
				RuleID:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
				Publisher:     "Publisher2",
				Domain:        "example.org",
				Factor:        1,
				Country:       null.StringFrom("CA"),
				Os:            null.String{},
				Device:        null.String{},
				PlacementType: null.StringFrom("Sidebar"),
				Browser:       null.String{},
			},
		},
		{
			name: "All fields empty",
			args: args{
				factor: &core.Factor{
					RuleId:        "966affd7-d087-57a2-baff-55b926f4c32d",
					Publisher:     "",
					Domain:        "",
					Factor:        1,
					OS:            "",
					Country:       "",
					Device:        "",
					PlacementType: "",
					Browser:       "",
				},
			},
			expected: &models.Factor{
				RuleID:        "966affd7-d087-57a2-baff-55b926f4c32d",
				Publisher:     "",
				Domain:        "",
				Factor:        1,
				Country:       null.String{},
				Os:            null.String{},
				Device:        null.String{},
				PlacementType: null.String{},
				Browser:       null.String{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mod := tt.args.factor.ToModel()
			assert.Equal(t, tt.expected, mod)
		})
	}
}

func TestFactorHistory(t *testing.T) {
	endpoint := "/factor"
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
		hasHistory bool
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
			requestBody:        `{"publisher":"333","domain":"3.com","country":"af","device":"tablet","os":"windowsphone","browser":"opera","placement_type":"rectangle","factor":0.02}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Bidder Targeting"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Bidder Targeting",
					Item:         "af_tablet_windowsphone_opera_rectangle",
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `{"publisher":"333","domain":"3.com","country":"af","device":"tablet","os":"windowsphone","browser":"opera","placement_type":"rectangle","factor":0.02}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Bidder Targeting"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Bidder Targeting",
					Item:         "af_tablet_windowsphone_opera_rectangle",
				},
			},
		},
		{
			name:               "validRequest_Updated",
			requestBody:        `{"publisher":"333","domain":"3.com","country":"af","device":"tablet","os":"windowsphone","browser":"opera","placement_type":"rectangle","factor":0.05}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Bidder Targeting"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Bidder Targeting",
					Item:         "af_tablet_windowsphone_opera_rectangle",
					Changes: []dto.Changes{
						{
							Property: "factor",
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

			var (
				got   []dto.History
				found bool
			)
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)
			for i := range got {
				got[i].ID = 0
				got[i].Date = time.Time{}
				for j := range got[i].Changes {
					got[i].Changes[j].ID = ""
				}

				if reflect.DeepEqual(tt.want.history, got[i]) {
					assert.Equal(t, tt.want.history, got[i])
					found = true
				}
			}
			assert.Equal(t, true, found)
		})
	}
}
