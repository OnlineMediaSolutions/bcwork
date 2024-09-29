package dpo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	app := fiber.New()

	app.Get("/test", ValidateQueryParams, func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Query parameters are valid",
		})
	})

	app.Post("/dpo", ValidateDPO, func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Request body is valid",
		})
	})

	return app
}

func TestValidateQueryParams(t *testing.T) {
	app := setupApp()

	tests := []struct {
		name       string
		params     url.Values
		statusCode int
		response   map[string]string
	}{
		{
			name: "Valid Params",
			params: url.Values{
				"rid":    {"123"},
				"factor": {"50"},
			},
			statusCode: http.StatusOK,
			response: map[string]string{
				"status":  "success",
				"message": "Query parameters are valid",
			},
		},
		{
			name:       "Missing rid",
			params:     url.Values{"factor": {"50"}},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "'rid' (rule id) is mandatory",
			},
		},
		{
			name:       "Missing factor",
			params:     url.Values{"rid": {"123"}},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "'Factor' must be a number between 0 and 100",
			},
		},
		{
			name:       "Invalid factor",
			params:     url.Values{"rid": {"123"}, "factor": {"150"}},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "'Factor' must be a number between 0 and 100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/test?%s", tt.params.Encode()), nil)
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

func TestValidateDPO(t *testing.T) {
	app := setupApp()

	tests := []struct {
		name       string
		body       map[string]interface{}
		statusCode int
		response   map[string]string
	}{
		{
			name: "Missing Country",
			body: map[string]interface{}{
				"factor":            50,
				"demand_partner_id": "rubicon",
			},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "Country is mandatory, validation failed",
			},
		},
		{
			name: "Missing Demand Partner",
			body: map[string]interface{}{
				"factor":  50,
				"country": "us",
			},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "DemandPartner is mandatory, validation failed",
			},
		},
		{
			name: "Invalid Country",
			body: map[string]interface{}{
				"country":           "XYZ",
				"factor":            50,
				"demand_partner_id": "rubicon",
			},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "Country code must be 2 characters long and should be in the allowed list",
			},
		},
		{
			name: "Missing Factor",
			body: map[string]interface{}{
				"country":           "us",
				"demand_partner_id": "rubicon",
			},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "Factor is mandatory, validation failed",
			},
		},
		{
			name: "Invalid Factor",
			body: map[string]interface{}{
				"country":           "us",
				"factor":            150,
				"demand_partner_id": "rubicon",
			},
			statusCode: http.StatusBadRequest,
			response: map[string]string{
				"status":  "error",
				"message": "Factor value not allowed, it should be >= 0 and <= 100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/dpo", bytes.NewBuffer(body))
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
