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

func TestGlobalFactorBulkPostHandler_InvalidJSON(t *testing.T) {
	endpoint := "/test/global/factor/bulk"

	invalidJSON := `{"key": "consultant_fee", "publisher_id": "id", "value": 5`

	req, err := http.NewRequest(fiber.MethodPost, baseURL+endpoint, strings.NewReader(invalidJSON))
	assert.NoError(t, err)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, `{"status":"error","message":"error parsing request body for global factor bulk update","error":"unexpected end of JSON input"}`, string(body))
}

func TestBulkGlobalFactorHistory(t *testing.T) {
	endpoint := "/bulk/global/factor"
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
			requestBody:        `[{"key":"tech_fee","value":2.1}]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Serving Fees"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Serving Fees",
					Item:         "Tech Fee",
					Changes: []dto.Changes{
						{Property: "value", OldValue: nil, NewValue: 2.1},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `[{"key":"tech_fee","value":2.1}]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Serving Fees"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Serving Fees",
					Item:         "Tech Fee",
					Changes: []dto.Changes{
						{Property: "value", OldValue: nil, NewValue: 2.1},
					},
				},
			},
		},
		{
			name:               "validRequest_Updated",
			requestBody:        `[{"key":"tech_fee","value":2.2}]`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Serving Fees"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Serving Fees",
					Item:         "Tech Fee",
					Changes: []dto.Changes{
						{
							Property: "value",
							OldValue: float64(2.1),
							NewValue: float64(2.2),
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
