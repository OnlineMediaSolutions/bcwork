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

func TestPublisherDomainGetHandler(t *testing.T) {
	endpoint := "/test/publisher/domain/get"

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
			requestBody: `{"filter": {"domain": ["direct.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"publisher_id":"666","publisher_name":"direct_publisher","domain":"direct.com","automation":true,"gpp_target":0.5,"integration_type":[],"created_at":"2024-10-01T13:51:28.407Z","confiant":{},"pixalate":{},"bid_caching":[],"refresh_cache":[],"updated_at":"2024-10-01T13:51:28.407Z","is_direct":true,"is_direct_publisher":true}]`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"filter": {"domain: ["direct.com"]}}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"Request body for publisher domain parsing error","error":"invalid character 'd' after object key"}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"filter": {"domain": ["unknown.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[]`,
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

func TestPublisherDomainHistory(t *testing.T) {
	t.Parallel()

	endpoint := "/publisher/domain"

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
			name:               "validRequest_Created",
			requestBody:        `{"automation":true,"gpp_target":20,"publisher_id":"1111111","domain":"1.com"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Domain"],"publisher_id": ["1111111"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Domain",
					Item:         "1.com (1111111)",
					Changes: []dto.Changes{
						{Property: "automation", OldValue: nil, NewValue: true},
						{Property: "domain", OldValue: nil, NewValue: "1.com"},
						{Property: "gpp_target", OldValue: nil, NewValue: float64(20)},
						{Property: "publisher_id", OldValue: nil, NewValue: "1111111"},
					},
				},
			},
		},
		{
			name:               "noNewChanges",
			requestBody:        `{"automation":true,"gpp_target":20,"publisher_id":"1111111","domain":"1.com"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Domain"],"publisher_id": ["1111111"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Domain",
					Item:         "1.com (1111111)",
					Changes: []dto.Changes{
						{Property: "automation", OldValue: nil, NewValue: true},
						{Property: "domain", OldValue: nil, NewValue: "1.com"},
						{Property: "gpp_target", OldValue: nil, NewValue: float64(20)},
						{Property: "publisher_id", OldValue: nil, NewValue: "1111111"},
					},
				},
			},
		},
		{
			name:               "validRequest_Updated",
			requestBody:        `{"automation":true,"gpp_target":25,"publisher_id":"1111111","domain":"1.com","is_direct":true}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["Domain"],"publisher_id": ["1111111"],"domain": ["1.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Domain",
					Item:         "1.com (1111111)",
					Changes: []dto.Changes{
						{Property: "gpp_target", OldValue: float64(20), NewValue: float64(25)},
						{Property: "is_direct", OldValue: nil, NewValue: true},
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
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "Factor Automation",
					Item:         "2.com (1111111)",
					Changes: []dto.Changes{
						{Property: "automation", OldValue: nil, NewValue: true},
						{Property: "gpp_target", OldValue: nil, NewValue: float64(20)},
					},
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
				history: dto.History{
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "Factor Automation",
					Item:         "2.com (1111111)",
					Changes: []dto.Changes{
						{Property: "gpp_target", OldValue: float64(20), NewValue: float64(25)},
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
