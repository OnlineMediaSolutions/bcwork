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
		return item{}, errors.New("not implemented")
	case publisherSubject:
		return getPublisherItem(value)
	case blockPublisherSubject:
		// Blocks - {pub name} ({pub id}) // TODO:
		return item{}, errors.New("not implemented")
	case confiantPublisherSubject:
		return getConfiantPublisherItem(value)
	case pixalatePublisherSubject:
		return getPixalatePublisherItem(value)
	case domainSubject:
		return getDomainItem(value)
	case blockDomainSubject:
		// Blocks - {domain} ({pub id}) // TODO:
		return item{}, errors.New("not implemented")
	case confiantDomainSubject:
		return getConfiantDomainItem(value)
	case pixalateDomainSubject:
		return getPixalateDomainItem(value)
	case factorSubject:
		// all rule dimensions joined by _ // TODO:
		return item{}, errors.New("not implemented")
	case jsTargetingSubject:
		return getJSTargetingItem(value)
	case floorSubject:
		// all rule dimensions joined by _ // TODO:
		return item{}, errors.New("not implemented")
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
