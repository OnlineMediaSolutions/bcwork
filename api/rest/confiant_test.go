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

func TestConfiantHistory(t *testing.T) {
	t.Parallel()

	endpoint := "/confiant"

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
			requestBody:        `{"confiant_key": "test-confiant","rate": 96,"publisher_id": "1111111"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Confiant - Publisher"],"publisher_id": ["1111111"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Confiant - Publisher",
					Item:         "Confiant - 1111111",
					Changes: []dto.Changes{
						{Property: "confiant_key", OldValue: nil, NewValue: "test-confiant"},
						{Property: "rate", OldValue: nil, NewValue: float64(96)},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `{"confiant_key": "test-confiant","rate": 96,"publisher_id": "1111111"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Confiant - Publisher"],"publisher_id": ["1111111"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Confiant - Publisher",
					Item:         "Confiant - 1111111",
					Changes: []dto.Changes{
						{Property: "confiant_key", OldValue: nil, NewValue: "test-confiant"},
						{Property: "rate", OldValue: nil, NewValue: float64(96)},
					},
				},
			},
		},
		{
			name:               "validRequest_Publisher_Updated",
			requestBody:        `{"confiant_key": "test-confiant","rate": 97,"publisher_id": "1111111"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Confiant - Publisher"],"publisher_id": ["1111111"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Confiant - Publisher",
					Item:         "Confiant - 1111111",
					Changes: []dto.Changes{
						{
							Property: "rate",
							OldValue: float64(96),
							NewValue: float64(97),
						},
					},
				},
			},
		},
		{
			name:               "validRequest_Domain_Created",
			requestBody:        `{"confiant_key": "test-confiant","rate": 96,"publisher_id": "1111111", "domain":"1.com"}`,
			query:              "?domain=true",
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Confiant - Domain"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Confiant - Domain",
					Item:         "Confiant - 1.com (1111111)",
					Changes: []dto.Changes{
						{Property: "confiant_key", OldValue: nil, NewValue: "test-confiant"},
						{Property: "rate", OldValue: nil, NewValue: float64(96)},
					},
				},
			},
		},
		{
			name:               "validRequest_Domain_Updated",
			requestBody:        `{"confiant_key": "test-confiant","rate": 97,"publisher_id": "1111111", "domain":"1.com"}`,
			query:              "?domain=true",
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Confiant - Domain"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Confiant - Domain",
					Item:         "Confiant - 1.com (1111111)",
					Changes: []dto.Changes{
						{
							Property: "rate",
							OldValue: float64(96),
							NewValue: float64(97),
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
