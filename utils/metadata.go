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
	AdsTxtMetaDataKeyTemplate     = "demand:%s:adtxtv2"
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
	return prefix + ":" + data.Publisher + ":" + data.Domain
}

func GetFormulaRegex(country, domain, device, placement_type, os, browser, publisher string) string {
	const (
		allValueString = "all"
		allValueRegexp = ".*"
	)

	if publisher == allValueString || publisher == "" {
		publisher = allValueRegexp
	}

	if domain == "" {
		domain = allValueRegexp
	}

	if country == allValueString || country == "" {
		country = allValueRegexp
	}

	if device == allValueString || device == "" {
		device = allValueRegexp
	} else if device != "mobile" {
		device = "desktop"
	}

	if placement_type == "" {
		placement_type = allValueRegexp
	}

	if os == "" {
		os = allValueRegexp
	}

	if browser == "" {
		browser = allValueRegexp
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
