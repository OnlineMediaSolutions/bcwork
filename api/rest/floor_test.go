package rest

import (
	"github.com/m6yf/bcwork/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/validations"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateFloors(t *testing.T) {
	app := fiber.New()
	app.Post("/floors", validations.ValidateFloors, func(c *fiber.Ctx) error {
		return c.SendString("Floor created successfully")
	})

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
			name:     "Missing device",
			body:     `{"publisher": "test", "country": "US", "floor": 1.0, "domain": "example.com"}`,
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
		req := httptest.NewRequest("POST", "/floors", strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
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
	app := fiber.New()
	app.Post("/floor", FloorPostHandler)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedJSON   string
	}{

		{
			name:           "error parsing body",
			body:           ``,
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":"invalid request","message":"Invalid JSON payload"}`,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("POST", "/floor", bytes.NewBufferString(tt.body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Errorf("Test %s: %v", tt.name, err)
			continue
		}

		if resp.StatusCode != tt.expectedStatus {
			t.Errorf("Test %s: expected status code %d, got %d", tt.name, tt.expectedStatus, resp.StatusCode)
			continue
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Test %s: %v", tt.name, err)
			continue
		}

		if string(respBody) != tt.expectedJSON {
			t.Errorf("Test %s: expected JSON response %s, got %s", tt.name, tt.expectedJSON, string(respBody))
		}
	}
}

func TestFloorGetAllHandler(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedResp string
	}{

		{
			name:         "empty request body",
			requestBody:  "",
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{status: "error", message: "error when parsing request body for /floor/get"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/floor/get", FloorGetAllHandler)

			req, err := http.NewRequest("POST", "/floor/get", bytes.NewBufferString(tt.requestBody))
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.expectedCode == http.StatusBadRequest {
				responseBody, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var responseBodyMap map[string]string
				err = json.Unmarshal(responseBody, &responseBodyMap)
				assert.NoError(t, err)
				assert.Equal(t, "error", responseBodyMap["Status"])
				assert.Equal(t, "invalid request body", responseBodyMap["Message"])
			}
		})
	}
}

func TestConvertingAllValues(t *testing.T) {
	tests := []struct {
		name     string
		data     FloorUpdateRequest
		expected FloorUpdateRequest
	}{
		{
			name: "device and country with all value",
			data: FloorUpdateRequest{
				Device:    "all",
				Country:   "all",
				Publisher: "345",
				Domain:    "bubu.com",
			},
			expected: FloorUpdateRequest{
				Device:    "",
				Country:   "",
				Publisher: "345",
				Domain:    "bubu.com",
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
