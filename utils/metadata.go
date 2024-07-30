package utils

import (
	"fmt"
)

type MetadataKey struct {
	Publisher string `json:"publisher"`
	Domain    string `json:"domain"`
	Device    string `json:"device"`
	Country   string `json:"country"`
}

type UpdateRequest interface {
	GetPublisher() string
	GetDomain() string
	GetDevice() string
	GetCountry() string
}

func CreateMetadataKey(data MetadataKey, prefix string) string {
	key := prefix + ":" + data.Publisher
	if data.Country != "" && data.Country != "all" {
		key = key + ":" + data.Country
	}
	if data.Device == "mobile" {
		key = "mobile:" + key
	}
	return key
}

func GetFormulaRegex(country, domain, device string, isConfiant bool) string {
	if domain == "" {
		domain = ".*"
	}

	if isConfiant {
		return fmt.Sprintf("(d=%s)", domain)
	}

	if country == "all" || country == "" {
		country = ".*"
	}

	if device == "all" || device == "" {
		device = ".*"
	}

	return fmt.Sprintf("(c=%s__d=%s__dt=%s)", country, domain, device)
}

func GetMetadataKey(updateRequest UpdateRequest) MetadataKey {
	key := MetadataKey{
		Publisher: updateRequest.GetPublisher(),
		Domain:    updateRequest.GetDomain(),
		Device:    updateRequest.GetDevice(),
		Country:   updateRequest.GetCountry(),
	}
	return key
}
