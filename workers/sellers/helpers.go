package sellers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

func FetchDataFromWebsite(url string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request for getting sellers: %w", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if sellers, ok := data["sellers"]; ok {
		if _, err = CheckSellersArray(sellers); err != nil {
			return nil, fmt.Errorf("invalid sellers format: %w", err)
		}
	} else {
		return nil, fmt.Errorf("sellers array not found in the response")
	}

	return data, nil
}

func CheckSellersArray(sellers interface{}) ([]interface{}, error) {
	if sellersArray, ok := sellers.([]interface{}); ok {
		return sellersArray, nil
	}

	return nil, fmt.Errorf("sellers should be an array, but got %T", sellers)
}

func TrimSellerIdByDemand(mappingValue, sellerId string) string {
	var prefixBuilder, suffixBuilder strings.Builder
	isPrefix := true

	for _, char := range mappingValue {
		if unicode.IsDigit(char) {
			if isPrefix {
				prefixBuilder.WriteRune(char)
			} else {
				suffixBuilder.WriteRune(char)
			}
		} else {
			isPrefix = false // Stop collecting prefix numbers when we hit a non-digit character
		}
	}

	numericPrefix := prefixBuilder.String()
	numericSuffix := suffixBuilder.String()

	// Trim numeric prefix if exists
	if numericPrefix != "" {
		sellerId = strings.TrimPrefix(sellerId, numericPrefix)
	}

	// Trim numeric suffix if exists
	if numericSuffix != "" {
		sellerId = strings.TrimSuffix(sellerId, numericSuffix)
	}

	return sellerId
}

func GetDPMap() map[string]string {
	return map[string]string{
		"GetMedia":   "12XXXXX",
		"Onomagic":   "XXXXX1",
		"Limpid":     "9XXXXX",
		"Audienciad": "XXXXX2",
		"OMS":        "XXXXXX",
	}
}
