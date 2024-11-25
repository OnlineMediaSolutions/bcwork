package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFactorAdjusterHandler(t *testing.T) {
	endpoint := "/test/adjust/factor"

	type want struct {
		statusCode int
		response   string
	}

	tests := []struct {
		name        string
		requestBody string
		query       string
		want        want
		wantErr     bool
	}{
		{
			name:        "validRequest_AdjustFactorRequest",
			requestBody: `{"domain":["oms.com"],"value":5}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"Adjusted Factor values"}`,
			},
		},
		{
			name:        "domainDoesntExistRequest",
			requestBody: `{"domain":["notExistInDB.com"],"value":5}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"Failed to fetch Factors","error":"%!s(\u003cnil\u003e)"}`,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodPost, endpoint, strings.NewReader(tt.requestBody))
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			resp, err := appTest.Test(req, -1)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}
