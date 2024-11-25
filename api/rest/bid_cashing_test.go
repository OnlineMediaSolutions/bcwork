package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestValidateBidCashings(t *testing.T) {
	app := fiber.New()
	app.Post("/bid_cashing", validations.ValidateFactor, func(c *fiber.Ctx) error {
		return c.SendString("Bid Cashing created successfully")
	})

	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{

		{
			name:         "Missing publisher",
			body:         `{"device": "test", "country": "US", "bid_cashing": 1, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Publisher is mandatory, validation failed","status":"error"}`,
		},
		{
			name:         "Invalid device",
			body:         `{"publisher":"1234", "device": "test", "country": "US", "bid_cashing": 1, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Device should be in the allowed list","status":"error"}`,
		},
		{
			name:         "Invalid country",
			body:         `{"publisher": "test", "device": "tablet", "country": "USA", "bid_cashing": 1, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Country code must be 2 characters long and should be in the allowed list","status":"error"}`,
		},
		{
			name:         "Missing bid cashing",
			body:         `{"publisher": "test", "device": "tablet", "country": "us", "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Bid cashing is mandatory, validation failed","status":"error"}`,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/bid_cashing", strings.NewReader(test.body))
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

func TestCreateBidCashingMetadataGeneration(t *testing.T) {
	tests := []struct {
		name         string
		modBC        models.BidCashingSlice
		finalRules   []core.BidCashingRealtimeRecord
		expectedJSON string
	}{
		{
			name: "Sort By Correct Order",
			modBC: models.BidCashingSlice{
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Device:     null.StringFrom("mobile"),
					BidCashing: 1,
				},
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Device:     null.StringFrom("mobile"),
					Country:    null.StringFrom("il"),
					BidCashing: 11,
				},
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Device:     null.StringFrom("mobile"),
					Country:    null.StringFrom("us"),
					BidCashing: 14,
				},
			},
			finalRules:   []core.BidCashingRealtimeRecord{},
			expectedJSON: `{"rules":[{"rule":"(p=20814__d=stream-together.org__c=il__os=.*__dt=mobile__pt=.*__b=.*)","bid_cashing":14,"rule_id":"cc11f229-1d4a-5bd2-a6d0-5fae8c7a9bf4"},{"rule":"(p=20814__d=stream-together.org__c=us__os=.*__dt=mobile__pt=.*__b=.*)","factor":0.14,"rule_id":"a0d406cd-bf98-50ab-9ff2-1b314b27da65"},{"rule":"(p=20814__d=stream-together.org__c=.*__os=.*__dt=mobile__pt=.*__b=.*)","factor":0.12,"rule_id":"cb45cb97-5ca2-503d-9008-317dbbe26d10"}]}`,
		},
		{
			name: "Device with null value",
			modBC: models.BidCashingSlice{
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Country:    null.StringFrom("us"),
					BidCashing: 11,
				},
			},
			finalRules:   []core.BidCashingRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=us__os=.*__dt=.*__pt=.*__b=.*)", "factor": 0.11, "rule_id": "ad18394a-ee20-58c2-bb9b-dd459550a9f7"}]}`,
		},
		{
			name: "Same ruleId different input factor",
			modBC: models.BidCashingSlice{
				{
					RuleID:     "",
					Publisher:  "20814",
					Domain:     "stream-together.org",
					Country:    null.StringFrom("us"),
					Device:     null.StringFrom("mobile"),
					BidCashing: 14,
				},
			},
			finalRules:   []core.BidCashingRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=us__os=.*__dt=mobile__pt=.*__b=.*)", "factor": 0.14, "rule_id": "a0d406cd-bf98-50ab-9ff2-1b314b27da65"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.CreateBidCashingMetadata(tt.modBC, tt.finalRules)

			resultJSON, err := json.Marshal(map[string]interface{}{"rules": result})
			if err != nil {
				t.Fatalf("Failed to marshal result to JSON: %v", err)
			}

			assert.JSONEq(t, tt.expectedJSON, string(resultJSON))
		})
	}
}

