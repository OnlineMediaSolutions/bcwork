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

func TestPublisherUpdateHandler(t *testing.T) {
	endpoint := "/test/publisher/update"

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
			requestBody: `{"publisher_id":"222","updates":{"publisher_id":"222","name":"publisher_for_test","status":"Active","office_location":"IL","integration_type":["oRTB"],"media_type":["Video"],"is_direct":true}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"updated"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{publisher_id":"222","updates":{"publisher_id":"222","name":"publisher_for_test","status":"Active","office_location":"IL","integration_type":["oRTB"]}}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"error when parsing request body","error":"invalid character 'p' looking for beginning of object key string"}`,
			},
		},
		{
			name:        "noPublisherFoundToUpdate",
			requestBody: `{"publisher_id":"9999999","updates":{"publisher_id":"9999999","name":"publisher_for_test","status":"Active","office_location":"IL","integration_type":["oRTB"]}}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update publisher fields","error":"failed to get publisher with id [9999999] to update: sql: no rows in result set"}`,
			},
		},
		{
			name:        "nothingToUpdate",
			requestBody: `{"publisher_id":"222","updates":{}}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update publisher fields","error":"applicaiton payload contains no vals for update (publisher_id:222)"}`,
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

func TestPublisherUpdateHistory(t *testing.T) {
	t.Parallel()

	endpoint := "/publisher/update"

	type want struct {
		statusCode int
		history    dto.History
	}

	tests := []struct {
		name               string
		requestBody        string
		historyRequestBody string
		want               want
		wantErr            bool
	}{
		{
			name:               "validRequest_Updated",
			requestBody:        `{"publisher_id":"333","updates":{"publisher_id":"333","name":"publisher_3","status":"Active","office_location":"IL","integration_type":["JS Tags (NP)"],"is_direct":true}}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Publisher"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Publisher",
					Item:         "333",
					Changes: []dto.Changes{
						{Property: "integration_type", OldValue: nil, NewValue: []any{"JS Tags (NP)"}},
						{Property: "is_direct", OldValue: false, NewValue: true},
						{Property: "office_location", OldValue: "LATAM", NewValue: "IL"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(fiber.MethodPost, baseURL+endpoint, strings.NewReader(tt.requestBody))
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

			time.Sleep(250 * time.Millisecond)

			historyReq, err := http.NewRequest(fiber.MethodPost, baseURL+constant.HistoryEndpoint, strings.NewReader(tt.historyRequestBody))
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
