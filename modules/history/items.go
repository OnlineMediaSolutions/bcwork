package history

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/m6yf/bcwork/utils/helpers"
)

const (
	dimensionSeparator      = "_"
	multipleValuesSeparator = "-"
)

type item struct {
	key             string
	publisherID     *string
	domain          *string
	demandPartnerID *string
	entityID        *string
}

func getItem(subject string, value any) (item, error) {
	switch subject {
	case GlobalFactorSubject:
		return getGlobalFactorItem(value)
	case UserSubject:
		return getUserItem(value)
	case DPOSubject:
		return getDPOItem(value)
	case PublisherSubject:
		return getPublisherItem(value)
	case BlockPublisherSubject:
		return getBlocksPublisherItem(value)
	case ConfiantPublisherSubject:
		return getConfiantPublisherItem(value)
	case PixalatePublisherSubject:
		return getPixalatePublisherItem(value)
	case DomainSubject:
		return getDomainItem(value)
	case BlockDomainSubject:
		return getBlocksDomainItem(value)
	case ConfiantDomainSubject:
		return getConfiantDomainItem(value)
	case PixalateDomainSubject:
		return getPixalateDomainItem(value)
	case FactorSubject:
		return getFactorItem(value)
	case JSTargetingSubject:
		return getJSTargetingItem(value)
	case FloorSubject:
		return getFloorItem(value)
	case FactorAutomationSubject:
		return getDomainItem(value)
	case BidCachingSubject:
		return getBidCachingItem(value)
	case BidCachingDomainSubject:
		return getBidCachingDomainItem(value)
	case RefreshCacheSubject:
		return getRefreshCacheItem(value)
	case RefreshCacheDomainSubject:
		return getRefreshCacheDomainItem(value)
	default:
		return item{}, errors.New("unknown item")
	}
}

func getPublisherItem(value any) (item, error) {
	publisher, ok := value.(*models.Publisher)
	if !ok {
		return item{}, errors.New("cannot cast value to publisher")
	}

	return item{
		key:         publisher.PublisherID,
		publisherID: helpers.GetPointerToString(publisher.PublisherID),
	}, nil
}

func getUserItem(value any) (item, error) {
	user, ok := value.(*models.User)
	if !ok {
		return item{}, errors.New("cannot cast value to user")
	}

	return item{
		key:      user.FirstName + " " + user.LastName,
		entityID: helpers.GetPointerToString(strconv.Itoa(user.ID)),
	}, nil
}

func getBlocksDomainItem(value any) (item, error) {
	block, ok := value.(*dto.BlockUpdateRequest)
	if !ok {
		return item{}, errors.New("cannot cast value to block")
	}

	return item{
		key:         "Blocks - " + block.Domain + " (" + block.Publisher + ")",
		publisherID: helpers.GetPointerToString(block.Publisher),
		domain:      helpers.GetPointerToString(block.Domain),
	}, nil
}

func getBlocksPublisherItem(value any) (item, error) {
	block, ok := value.(*dto.BlockUpdateRequest)
	if !ok {
		return item{}, errors.New("cannot cast value to block")
	}

	return item{
		key:         "Blocks - " + block.Publisher,
		publisherID: helpers.GetPointerToString(block.Publisher),
	}, nil
}

func getConfiantDomainItem(value any) (item, error) {
	confiant, ok := value.(*models.Confiant)
	if !ok {
		return item{}, errors.New("cannot cast value to confiant")
	}

	return item{
		key:         "Confiant - " + confiant.Domain + " (" + confiant.PublisherID + ")",
		publisherID: helpers.GetPointerToString(confiant.PublisherID),
		domain:      helpers.GetPointerToString(confiant.Domain),
	}, nil
}

func getConfiantPublisherItem(value any) (item, error) {
	confiant, ok := value.(*models.Confiant)
	if !ok {
		return item{}, errors.New("cannot cast value to confiant")
	}

	return item{
		key:         "Confiant - " + confiant.PublisherID,
		publisherID: helpers.GetPointerToString(confiant.PublisherID),
	}, nil
}

