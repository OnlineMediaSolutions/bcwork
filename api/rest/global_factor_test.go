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

func TestGlobalFactorHistory(t *testing.T) {
	t.Parallel()

	endpoint := "/global/factor"

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
			requestBody:        `{"key":"tam_fee","value":0.249}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Serving Fees"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Serving Fees",
					Item:         "Amazon TAM Fee",
					Changes: []dto.Changes{
						{Property: "value", OldValue: nil, NewValue: float64(0.249)},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `{"key":"tam_fee","value":0.249}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Serving Fees"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Serving Fees",
					Item:         "Amazon TAM Fee",
					Changes: []dto.Changes{
						{Property: "value", OldValue: nil, NewValue: float64(0.249)},
					},
				},
			},
		},
		{
			name:               "validRequest_Updated",
			requestBody:        `{"key":"tam_fee","value":0.25}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Serving Fees"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Serving Fees",
					Item:         "Amazon TAM Fee",
					Changes: []dto.Changes{
						{
							Property: "value",
							OldValue: float64(0.249),
							NewValue: float64(0.25),
						},
					},
				},
			},
		},
		{
			name:               "validRequest_Created_ConsultantFee",
			requestBody:        `{"key":"consultant_fee","publisher_id":"1111111","value":0.249}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Serving Fees"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Serving Fees",
					Item:         "Consultant Fee - 1111111",
					Changes: []dto.Changes{
						{Property: "value", OldValue: nil, NewValue: float64(0.249)},
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
