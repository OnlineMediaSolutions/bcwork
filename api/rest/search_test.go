package rest

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSearchHandler(t *testing.T) {
	endpoint := "/test/search"

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
			name:        "invalidRequest",
			requestBody: `{section_type": "Floors list","query": "1"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse search request","error":"invalid character 's' looking for beginning of object key string"}`,
			},
		},
		{
			name:        "allSections",
			requestBody: `{"query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Bidder Targetings":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"DPO Rules":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"JS Targetings":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Domains list":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Floors list":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Publishers list":[],"Domain - Dashboard":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Domain - Demand":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "PublisherSectionType",
			requestBody: `{"section_type": "Publishers list","query": "online"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Publishers list":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":""}]}`,
			},
		},
		{
			name:        "DomainSectionType",
			requestBody: `{"section_type": "Domains list","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Domains list":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "DomainDashboardSectionType",
			requestBody: `{"section_type": "Domain - Dashboard","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Domain - Dashboard":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "FactorSectionType",
			requestBody: `{"section_type": "Bidder Targetings","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Bidder Targetings":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "FactorSectionType_noActiveRules",
			requestBody: `{"section_type": "Bidder Targetings","query": "brightcom"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Bidder Targetings":[]}`,
			},
		},
		{
			name:        "JSTargetingSectionType",
			requestBody: `{"section_type": "JS Targetings","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"JS Targetings":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "FloorsSectionType",
			requestBody: `{"section_type": "Floors list","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Floors list":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "FloorsSectionType_noActiveRules",
			requestBody: `{"section_type": "Floors list","query": "brightcom"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Floors list":[]}`,
			},
		},
		{
			name:        "PublisherDemandSectionType",
			requestBody: `{"section_type": "Domain - Demand","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Domain - Demand":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "DPOSectionType",
			requestBody: `{"section_type": "DPO Rules","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"DPO Rules":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "DPOSectionType_noActiveRules",
			requestBody: `{"section_type": "DPO Rules","query": "active"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"DPO Rules":[]}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"query": "verylonguselessquery"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Bidder Targetings":[],"DPO Rules":[],"JS Targetings":[],"Domains list":[],"Floors list":[],"Publishers list":[],"Domain - Dashboard":[],"Domain - Demand":[]}`,
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
