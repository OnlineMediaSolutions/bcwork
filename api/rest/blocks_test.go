package rest

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestBlockGetAllHandler(t *testing.T) {
	endpoint := "/block/get"

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
			requestBody: `{"types": ["badv"], "publisher": "20356", "domain": "playpilot.com"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response: `[` +
					`{` +
					`"transaction_id":"c53c4dd2-6f68-5b62-b613-999a5239ad36",` +
					`"key":"badv:20356:playpilot.com",` +
					`"version":null,` +
					`"value":["fraction-content.com"],` +
					`"commited_instances":0,` +
					`"created_at":"2024-09-20T10:10:10.1Z",` +
					`"updated_at":"2024-09-26T10:10:10.1Z"` +
					`}` +
					`]`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"types: ["badv"], "publisher": "20356", "domain": "playpilot.com"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"Failed to parse metadata update payload"}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"types": ["badv"], "publisher": "20357", "domain": "playpilot.com"}`,
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

func TestCreateKeyForQuery(t *testing.T) {
	testCases := []struct {
		name     string
		request  BlockGetRequest
		expected string
	}{
		{
			name:     "Empty request",
			request:  BlockGetRequest{},
			expected: " and 1=1 ",
		},
		{
			name: "Publisher and types provided",
			request: BlockGetRequest{
				Types:     []string{"badv", "bcat"},
				Publisher: "publisher",
				Domain:    "domain",
			},
			expected: "AND ( (metadata_queue.key = 'badv:publisher:domain') OR (metadata_queue.key = 'bcat:publisher:domain'))",
		},
		{
			name: "Publisher provided, no domain",
			request: BlockGetRequest{
				Types:     []string{"badv", "cbat"},
				Publisher: "publisher",
			},
			expected: "AND ( (metadata_queue.key = 'badv:publisher') OR (metadata_queue.key = 'cbat:publisher'))",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := createKeyForQuery(&tc.request)
			if result != tc.expected {
				t.Errorf("Test %s failed: expected '%s', got '%s'", tc.name, tc.expected, result)
			}
		})
	}
}
