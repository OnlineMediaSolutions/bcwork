package validations

import (
	"github.com/go-playground/validator/v10"
	"github.com/m6yf/bcwork/utils/constant"
)

var Validator = validator.New()

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
}

func floorValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= 0
}

func factorValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= constant.MinFactorValue && val <= 20.0
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
