package utils

import (
	"fmt"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"time"
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
	GetPlacementType() string
	GetOS() string
	GetBrowser() string
}

func CreateMetadataKey(data MetadataKey, prefix string) string {
	key := prefix + ":" + data.Publisher + ":" + data.Domain
	return key
}

func GetFormulaRegex(country, domain, device, placement_type, os, browser, publisher string, isConfiant bool) string {

	if isConfiant {
		return fmt.Sprintf("(d=%s)", domain)
	}

	if publisher == "all" || publisher == "" {
		publisher = ".*"
	}

	if country == "all" || country == "" {
		country = ".*"
	}

	if device == "all" || device == "" {
		device = ".*"
	} else if device != "mobile" {
		device = "desktop"
	}

	if placement_type == "" {
		placement_type = ".*"
	}

	if os == "" {
		os = ".*"
	}

	if browser == "" {
		browser = ".*"
	}

	return fmt.Sprintf("(p=%s__c=%s__d=%s__dt=%s__pt=%s__os=%s__b=%s,)", publisher, country, domain, device, placement_type, os, browser)
}

func GetMetadataObject(updateRequest UpdateRequest) MetadataKey {
	key := MetadataKey{
		Publisher: updateRequest.GetPublisher(),
		Domain:    updateRequest.GetDomain(),
		Device:    updateRequest.GetDevice(),
		Country:   updateRequest.GetCountry(),
	}
	return key
}

func CreateMetadataObject(updateRequest UpdateRequest, key string, b []byte) models.MetadataQueue {
	modMeta := models.MetadataQueue{
		TransactionID: bcguid.NewFromf(updateRequest.GetPublisher(), updateRequest.GetDomain(), time.Now()),
		Key:           key,
		Value:         b,
	}
	return modMeta
}
