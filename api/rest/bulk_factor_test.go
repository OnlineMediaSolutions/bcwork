package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/gofiber/fiber/v2"
)

func TestBulkFactorForAutomation(t *testing.T) {
	endpoint := "/test/bulk/factor"

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
			requestBody: `[{"publisher":"publisher1","domain":"domain1","device":"desktop","factor":1.23,"country":"us"},{"publisher":"publisher2","domain":"domain2","device":"mobile","factor":3,"country":"il"},{"publisher":"publisher2","domain":"domain1","device":"mobile","factor":3,"country":"il"},{"publisher":"publisher1","domain":"domain1","device":"mobile","factor":10,"country":"uk"}]`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"factor bulk update successfully processed"}`,
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

func TestBulkFactorHistory(t *testing.T) {
	t.Parallel()

	endpoint := "/bulk/factor"
	historyEndpoint := "/history/get"

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
			requestBody:        `[{"publisher":"publisher_4","domain":"4.com","country":"il","device":"mobile","factor":0.1}]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Bidder Targeting"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Bidder Targeting",
					Item:         "publisher_4_4.com_il_mobile_all_all_all",
					Changes: []dto.Changes{
						{Property: "factor", OldValue: nil, NewValue: float64(0.1)},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `[{"publisher":"publisher_4","domain":"4.com","country":"il","device":"mobile","factor":0.1}]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Bidder Targeting"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Bidder Targeting",
					Item:         "publisher_4_4.com_il_mobile_all_all_all",
					Changes: []dto.Changes{
						{Property: "factor", OldValue: nil, NewValue: float64(0.1)},
					},
				},
			},
		},
		{
			name:               "validRequest_Updated",
			requestBody:        `[{"publisher":"publisher_4","domain":"4.com","country":"il","device":"mobile","factor":0.15}]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Bidder Targeting"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Bidder Targeting",
					Item:         "publisher_4_4.com_il_mobile_all_all_all",
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
