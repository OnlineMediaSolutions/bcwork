package history

import (
	"errors"
	"fmt"
	"github.com/volatiletech/null/v8"
	"strconv"
	"strings"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
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
		key: publisher.PublisherID,
		publisherID: func() *string {
			s := publisher.PublisherID
			return &s
		}(),
	}, nil
}

func getUserItem(value any) (item, error) {
	user, ok := value.(*models.User)
	if !ok {
		return item{}, errors.New("cannot cast value to user")
	}
	return item{
		key: user.FirstName + " " + user.LastName,
		entityID: func() *string {
			s := strconv.Itoa(user.ID)
			return &s
		}(),
	}, nil
}

func getBlocksDomainItem(value any) (item, error) {
	block, ok := value.(*dto.BlockUpdateRequest)
	if !ok {
		return item{}, errors.New("cannot cast value to block")
	}
	return item{
		key: "Blocks - " + block.Domain + " (" + block.Publisher + ")",
		publisherID: func() *string {
			s := block.Publisher
			return &s
		}(),
		domain: func() *string {
			s := block.Domain
			return &s
		}(),
	}, nil
}

func getBlocksPublisherItem(value any) (item, error) {
	block, ok := value.(*dto.BlockUpdateRequest)
	if !ok {
		return item{}, errors.New("cannot cast value to block")
	}
	return item{
		key: "Blocks - " + block.Publisher,
		publisherID: func() *string {
			s := block.Publisher
			return &s
		}(),
	}, nil
}

func getConfiantDomainItem(value any) (item, error) {
	confiant, ok := value.(*models.Confiant)
	if !ok {
		return item{}, errors.New("cannot cast value to confiant")
	}
	return item{
		key: "Confiant - " + confiant.Domain + " (" + confiant.PublisherID + ")",
		publisherID: func() *string {
			s := confiant.PublisherID
			return &s
		}(),
		domain: func() *string {
			s := confiant.Domain
			return &s
		}(),
	}, nil
}

func getConfiantPublisherItem(value any) (item, error) {
	confiant, ok := value.(*models.Confiant)
	if !ok {
		return item{}, errors.New("cannot cast value to confiant")
	}
	return item{
		key: "Confiant - " + confiant.PublisherID,
		publisherID: func() *string {
			s := confiant.PublisherID
			return &s
		}(),
	}, nil
}

func getPixalateDomainItem(value any) (item, error) {
	pixalate, ok := value.(*models.Pixalate)
	if !ok {
		return item{}, errors.New("cannot cast value to pixalate")
	}
	return item{
		key: "Pixalate - " + pixalate.Domain + " (" + pixalate.PublisherID + ")",
		publisherID: func() *string {
			s := pixalate.PublisherID
			return &s
		}(),
		domain: func() *string {
			s := pixalate.Domain
			return &s
		}(),
	}, nil
}

func getPixalatePublisherItem(value any) (item, error) {
	pixalate, ok := value.(*models.Pixalate)
	if !ok {
		return item{}, errors.New("cannot cast value to pixalate")
	}
	return item{
		key: "Pixalate - " + pixalate.PublisherID,
		publisherID: func() *string {
			s := pixalate.PublisherID
			return &s
		}(),
	}, nil
}

func getDomainItem(value any) (item, error) {
	domain, ok := value.(*models.PublisherDomain)
	if !ok {
		return item{}, errors.New("cannot cast value to domain")
	}
	return item{
		key: domain.Domain + " (" + domain.PublisherID + ")",
		publisherID: func() *string {
			s := domain.PublisherID
			return &s
		}(),
		domain: func() *string {
			s := domain.Domain
			return &s
		}(),
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
				s := globalFactor.PublisherID
				return &s
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
		publisherID: func() *string {
			s := dpo.Publisher.String
			return &s
		}(),
		domain: func() *string {
			s := dpo.Domain.String
			return &s
		}(),
		demandPartnerID: func() *string {
			s := dpo.DemandPartnerID
			return &s
		}(),
		entityID: func() *string {
			s := dpo.RuleID
			return &s
		}(),
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
		publisherID: func() *string {
			s := factor.Publisher
			return &s
		}(),
		domain: func() *string {
			s := factor.Domain
			return &s
		}(),
		entityID: func() *string {
			s := factor.RuleID
			return &s
		}(),
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
		publisherID: func() *string {
			s := floor.Publisher
			return &s
		}(),
		domain: func() *string {
			s := floor.Domain
			return &s
		}(),
		entityID: func() *string {
			s := floor.RuleID
			return &s
		}(),
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
		publisherID: func() *string {
			s := targeting.PublisherID
			return &s
		}(),
		domain: func() *string {
			s := targeting.Domain
			return &s
		}(),
		entityID: func() *string {
			s := strconv.Itoa(targeting.ID)
			return &s
		}(),
	}, nil
}

func getBidCachingItem(value any) (item, error) {
	bc, ok := value.(*models.BidCaching)
	if !ok {
		return item{}, errors.New("cannot cast value to Bid Caching")
	}
	return item{
		key: "Max Creatives in Cache - " + bc.Publisher,
		publisherID: func() *string {
			s := bc.Publisher
			return &s
		}(),
		entityID: func() *string {
			s := bc.RuleID
			return &s
		}(),
	}, nil
}

func getBidCachingDomainItem(value any) (item, error) {
	bc, ok := value.(*models.BidCaching)
	if !ok {
		return item{}, errors.New("cannot cast value to bid caching")
	}
	return item{
		key: "Max Creatives in Cache -  (" + bc.Publisher + ")",
		publisherID: func() *string {
			s := bc.Publisher
			return &s
		}(),
		entityID: func() *string {
			s := bc.RuleID
			return &s
		}(),
	}, nil
}

func getRefreshCacheDomainItem(value any) (item, error) {
	rc, ok := value.(*models.RefreshCache)
	if !ok {
		return item{}, errors.New("cannot cast value to refresh cache")
	}

	domainValue := rc.Domain
	if !domainValue.Valid {
		domainValue = null.StringFrom("*")
	}

	domainStr := domainValue.String

	return item{
		key: "Max Client Refresh - " + domainStr + " (" + rc.Publisher + ")",
		publisherID: func() *string {
			s := rc.Publisher
			return &s
		}(),
		domain: func() *string {
			return &domainStr
		}(),
		entityID: func() *string {
			s := rc.RuleID
			return &s
		}(),
	}, nil
}

func getRefreshCacheItem(value any) (item, error) {
	rc, ok := value.(*models.RefreshCache)
	if !ok {
		return item{}, errors.New("cannot cast value to Refresh Cache")
	}

	domainValue := rc.Domain
	if !domainValue.Valid {
		domainValue = null.StringFrom("*")
	}

	domainStr := domainValue.String

	return item{
		key: "Max Client Refresh - " + rc.Publisher,
		publisherID: func() *string {
			s := rc.Publisher
			return &s
		}(),
		domain: func() *string {
			return &domainStr
		}(),
		entityID: func() *string {
			s := rc.RuleID
			return &s
		}(),
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
