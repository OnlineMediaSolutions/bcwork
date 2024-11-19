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

func TestDPORuleHistory(t *testing.T) {
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
		hasHistory bool
		history    dto.History
	}

	tests := []struct {
		name               string
		endpoint           string
		requestBody        string
		historyRequestBody string
		want               want
		wantErr            bool
	}{
		{
			name:               "validRequest_Created",
			endpoint:           "/dpo/set",
			requestBody:        `{"demand_partner_name":"dp_1","publisherName":"publisher_3 (333)","domain":"3.com","country":"de","device_type":"mobile","os":"android","browser":"chrome","placement_type":"leaderboard","factor":4,"publisher":"333","demand_partner_id":"dp_1"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["DPO"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "DPO",
					Item:         "de_mobile_android_chrome_leaderboard",
				},
			},
		},
		{
			name:               "noNewChanges",
			endpoint:           "/dpo/set",
			requestBody:        `{"demand_partner_name":"dp_1","publisherName":"publisher_3 (333)","domain":"3.com","country":"de","device_type":"mobile","os":"android","browser":"chrome","placement_type":"leaderboard","factor":4,"publisher":"333","demand_partner_id":"dp_1"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["DPO"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "DPO",
					Item:         "de_mobile_android_chrome_leaderboard",
				},
			},
		},
		{
			name:               "validRequest_Updated",
			endpoint:           "/dpo/set",
			requestBody:        `{"demand_partner_name":"dp_1","publisherName":"publisher_3 (333)","domain":"3.com","country":"de","device_type":"mobile","os":"android","browser":"chrome","placement_type":"leaderboard","factor":5,"publisher":"333","demand_partner_id":"dp_1"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["DPO"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "DPO",
					Item:         "de_mobile_android_chrome_leaderboard",
					Changes: []dto.Changes{
						{
							Property: "factor",
							OldValue: float64(4),
							NewValue: float64(5),
						},
					},
				},
			},
		},
		{
			name:               "validRequest_Updated_SoftDeleted",
			endpoint:           "/dpo/delete",
			requestBody:        `["6986794e-419e-517d-808f-82f79fbaac0b"]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["DPO"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "DPO",
					Item:         "de_mobile_android_chrome_leaderboard",
					Changes: []dto.Changes{
						{
							Property: "active",
							OldValue: true,
							NewValue: false,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(fiber.MethodPost, baseURL+tt.endpoint, strings.NewReader(tt.requestBody))
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
					assert.Equal(t, tt.want.history, got[i])
					found = true
				}
			}
			assert.Equal(t, true, found)
		})
	}
}
