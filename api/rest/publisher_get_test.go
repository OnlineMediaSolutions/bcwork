package rest

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPublisherGetHandler(t *testing.T) {
	endpoint := "/test/publisher/get"

	type want struct {
		statusCode int
		response   string
	}

	tests := []struct {
		name        string
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "validRequest",
			requestBody: `{"filter": {"publisher_id": ["555"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"publisher_id":"555","created_at":"2024-10-01T13:46:41.302Z","name":"test_publisher","account_manager_id":"1","account_manager_full_name":"name_1 surname_1","media_buyer_id":"2","media_buyer_full_name":"name_2 surname_2","campaign_manager_id":"3","campaign_manager_full_name":"name_temp surname_temp","office_location":"IL","integration_type":[],"status":"Active","confiant":{},"pixalate":{},"bid_caching":[],"refresh_cache":[]}]`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"filter": {"publisher_id: ["555"]}}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"error when parsing request body"}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"filter": {"publisher_id": ["xxxxxxxx"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[]`,
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
