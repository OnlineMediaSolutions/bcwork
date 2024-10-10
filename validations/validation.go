package validations

import (
	"net/url"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/m6yf/bcwork/utils/constant"
)

const (
	// Keys
	targetingPriceModelValidationKey = "targetingPriceModel"
	targetingStatusValidationKey     = "targetingStatus"
	countriesValidationKey           = "countries"
	devicesValidationKey             = "devices"
	// Error messages
	countryValidationErrorMessage            = "country code must be 2 characters long and should be in the allowed list"
	deviceValidationErrorMessage             = "device should be in the allowed list"
	targetingCostModelValidationErrorMessage = "targeting price model should be 'CPM' or 'Rev Share'"
	targetingStatusValidationErrorMessage    = "targeting status should be 'Active', 'Paused' or 'Archived'"
)

var (
	Validator = validator.New()

	integrationTypes     = []string{"JS Tags (Compass)", "JS Tags (NP)", "Prebid.js", "Prebid Server", "oRTB EP"}
	globalFactorKeyTypes = []string{"tech_fee", "consultant_fee", "tam_fee"}
	targetingCostModels  = []string{constant.TargetingPriceModelCPM, constant.TargetingPriceModelRevShare}
	targetingStatuses    = []string{"", constant.TargetingStatusActive, constant.TargetingStatusPaused, constant.TargetingStatusArchived}
)

func init() {
	err := Validator.RegisterValidation("country", countryValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("device", deviceValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("floor", floorValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("factor", factorValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("factorDpo", factorDpoValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("rate", rateValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("active", activeValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("all", notAllowedTheWordAllValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("integrationType", integrationTypeValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("url", validateURL)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation("globalFactorKey", globalFactorKeyValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(targetingPriceModelValidationKey, targetingCostModelValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(targetingStatusValidationKey, targetingStatusValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(countriesValidationKey, countriesValidation)
	if err != nil {
		return
	}
	err = Validator.RegisterValidation(devicesValidationKey, devicesValidation)
	if err != nil {
		return
	}
}

func floorValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= 0
}

func factorValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= constant.MinFactorValue && val <= constant.MaxFactorValue
}

func factorDpoValidation(fieldLevel validator.FieldLevel) bool {
	val := fieldLevel.Field().Float()
	return val >= constant.MinDPOFactorValue && val <= constant.MaxDPOFactorValue
}

func rateValidation(fieldLevel validator.FieldLevel) bool {
	val := fieldLevel.Field().Float()
	return val >= constant.MinDPOFactorValue && val <= constant.MaxDPOFactorValue
}

func activeValidation(fieldLevel validator.FieldLevel) bool {
	if len(fieldLevel.Field().String()) == 0 {
		return false
	}
	active := fieldLevel.Field().Bool()
	if active != true && active != false {
		return false
	}
	return true
}

func countryValidation(fl validator.FieldLevel) bool {
	country := fl.Field().String()
	return validateCountry(country)
}

func countriesValidation(fl validator.FieldLevel) bool {
	countries, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	for _, country := range countries {
		isValid := validateCountry(country)
		if !isValid {
			return false
		}
	}

	return true
}

func validateCountry(country string) bool {
	if country == "all" || len(country) == 0 {
		return true
	}
	if len(country) != constant.MaxCountryCodeLength {
		return false
	}
	if _, ok := constant.AllowedCountries[country]; !ok {
		return false
	}
	return true
}

func deviceValidation(fl validator.FieldLevel) bool {
	device := fl.Field().String()
	return validateDevice(device)
}

func devicesValidation(fl validator.FieldLevel) bool {
	devices, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	for _, device := range devices {
		isValid := validateDevice(device)
		if !isValid {
			return false
		}
	}

	return true
}

func validateDevice(device string) bool {
	if device == "all" {
		return true
	}

	if _, ok := constant.AllowedDevices[device]; !ok {
		return false
	}
	return true
}

func notAllowedTheWordAllValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "all" {
		return false
	}
	return true
}

func integrationTypeValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	integTypes, ok := field.Interface().([]string)
	if !ok {
		return false
	}
	for _, integType := range integTypes {
		found := false
		for _, integrationType := range integrationTypes {
			if integType == integrationType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func validateURL(fl validator.FieldLevel) bool {
	u, err := url.ParseRequestURI(fl.Field().String())
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func globalFactorKeyValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(globalFactorKeyTypes, field.String())
}

func targetingCostModelValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(targetingCostModels, field.String())
}

func targetingStatusValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(targetingStatuses, field.String())
}
