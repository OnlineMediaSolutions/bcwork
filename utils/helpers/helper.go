package helpers

import (
	"github.com/volatiletech/null/v8"
	"reflect"
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
