package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestTargetingGetHandler(t *testing.T) {
	endpoint := "/test/targeting/get"

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
			requestBody: `{"filter": {"publisher_id": ["22222222"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"id":10,"publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"os":[],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"daily_cap":null,"status":"Active"},{"id":30,"publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["al"],"device_type":["mobile"],"browser":["firefox"],"os":[],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"daily_cap":null,"status":"Active"}]`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"filter": {"publisher_id: ["22222222"]}}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse request for getting targeting data","error":"invalid character '2' after object key"}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"filter": {"publisher_id": ["xxxxxxxx"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[]`,
			},
		},
		{
			name:        "validRequest_withDailyCap",
			requestBody: `{"filter": {"publisher_id": ["333"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"id":20,"publisher_id":"333","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["ru","us"],"device_type":["mobile"],"browser":["firefox"],"os":[],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"","value":0,"daily_cap":1000,"status":"Active"}]`,
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

func TestTargetingSetHandler(t *testing.T) {
	endpoint := "/test/targeting/set"

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
			requestBody: `{"publisher_id":"22222222","domain":"3.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"status":"Active"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"targeting successfully added"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"publisher_id: "22222222"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"message":"Invalid request body for Targeting. Please ensure it's a valid JSON.","status":"error"}`,
			},
		},
		{
			name:        "hasDuplicate",
			requestBody: `{"publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","ru"],"device_type":["mobile","desktop"],"browser":["firefox","chrome"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"status":"Active"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"found duplicate while creating targeting","error":"checking for duplicates: found duplicate: there is targeting with such parameters","duplicate":{"id":10,"publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"os":[],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"daily_cap":null,"status":"Active"}}`,
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

func TestTargetingUpdateHandler(t *testing.T) {
	endpoint := "/test/targeting/update"

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
			requestBody: `{"id":10, "publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"Active"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"targeting successfully updated"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"publisher_id: "22222222"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"message":"Invalid request body for Targeting. Please ensure it's a valid JSON.","status":"error"}`,
			},
		},
		{
			name:        "noTargetingFoundToUpdate",
			requestBody: `{"id":12, "publisher_id":"33333333","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"Active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update targeting","error":"failed to get targeting with id [12] to update: sql: no rows in result set"}`,
			},
		},
		{
			// based on results of "validRequest"
			name:        "nothingToUpdate",
			requestBody: `{"id":10, "publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"Active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update targeting","error":"there are no new values to update targeting"}`,
			},
		},
		{
			name:        "duplicateConflictOnUpdatedEntity",
			requestBody: `{"id":11, "publisher_id":"1111111","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"Active"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"found duplicate while updating targeting","error":"checking for duplicates: found duplicate: there is targeting with such parameters","duplicate":{"id":9,"publisher_id":"1111111","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["ru","us"],"device_type":["mobile"],"browser":["firefox"],"os":[],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"","value":0,"daily_cap":null,"status":"Active"}}`,
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

func TestTargetingExportTagsHandler(t *testing.T) {
	endpoint := "/test/targeting/tags"

	now := time.Now().Format(time.DateOnly)

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
			requestBody: `{"ids": [9, 10]}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   "{\"status\":\"success\",\"message\":\"tags successfully exported\",\"tags\":[{\"id\":9,\"tag\":\"\\u003c!-- HTML Tag for publisher='publisher_1', domain='2.com', size='300X250', key_1='value_1', key_2='value_2', key_3='value_3', exported='" + now + "' --\\u003e\\n\\u003cscript src=\\\"https://rt.marphezis.com/js?pid=1111111\\u0026size=300X250\\u0026dom=2.com\\u0026key_1=value_1\\u0026key_2=value_2\\u0026key_3=value_3\\\"\\u003e\\u003c/script\\u003e\"},{\"id\":10,\"tag\":\"\\u003c!-- HTML Tag for publisher='publisher_2', domain='2.com', size='300X250', key_1='value_1', key_2='value_2', key_3='value_3', exported='" + now + "' --\\u003e\\n\\u003cscript src=\\\"https://rt.marphezis.com/js?pid=22222222\\u0026size=300X250\\u0026dom=2.com\\u0026key_1=value_1\\u0026key_2=value_2\\u0026key_3=value_3\\\"\\u003e\\u003c/script\\u003e\"}]}",
			},
		},
		{
			name:        "validRequest_withGDPR",
			requestBody: `{"ids": [9, 10], "add_gdpr": true}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   "{\"status\":\"success\",\"message\":\"tags successfully exported\",\"tags\":[{\"id\":9,\"tag\":\"\\u003c!-- HTML Tag for publisher='publisher_1', domain='2.com', size='300X250', key_1='value_1', key_2='value_2', key_3='value_3', exported='" + now + "' --\\u003e\\n\\u003cscript src=\\\"https://rt.marphezis.com/js?pid=1111111\\u0026size=300X250\\u0026dom=2.com\\u0026key_1=value_1\\u0026key_2=value_2\\u0026key_3=value_3\\u0026gdpr=${GDPR}\\u0026gdpr_concent=${GDPR_CONSENT_883}\\\"\\u003e\\u003c/script\\u003e\"},{\"id\":10,\"tag\":\"\\u003c!-- HTML Tag for publisher='publisher_2', domain='2.com', size='300X250', key_1='value_1', key_2='value_2', key_3='value_3', exported='" + now + "' --\\u003e\\n\\u003cscript src=\\\"https://rt.marphezis.com/js?pid=22222222\\u0026size=300X250\\u0026dom=2.com\\u0026key_1=value_1\\u0026key_2=value_2\\u0026key_3=value_3\\u0026gdpr=${GDPR}\\u0026gdpr_concent=${GDPR_CONSENT_883}\\\"\\u003e\\u003c/script\\u003e\"}]}",
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"ids: [9, 10]}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse request for export tags","error":"unexpected end of JSON input"}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"ids": [100, 101]}`,
			want: want{
				statusCode: fiber.StatusNotFound,
				response:   `{"status":"error","message":"failed to export tags","error":"no tags found for ids [100 101]"}`,
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

func TestTargetingUpdate_History(t *testing.T) {
	endpoint := "/targeting/update"
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
		hasHistory bool
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
			name:               "noChanges",
			requestBody:        `{"id":30, "publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["al"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"Active"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["JS Targeting"],"publisher_id":["22222222"],"domain":["2.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: false,
			},
		},
		{
			name:               "validRequest",
			requestBody:        `{"id":30, "publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["al"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":3,"status":"Active"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["JS Targeting"],"publisher_id":["22222222"],"domain":["2.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "JS Targeting",
					Item:         "al_mobile__firefox_top",
					Changes: []dto.Changes{
						{
							Property: "value",
							OldValue: float64(2),
							NewValue: float64(3),
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
			if !tt.want.hasHistory {
				assert.Equal(t, []dto.History{}, got)
				return
			}

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

func TestTargetingSet_History(t *testing.T) {
	endpoint := "/targeting/set"
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
		hasHistory bool
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
			name:               "validRequest",
			requestBody:        `{"publisher_id":"22222222","domain":"3.com","unit_size":"300X250","placement_type":"top","country":["by"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"status":"Active"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["JS Targeting"],"publisher_id":["22222222"],"domain":["3.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "JS Targeting",
					Item:         "by_mobile__firefox_top",
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
			if !tt.want.hasHistory {
				assert.Equal(t, []dto.History{}, got)
				return
			}

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
