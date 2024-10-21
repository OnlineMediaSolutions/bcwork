package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/m6yf/bcwork/validations"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
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

func TestConvertingAllValuesFactor(t *testing.T) {
	tests := []struct {
		name     string
		data     constant.FactorUpdateRequest
		expected constant.FactorUpdateRequest
	}{
		{
			name: "device and country with empty values",
			data: constant.FactorUpdateRequest{
				Device:    "",
				Country:   "",
				Publisher: "345",
				Domain:    "active.com",
			},
			expected: constant.FactorUpdateRequest{
				Device:    "",
				Country:   "",
				Publisher: "345",
				Domain:    "active.com",
			},
		},
	}

	for _, tt := range tests {
		utils.ConvertingAllValues(&tt.data)
		if !reflect.DeepEqual(tt.data, tt.expected) {
			t.Errorf("Test %s failed: got %+v, expected %+v", tt.name, tt.data, tt.expected)
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
					RuleID:    "cb45cb97-5ca2-503d-9008-317dbbe26d10",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Country:   "",
					Device:    "mobile",
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
					RuleID:    "4f63927a-2497-5496-82c1-e748277afe24",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Country:   "",
					Device:    "",
					Factor:    0.11,
				},
			},
			finalRules:   []core.FactorRealtimeRecord{},
			expectedJSON: `{"rules": [{"rule": "(p=20814__d=stream-together.org__c=.*__os=.*__dt=.*__pt=.*__b=.*)", "factor": 0.11, "rule_id": "4f63927a-2497-5496-82c1-e748277afe24"}]}`,
		},
		{
			name: "Same ruleId different input factor",
			modFactor: models.FactorSlice{
				{
					RuleID:    "a0d406cd-bf98-50ab-9ff2-1b314b27da65",
					Publisher: "20814",
					Domain:    "stream-together.org",
					Country:   "us",
					Device:    "mobile",
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
