package rest

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBlockGetAllHandler(t *testing.T) {
	endpoint := "/block/get"

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					BlockGetAllHandler,
				},
			},
		},
	})
	defer app.Shutdown()

	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	defer func() {
		pool.Purge(pg)
	}()

	createBlockTables(bcdb.DB())

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

			resp, err := app.Test(req, -1)
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

func createBlockTables(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("CREATE TABLE metadata_queue (transaction_id varchar(36), key varchar(256), version varchar(16),value varchar(512),commited_instances integer, created_at timestamp, updated_at timestamp)")
	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)",
		"f2b8833e-e0e4-57e0-a68b-6792e337ab4d", "badv:20223:realgm.com", nil, "[\"safesysdefender.xyz\"]", 0, "2024-09-20T10:10:10.100", "2024-09-26T10:10:10.100")
	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)",
		"c53c4dd2-6f68-5b62-b613-999a5239ad36", "badv:20356:playpilot.com", nil, "[\"fraction-content.com\"]", 0, "2024-09-20T10:10:10.100", "2024-09-26T10:10:10.100")
	tx.Commit()
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
