package rest

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestValidateInputs(t *testing.T) {
	tests := []struct {
		name           string
		data           FactorUpdateRequest
		expectedErrMsg string
		expectedStatus int
		expectedErr    error
	}{
		{
			name: "valid input",
			data: FactorUpdateRequest{
				Country:   "us",
				Publisher: "example",
				Domain:    "example.com",
				Factor:    5.0,
				Device:    "tablet",
			},
			expectedErr:    nil,
			expectedStatus: http.StatusOK,
			expectedErrMsg: "ok",
		},
		{
			name: "invalid country",
			data: FactorUpdateRequest{
				Country:   "USA",
				Publisher: "example",
				Domain:    "example.com",
				Factor:    5.0,
				Device:    "",
			},
			expectedErr:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: "Country must be a 2-letter country code",
		},
		{
			name: "missing publisher",
			data: FactorUpdateRequest{
				Country:   "US",
				Publisher: "",
				Domain:    "example.com",
				Factor:    5.0,
				Device:    "",
			},
			expectedErr:    nil,
			expectedErrMsg: "Publisher is mandatory",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing country",
			data: FactorUpdateRequest{
				Country:   "US",
				Publisher: "",
				Domain:    "example.com",
				Factor:    5.0,
				Device:    "",
			},
			expectedErr:    nil,
			expectedErrMsg: "Country is mandatory",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing device",
			data: FactorUpdateRequest{
				Country:   "US",
				Publisher: "active",
				Domain:    "example.com",
				Factor:    5.0,
				Device:    "",
			},
			expectedErr:    nil,
			expectedErrMsg: "Device is mandatory",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing factor",
			data: FactorUpdateRequest{
				Country:   "US",
				Publisher: "active",
				Domain:    "example.com",
				Device:    "tablet",
			},
			expectedErr:    nil,
			expectedErrMsg: "Factor is mandatory and must be between 0.1 and 10",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid factor",
			data: FactorUpdateRequest{
				Country:   "US",
				Publisher: "example",
				Domain:    "example.com",
				Factor:    -1.0,
				Device:    "",
			},
			expectedErr:    nil,
			expectedErrMsg: "Factor is mandatory and must be between 0.1 and 10",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid device",
			data: FactorUpdateRequest{
				Country:   "us",
				Publisher: "example",
				Domain:    "example.com",
				Factor:    5.0,
				Device:    "someOtherDevice",
			},
			expectedErr:    nil,
			expectedErrMsg: "Not allowed as device  name",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			req := fasthttp.AcquireRequest()
			req.SetRequestURI("/")
			rc := &fasthttp.RequestCtx{
				Request: *req,
			}
			c := app.AcquireCtx(rc)
			err, isValid := validateInputs(c, &tt.data)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if c.Response().StatusCode() != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, c.Response().StatusCode())
			}
			if isValid != (tt.expectedStatus == http.StatusBadRequest) {
				t.Errorf("expected isValid %v, got %v", tt.expectedStatus == http.StatusBadRequest, isValid)
			}
			fasthttp.ReleaseRequest(req)
		})
	}
}

func TestFactorPostHandler(t *testing.T) {
	app := fiber.New()
	app.Post("/factor", FactorPostHandler)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:           "valid request",
			body:           `{"publisher":"Active Network_Display","device":"all","country":"all","factor":10,"domain":"active.com"}`,
			expectedStatus: http.StatusOK,
			expectedJSON:   `{ "status": "ok","message": "Factor and metadata tables successfully updated"}`,
		},
		{
			name:           "error parsing body",
			body:           ``,
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":"unexpected end of JSON input","message":"Error when parsing factor payload"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/factor", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var response Response
			err = json.NewDecoder(resp.Body).Decode(&response)
			if err != nil {
				t.Fatal(err)
			}
			return

		})
	}
}

func TestFactorGetAllHandler(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedResp string
	}{
		{
			name: "valid request",
			requestBody: `{
				"test": "test"
			}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "empty request body",
			requestBody:  "",
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{status: "error", message: "error when parsing request body for /factor/get"}`,
		},
		{
			name:         "empty request body",
			requestBody:  "{test",
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{status: "error", message: "error when parsing request body for /factor/get"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/factor/get", FactorGetAllHandler)

			req, err := http.NewRequest("POST", "/factor/get", bytes.NewBufferString(tt.requestBody))
			assert.NoError(t, err)

			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			// Check if the error is being returned correctly
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
