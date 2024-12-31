package rest

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
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
				response:   `[{"demand_partner_id":"Finkiel","demand_partner_name":"Finkiel DP","dp_domain":"finkiel.com","children":[{"id":1,"parent_id":"Finkiel","dp_child_name":"Open X","dp_child_domain":"openx.com","publisher_account":"88888","certification_authority_id":null,"is_required_for_ads_txt":false,"active":true,"created_at":"2024-10-01T13:51:28.407Z","updated_at":null}],"connection":[{"id":4,"demand_partner_id":"Finkiel","publisher_account":"11111","integration_type":["js","s2s"],"active":true,"created_at":"2024-10-01T13:51:28.407Z","updated_at":null}],"certification_authority_id":"jtfliy6893gfc","approval_process":"Other","dp_blocks":"Other","poc_name":"","poc_email":"","seat_owner_id":10,"manager_id":1,"is_include":false,"active":true,"is_direct":false,"is_approval_needed":true,"approval_before_going_live":false,"is_required_for_ads_txt":true,"score":3,"comments":null,"created_at":"2024-06-25T14:51:57Z","updated_at":"2024-06-25T14:51:57Z"}]`,
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
					"connection": [
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
							"is_required_for_ads_txt": false
						},
						{
							"dp_child_name": "Appnexus",
							"dp_child_domain": "appnexus.com",
							"publisher_account": "EFGH5678",
							"certification_authority_id": "appnexus_id",
							"is_required_for_ads_txt": true
						}
					],
					"connection": [
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
							"is_required_for_ads_txt": false
						},
						{
							"dp_child_name": "Appnexus",
							"dp_child_domain": "appnexus.com",
							"publisher_account": "EFGH5678",
							"certification_authority_id": "appnexus_id",
							"is_required_for_ads_txt": true
						}
					],
					"connection": [
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
