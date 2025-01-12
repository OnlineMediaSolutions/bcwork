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
	"github.com/stretchr/testify/assert"
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
				response:   `[{"id":10,"seat_owner_name":"OMS","seat_owner_domain":"onlinemediasolutions.com","publisher_account":"%s","certification_authority_id":"","created_at":"2024-10-01T13:51:28.407Z","updated_at":null}]`,
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
				response:   `[{"demand_partner_id":"Finkiel","demand_partner_name":"Finkiel DP","dp_domain":"finkiel.com","children":[{"id":1,"parent_id":"Finkiel","dp_child_name":"Open X","dp_child_domain":"openx.com","publisher_account":"88888","certification_authority_id":null,"is_required_for_ads_txt":false,"is_direct":false,"active":true,"created_at":"2024-10-01T13:51:28.407Z","updated_at":null}],"connections":[{"id":4,"demand_partner_id":"Finkiel","publisher_account":"11111","integration_type":["js","s2s"],"active":true,"created_at":"2024-10-01T13:51:28.407Z","updated_at":null}],"certification_authority_id":"jtfliy6893gfc","approval_process":"Other","dp_blocks":"Other","poc_name":"","poc_email":"","seat_owner_id":10,"manager_id":1,"is_include":false,"active":true,"is_direct":false,"is_approval_needed":true,"approval_before_going_live":false,"is_required_for_ads_txt":true,"automation":false,"automation_name":"","threshold":0,"score":3,"comments":null,"created_at":"2024-06-25T14:51:57Z","updated_at":"2024-06-25T14:51:57Z"}]`,
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
					],
					"connections": [
						{
							"publisher_account": "77777",
							"integration_type": [
								"js",
								"s2s"
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
					],
					"connections": [
						{
							"id": 4,
							"publisher_account": "11111",
							"integration_type": [
								"js", "s2s"
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
			name:        "noTargetingFoundToUpdate",
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
					],
					"connections": [
						{
							"id": 4,
							"publisher_account": "11111",
							"integration_type": [
								"js", "s2s"
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
	getEndpoint := "/test/dp/get"
	getRequestBody := `{"filter": {"demand_partner_name": ["Flow"]}}`

	setEndpoint := "/test/dp/set"
	setRequestBody := `
		{
			"demand_partner_name": "Flow",
			"dp_domain": "flow.com",
			"children": [
				{
					"dp_child_name": "Index",
					"dp_child_domain": "index.com",
					"publisher_account": "12345678",
					"certification_authority_id": "index_id",
					"is_required_for_ads_txt": false,
					"active": true
				},
				{
					"dp_child_name": "OpenX",
					"dp_child_domain": "openx.com",
					"publisher_account": "87654321",
					"certification_authority_id": "openx_id",
					"is_required_for_ads_txt": true,
					"active": true
				}
			],
			"connections": [
				{
					"publisher_account": "a1b2c3d4",
					"integration_type": [
						"js",
						"s2s"
					],
					"active": true
				}
			],
			"certification_authority_id": "flow_ca_id",
			"seat_owner_id": 10,
			"manager_id": 1,
			"approval_process": "GDoc",
			"dp_blocks": "Other",
			"is_include": false,
			"active": true,
			"is_direct": false,
			"is_approval_needed": true,
			"approval_before_going_live": true,
			"is_required_for_ads_txt": true,
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
			"children": [
				{
					"dp_child_name": "Index",
					"dp_child_domain": "index.com",
					"publisher_account": "12345678",
					"certification_authority_id": "index_id",
					"is_required_for_ads_txt": true,
					"active": true
				},
				{
					"dp_child_name": "OpenX",
					"dp_child_domain": "openx.com",
					"publisher_account": "87654321",
					"certification_authority_id": "openx_id",
					"is_required_for_ads_txt": false,
					"active": true
				}
			],
			"connections": [
				{
					"publisher_account": "a1b2c3d4",
					"integration_type": [
						"s2s"
					],
					"active": true
				},
				{
					"publisher_account": "e5f6g7h8",
					"integration_type": [
						"js"
					],
					"active": true
				}
			],
			"certification_authority_id": "flow_ca_id",
			"seat_owner_id": 10,
			"manager_id": 1,
			"approval_process": "GDoc",
			"dp_blocks": "Other",
			"is_include": false,
			"active": true,
			"is_direct": false,
			"is_approval_needed": true,
			"approval_before_going_live": true,
			"is_required_for_ads_txt": true,
			"poc_name": "poc_name_2",
			"poc_email": "poc@mail.com"
		}
	`

	updateRequestBody2 := `
		{
			"demand_partner_id": "flow",
			"demand_partner_name": "Flow",
			"dp_domain": "flow.com",
			"children": [
				{
					"dp_child_name": "OpenX",
					"dp_child_domain": "openx.com",
					"publisher_account": "87654321",
					"certification_authority_id": "openx_id",
					"is_required_for_ads_txt": false,
					"active": true
				}
			],
			"connections": [
				{
					"publisher_account": "e5f6g7h8",
					"integration_type": [
						"js"
					],
					"active": true
				}
			],
			"certification_authority_id": "flow_ca_id",
			"seat_owner_id": 10,
			"manager_id": 1,
			"approval_process": "GDoc",
			"dp_blocks": "Other",
			"is_include": false,
			"active": true,
			"is_direct": false,
			"is_approval_needed": true,
			"approval_before_going_live": true,
			"is_required_for_ads_txt": true,
			"poc_name": "poc_name_2",
			"poc_email": "poc@mail.com"
		}
	`

	mockDP := getMockDemandPartner()

	// creating new demand partner with 2 childs and 1 connection
	setReq, err := http.NewRequest(fiber.MethodPost, baseURL+setEndpoint, strings.NewReader(setRequestBody))
	assert.NoError(t, err)
	setReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	setResp, err := http.DefaultClient.Do(setReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, setResp.StatusCode)
	// checking result
	getReq, err := http.NewRequest(fiber.MethodPost, baseURL+getEndpoint, strings.NewReader(getRequestBody))
	assert.NoError(t, err)
	getReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	getResp, err := http.DefaultClient.Do(getReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, getResp.StatusCode)
	getRespBody, err := io.ReadAll(getResp.Body)
	assert.NoError(t, err)
	defer getResp.Body.Close()
	dps, err := getDPFromResponse(getRespBody)
	assert.NoError(t, err)
	assert.Equal(t, mockDP, dps)

	// updating demand partner: update poc_name, poc_email, children, connection[0] and add 1 new connection
	updateReq, err := http.NewRequest(fiber.MethodPost, baseURL+updateEndpoint, strings.NewReader(updateRequestBody1))
	assert.NoError(t, err)
	updateReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	updateResp, err := http.DefaultClient.Do(updateReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, updateResp.StatusCode)
	// checking result
	getReq2, err := http.NewRequest(fiber.MethodPost, baseURL+getEndpoint, strings.NewReader(getRequestBody))
	assert.NoError(t, err)
	getReq2.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	getResp2, err := http.DefaultClient.Do(getReq2)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, getResp2.StatusCode)
	getRespBody2, err := io.ReadAll(getResp2.Body)
	assert.NoError(t, err)
	defer getResp2.Body.Close()
	dps, err = getDPFromResponse(getRespBody2)
	assert.NoError(t, err)
	// add changes to mock dp
	mockDP[0].POCName = "poc_name_2"
	mockDP[0].Connections[0].IntegrationType = []string{"s2s"}
	mockDP[0].Connections = append(mockDP[0].Connections,
		&dto.DemandPartnerConnection{
			DemandPartnerID:  "flow",
			PublisherAccount: "e5f6g7h8",
			IntegrationType:  []string{"js"},
			Active:           true,
		})
	mockDP[0].Children[0].IsRequiredForAdsTxt = true
	mockDP[0].Children[1].IsRequiredForAdsTxt = false
	assert.Equal(t, mockDP, dps)

	// updating demand partner: turn off first child and first connection
	updateReq2, err := http.NewRequest(fiber.MethodPost, baseURL+updateEndpoint, strings.NewReader(updateRequestBody2))
	assert.NoError(t, err)
	updateReq2.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	updateResp2, err := http.DefaultClient.Do(updateReq2)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, updateResp2.StatusCode)
	// checking result
	getReq3, err := http.NewRequest(fiber.MethodPost, baseURL+getEndpoint, strings.NewReader(getRequestBody))
	assert.NoError(t, err)
	getReq3.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	getResp3, err := http.DefaultClient.Do(getReq3)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, getResp3.StatusCode)
	getRespBody3, err := io.ReadAll(getResp3.Body)
	assert.NoError(t, err)
	defer getResp3.Body.Close()
	dps, err = getDPFromResponse(getRespBody3)
	assert.NoError(t, err)
	// add changes to mock dp
	mockDP[0].Connections[0].Active = false
	mockDP[0].Children[0].Active = false
	mockDP[0].Connections[0], mockDP[0].Connections[1] = mockDP[0].Connections[1], mockDP[0].Connections[0]
	mockDP[0].Children[0], mockDP[0].Children[1] = mockDP[0].Children[1], mockDP[0].Children[0]
	assert.Equal(t, mockDP, dps)
}

func getMockDemandPartner() []*dto.DemandPartner {
	return []*dto.DemandPartner{
		{
			DemandPartnerID:   "flow",
			DemandPartnerName: "Flow",
			DPDomain:          "flow.com",
			Children: []*dto.DemandPartnerChild{
				{
					ParentID:         "flow",
					DPChildName:      "Index",
					DPChildDomain:    "index.com",
					PublisherAccount: "12345678",
					CertificationAuthorityID: func() *string {
						s := "index_id"
						return &s
					}(),
					Active: true,
				},
				{
					ParentID:         "flow",
					DPChildName:      "OpenX",
					DPChildDomain:    "openx.com",
					PublisherAccount: "87654321",
					CertificationAuthorityID: func() *string {
						s := "openx_id"
						return &s
					}(),
					IsRequiredForAdsTxt: true,
					Active:              true,
				},
			},
			Connections: []*dto.DemandPartnerConnection{
				{
					DemandPartnerID:  "flow",
					PublisherAccount: "a1b2c3d4",
					IntegrationType:  []string{"js", "s2s"},
					Active:           true,
				},
			},
			CertificationAuthorityID: func() *string {
				s := "flow_ca_id"
				return &s
			}(),
			ApprovalProcess: "GDoc",
			DPBlocks:        "Other",
			POCName:         "pocname",
			POCEmail:        "poc@mail.com",
			SeatOwnerID: func() *int {
				n := 10
				return &n
			}(),
			ManagerID: func() *int {
				n := 1
				return &n
			}(),
			Active:                  true,
			IsApprovalNeeded:        true,
			ApprovalBeforeGoingLive: true,
			IsRequiredForAdsTxt:     true,
			Score:                   1000,
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

		for _, child := range dp.Children {
			child.ID = 0
			child.CreatedAt = time.Time{}
			child.UpdatedAt = nil
		}

		for _, connection := range dp.Connections {
			connection.ID = 0
			connection.CreatedAt = time.Time{}
			connection.UpdatedAt = nil
		}
	}

	return dps, nil
}
