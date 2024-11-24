package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestBlockGetAllHandler(t *testing.T) {
	endpoint := "/test/block/get"

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
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse request for getting blocks","error":"invalid character 'b' after object key"}`,
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

func TestBlockHistory(t *testing.T) {
	endpoint := "/block"
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
		history    dto.History
	}

	tests := []struct {
		name               string
		requestBody        string
		query              string
		historyRequestBody string
		want               want
		wantErr            bool
	}{
		{
			name:               "validRequest_Publisher_Created",
			requestBody:        `{"badv": [],"bcat": ["IAB1-1"],"publisher": "1111111"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Blocks - Publisher"],"publisher_id": ["1111111"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Blocks - Publisher",
					Item:         "Blocks - 1111111",
					Changes: []dto.Changes{
						{Property: "badv", OldValue: nil, NewValue: []interface{}{}},
						{Property: "bcat", OldValue: nil, NewValue: []interface{}{"IAB1-1"}},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `{"badv": [],"bcat": ["IAB1-1"],"publisher": "1111111"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Blocks - Publisher"],"publisher_id": ["1111111"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Blocks - Publisher",
					Item:         "Blocks - 1111111",
					Changes: []dto.Changes{
						{Property: "badv", OldValue: nil, NewValue: []interface{}{}},
						{Property: "bcat", OldValue: nil, NewValue: []interface{}{"IAB1-1"}},
					},
				},
			},
		},
		{
			name:               "validRequest_Publisher_Updated",
			requestBody:        `{"badv": [],"bcat": [],"publisher": "1111111"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Blocks - Publisher"],"publisher_id": ["1111111"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Blocks - Publisher",
					Item:         "Blocks - 1111111",
					Changes: []dto.Changes{
						{
							Property: "bcat",
							OldValue: []any{"IAB1-1"},
							NewValue: []any{},
						},
					},
				},
			},
		},
		{
			name:               "validRequest_Domain_Created",
			requestBody:        `{"badv": [],"bcat": [],"publisher": "1111111", "domain":"1.com"}`,
			query:              "?domain=true",
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Blocks - Domain"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Blocks - Domain",
					Item:         "Blocks - 1.com (1111111)",
					Changes: []dto.Changes{
						{Property: "badv", OldValue: nil, NewValue: []interface{}{}},
						{Property: "bcat", OldValue: nil, NewValue: []interface{}{}},
					},
				},
			},
		},
		{
			name:               "validRequest_Domain_Updated",
			requestBody:        `{"badv": [],"bcat": ["IAB1-1"],"publisher": "1111111", "domain":"1.com"}`,
			query:              "?domain=true",
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Blocks - Domain"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Blocks - Domain",
					Item:         "Blocks - 1.com (1111111)",
					Changes: []dto.Changes{
						{
							Property: "bcat",
							OldValue: []any{},
							NewValue: []any{"IAB1-1"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(fiber.MethodPost, baseURL+endpoint+tt.query, strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))

			_, err = http.DefaultClient.Do(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			historyReq, err := http.NewRequest(fiber.MethodPost, baseURL+historyEndpoint, strings.NewReader(tt.historyRequestBody))
			if err != nil {
				t.Fatal(err)
			}
			historyReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			historyReq.Header.Set(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))

			historyResp, err := http.DefaultClient.Do(historyReq)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, historyResp.StatusCode)

			body, err := io.ReadAll(historyResp.Body)
			assert.NoError(t, err)
			defer historyResp.Body.Close()

			var got []dto.History
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)

			for i := range got {
				got[i].ID = 0
				got[i].Date = time.Time{}
				for j := range got[i].Changes {
					got[i].Changes[j].ID = ""
				}
			}

			assert.Contains(t, got, tt.want.history)
		})
	}
}