func getPixalateDomainItem(value any) (item, error) {
	pixalate, ok := value.(*models.Pixalate)
	if !ok {
		return item{}, errors.New("cannot cast value to pixalate")
	}

	return item{
		key:         "Pixalate - " + pixalate.Domain + " (" + pixalate.PublisherID + ")",
		publisherID: helpers.GetPointerToString(pixalate.PublisherID),
		domain:      helpers.GetPointerToString(pixalate.Domain),
	}, nil
}

func getPixalatePublisherItem(value any) (item, error) {
	pixalate, ok := value.(*models.Pixalate)
	if !ok {
		return item{}, errors.New("cannot cast value to pixalate")
	}

	return item{
		key:         "Pixalate - " + pixalate.PublisherID,
		publisherID: helpers.GetPointerToString(pixalate.PublisherID),
	}, nil
}

func getDomainItem(value any) (item, error) {
	domain, ok := value.(*models.PublisherDomain)
	if !ok {
		return item{}, errors.New("cannot cast value to domain")
	}

	return item{
		key:         domain.Domain + " (" + domain.PublisherID + ")",
		publisherID: helpers.GetPointerToString(domain.PublisherID),
		domain:      helpers.GetPointerToString(domain.Domain),
	}, nil
}

func getGlobalFactorItem(value any) (item, error) {
	const (
		techFeeKey       = "Tech Fee"
		tamFeeKey        = "Amazon TAM Fee"
		consultantFeeKey = "Consultant Fee"
	)

	globalFactor, ok := value.(*models.GlobalFactor)
	if !ok {
		return item{}, errors.New("cannot cast value to global factor")
	}

	var key string
	switch globalFactor.Key {
	case constant.GlobalFactorTechFeeType:
		key = techFeeKey
	case constant.GlobalFactorTAMFeeType:
		key = tamFeeKey
	case constant.GlobalFactorConsultantFeeType:
		key = consultantFeeKey + " - " + globalFactor.PublisherID
	}

	if key == "" {
		return item{}, fmt.Errorf("cannot get key from global factor fee type [%v]", globalFactor.Key)
	}

	return item{
		key: key,
		publisherID: func(key string) *string {
			if globalFactor.Key == constant.GlobalFactorConsultantFeeType {
				return helpers.GetPointerToString(globalFactor.PublisherID)
			}

			return nil
		}(key),
	}, nil
}

func getDPOItem(value any) (item, error) {
	dpo, ok := value.(*models.DpoRule)
	if !ok {
		return item{}, errors.New("cannot cast value to dpo rule")
	}

	return item{
		key: dpo.DemandPartnerID + dimensionSeparator +
			getDimensionString(
				dpo.Publisher.String,
				dpo.Domain.String,
				dpo.Country.String,
				"",
				dpo.DeviceType.String,
				dpo.Os.String,
				dpo.Browser.String,
				dpo.PlacementType.String,
			),
		publisherID:     dpo.Publisher.Ptr(),
		domain:          dpo.Domain.Ptr(),
		demandPartnerID: helpers.GetPointerToString(dpo.DemandPartnerID),
		entityID:        helpers.GetPointerToString(dpo.RuleID),
	}, nil
}

func getFactorItem(value any) (item, error) {
	factor, ok := value.(*models.Factor)
	if !ok {
		return item{}, errors.New("cannot cast value to factor")
	}

	return item{
		key: getDimensionString(
			factor.Publisher,
			factor.Domain,
			factor.Country.String,
			"",
			factor.Device.String,
			factor.Os.String,
			factor.Browser.String,
			factor.PlacementType.String,
		),
		publisherID: helpers.GetPointerToString(factor.Publisher),
		domain:      helpers.GetPointerToString(factor.Domain),
		entityID:    helpers.GetPointerToString(factor.RuleID),
	}, nil
}

func getFloorItem(value any) (item, error) {
	floor, ok := value.(*models.Floor)
	if !ok {
		return item{}, errors.New("cannot cast value to floor")
	}

	return item{
		key: getDimensionString(
			floor.Publisher,
			floor.Domain,
			floor.Country.String,
			"",
			floor.Device.String,
			floor.Os.String,
			floor.Browser.String,
			floor.PlacementType.String,
		),
		publisherID: helpers.GetPointerToString(floor.Publisher),
		domain:      helpers.GetPointerToString(floor.Domain),
		entityID:    helpers.GetPointerToString(floor.RuleID),
	}, nil
}

