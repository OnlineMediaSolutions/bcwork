package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

const (
	jsonTagName              = "json"
	updatedAtFieldJsonName   = "updated_at"
	nonImportedFieldJsonName = "-"
)

type change struct {
	Property string `json:"property"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

func getChanges(oldValue, newValue any) ([]byte, error) {
	if oldValue == nil {
		return nil, errors.New("old value is nil")
	}

	oldValueReflection, err := getValue(oldValue)
	if err != nil {
		return nil, fmt.Errorf("cannot get reflection of old value: %w", err)
	}

	newValueReflection, err := getValue(newValue)
	if err != nil {
		return nil, fmt.Errorf("cannot get reflection of new value: %w", err)
	}

	if oldValueReflection.Type().Name() != newValueReflection.Type().Name() {
		return nil, fmt.Errorf(
			"provided different structs: old [%v], new [%v]",
			oldValueReflection.Type().Name(), newValueReflection.Type().Name(),
		)
	}

	var changes []change
	for i := 0; i < oldValueReflection.NumField(); i++ {
		field := oldValueReflection.Type().Field(i)
		property := strings.Split(field.Tag.Get(jsonTagName), ",")[0]
		oldFieldValue := oldValueReflection.Field(i)
		newFieldValue := newValueReflection.Field(i)

		if !reflect.DeepEqual(oldFieldValue.Interface(), newFieldValue.Interface()) &&
			!slices.Contains([]string{updatedAtFieldJsonName, nonImportedFieldJsonName}, property) {
			changes = append(changes, change{
				Property: property,
				OldValue: oldFieldValue.Interface(),
				NewValue: newFieldValue.Interface(),
			})
		}
	}

	if len(changes) == 0 {
		return nil, errors.New("no changes found")
	}

	data, err := json.Marshal(changes)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal changes: %w", err)
	}

	return data, nil
}

func getValue(i any) (reflect.Value, error) {
	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("old value is not a struct")
	}

	return value, nil
}
