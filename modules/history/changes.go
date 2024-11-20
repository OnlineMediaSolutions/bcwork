package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
)

const jsonTagName = "json"

type change struct {
	Property string `json:"property"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

func getChanges(action string, oldValue, newValue any) ([]byte, error) {
	switch action {
	case createdAction:
		return getChangesForActionCreated(newValue)
	case updatedAction:
		return getChangesForActionUpdated(oldValue, newValue)
	case deletedAction:
		return getChangesForActionDeleted(oldValue)
	}

	return nil, errors.New("cannot get changes from unknown action")
}

func getChangesForActionCreated(newValue any) ([]byte, error) {
	newValueData, err := json.Marshal(newValue)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal newValue: %w", err)
	}

	var newValueMap map[string]interface{}
	err = json.Unmarshal(newValueData, &newValueMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal newValue to map[string]interface{}: %w", err)
	}

	var changes []change
	for property, value := range newValueMap {
		if !isInternalField(property) && isValueContainsData(value) {
			changes = append(changes, change{
				Property: property,
				OldValue: nil,
				NewValue: value,
			})
		}
	}

	sort.SliceStable(changes, func(i, j int) bool { return changes[i].Property < changes[j].Property })

	data, err := json.Marshal(changes)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal changes: %w", err)
	}

	return data, nil
}

func getChangesForActionUpdated(oldValue, newValue any) ([]byte, error) {
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
			!isInternalField(property) {
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

	sort.SliceStable(changes, func(i, j int) bool { return changes[i].Property < changes[j].Property })

	data, err := json.Marshal(changes)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal changes: %w", err)
	}

	return data, nil
}

func getChangesForActionDeleted(oldValue any) ([]byte, error) {
	oldValueData, err := json.Marshal(oldValue)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal old value: %w", err)
	}

	var oldValueMap map[string]interface{}
	err = json.Unmarshal(oldValueData, &oldValueMap)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal old value to map[string]interface{}: %w", err)
	}

	var changes []change
	for property, value := range oldValueMap {
		if !isInternalField(property) && isValueContainsData(value) {
			changes = append(changes, change{
				Property: property,
				OldValue: value,
				NewValue: nil,
			})
		}
	}

	sort.SliceStable(changes, func(i, j int) bool { return changes[i].Property < changes[j].Property })

	data, err := json.Marshal(changes)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal changes: %w", err)
	}

	return data, nil
}

func isInternalField(property string) bool {
	const (
		idFieldJsonName          = "id"
		ruleIDFieldJsonName      = "rule_id"
		createdAtFieldJsonName   = "created_at"
		updatedAtFieldJsonName   = "updated_at"
		disabledAtFieldJsonName  = "disabled_at"
		nonImportedFieldJsonName = "-"
	)

	return slices.Contains(
		[]string{
			idFieldJsonName, ruleIDFieldJsonName, createdAtFieldJsonName,
			updatedAtFieldJsonName, disabledAtFieldJsonName, nonImportedFieldJsonName,
		},
		property,
	)
}

func isValueContainsData(value any) bool {
	if value == nil {
		return false
	}

	data, ok := value.(string)
	if ok && data == "" {
		return false
	}

	return true
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
