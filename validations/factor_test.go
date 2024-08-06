package validations

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupFactorApp() *fiber.App {
	app := fiber.New()

	app.Post("/factor", ValidateFactor, func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "error",
			"message": "Request body is valid",
		})
	})

	return app
}

func TestValidateFactor(t *testing.T) {
	app := setupFactorApp()

	tests := []struct {
		name       string
		body       map[string]interface{}
		statusCode int
		response   map[string]string
	}{

		{
			name: "Missing Publisher",
			body: map[string]interface{}{
				"domain":  "testdomain.com",
				"device":  "mobile",
				"factor":  1.23,
				"country": "us",
			},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "Publisher is mandatory, validation failed",
			},
		},
		{
			name: "Missing Factor",
			body: map[string]interface{}{
				"device":    "mobile",
				"country":   "us",
				"publisher": "somePublisher",
			},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "Factor is mandatory, validation failed",
			},
		},
		{
			name: "Wrong value for Factor",
			body: map[string]interface{}{
				"device":    "mobile",
				"factor":    100,
				"country":   "us",
				"publisher": "somePublisher",
			},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "Factor value not allowed, it should be >= 0.01 and <= 10.00",
			},
		},
		{
			name: "Wrong value for Country",
			body: map[string]interface{}{
				"device":    "mobile",
				"factor":    1,
				"country":   "usa",
				"publisher": "somePublisher",
			},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "Country code must be 2 characters long and should be in the allowed list",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/factor", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode)

			var response map[string]string
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.response, response)
		})
	}
}
