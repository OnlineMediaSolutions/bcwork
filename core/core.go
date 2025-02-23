package core

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/m6yf/bcwork/utils/helpers"
)

func getModelsColumnsToUpdate(oldData, newData any, blacklistColumns []string) ([]string, error) {
	const boilTagName = "boil"

	oldValueReflection, err := helpers.GetStructReflectValue(oldData)
	if err != nil {
		return nil, fmt.Errorf("cannot get reflection of old data: %w", err)
	}

	newValueReflection, err := helpers.GetStructReflectValue(newData)
	if err != nil {
		return nil, fmt.Errorf("cannot get reflection of new data: %w", err)
	}

	if oldValueReflection.Type().Name() != newValueReflection.Type().Name() {
		return nil, fmt.Errorf(
			"provided different structs: old [%v], new [%v]",
			oldValueReflection.Type().Name(), newValueReflection.Type().Name(),
		)
	}

	columns := make([]string, 0, oldValueReflection.NumField())

	for i := 0; i < oldValueReflection.NumField(); i++ {
		field := oldValueReflection.Type().Field(i)
		property := strings.Split(field.Tag.Get(boilTagName), ",")[0]
		oldFieldValue := oldValueReflection.Field(i)
		newFieldValue := newValueReflection.Field(i)

		if !reflect.DeepEqual(oldFieldValue.Interface(), newFieldValue.Interface()) &&
			!slices.Contains(blacklistColumns, property) {
			columns = append(columns, property)
		}
	}

	return columns, nil
}
