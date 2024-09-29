package validations

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
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
		payload      core.CompetitorUpdateRequest
		expectedCode int
		expectedBody string
	}{
		{
			name: "Invalid URL format",
			payload: core.CompetitorUpdateRequest{
				URL: "invalid-url.com",
			},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: "URL must be valid and start with either 'http' or 'https'.",
		},
		{
			name:         "Missing URL",
			payload:      core.CompetitorUpdateRequest{},
			expectedCode: fiber.StatusBadRequest,
			expectedBody: "URL must be valid and start with either 'http' or 'https'.",
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

			if tt.expectedCode != fiber.StatusOK {
				var respBody map[string]string
				err := json.NewDecoder(resp.Body).Decode(&respBody)
				assert.NoError(t, err)

				assert.Contains(t, respBody["message"], tt.expectedBody)
			}
		})
	}
}
