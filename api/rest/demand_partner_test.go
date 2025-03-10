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
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
)

func TestDemandPartnerGetSeatOwnersHandler(t *testing.T) {
	endpoint := "/test/dp/seat_owner/get"

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
			requestBody: `{"filter": {"seat_owner_name": ["OMS"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"id":4,"seat_owner_name":"OMS","seat_owner_domain":"onlinemediasolutions.com","publisher_account":"%s","certification_authority_id":"","ads_txt_line":"onlinemediasolutions.com, XXXXX, DIRECT","line_name":"OMS - Direct","created_at":"2024-10-01T13:51:28.407Z","updated_at":null}]`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"filter": {"seat_owner_name: ["OMS"]}}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse request for getting seat owners data","error":"invalid character 'O' after object key"}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"filter": {"seat_owner_name": ["xxxxxxxx"]}}`,
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

func TestDemandPartnerGetHandler(t *testing.T) {
	endpoint := "/test/dp/get"

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
			requestBody: `{"filter": {"demand_partner_name": ["Finkiel DP"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"demand_partner_id":"Finkiel","demand_partner_name":"Finkiel DP","dp_domain":"finkiel.com","connections":[{"id":2,"demand_partner_id":"Finkiel","publisher_account":"11111","media_type":["Web Banners"],"is_direct":false,"is_required_for_ads_txt":true,"children":[],"ads_txt_line":"finkiel.com, 11111, RESELLER, jtfliy6893gfc","line_name":"Finkiel DP - Finkiel DP","created_at":"2024-10-01T13:51:28.407Z","updated_at":null}],"certification_authority_id":"jtfliy6893gfc","approval_process":"Other","dp_blocks":"Other","poc_name":"","poc_email":"","seat_owner_id":2,"seat_owner_name":"GetMedia","manager_id":1,"manager_full_name":"name_1 surname_1","integration_type":["oRTB","Prebid Server"],"media_type_list":["Web Banners"],"is_include":false,"active":true,"is_approval_needed":true,"approval_before_going_live":false,"automation":false,"automation_name":"","threshold":0,"score":3,"comments":null,"created_at":"2024-06-25T14:51:57Z","updated_at":"2024-06-25T14:51:57Z"}]`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"filter": {"demand_partner_name: ["Finkiel DP"]}}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse request for getting demand partners data","error":"invalid character 'F' after object key"}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"filter": {"demand_partner_name": ["xxxxxxxx"]}}`,
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

func TestDemandPartnerSetHandler(t *testing.T) {
	endpoint := "/test/dp/set"

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
			name: "validRequest",
			requestBody: `
				{
					"demand_partner_name": "New Demand Partner",
					"dp_domain": "newdemandpartner.com",
					"integration_type": [
						"oRTB",
						"Prebid Server"
					],
					"connections": [
						{
							"publisher_account": "77777",
							"media_type": [
								"Video",
								"Web Banners"
							],
							"children": [
								{
									"dp_child_name": "Pubmatic",
									"dp_child_domain": "pubmatic.com",
									"publisher_account": "abcd1234",
									"certification_authority_id": "pubmatic_id",
									"is_required_for_ads_txt": false
								},
								{
									"dp_child_name": "Appnexus",
									"dp_child_domain": "appnexus.com",
									"publisher_account": "efgh5678",
									"certification_authority_id": "appnexus_id",
									"is_required_for_ads_txt": true
								}
							]
						}
					],
					"certification_authority_id": "new_demand_partner_ca_id",
					"seat_owner_id": 10,
					"manager_id": 1,
					"approval_process": "GDoc",
					"is_include": false,
					"active": true,
					"is_direct": false,
					"is_approval_needed": true,
					"approval_before_going_live": true,
					"is_required_for_ads_txt": true,
					"poc_name": "pocname",
					"poc_email": "poc@mail.com"
				}
			`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"demand partner successfully created"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"demand_partner_name: "New Demand Partner 2"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse request for creating demand partner","error":"invalid character 'N' after object key"}`,
			},
		},
		{
			name:        "hasDuplicate",
			requestBody: `{"demand_partner_name": "Amazon"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to create demand partner","error":"demand partner with name [Amazon] already exists"}`,
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

func TestDemandPartnerUpdateHandler(t *testing.T) {
	endpoint := "/test/dp/update"

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
			name: "validRequest",
			requestBody: `
				{
					"demand_partner_id": "Finkiel",
					"demand_partner_name": "Finkiel DP",
					"dp_domain": "finkiel.com",
					"integration_type": [
						"oRTB", "Prebid Server"
					],
					"connections": [
						{
							"id": 2,
							"publisher_account": "11111",
							"media_type": [
								"Web Banners", "Video"
							],
							"children": [
								{
									"dp_child_name": "Pubmatic",
									"dp_child_domain": "pubmatic.com",
									"publisher_account": "ABCD1234",
									"certification_authority_id": "pubmatic_id",
									"active": true,
									"is_required_for_ads_txt": false
								},
								{
									"dp_child_name": "Appnexus",
									"dp_child_domain": "appnexus.com",
									"publisher_account": "EFGH5678",
									"certification_authority_id": "appnexus_id",
									"active": true,
									"is_required_for_ads_txt": true
								}
							]
						}
					],
					"certification_authority_id": "new_demand_partner_ca_id",
					"seat_owner_id": 10,
					"manager_id": 1,
					"approval_process": "GDoc",
					"dp_blocks": "Other",
					"is_include": false,
					"active": true,
					"is_direct": false,
					"is_approval_needed": false,
					"approval_before_going_live": false,
					"is_required_for_ads_txt": true,
					"poc_name": "pocnamenew",
					"poc_email": "pocnew@mail.com"
				}
			`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"demand partner successfully updated"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"demand_partner_name: "New Demand Partner"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse request for updating demand partner","error":"invalid character 'N' after object key"}`,
			},
		},
		{
			name:        "noDemandPartnerFoundToUpdate",
			requestBody: `{"demand_partner_id": "unknowndemandpartner"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update demand partner","error":"failed to get demand partner with id [unknowndemandpartner] to update: sql: no rows in result set"}`,
			},
		},
		{
			// based on results of "validRequest"
			name: "nothingToUpdate",
			requestBody: `
				{
					"demand_partner_id": "Finkiel",
					"demand_partner_name": "Finkiel DP",
					"dp_domain": "finkiel.com",
					"integration_type": [
						"oRTB", "Prebid Server"
					],
					"connections": [
						{
							"id": 2,
							"publisher_account": "11111",
							"media_type": [
								"Web Banners", "Video"
							],
							"children": [
								{
									"dp_child_name": "Pubmatic",
									"dp_child_domain": "pubmatic.com",
									"publisher_account": "ABCD1234",
									"certification_authority_id": "pubmatic_id",
									"active": true,
									"is_required_for_ads_txt": false
								},
								{
									"dp_child_name": "Appnexus",
									"dp_child_domain": "appnexus.com",
									"publisher_account": "EFGH5678",
									"certification_authority_id": "appnexus_id",
									"active": true,
									"is_required_for_ads_txt": true
								}
							]
						}
					],
					"certification_authority_id": "new_demand_partner_ca_id",
					"seat_owner_id": 10,
					"manager_id": 1,
					"approval_process": "GDoc",
					"dp_blocks": "Other",
					"is_include": false,
					"active": true,
					"is_direct": false,
					"is_approval_needed": false,
					"approval_before_going_live": false,
					"is_required_for_ads_txt": true,
					"poc_name": "pocnamenew",
					"poc_email": "pocnew@mail.com"
				}
			`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update demand partner","error":"there are no new values to update demand partner"}`,
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

func TestDemandPartnerFlow(t *testing.T) {
	adsTxtEndpoint := "/test/ads_txt/main"
	getEndpoint := "/test/dp/get"
	setEndpoint := "/test/dp/set"

	getRequestBody := `{"filter": {"demand_partner_name": ["Flow"]}}`

	setRequestBody := `
		{
			"demand_partner_name": "Flow",
			"dp_domain": "flow.com",
			"integration_type": [
				"oRTB",
				"Prebid Server"
			],
			"connections": [
				{
					"publisher_account": "a1b2c3d4",
					"media_type": [
						"Web Banners",
						"Video"
					],
					"children": [
						{
							"dp_child_name": "Index",
							"dp_child_domain": "index.com",
							"publisher_account": "12345678",
							"certification_authority_id": "index_id",
							"is_required_for_ads_txt": false
						}
					],
					"is_direct": false,
					"is_required_for_ads_txt": true
				}
			],
			"certification_authority_id": "flow_ca_id",
			"seat_owner_id": 5,
			"manager_id": 1,
			"approval_process": "GDoc",
			"dp_blocks": "Other",
			"is_include": false,
			"active": true,
			"is_approval_needed": true,
			"approval_before_going_live": true,
			"poc_name": "pocname",
			"poc_email": "poc@mail.com"
		}
	`

	updateEndpoint := "/test/dp/update"
	updateRequestBody1 := `
		{
			"demand_partner_id": "flow",
			"demand_partner_name": "Flow",
			"dp_domain": "flow.com",
			"integration_type": [
				"oRTB",
				"Prebid Server"
			],
			"connections": [
				{
					"publisher_account": "a1b2c3d4",
					"media_type": [
						"Video"
					],
					"children": [
						{
							"dp_child_name": "Index",
							"dp_child_domain": "index.com",
							"publisher_account": "12345678",
							"certification_authority_id": "index_id",
							"is_required_for_ads_txt": true
						},
						{
							"dp_child_name": "OpenX",
							"dp_child_domain": "openx.com",
							"publisher_account": "87654321",
							"certification_authority_id": "openx_id",
							"is_required_for_ads_txt": false
						}
					],
					"is_direct": false,
					"is_required_for_ads_txt": true
				},
				{
					"publisher_account": "e5f6g7h8",
					"media_type": [
						"Web Banners"
					],
					"is_direct": false,
					"is_required_for_ads_txt": true
				}
			],
			"certification_authority_id": "flow_ca_id",
			"seat_owner_id": 5,
			"manager_id": 1,
			"approval_process": "GDoc",
			"dp_blocks": "Other",
			"is_include": false,
			"active": true,
			"is_approval_needed": true,
			"approval_before_going_live": true,
			"poc_name": "poc_name_2",
			"poc_email": "poc@mail.com"
		}
	`

	updateRequestBody2 := `
		{
			"demand_partner_id": "flow",
			"demand_partner_name": "Flow",
			"dp_domain": "flow.com",
			"integration_type": [
				"oRTB",
				"Prebid Server"
			],
			"connections": [
				{
					"publisher_account": "e5f6g7h8",
					"media_type": [
						"Web Banners"
					],
					"children": [
						{
							"dp_child_name": "OpenX",
							"dp_child_domain": "openx.com",
							"publisher_account": "87654321",
							"certification_authority_id": "openx_id",
							"is_required_for_ads_txt": false
						}
					],
					"is_direct": false,
					"is_required_for_ads_txt": true
				}
			],
			"certification_authority_id": "flow_ca_id",
			"seat_owner_id": 6,
			"manager_id": 1,
			"approval_process": "GDoc",
			"dp_blocks": "Other",
			"is_include": false,
			"active": true,
			"is_approval_needed": true,
			"approval_before_going_live": true,
			"poc_name": "poc_name_2",
			"poc_email": "poc@mail.com"
		}
	`

	mockDP := getMockDemandPartner()
	mockAdsTxtLines := getMockAdsTxtLines()

	// creating new demand partner with 1 connection (1 child)
	setReq, err := http.NewRequest(fiber.MethodPost, baseURL+setEndpoint, strings.NewReader(setRequestBody))
	require.NoError(t, err)
	setReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	setResp, err := http.DefaultClient.Do(setReq)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, setResp.StatusCode)
	// checking result
	getReq, err := http.NewRequest(fiber.MethodPost, baseURL+getEndpoint, strings.NewReader(getRequestBody))
	require.NoError(t, err)
	getReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	getResp, err := http.DefaultClient.Do(getReq)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, getResp.StatusCode)
	defer getResp.Body.Close()
	getRespBody, err := io.ReadAll(getResp.Body)
	require.NoError(t, err)
	dps, err := getDPFromResponse(getRespBody)
	require.NoError(t, err)
	require.Equal(t, mockDP, dps)
	// getting ads.txt lines
	adsTxtReq, err := http.NewRequest(fiber.MethodPost, baseURL+adsTxtEndpoint, strings.NewReader("{}"))
	require.NoError(t, err)
	adsTxtReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	adsTxtResp, err := http.DefaultClient.Do(adsTxtReq)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, adsTxtResp.StatusCode)
	defer adsTxtResp.Body.Close()
	adsTxtRespBody, err := io.ReadAll(adsTxtResp.Body)
	require.NoError(t, err)
	adsTxtLines, err := getAdsTxtFromResponse(adsTxtRespBody)
	require.NoError(t, err)
	for _, mockAdsTxtLine := range mockAdsTxtLines {
		require.Contains(t, adsTxtLines, mockAdsTxtLine)
	}

	// updating demand partner:
	// update poc_name, poc_email, connection[0] (add new child), connection[0]child[0].IsRequired = true and add 1 new connection (without children)
	updateReq, err := http.NewRequest(fiber.MethodPost, baseURL+updateEndpoint, strings.NewReader(updateRequestBody1))
	require.NoError(t, err)
	updateReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	updateResp, err := http.DefaultClient.Do(updateReq)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, updateResp.StatusCode)
	// checking result
	getReq2, err := http.NewRequest(fiber.MethodPost, baseURL+getEndpoint, strings.NewReader(getRequestBody))
	require.NoError(t, err)
	getReq2.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	getResp2, err := http.DefaultClient.Do(getReq2)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, getResp2.StatusCode)
	defer getResp2.Body.Close()
	getRespBody2, err := io.ReadAll(getResp2.Body)
	require.NoError(t, err)
	dps, err = getDPFromResponse(getRespBody2)
	require.NoError(t, err)
	// add changes to mock dp
	mockDP[0].POCName = "poc_name_2"
	mockDP[0].Connections[0].MediaType = []string{"Video"}
	mockDP[0].Connections = append(mockDP[0].Connections,
		&dto.DemandPartnerConnection{
			DemandPartnerID:     "flow",
			PublisherAccount:    "e5f6g7h8",
			MediaType:           []string{"Web Banners"},
			Children:            []*dto.DemandPartnerChild{},
			AdsTxtLine:          "flow.com, e5f6g7h8, RESELLER, flow_ca_id",
			LineName:            "Flow - Flow",
			IsRequiredForAdsTxt: true,
		})
	mockDP[0].Connections[0].Children[0].IsRequiredForAdsTxt = true
	mockDP[0].Connections[0].Children = append(mockDP[0].Connections[0].Children,
		&dto.DemandPartnerChild{
			DPChildName:              "OpenX",
			DPChildDomain:            "openx.com",
			PublisherAccount:         "87654321",
			CertificationAuthorityID: helpers.GetPointerToString("openx_id"),
			AdsTxtLine:               "openx.com, 87654321, RESELLER, openx_id",
			LineName:                 "Flow - OpenX",
			IsRequiredForAdsTxt:      false,
		})
	require.Equal(t, mockDP, dps)
	// getting ads.txt lines
	adsTxtReq2, err := http.NewRequest(fiber.MethodPost, baseURL+adsTxtEndpoint, strings.NewReader("{}"))
	require.NoError(t, err)
	adsTxtReq2.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	adsTxtResp2, err := http.DefaultClient.Do(adsTxtReq2)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, adsTxtResp2.StatusCode)
	defer adsTxtResp2.Body.Close()
	adsTxtRespBody2, err := io.ReadAll(adsTxtResp2.Body)
	require.NoError(t, err)
	adsTxtLines2, err := getAdsTxtFromResponse(adsTxtRespBody2)
	require.NoError(t, err)
	// add changes to mock ads.txt
	mockAdsTxtLines[1].IsRequired = true
	mockAdsTxtLines[0].MediaType = []string{dto.VideoMediaType}
	mockAdsTxtLines[1].MediaType = []string{dto.VideoMediaType}
	mockAdsTxtLines = append(mockAdsTxtLines,
		&dto.AdsTxt{
			PublisherID:               "999",
			PublisherName:             "online-media-soluctions",
			Domain:                    "oms.com",
			DomainStatus:              dto.DomainStatusNew,
			DemandStatus:              dto.DPStatusNotSent,
			MediaType:                 []string{dto.WebBannersMediaType},
			DemandPartnerID:           "flow",
			DemandPartnerName:         "Flow",
			DemandPartnerNameExtended: "Flow - Flow",
			DemandManagerID:           null.StringFrom("1"),
			DemandManagerFullName:     "name_1 surname_1",
			Status:                    dto.AdsTxtStatusNotScanned,
			IsRequired:                true,
			IsDemandPartnerActive:     true,
			AdsTxtLine:                "flow.com, e5f6g7h8, RESELLER, flow_ca_id",
		},
		&dto.AdsTxt{
			PublisherID:               "999",
			PublisherName:             "online-media-soluctions",
			Domain:                    "oms.com",
			DomainStatus:              dto.DomainStatusNew,
			DemandStatus:              dto.DPStatusNotSent,
			MediaType:                 []string{dto.VideoMediaType},
			DemandPartnerID:           "flow",
			DemandPartnerName:         "Flow",
			DemandPartnerNameExtended: "Flow - OpenX",
			DemandManagerID:           null.StringFrom("1"),
			DemandManagerFullName:     "name_1 surname_1",
			Status:                    dto.AdsTxtStatusNotScanned,
			IsDemandPartnerActive:     true,
			AdsTxtLine:                "openx.com, 87654321, RESELLER, openx_id",
		},
	)
	for _, mockAdsTxtLine := range mockAdsTxtLines {
		require.Contains(t, adsTxtLines2, mockAdsTxtLine)
	}

	// updating demand partner: turn off first connection (with children), add child to second connection, changed seat owner
	updateReq2, err := http.NewRequest(fiber.MethodPost, baseURL+updateEndpoint, strings.NewReader(updateRequestBody2))
	require.NoError(t, err)
	updateReq2.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	updateResp2, err := http.DefaultClient.Do(updateReq2)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, updateResp2.StatusCode)
	// checking result
	getReq3, err := http.NewRequest(fiber.MethodPost, baseURL+getEndpoint, strings.NewReader(getRequestBody))
	require.NoError(t, err)
	getReq3.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	getResp3, err := http.DefaultClient.Do(getReq3)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, getResp3.StatusCode)
	defer getResp3.Body.Close()
	getRespBody3, err := io.ReadAll(getResp3.Body)
	require.NoError(t, err)
	dps, err = getDPFromResponse(getRespBody3)
	require.NoError(t, err)
	// add changes to mock dp
	mockDP[0].SeatOwnerID = helpers.GetPointerToInt(6)
	mockDP[0].SeatOwnerName = "TSO2"
	children1Copy := *mockDP[0].Connections[0].Children[1]
	mockDP[0].Connections[1].Children = append(mockDP[0].Connections[1].Children, &children1Copy)
	mockDP[0].Connections[0], mockDP[0].Connections[1] = mockDP[0].Connections[1], mockDP[0].Connections[0]
	mockDP[0].Connections = mockDP[0].Connections[:1]
	mockDP[0].MediaTypeList = []string{"Web Banners"}
	require.Equal(t, mockDP, dps)
	// getting ads.txt lines
	adsTxtReq3, err := http.NewRequest(fiber.MethodPost, baseURL+adsTxtEndpoint, strings.NewReader("{}"))
	require.NoError(t, err)
	adsTxtReq3.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	adsTxtResp3, err := http.DefaultClient.Do(adsTxtReq3)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, adsTxtResp3.StatusCode)
	defer adsTxtResp3.Body.Close()
	adsTxtRespBody3, err := io.ReadAll(adsTxtResp3.Body)
	require.NoError(t, err)
	adsTxtLines3, err := getAdsTxtFromResponse(adsTxtRespBody3)
	require.NoError(t, err)
	// add changes to mock ads.txt
	mockDeletedAdsTxtLines := []*dto.AdsTxt{mockAdsTxtLines[0], mockAdsTxtLines[1], mockAdsTxtLines[2]}
	mockAdsTxtLines = []*dto.AdsTxt{mockAdsTxtLines[3], mockAdsTxtLines[4]}
	mockAdsTxtLines[0].MediaType = []string{dto.WebBannersMediaType}
	mockAdsTxtLines[1].MediaType = []string{dto.WebBannersMediaType}
	mockAdsTxtLines = append(mockAdsTxtLines,
		&dto.AdsTxt{
			PublisherID:               "999",
			PublisherName:             "online-media-soluctions",
			Domain:                    "oms.com",
			DomainStatus:              dto.DomainStatusNew,
			DemandStatus:              dto.DPStatusApproved,
			DemandPartnerID:           "flow",
			DemandPartnerName:         "Flow",
			DemandPartnerNameExtended: "TSO2 - Direct",
			Status:                    dto.AdsTxtStatusNotScanned,
			IsRequired:                true,
			IsDemandPartnerActive:     true,
			AdsTxtLine:                "testseatowner2.com, 52999, DIRECT",
		},
	)
	for _, mockAdsTxtLine := range mockAdsTxtLines {
		require.Contains(t, adsTxtLines3, mockAdsTxtLine)
	}
	// also checking deleted lines
	for _, mockDeletedAdsTxtLine := range mockDeletedAdsTxtLines {
		require.NotContains(t, adsTxtLines3, mockDeletedAdsTxtLine)
	}
}

