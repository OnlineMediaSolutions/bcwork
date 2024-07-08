package utils

import (
	"reflect"
)

func ConvertingAllValues(data interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.Kind() == reflect.String {
				if field.String() == "all" {
					field.SetString("")
				}
			}
		}
	}
}
