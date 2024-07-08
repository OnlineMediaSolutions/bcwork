package utils

import (
	"reflect"
)

func ConvertingAllValues(data interface{}) {
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
