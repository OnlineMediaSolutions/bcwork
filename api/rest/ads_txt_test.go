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
		// {
		// 	name:        "mainTable_validRequest",
		// 	endpoint:    "/test/ads_txt/main",
		// 	requestBody: `{}`,
		// 	want: want{
		// 		statusCode: fiber.StatusOK,
		// 		response:   getAdsTxtData(t, "./testdata/ads_txt_main_table.json"),
		// 	},
		// },
		{
			name:        "cmTable_validRequest",
			endpoint:    "/test/ads_txt/cm",
			requestBody: `{}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   getAdsTxtData(t, "./testdata/ads_txt_cm_table.json"),
			},
		},
		// {
		// 	name:        "mbTable_validRequest",
		// 	endpoint:    "/test/ads_txt/mb",
		// 	requestBody: `{}`,
		// 	want: want{
		// 		statusCode: fiber.StatusOK,
		// 		response:   getAdsTxtData(t, "./testdata/ads_txt_mb_table.json"),
		// 	},
		// },
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

func TestAdsTxtGroupByDPTable(t *testing.T) {
	endpoint := "/test/ads_txt/group_by_dp"

	type want struct {
		statusCode int
		response   map[string]*dto.AdsTxtGroupedByDPData
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
				response:   getAdsTxtGroupByDPData(t, "./testdata/ads_txt_group_by_dp_table.json"),
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

			var got map[string]*dto.AdsTxtGroupedByDPData
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

func getAdsTxtGroupByDPData(t *testing.T, datapath string) map[string]*dto.AdsTxtGroupedByDPData {
	data, err := getFileData(datapath)
	if err != nil {
		t.Fatal(err)
	}

	var want map[string]*dto.AdsTxtGroupedByDPData
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
