package rest

import (
	"encoding/json"
	"io"
	"net/http"
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

func TestPixalateHistory(t *testing.T) {
	endpoint := "/pixalate"
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
			requestBody:        `{"active": true,"rate": 8.2,"publisher_id": "1111111"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Pixalate - Publisher"],"publisher_id": ["1111111"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Pixalate - Publisher",
					Item:         "Pixalate - 1111111",
					Changes: []dto.Changes{
						{Property: "active", OldValue: nil, NewValue: true},
						{Property: "rate", OldValue: nil, NewValue: 8.2},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `{"active": true,"rate": 8.2,"publisher_id": "1111111"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Pixalate - Publisher"],"publisher_id": ["1111111"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Pixalate - Publisher",
					Item:         "Pixalate - 1111111",
					Changes: []dto.Changes{
						{Property: "active", OldValue: nil, NewValue: true},
						{Property: "rate", OldValue: nil, NewValue: 8.2},
					},
				},
			},
		},
		{
			name:               "validRequest_Publisher_Updated",
			requestBody:        `{"active": true,"rate": 8.3,"publisher_id": "1111111"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Pixalate - Publisher"],"publisher_id": ["1111111"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Pixalate - Publisher",
					Item:         "Pixalate - 1111111",
					Changes: []dto.Changes{
						{
							Property: "rate",
							OldValue: float64(8.2),
							NewValue: float64(8.3),
						},
					},
				},
			},
		},
		{
			name:               "validRequest_Domain_Created",
			requestBody:        `{"active": true,"rate": 8.2,"publisher_id": "1111111", "domain":"1.com"}`,
			query:              "?domain=true",
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Pixalate - Domain"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Pixalate - Domain",
					Item:         "Pixalate - 1.com (1111111)",
					Changes: []dto.Changes{
						{Property: "active", OldValue: nil, NewValue: true},
						{Property: "rate", OldValue: nil, NewValue: 8.2},
					},
				},
			},
		},
		{
			name:               "validRequest_Domain_Updated",
			requestBody:        `{"active": true,"rate": 8.3,"publisher_id": "1111111", "domain":"1.com"}`,
			query:              "?domain=true",
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Pixalate - Domain"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Pixalate - Domain",
					Item:         "Pixalate - 1.com (1111111)",
					Changes: []dto.Changes{
						{
							Property: "rate",
							OldValue: float64(8.2),
							NewValue: float64(8.3),
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

			time.Sleep(250 * time.Millisecond)

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