func getMockDemandPartner() []*dto.DemandPartner {
	return []*dto.DemandPartner{
		{
			DemandPartnerID:   "flow",
			DemandPartnerName: "Flow",
			DPDomain:          "flow.com",
			IntegrationType:   []string{"Prebid Server", "oRTB"},
			MediaTypeList:     []string{"Video", "Web Banners"},
			Connections: []*dto.DemandPartnerConnection{
				{
					DemandPartnerID:     "flow",
					PublisherAccount:    "a1b2c3d4",
					MediaType:           []string{"Video", "Web Banners"},
					IsRequiredForAdsTxt: true,
					AdsTxtLine:          "flow.com, a1b2c3d4, RESELLER, flow_ca_id",
					LineName:            "Flow - Flow",
					Children: []*dto.DemandPartnerChild{
						{
							DPChildName:              "Index",
							DPChildDomain:            "index.com",
							PublisherAccount:         "12345678",
							CertificationAuthorityID: helpers.GetPointerToString("index_id"),
							AdsTxtLine:               "index.com, 12345678, RESELLER, index_id",
							LineName:                 "Flow - Index",
						},
					},
				},
			},
			CertificationAuthorityID: helpers.GetPointerToString("flow_ca_id"),
			ApprovalProcess:          "GDoc",
			DPBlocks:                 "Other",
			POCName:                  "pocname",
			POCEmail:                 "poc@mail.com",
			SeatOwnerID:              helpers.GetPointerToInt(5),
			SeatOwnerName:            "TSO",
			ManagerID:                helpers.GetPointerToInt(1),
			ManagerFullName:          "name_1 surname_1",
			Active:                   true,
			IsApprovalNeeded:         true,
			ApprovalBeforeGoingLive:  true,
			Score:                    1000,
		},
	}
}

