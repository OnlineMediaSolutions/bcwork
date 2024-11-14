package helpers

import (
	"math"
	"reflect"
	"strings"

	"github.com/volatiletech/null/v8"
)

func ReplaceWildcardValues(data interface{}) {
	val := reflect.ValueOf(data).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.CanSet() && field.Kind() == reflect.String {
			if field.String() == "all" {
				field.SetString("")
			}
		}
	}
}

func GetStringOrEmpty(ns null.String) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func GetStringWithDefaultValue(str, defaultValue string) string {
	if str == "" {
		return defaultValue
	}
	return str
}

func GetStringFromSliceWithDefaultValue(elems []string, sep, defaultValue string) string {
	if len(elems) == 0 {
		return defaultValue
	}

	return "(" + strings.Join(elems, sep) + ")"
}

func GetNullString(s string) null.String {
	if s == "" {
		return null.NewString("", false)
	}
	return null.StringFrom(s)
}

func RoundFloat(value float64) float64 {
	return math.Round(value*100) / 100
}
