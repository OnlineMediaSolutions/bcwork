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

func TestDPORuleHistory(t *testing.T) {
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
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
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "DPO",
					Item:         "dp_1_333_3.com_de_mobile_android_chrome_leaderboard",
					Changes: []dto.Changes{
						{Property: "browser", OldValue: nil, NewValue: "chrome"},
						{Property: "country", OldValue: nil, NewValue: "de"},
						{Property: "demand_partner_id", OldValue: nil, NewValue: "dp_1"},
						{Property: "device_type", OldValue: nil, NewValue: "mobile"},
						{Property: "domain", OldValue: nil, NewValue: "3.com"},
						{Property: "factor", OldValue: nil, NewValue: float64(4)},
						{Property: "os", OldValue: nil, NewValue: "android"},
						{Property: "placement_type", OldValue: nil, NewValue: "leaderboard"},
						{Property: "publisher", OldValue: nil, NewValue: "333"},
					},
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
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "DPO",
					Item:         "dp_1_333_3.com_de_mobile_android_chrome_leaderboard",
					Changes: []dto.Changes{
						{Property: "browser", OldValue: nil, NewValue: "chrome"},
						{Property: "country", OldValue: nil, NewValue: "de"},
						{Property: "demand_partner_id", OldValue: nil, NewValue: "dp_1"},
						{Property: "device_type", OldValue: nil, NewValue: "mobile"},
						{Property: "domain", OldValue: nil, NewValue: "3.com"},
						{Property: "factor", OldValue: nil, NewValue: float64(4)},
						{Property: "os", OldValue: nil, NewValue: "android"},
						{Property: "placement_type", OldValue: nil, NewValue: "leaderboard"},
						{Property: "publisher", OldValue: nil, NewValue: "333"},
					},
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
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "DPO",
					Item:         "dp_1_333_3.com_de_mobile_android_chrome_leaderboard",
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
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Deleted",
					Subject:      "DPO",
					Item:         "dp_1_333_3.com_de_mobile_android_chrome_leaderboard",
					Changes: []dto.Changes{
						{Property: "browser", NewValue: nil, OldValue: "chrome"},
						{Property: "country", NewValue: nil, OldValue: "de"},
						{Property: "demand_partner_id", NewValue: nil, OldValue: "dp_1"},
						{Property: "device_type", NewValue: nil, OldValue: "mobile"},
						{Property: "domain", NewValue: nil, OldValue: "3.com"},
						{Property: "factor", NewValue: nil, OldValue: float64(5)},
						{Property: "os", NewValue: nil, OldValue: "android"},
						{Property: "placement_type", NewValue: nil, OldValue: "leaderboard"},
						{Property: "publisher", NewValue: nil, OldValue: "333"},
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
