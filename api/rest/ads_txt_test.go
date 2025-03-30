package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/stretchr/testify/assert"
)

func TestAdsTxtTables(t *testing.T) {
	type want struct {
		statusCode int
		response   []*dto.AdsTxt
	}

	tests := []struct {
		name        string
		endpoint    string
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "cmTable_validRequest",
			endpoint:    "/test/ads_txt/cm",
			requestBody: `{}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   getAdsTxtData(t, "./testdata/ads_txt_cm_table.json"),
			},
		},
		{
			name:        "mbTable_validRequest",
			endpoint:    "/test/ads_txt/mb",
			requestBody: `{}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   getAdsTxtData(t, "./testdata/ads_txt_mb_table.json"),
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

			var got []*dto.AdsTxt
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.response, got)
		})
	}
}

func TestAdsTxtMainTable(t *testing.T) {
	type want struct {
		statusCode int
		response   *dto.AdsTxtResponse
	}

	tests := []struct {
		name        string
		endpoint    string
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "mainTable_validRequest",
			endpoint:    "/test/ads_txt/main",
			requestBody: `{}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   getAdsTxtResponse(t, "./testdata/ads_txt_main_table.json"),
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

			var got *dto.AdsTxtResponse
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.response, got)
		})
	}
}

func TestAdsTxtGroupByDPTable(t *testing.T) {
	endpoint := "/test/ads_txt/group_by_dp"

	type want struct {
		statusCode int
		response   *dto.AdsTxtGroupByDPResponse
	}

	tests := []struct {
		name        string
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "validRequest",
			requestBody: `{}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   getAdsTxtGroupByDPResponse(t, "./testdata/ads_txt_group_by_dp_table.json"),
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

			var got *dto.AdsTxtGroupByDPResponse
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.response, got)
		})
	}
}

func TestAdsTxtUpdateHandler(t *testing.T) {
	type want struct {
		statusCode int
		response   *dto.AdsTxtResponse
	}

	tests := []struct {
		name           string
		endpoint       string
		requestBody    string
		getEndpoint    string
		getRequestBody string
		want           want
		wantErr        bool
	}{
		{
			name:           "valid",
			endpoint:       "/test/ads_txt/update",
			requestBody:    `{"domain":["test2.com"],"demand_partner_id":"index","demand_status":"approved"}`,
			getEndpoint:    "/test/ads_txt/main",
			getRequestBody: `{"filter":{"publisher_id":["1111111"],"domain":["test2.com"],"demand_partner_name_extended":["Index - Index"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response: func() *dto.AdsTxtResponse {
					row := []byte("{\"data\":[{\"id\":15,\"group_by_dp_id\":0,\"cursor_id\":1,\"publisher_id\":\"1111111\",\"publisher_name\":\"publisher_1\",\"mirror_publisher_id\":\"\",\"account_manager_id\":null,\"account_manager_full_name\":null,\"campaign_manager_id\":null,\"campaign_manager_full_name\":null,\"domain\":\"test2.com\",\"domain_status\":\"New\",\"demand_partner_id\":\"index\",\"demand_partner_name\":\"Index\",\"demand_partner_name_extended\":\"Index - Index\",\"demand_partner_connection_id\":1,\"media_type\":[\"Web Banners\"],\"demand_manager_id\":\"1\",\"demand_manager_full_name\":\"name_1 surname_1\",\"demand_status\":\"Approved\",\"is_demand_partner_active\":false,\"seat_owner_name\":\"\",\"score\":0,\"action\":\"\",\"status\":\"Not Scanned\",\"is_required\":true,\"ads_txt_line\":\"indexexchange.com, 181818, RESELLER\",\"added\":0,\"total\":0,\"dp_enabled\":false,\"last_scanned_at\":null,\"error_message\":null,\"is_mirror_used\":false}],\"total\":1}")
					var want *dto.AdsTxtResponse
					err := json.Unmarshal(row, &want)
					if err != nil {
						t.Fatal(err)
					}

					return want
				}(),
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

			resp, err := http.DefaultClient.Do(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			getReq, err := http.NewRequest(fiber.MethodPost, baseURL+tt.getEndpoint, strings.NewReader(tt.getRequestBody))
			if err != nil {
				t.Fatal(err)
			}
			getReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			getResp, err := http.DefaultClient.Do(getReq)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			body, err := io.ReadAll(getResp.Body)
			assert.NoError(t, err)
			defer getResp.Body.Close()

			var got *dto.AdsTxtResponse
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.response, got)
		})
	}
}

func getAdsTxtData(t *testing.T, datapath string) []*dto.AdsTxt {
	data, err := getFileData(datapath)
	if err != nil {
		t.Fatal(err)
	}

	var want []*dto.AdsTxt
	err = json.Unmarshal(data, &want)
	if err != nil {
		t.Fatal(err)
	}

	return want
}

func getAdsTxtResponse(t *testing.T, datapath string) *dto.AdsTxtResponse {
	data, err := getFileData(datapath)
	if err != nil {
		t.Fatal(err)
	}

	var want *dto.AdsTxtResponse
	err = json.Unmarshal(data, &want)
	if err != nil {
		t.Fatal(err)
	}

	return want
}

func getAdsTxtGroupByDPResponse(t *testing.T, datapath string) *dto.AdsTxtGroupByDPResponse {
	data, err := getFileData(datapath)
	if err != nil {
		t.Fatal(err)
	}

	var want *dto.AdsTxtGroupByDPResponse
	err = json.Unmarshal(data, &want)
	if err != nil {
		t.Fatal(err)
	}

	return want
}

func getFileData(datapath string) ([]byte, error) {
	f, err := os.Open(datapath)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}
