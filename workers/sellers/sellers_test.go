package sellers

import (
	"encoding/json"
	"github.com/m6yf/bcwork/utils/constant"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestCheckSellersArray(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"Valid Array", []interface{}{"seller1", "seller2"}, false},
		{"Invalid Type (String)", "invalidString", true},
		{"Invalid Type (Map)", map[string]interface{}{"seller": "value"}, true},
		{"Empty Array", []interface{}{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CheckSellersArray(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSellersArray() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchDataFromWebsite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"sellers": []interface{}{"seller1", "seller2"},
		}
		if strings.Contains(r.URL.Path, "nosellers") {
			delete(response, "sellers")
		} else if strings.Contains(r.URL.Path, "invalidsellers") {
			response["sellers"] = "invalid"
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Missing Sellers", server.URL + "/nosellers", true},
		{"Invalid Sellers Format", server.URL + "/invalidsellers", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FetchDataFromWebsite(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchDataFromWebsite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAdsTxtStatus(t *testing.T) {
	tests := []struct {
		name           string
		domain         string
		sellerId       string
		competitorType string
		expected       string
	}{
		{
			name:           "Valid ads.txt with seller ID",
			domain:         "dailydot.com",
			sellerId:       "12754",
			competitorType: "web",
			expected:       constant.AdsTxtIncludedStatus,
		},
		{
			name:           "in - valid ads.txt with seller ID",
			domain:         "dailydot.com",
			sellerId:       "44111",
			competitorType: "web",
			expected:       constant.AdsTxtNotIncludedStatus,
		},
		{
			name:           "Empty domain",
			domain:         "",
			sellerId:       "seller123",
			competitorType: "inapp",
			expected:       constant.AdsTxtNotVerifiedStatus,
		},
		{
			name:           "HTTP error",
			domain:         "example.com",
			sellerId:       "seller123",
			competitorType: "web",
			expected:       constant.AdsTxtNotVerifiedStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := &Worker{}
			got := worker.GetAdsTxtStatus(tt.domain, tt.sellerId, tt.competitorType)
			if got != tt.expected {
				t.Errorf("GetAdsTxtStatus(%q, %q) = %q; want %q", tt.domain, tt.sellerId, got, tt.expected)
			}
		})
	}
}

func TestPrepareDeletedData(t *testing.T) {
	worker := &Worker{}

	tests := []struct {
		name              string
		deletedPublishers []string
		deletedDomains    []string
		sellerTypes       []string
		expectedResult    []PublisherDomain
		expectPanic       bool
	}{
		{
			name:              "Valid Input",
			deletedPublishers: []string{"Publisher1", "Publisher2"},
			deletedDomains:    []string{"Domain1", "Domain2"},
			sellerTypes:       []string{"Type1", "Type2"},
			expectedResult: []PublisherDomain{
				{Publisher: "Publisher1", Domain: "Domain1", SellerType: "Type1"},
				{Publisher: "Publisher2", Domain: "Domain2", SellerType: "Type2"},
			},
			expectPanic: false,
		},
		{
			name:              "Nil Input",
			deletedPublishers: nil,
			deletedDomains:    nil,
			sellerTypes:       nil,
			expectedResult:    []PublisherDomain{},
			expectPanic:       false,
		},
		{
			name:              "Empty Input",
			deletedPublishers: []string{},
			deletedDomains:    []string{},
			sellerTypes:       []string{},
			expectedResult:    []PublisherDomain{},
			expectPanic:       false,
		},
		{
			name:              "Mismatched Input Sizes",
			deletedPublishers: []string{"Publisher1", "Publisher2"},
			deletedDomains:    []string{"Domain1"},
			sellerTypes:       []string{"Type1"},
			expectPanic:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic, but no panic occurred")
					}
				}()
			}

			result := worker.prepareDeletedData(test.deletedPublishers, test.deletedDomains, test.sellerTypes)
			if !reflect.DeepEqual(result, test.expectedResult) && !test.expectPanic {
				t.Errorf("Expected %v, got %v", test.expectedResult, result)
			}
		})
	}
}

func TestGetAdsTxtUrl(t *testing.T) {
	tests := []struct {
		name           string
		domain         string
		competitorType string
		expected       string
	}{
		{
			name:           "Ads txt url for web",
			domain:         "dailydot.com",
			competitorType: "web",
			expected:       "https://dailydot.com/ads.txt",
		},
		{
			name:           "Ads txt url for in-app",
			domain:         "dailydot.com",
			competitorType: "inapp",
			expected:       "https://dailydot.com/app-ads.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := &Worker{}
			got := worker.GetAdsTxtUrl(tt.domain, tt.competitorType)
			if got != tt.expected {
				t.Errorf("GetAdsTxtUrl(%q) = %q; want %q", tt.domain, got, tt.expected)
			}
		})
	}
}
