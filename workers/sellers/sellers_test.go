package sellers

import (
	"encoding/json"
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
