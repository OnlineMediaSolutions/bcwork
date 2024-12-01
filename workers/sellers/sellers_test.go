package sellers

import (
	"encoding/json"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
		mockResponse   string
		mockStatusCode int
		expectedStatus string
		expectedError  string
	}{
		{
			name:           "Seller ID not included",
			domain:         "example.com",
			sellerId:       "12345",
			mockResponse:   "example.com, 12345\nother.com, 67890\n",
			mockStatusCode: http.StatusOK,
			expectedStatus: constant.AdsTxtNotVerifiedStatus,
			expectedError:  "ads.txt not found or invalid for domain example.com",
		},
		{
			name:           "Seller ID not included",
			domain:         "google.com",
			sellerId:       "105199474",
			mockResponse:   "google.com, 105199474\nother.com, 67890\n",
			mockStatusCode: http.StatusOK,
			expectedStatus: constant.AdsTxtNotVerifiedStatus,
			expectedError:  "ads.txt not found or invalid for domain google.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := &Worker{}
			status, err := worker.GetAdsTxtStatus(tt.domain, tt.sellerId)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, status)
			}
		})
	}
}
