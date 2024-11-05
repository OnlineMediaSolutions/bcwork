package history

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
)

type item struct {
	key         string
	publisherID *string
	domain      *string
	entityID    *string
}

func getItem(subject string, value any) (item, error) {
	switch subject {
	case globalFactorSubject:
		return getGlobalFactorItem(value)
	case userSubject:
		return getUserItem(value)
	case dpoSubject:
		// all rule dimensions joined by _ // TODO:
		return item{}, nil
	case publisherSubject:
		return getPublisherItem(value)
	case blockPublisherSubject:
		// Blocks - {pub name} ({pub id}) // TODO:
		return item{}, nil
	case confiantPublisherSubject:
		// Confiant - {pub name} ({pub id}) // TODO:
		return item{}, nil
	case pixalatePublisherSubject:
		// Pixalate - {pub name} ({pub id}) // TODO:
		return item{}, nil
	case domainSubject:
		return getDomainItem(value)
	case blockDomainSubject:
		// Blocks - {domain} ({pub id}) // TODO:
		return item{}, nil
	case confiantDomainSubject:
		// Confiant - {domain} ({pub id}) // TODO:
		return item{}, nil
	case pixalateDomainSubject:
		// Pixalate - {domain} ({pub id}) // TODO:
		return item{}, nil
	case factorSubject:
		// all rule dimensions joined by _ // TODO:
		return item{}, nil
	case jsTargetingSubject:
		return getJSTargetingItem(value)
	case floorSubject:
		// all rule dimensions joined by _ // TODO:
		return item{}, nil
	case factorAutomationSubject:
		return getDomainItem(value)
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
		key = consultantFeeKey
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

func getJSTargetingItem(value any) (item, error) {
	targeting, ok := value.(*models.Targeting)
	if !ok {
		return item{}, errors.New("cannot cast value to targeting")
	}
	return item{
		key: fmt.Sprintf(
			"%v_%v_%v_%v_%v_%v_%v_%v_%v_%v",
			targeting.Country, targeting.UnitSize, targeting.DeviceType,
			targeting.Os, targeting.Browser, targeting.PlacementType,
			targeting.PriceModel, targeting.Value, targeting.DailyCap,
			targeting.KV.JSON,
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
