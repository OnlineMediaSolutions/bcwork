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

func GetFormulaRegex(country, domain, device string) string {
	c := country
	if country == "all" {
		c = ".*"
	}

	d := domain
	if d == "" {
		d = ".*"
	}

	dt := device
	if dt == "all" {
		dt = ".*"
	}
	return fmt.Sprintf("(c=%s__d=%s__dt=%s)", c, d, dt)
}
