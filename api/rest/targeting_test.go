package rest

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestTargetingGetHandler(t *testing.T) {
	endpoint := "/targeting/get"

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
				response:   `[{"id":10,"publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"os":[],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"daily_cap":null,"status":"Active"}]`,
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
	endpoint := "/targeting/set"

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
	endpoint := "/targeting/update"

	type want struct {
		statusCode int
		response   string
	}

	tests := []struct {
		name        string
		endpoint    string
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "validRequest",
			endpoint:    endpoint,
			requestBody: `{"id":10, "publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"Active"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"targeting successfully updated"}`,
			},
		},
		{
			name:        "invalidRequest",
			endpoint:    endpoint,
			requestBody: `{"publisher_id: "22222222"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"message":"Invalid request body for Targeting. Please ensure it's a valid JSON.","status":"error"}`,
			},
		},
		{
			name:        "noTargetingFoundToUpdate",
			endpoint:    endpoint,
			requestBody: `{"id":12, "publisher_id":"33333333","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"Active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update targeting","error":"failed to get targeting with id [12] to update: sql: no rows in result set"}`,
			},
		},
		{
			// based on results of "validRequest"
			name:        "nothingToUpdate",
			endpoint:    endpoint,
			requestBody: `{"id":10, "publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"Active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update targeting","error":"there are no new values to update targeting"}`,
			},
		},
		{
			name:        "duplicateConflictOnUpdatedEntity",
			endpoint:    endpoint,
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
			req := httptest.NewRequest(fiber.MethodPost, tt.endpoint, strings.NewReader(tt.requestBody))
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
	endpoint := "/targeting/tags"

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
