package rest

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestValidateInputs(t *testing.T) {
	tests := []struct {
		name           string
		data           FactorUpdateRequest
		expectedErr    error
		expectedStatus int
	}{
		{
			name: "valid input",
			data: FactorUpdateRequest{
				Country:   "US",
				Publisher: "example",
				Domain:    "example.com",
				Factor:    5.0,
				Device:    "",
			},
			expectedErr:    nil,
			expectedStatus: http.StatusOK,
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
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing domain",
			data: FactorUpdateRequest{
				Country:   "US",
				Publisher: "example",
				Domain:    "",
				Factor:    5.0,
				Device:    "",
			},
			expectedErr:    nil,
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
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid device",
			data: FactorUpdateRequest{
				Country:   "US",
				Publisher: "example",
				Domain:    "example.com",
				Factor:    5.0,
				Device:    "invalid",
			},
			expectedErr:    nil,
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