func getMockAdsTxtLines() []*dto.AdsTxt {
	return []*dto.AdsTxt{
		{
			PublisherID:               "999",
			PublisherName:             "online-media-soluctions",
			Domain:                    "oms.com",
			DomainStatus:              dto.DomainStatusNew,
			DemandStatus:              dto.DPStatusNotSent,
			DemandPartnerID:           "flow",
			DemandPartnerName:         "Flow",
			DemandPartnerNameExtended: "Flow - Flow",
			MediaType:                 []string{dto.VideoMediaType, dto.WebBannersMediaType},
			DemandManagerID:           null.StringFrom("1"),
			DemandManagerFullName:     "name_1 surname_1",
			Status:                    dto.AdsTxtStatusNotScanned,
			IsRequired:                true,
			IsDemandPartnerActive:     true,
			AdsTxtLine:                "flow.com, a1b2c3d4, RESELLER, flow_ca_id",
		},
		{
			PublisherID:               "999",
			PublisherName:             "online-media-soluctions",
			Domain:                    "oms.com",
			DomainStatus:              dto.DomainStatusNew,
			DemandStatus:              dto.DPStatusNotSent,
			DemandPartnerID:           "flow",
			DemandPartnerName:         "Flow",
			DemandPartnerNameExtended: "Flow - Index",
			MediaType:                 []string{dto.VideoMediaType, dto.WebBannersMediaType},
			DemandManagerID:           null.StringFrom("1"),
			DemandManagerFullName:     "name_1 surname_1",
			Status:                    dto.AdsTxtStatusNotScanned,
			IsDemandPartnerActive:     true,
			AdsTxtLine:                "index.com, 12345678, RESELLER, index_id",
		},
		{
			PublisherID:               "999",
			PublisherName:             "online-media-soluctions",
			Domain:                    "oms.com",
			DomainStatus:              dto.DomainStatusNew,
			DemandStatus:              dto.DPStatusApproved,
			DemandPartnerID:           "flow",
			DemandPartnerName:         "Flow",
			DemandPartnerNameExtended: "TSO - Direct",
			Status:                    dto.AdsTxtStatusNotScanned,
			IsRequired:                true,
			IsDemandPartnerActive:     true,
			AdsTxtLine:                "testseatowner.com, 5999, DIRECT",
		},
	}
}

func getDPFromResponse(body []byte) ([]*dto.DemandPartner, error) {
	var dps []*dto.DemandPartner

	err := json.Unmarshal(body, &dps)
	if err != nil {
		return nil, err
	}

	for _, dp := range dps {
		dp.CreatedAt = time.Time{}
		dp.UpdatedAt = nil

		for _, connection := range dp.Connections {
			connection.ID = 0
			connection.CreatedAt = time.Time{}
			connection.UpdatedAt = nil

			for _, child := range connection.Children {
				child.ID = 0
				child.DPConnectionID = 0
				child.CreatedAt = time.Time{}
				child.UpdatedAt = nil
			}
		}
	}

	return dps, nil
}

func getAdsTxtFromResponse(body []byte) ([]*dto.AdsTxt, error) {
	var adsTxtLines []*dto.AdsTxt

	err := json.Unmarshal(body, &adsTxtLines)
	if err != nil {
		return nil, err
	}

	for _, adsTxtLine := range adsTxtLines {
		adsTxtLine.ID = 0
	}

	return adsTxtLines, nil
}
