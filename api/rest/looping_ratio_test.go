package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/validations"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateLoopingRatio(t *testing.T) {
	app := fiber.New()
	app.Post("/looping_ratio", validations.ValidateLoopingRatio, func(c *fiber.Ctx) error {
		return c.SendString("Looping ratio created successfully")
	})

	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{

		{
			name:         "Invalid device",
			body:         `{"publisher":"1234", "device": "mm", "country": "US", "looping_ratio": 1, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Device should be in the allowed list","status":"error"}`,
		},
		{
			name:         "Invalid country",
			body:         `{"publisher": "test", "device": "tablet", "country": "USA", "looping_ratio": 1, "domain": "example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Country code must be 2 characters long and should be in the allowed list","status":"error"}`,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/looping_ratio", strings.NewReader(test.body))
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

func TestCreateLoopingRatioMetadataGeneration(t *testing.T) {
	tests := []struct {
		name         string
		modBC        models.LoopingRatioSlice
		finalRules   []core.LoopingRatioRealtimeRecord
		expectedJSON string
	}{
		{
			name: "Sort By Correct Order",
			modBC: models.LoopingRatioSlice{
				{
					RuleID:       "",
					Publisher:    "20814",
					Domain:       "stream-together.org",
					Device:       null.StringFrom("mobile"),
					LoopingRatio: 12,
				},
				{
					RuleID:       "",
					Publisher:    "20814",
					Domain:       "stream-together.org",
					Device:       null.StringFrom("mobile"),
					Country:      null.StringFrom("il"),
					LoopingRatio: 11,
				},
				{
					RuleID:       "",
					Publisher:    "20814",
					Domain:       "stream-together.org",
					Device:       null.StringFrom("mobile"),
					Country:      null.StringFrom("us"),
					LoopingRatio: 14,
				},
			},
			finalRules:   []core.LoopingRatioRealtimeRecord{},
			expectedJSON: `{"rules":[{"rule":"(p=20814__d=stream-together.org__c=il__os=.*__dt=mobile__pt=.*__b=.*)","looping_ratio":11,"rule_id":"cc11f229-1d4a-5bd2-a6d0-5fae8c7a9bf4"},{"rule":"(p=20814__d=stream-together.org__c=us__os=.*__dt=mobile__pt=.*__b=.*)","looping_ratio":14,"rule_id":"a0d406cd-bf98-50ab-9ff2-1b314b27da65"},{"rule":"(p=20814__d=stream-together.org__c=.*__os=.*__dt=mobile__pt=.*__b=.*)","looping_ratio":12,"rule_id":"cb45cb97-5ca2-503d-9008-317dbbe26d10"}]}`,
		},
		{
			name: "Device with null value",
			modBC: models.LoopingRatioSlice{
				{
					RuleID:       "",
					Publisher:    "20814",
					Domain:       "stream-together.org",
					Country:      null.StringFrom("us"),
					LoopingRatio: 11,
				},
			},
			finalRules:   []core.LoopingRatioRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=us__os=.*__dt=.*__pt=.*__b=.*)", "looping_ratio": 11, "rule_id": "ad18394a-ee20-58c2-bb9b-dd459550a9f7"}]}`,
		},
		{
			name: "Same ruleId different input looping ratio",
			modBC: models.LoopingRatioSlice{
				{
					RuleID:       "",
					Publisher:    "20814",
					Domain:       "stream-together.org",
					Country:      null.StringFrom("us"),
					Device:       null.StringFrom("mobile"),
					LoopingRatio: 14,
				},
			},
			finalRules:   []core.LoopingRatioRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=us__os=.*__dt=mobile__pt=.*__b=.*)", "looping_ratio": 14, "rule_id": "a0d406cd-bf98-50ab-9ff2-1b314b27da65"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.CreateLoopingRatioMetadata(tt.modBC, tt.finalRules)

			resultJSON, err := json.Marshal(map[string]interface{}{"rules": result})
			if err != nil {
				t.Fatalf("Failed to marshal result to JSON: %v", err)
			}

			assert.JSONEq(t, tt.expectedJSON, string(resultJSON))
		})
	}
}

func Test_LR_ToModel(t *testing.T) {
	t.Parallel()

	type args struct {
		loopingRatio *core.LoopingRatio
	}

	tests := []struct {
		name     string
		args     args
		expected *models.LoopingRatio
	}{
		{
			name: "All fields populated",
			args: args{
				loopingRatio: &core.LoopingRatio{
					RuleId:        "50afedac-d41a-53b0-a922-2c64c6e80623",
					Publisher:     "Publisher1",
					Domain:        "example.com",
					LoopingRatio:  1,
					OS:            "Windows",
					Country:       "US",
					Device:        "Desktop",
					PlacementType: "Banner",
					Browser:       "Chrome",
				},
			},
			expected: &models.LoopingRatio{
				RuleID:        "50afedac-d41a-53b0-a922-2c64c6e80623",
				Publisher:     "Publisher1",
				Domain:        "example.com",
				LoopingRatio:  1,
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
				loopingRatio: &core.LoopingRatio{
					RuleId:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
					Publisher:     "Publisher2",
					Domain:        "example.org",
					LoopingRatio:  1,
					OS:            "",
					Country:       "CA",
					Device:        "",
					PlacementType: "Sidebar",
					Browser:       "",
				},
			},
			expected: &models.LoopingRatio{
				RuleID:        "d823a92a-83e5-5c2b-a067-b982d6cdfaf8",
				Publisher:     "Publisher2",
				Domain:        "example.org",
				LoopingRatio:  1,
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
				loopingRatio: &core.LoopingRatio{
					RuleId:        "966affd7-d087-57a2-baff-55b926f4c32d",
					Publisher:     "",
					Domain:        "",
					LoopingRatio:  1,
					OS:            "",
					Country:       "",
					Device:        "",
					PlacementType: "",
					Browser:       "",
				},
			},
			expected: &models.LoopingRatio{
				RuleID:        "966affd7-d087-57a2-baff-55b926f4c32d",
				Publisher:     "",
				Domain:        "",
				LoopingRatio:  1,
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

			mod := tt.args.loopingRatio.ToModel()
			assert.Equal(t, tt.expected, mod)
		})
	}
}
