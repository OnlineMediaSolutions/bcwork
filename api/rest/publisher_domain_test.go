package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
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

func TestPublisherDomainHistory(t *testing.T) {
	endpoint := "/publisher/domain"
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
		hasHistory bool
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
			name:               "validRequest_Created",
			requestBody:        `{"automation":true,"gpp_target":20,"publisher_id":"1111111","domain":"1.com"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Domain"],"publisher_id": ["1111111"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Domain",
					Item:         "1.com (1111111)",
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `{"automation":true,"gpp_target":20,"publisher_id":"1111111","domain":"1.com"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Domain"],"publisher_id": ["1111111"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Domain",
					Item:         "1.com (1111111)",
				},
			},
		},
		{
			name:               "validRequest_Updated",
			requestBody:        `{"automation":true,"gpp_target":25,"publisher_id":"1111111","domain":"1.com"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Domain"],"publisher_id": ["1111111"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Domain",
					Item:         "1.com (1111111)",
					Changes: []dto.Changes{
						{
							Property: "gpp_target",
							OldValue: float64(20),
							NewValue: float64(25),
						},
					},
				},
			},
		},
		{
			name:               "validRequest_Automation_Created",
			requestBody:        `{"automation":true,"gpp_target":20,"publisher_id":"1111111","domain":"2.com"}`,
			query:              "?automation=true",
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Factor Automation"],"publisher_id": ["1111111"],"domain": ["2.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Factor Automation",
					Item:         "2.com (1111111)",
				},
			},
		},
		{
			name:               "validRequest_Automation_Updated",
			requestBody:        `{"automation":true,"gpp_target":25,"publisher_id":"1111111","domain":"2.com"}`,
			query:              "?automation=true",
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Factor Automation"],"publisher_id": ["1111111"],"domain": ["2.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Factor Automation",
					Item:         "2.com (1111111)",
					Changes: []dto.Changes{
						{
							Property: "gpp_target",
							OldValue: float64(20),
							NewValue: float64(25),
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

			var (
				got   []dto.History
				found bool
			)
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)
			for i := range got {
				got[i].ID = 0
				got[i].Date = time.Time{}
				for j := range got[i].Changes {
					got[i].Changes[j].ID = ""
				}

				if reflect.DeepEqual(tt.want.history, got[i]) {
					found = true
				}
			}
			assert.Equal(t, true, found)
		})
	}
}
