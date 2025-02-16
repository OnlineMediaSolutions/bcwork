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

func TestBulkDPORule(t *testing.T) {
	endpoint := "/test/bulk/dpo"

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
			requestBody: `[{"demand_partner_id":"dp_1","publisher":"publisher_1","domain":"1.com","country":"il","browser":"firefox","os":"linux","device_type":"mobile","placement_type":"top","factor":0.1},{"demand_partner_id":"dp_2","publisher":"publisher_2","domain":"2.com","country":"us","browser":"chrome","os":"macos","device_type":"mobile","placement_type":"bottom","factor":0.05},{"demand_partner_id":"dp_3","publisher":"publisher_3","domain":"3.com","country":"ru","browser":"opera","os":"windows","device_type":"mobile","placement_type":"side","factor":0.15}]`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"Dpo_rule with Metadata_queue bulk update successfully processed"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `["demand_partner_id":"dp_1","publisher":"publisher_1","domain":"1.com","country":"il","browser":"firefox","os":"linux","device_type":"mobile","placement_type":"top","factor":0.1},{"demand_partner_id":"dp_2","publisher":"publisher_2","domain":"2.com","country":"us","browser":"chrome","os":"macos","device_type":"mobile","placement_type":"bottom","factor":0.05},{"demand_partner_id":"dp_3","publisher":"publisher_3","domain":"3.com","country":"ru","browser":"opera","os":"windows","device_type":"mobile","placement_type":"side","factor":0.15}]`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"Failed to parse metadata for DPO bulk","error":"invalid character ':' after array element"}`,
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

			resp, err := http.DefaultClient.Do(req)
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

func TestBulkDPORuleHistory(t *testing.T) {
	t.Parallel()

	endpoint := "/bulk/dpo"

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
			name:               "validRequest_Created",
			requestBody:        `[{"demand_partner_id":"dp_4","publisher":"publisher_4","domain":"4.com","country":"il","browser":"firefox","os":"linux","device_type":"mobile","placement_type":"top","factor":0.1}]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["DPO"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "DPO",
					Item:         "dp_4_publisher_4_4.com_il_mobile_linux_firefox_top",
					Changes: []dto.Changes{
						{Property: "factor", OldValue: nil, NewValue: 0.1},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `[{"demand_partner_id":"dp_4","publisher":"publisher_4","domain":"4.com","country":"il","browser":"firefox","os":"linux","device_type":"mobile","placement_type":"top","factor":0.1}]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["DPO"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "DPO",
					Item:         "dp_4_publisher_4_4.com_il_mobile_linux_firefox_top",
					Changes: []dto.Changes{
						{Property: "factor", OldValue: nil, NewValue: 0.1},
					},
				},
			},
		},
		{
			name:               "validRequest_Updated",
			requestBody:        `[{"demand_partner_id":"dp_4","publisher":"publisher_4","domain":"4.com","country":"il","browser":"firefox","os":"linux","device_type":"mobile","placement_type":"top","factor":0.15}]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["DPO"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "DPO",
					Item:         "dp_4_publisher_4_4.com_il_mobile_linux_firefox_top",
					Changes: []dto.Changes{
						{
							Property: "factor",
							OldValue: float64(0.1),
							NewValue: float64(0.15),
						},
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
