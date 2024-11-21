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
			requestBody: `{section_type": "Floors","query": "1"}`,
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
				response:   `{"DPO Rule":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Floors":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Publisher / domain - Dashboard":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Publisher / domain - Demand":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Publisher / domain list":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Publisher list":[],"Targeting - Bidder":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}],"Targeting - JS":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "PublisherSectionType",
			requestBody: `{"section_type": "Publisher list","query": "online"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Publisher list":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":null}]}`,
			},
		},
		{
			name:        "DomainSectionType",
			requestBody: `{"section_type": "Publisher / domain list","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Publisher / domain list":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "DomainDashboardSectionType",
			requestBody: `{"section_type": "Publisher / domain - Dashboard","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Publisher / domain - Dashboard":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "FactorSectionType",
			requestBody: `{"section_type": "Targeting - Bidder","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Targeting - Bidder":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "JSTargetingSectionType",
			requestBody: `{"section_type": "Targeting - JS","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Targeting - JS":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "FloorsSectionType",
			requestBody: `{"section_type": "Floors","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Floors":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "PublisherDemandSectionType",
			requestBody: `{"section_type": "Publisher / domain - Demand","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"Publisher / domain - Demand":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "DPOSectionType",
			requestBody: `{"section_type": "DPO Rule","query": "oms"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"DPO Rule":[{"publisher_id":"999","publisher_name":"online-media-soluctions","domain":"oms.com"}]}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"query": "verylonguselessquery"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"DPO Rule":[],"Floors":[],"Publisher / domain - Dashboard":[],"Publisher / domain - Demand":[],"Publisher / domain list":[],"Publisher list":[],"Targeting - Bidder":[],"Targeting - JS":[]}`,
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