func getJSTargetingItem(value any) (item, error) {
	targeting, ok := value.(*models.Targeting)
	if !ok {
		return item{}, errors.New("cannot cast value to targeting")
	}

	return item{
		key: getDimensionString(
			targeting.PublisherID,
			targeting.Domain,
			strings.Join(targeting.Country, multipleValuesSeparator),
			targeting.UnitSize,
			strings.Join(targeting.DeviceType, multipleValuesSeparator),
			strings.Join(targeting.Os, multipleValuesSeparator),
			strings.Join(targeting.Browser, multipleValuesSeparator),
			targeting.PlacementType.String,
		),
		publisherID: helpers.GetPointerToString(targeting.PublisherID),
		domain:      helpers.GetPointerToString(targeting.Domain),
		entityID:    helpers.GetPointerToString(strconv.Itoa(targeting.ID)),
	}, nil
}

func getBidCachingItem(value any) (item, error) {
	bc, ok := value.(*models.BidCaching)
	if !ok {
		return item{}, errors.New("cannot cast value to Bid Caching")
	}

	return item{
		key: getDimensionString(
			bc.Publisher,
			bc.Domain.String,
			bc.Country.String,
			"",
			bc.Device.String,
			bc.Os.String,
			bc.Browser.String,
			bc.PlacementType.String,
		),
		publisherID: helpers.GetPointerToString(bc.Publisher),
		entityID:    helpers.GetPointerToString(bc.RuleID),
	}, nil
}

func getBidCachingDomainItem(value any) (item, error) {
	bc, ok := value.(*models.BidCaching)
	if !ok {
		return item{}, errors.New("cannot cast value to bid caching")
	}

	return item{
		key: getDimensionString(
			bc.Publisher,
			bc.Domain.String,
			bc.Country.String,
			"",
			bc.Device.String,
			bc.Os.String,
			bc.Browser.String,
			bc.PlacementType.String,
		),
		publisherID: helpers.GetPointerToString(bc.Publisher),
		domain:      bc.Domain.Ptr(),
		entityID:    helpers.GetPointerToString(bc.RuleID),
	}, nil
}

func getRefreshCacheItem(value any) (item, error) {
	rc, ok := value.(*models.RefreshCache)
	if !ok {
		return item{}, errors.New("cannot cast value to Refresh Cache")
	}

	return item{
		key: getDimensionString(
			rc.Publisher,
			rc.Domain.String,
			rc.Country.String,
			"",
			rc.Device.String,
			rc.Os.String,
			rc.Browser.String,
			rc.PlacementType.String,
		),
		publisherID: helpers.GetPointerToString(rc.Publisher),
		entityID:    helpers.GetPointerToString(rc.RuleID),
	}, nil
}

func getRefreshCacheDomainItem(value any) (item, error) {
	rc, ok := value.(*models.RefreshCache)
	if !ok {
		return item{}, errors.New("cannot cast value to refresh cache")
	}

	return item{
		key: getDimensionString(
			rc.Publisher,
			rc.Domain.String,
			rc.Country.String,
			"",
			rc.Device.String,
			rc.Os.String,
			rc.Browser.String,
			rc.PlacementType.String,
		),
		publisherID: helpers.GetPointerToString(rc.Publisher),
		domain:      rc.Domain.Ptr(),
		entityID:    helpers.GetPointerToString(rc.RuleID),
	}, nil
}

func getDimensionString(publisherID, domain, country, unitSize, device, os, browser, placementType string) string {
	return getDimensionValue(publisherID) + dimensionSeparator +
		getDimensionValue(domain) + dimensionSeparator +
		getDimensionValue(country) + dimensionSeparator +
		func() string {
			if unitSize != "" {
				return unitSize + dimensionSeparator
			}

			return ""
		}() +
		getDimensionValue(device) + dimensionSeparator +
		getDimensionValue(os) + dimensionSeparator +
		getDimensionValue(browser) + dimensionSeparator +
		getDimensionValue(placementType)
}

func getDimensionValue(value string) string {
	if value == "" {
		return "all"
	}

	return value
}
