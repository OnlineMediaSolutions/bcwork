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
}

func floorValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= 0
}

func factorValidation(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= constant.MinFactorValue && val <= constant.MaxFactorValue
}

func countryValidation(fl validator.FieldLevel) bool {
	country := fl.Field().String()
	if country == "all" {
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
