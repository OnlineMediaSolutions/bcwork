package validations

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	app := fiber.New()
	app.Post("/validate", ValidateCompetitorURL)
	return app
}

func TestValidateCompetitorURL(t *testing.T) {
	app := setupApp()

	tests := []struct {
		name         string
		payload      []constant.CompetitorUpdateRequest
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name: "Invalid URL format",
			payload: []constant.CompetitorUpdateRequest{
				{Name: "Invalid URL Competitor", URL: "invalid-url.com"},
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"errors": []map[string]interface{}{
					{
						"competitor": "Invalid URL Competitor",
						"field":      "URL",
						"message":    "Competitor 'Invalid URL Competitor': URL must be valid and start with either 'http' or 'https'.",
					},
				},
			},
		},
		{
			name: "Missing URL",
			payload: []constant.CompetitorUpdateRequest{
				{Name: "Missing URL Competitor", URL: ""},
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status": "error",
				"errors": []map[string]interface{}{
					{
						"competitor": "Missing URL Competitor",
						"field":      "URL",
						"message":    "Competitor 'Missing URL Competitor': URL must be valid and start with either 'http' or 'https'.",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, err := json.Marshal(tt.payload)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/validate", bytes.NewReader(payloadBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			var respBody map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBody["status"], respBody["status"])
			assert.ElementsMatch(t, tt.expectedBody["errors"], respBody["errors"])
		})
	}
}
