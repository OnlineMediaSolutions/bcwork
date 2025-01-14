package utils

import (
	"fmt"
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
)

const (
	FactorMetaDataKeyPrefix       = "price:factor:v2"
	FloorMetaDataKeyPrefix        = "price:floor:v2"
	DPOMetaDataKeyPrefix          = "dpo"
	JSTagMetaDataKeyPrefix        = "jstag"
	BidCachingMetaDataKeyPrefix   = "bid:cache"
	RefreshCacheMetaDataKeyPrefix = "refresh:cache"
	ConfiantMetaDataKeyPrefix     = "confiant:v2"
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

func GetFormulaRegex(country, domain, device, placement_type, os, browser, publisher string) string {
	if publisher == "all" || publisher == "" {
		publisher = ".*"
	}

	if domain == "" {
		domain = ".*"
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

	return fmt.Sprintf("(p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s)", publisher, domain, country, os, device, placement_type, browser)
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
