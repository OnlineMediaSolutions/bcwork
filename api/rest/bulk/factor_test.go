package bulk

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

func TestValidateBulkFactors(t *testing.T) {
	app := fiber.New()
	app.Post("/bulkFactors", validations.ValidateBulkFactors, func(c *fiber.Ctx) error {
		return c.SendString("Bulk factors created successfully")
	})

	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid bulk request",
			body:         `[{"publisher": "1234", "device": "mobile", "country": "us", "factor": 1.0, "domain": "example.com"}, {"publisher": "5678", "device": "tablet", "country": "ca", "factor": 1.2, "domain": "example.org"}]`,
			expectedCode: http.StatusOK,
			expectedBody: `Bulk factors created successfully`,
		},
		{
			name:         "Missing publisher in one request",
			body:         `[{"device": "mobile", "country": "us", "factor": 1.0, "domain": "example.com","domain": "example.com"}, {"publisher": "5678", "device": "tablet", "country": "ca", "factor": 1.2, "domain": "example.org"}]`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Publisher is mandatory, validation failed","status":"error"}`,
		},
		{
			name:         "Missing domain in one of  requests",
			body:         `[{"device": "mobile", "country": "us", "publisher":"1234", "factor": 1.0}, {"publisher": "5678", "device": "tablet", "country": "ca", "factor": 1.2, "domain": "example.org"}]`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"message":"Domain is mandatory, validation failed","status":"error"}`,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/bulkFactors", strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Errorf("Test %s failed: %s", test.name, err)
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := strings.TrimSpace(string(bodyBytes))
		assert.Equal(t, test.expectedCode, resp.StatusCode, "Unexpected status code in %s", test.name)
		assert.Equal(t, test.expectedBody, bodyString, "Unexpected response body in %s", test.name)
	}
}

func TestCreateBulkFactorMetadataGeneration(t *testing.T) {
	tests := []struct {
		name         string
		modFactor    models.FactorSlice
		finalRules   []core.FactorRealtimeRecord
		expectedJSON string
	}{
		{
			name: "Multiple factors",
			modFactor: models.FactorSlice{
				{
					Publisher: "20814",
					Domain:    "example.com",
					Country:   null.StringFrom("US"),
					Device:    null.StringFrom("mobile"),
					Factor:    0.11,
				},
				{
					Publisher: "20814",
					Domain:    "example.com",
					Country:   null.StringFrom("CA"),
					Device:    null.StringFrom("tablet"),
					Factor:    1.2,
				},
			},
			finalRules: []core.FactorRealtimeRecord{},
			expectedJSON: `{"rules": [
                {"rule": "(p=20814__d=example.com__c=US__os=.*__dt=mobile__pt=.*__b=.*)", "factor": 0.11, "rule_id": "81e33064-b91f-5ad8-88cb-3375deb3a8bd"},
                {"rule": "(p=20814__d=example.com__c=CA__os=.*__dt=desktop__pt=.*__b=.*)", "factor": 1.2, "rule_id": "bde23553-cbd7-51c0-9116-b20f0dd54e28"}
            ]}`,
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