func Test_BC_ToModel(t *testing.T) {
	t.Parallel()

	type args struct {
		bidCashing *core.BidCashing
	}

	tests := []struct {
		name     string
		args     args
		expected *models.BidCashing
	}{
		{
			name: "All fields populated",
			args: args{
				bidCashing: &core.BidCashing{
					RuleId:        "50afedac-d41a-53b0-a922-2c64c6e80623",
					Publisher:     "Publisher1",
					Domain:        "example.com",
					BidCashing:    1,
					OS:            "Windows",
					Country:       "US",
					Device:        "Desktop",
					PlacementType: "Banner",
					Browser:       "Chrome",
				},
			},
			expected: &models.BidCashing{
				RuleID:        "50afedac-d41a-53b0-a922-2c64c6e80623",
				Publisher:     "Publisher1",
				Domain:        "example.com",
				BidCashing:    1,
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
				bidCashing: &core.BidCashing{
					RuleId:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
					Publisher:     "Publisher2",
					Domain:        "example.org",
					BidCashing:    1,
					OS:            "",
					Country:       "CA",
					Device:        "",
					PlacementType: "Sidebar",
					Browser:       "",
				},
			},
			expected: &models.BidCashing{
				RuleID:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
				Publisher:     "Publisher2",
				Domain:        "example.org",
				BidCashing:    1,
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
				bidCashing: &core.BidCashing{
					RuleId:        "966affd7-d087-57a2-baff-55b926f4c32d",
					Publisher:     "",
					Domain:        "",
					BidCashing:    1,
					OS:            "",
					Country:       "",
					Device:        "",
					PlacementType: "",
					Browser:       "",
				},
			},
			expected: &models.BidCashing{
				RuleID:        "966affd7-d087-57a2-baff-55b926f4c32d",
				Publisher:     "",
				Domain:        "",
				BidCashing:    1,
				Country:       null.String{},
				Os:            null.String{},
				Device:        null.String{},
				PlacementType: null.String{},
				Browser:       null.String{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mod := tt.args.bidCashing.ToModel()
			assert.Equal(t, tt.expected, mod)
		})
	}
}

func TestBIdCashingHistory(t *testing.T) {
	endpoint := "/bid_cashing"
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
			requestBody:        `{"publisher":"333","domain":"3.com","country":"af","device":"tablet","os":"windowsphone","browser":"opera","placement_type":"rectangle","bid_cashing":2}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Bidder Targeting"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Bidder Targeting",
					Item:         "333_3.com_af_tablet_windowsphone_opera_rectangle",
					Changes: []dto.Changes{
						{Property: "browser", OldValue: nil, NewValue: "opera"},
						{Property: "country", OldValue: nil, NewValue: "af"},
						{Property: "device", OldValue: nil, NewValue: "tablet"},
						{Property: "domain", OldValue: nil, NewValue: "3.com"},
						{Property: "bid_cashing", OldValue: nil, NewValue: 2},
						{Property: "os", OldValue: nil, NewValue: "windowsphone"},
						{Property: "placement_type", OldValue: nil, NewValue: "rectangle"},
						{Property: "publisher", OldValue: nil, NewValue: "333"},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `{"publisher":"333","domain":"3.com","country":"af","device":"tablet","os":"windowsphone","browser":"opera","placement_type":"rectangle","bid_cashing":2}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Bidder Targeting"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Bidder Targeting",
					Item:         "333_3.com_af_tablet_windowsphone_opera_rectangle",
					Changes: []dto.Changes{
						{Property: "browser", OldValue: nil, NewValue: "opera"},
						{Property: "country", OldValue: nil, NewValue: "af"},
						{Property: "device", OldValue: nil, NewValue: "tablet"},
						{Property: "domain", OldValue: nil, NewValue: "3.com"},
						{Property: "factor", OldValue: nil, NewValue: int64(2)},
						{Property: "os", OldValue: nil, NewValue: "windowsphone"},
						{Property: "placement_type", OldValue: nil, NewValue: "rectangle"},
						{Property: "publisher", OldValue: nil, NewValue: "333"},
					},
				},
			},
		},
		{
			name:               "validRequest_Updated",
			requestBody:        `{"publisher":"333","domain":"3.com","country":"af","device":"tablet","os":"windowsphone","browser":"opera","placement_type":"rectangle","bid_cashing":5}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Bidder Targeting"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Bidder Targeting",
					Item:         "333_3.com_af_tablet_windowsphone_opera_rectangle",
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
