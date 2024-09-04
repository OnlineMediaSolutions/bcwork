package validations

import (
	"github.com/go-playground/validator/v10"
	"github.com/m6yf/bcwork/utils/constant"
	"slices"
)

var Validator = validator.New()
var integrationTypes = []string{"JS Tags (Compass)", "JS Tags (NP)", "Prebid.js", "Prebid Server", "oRTB EP"}
var globalFactorKeyTypes = []string{"NP Tech Fee", "Consultant Fee", "Amazon TAM Fee"}

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
	err = Validator.RegisterValidation("globalFactorKey", globalFactorKeyValidation)
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
	if len(device) == 0 || device == "all" {
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

func globalFactorKeyValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	return slices.Contains(globalFactorKeyTypes, field.String())
}
