package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/m6yf/bcwork/utils/helpers"
)

const jsonTagName = "json"

type change struct {
	Property string `json:"property"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

func getChanges(action, subject string, oldValue, newValue any) ([]byte, error) {
	switch action {
	case createdAction:
		return getChangesForActionCreated(subject, newValue)
	case updatedAction:
		return getChangesForActionUpdated(subject, oldValue, newValue)
	case deletedAction:
		return getChangesForActionDeleted(subject, oldValue)
	}

	return nil, errors.New("cannot get changes from unknown action")
}

func getChangesForActionCreated(subject string, newValue any) ([]byte, error) {
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
		if !isInternalField(subject, property) && isValueContainsData(value) {
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

func getChangesForActionUpdated(subject string, oldValue, newValue any) ([]byte, error) {
	if oldValue == nil {
		return nil, errors.New("old value is nil")
	}

	oldValueReflection, err := helpers.GetStructReflectValue(oldValue)
	if err != nil {
		return nil, fmt.Errorf("cannot get reflection of old value: %w", err)
	}

	newValueReflection, err := helpers.GetStructReflectValue(newValue)
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
			!isInternalField(subject, property) {
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

func getChangesForActionDeleted(subject string, oldValue any) ([]byte, error) {
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
		if !isInternalField(subject, property) && isValueContainsData(value) {
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

func isInternalField(subject, property string) bool {
	return slices.Contains(
		getInternalFields(subject),
		property,
	)
}

func getInternalFields(subject string) []string {
	internalFields := make([]string, 0, 12)
	internalFields = append(internalFields, getCommonInternalFields()...)
	internalFields = append(internalFields, getInternalFieldsBasedOnSubject(subject)...)

	return internalFields
}

func getCommonInternalFields() []string {
	const (
		idFieldJsonName          = "id"
		ruleIDFieldJsonName      = "rule_id"
		createdAtFieldJsonName   = "created_at"
		updatedAtFieldJsonName   = "updated_at"
		disabledAtFieldJsonName  = "disabled_at"
		nonImportedFieldJsonName = "-"
	)

	return []string{
		idFieldJsonName, ruleIDFieldJsonName, createdAtFieldJsonName,
		updatedAtFieldJsonName, disabledAtFieldJsonName, nonImportedFieldJsonName,
	}
}

func getInternalFieldsBasedOnSubject(subject string) []string {
	const (
		publisherIDFieldJsonName     = "publisher_id"
		publisherFieldJsonName       = "publisher"
		domainFieldJsonName          = "domain"
		keyFieldJsonName             = "key"
		activeFieldJsonName          = "active"
		userIDFieldJsonName          = "user_id"
		passwordChangedFieldJsonName = "password_changed"

		browserFieldJsonName         = "browser"
		countryFieldJsonName         = "country"
		demandPartnerIDFieldJsonName = "demand_partner_id"
		deviceTypeFieldJsonName      = "device_type"
		deviceFieldJsonName          = "device"
		osFieldJsonName              = "os"
		placementTypeFieldJsonName   = "placement_type"
	)

	switch subject {
	case GlobalFactorSubject:
		return []string{publisherIDFieldJsonName, keyFieldJsonName}
	case DPOSubject, FloorSubject, FactorAutomationSubject, FactorSubject,
		RefreshCacheSubject, RefreshCacheDomainSubject, BidCachingSubject, BidCachingDomainSubject:
		return []string{
			publisherIDFieldJsonName, publisherFieldJsonName, domainFieldJsonName,
			activeFieldJsonName, browserFieldJsonName, countryFieldJsonName,
			demandPartnerIDFieldJsonName, deviceTypeFieldJsonName, deviceFieldJsonName,
			osFieldJsonName, placementTypeFieldJsonName,
		}
	case UserSubject:
		return []string{userIDFieldJsonName, passwordChangedFieldJsonName}
	case JSTargetingSubject,
		BlockPublisherSubject, BlockDomainSubject,
		PixalatePublisherSubject, PixalateDomainSubject,
		ConfiantPublisherSubject, ConfiantDomainSubject:
		return []string{publisherIDFieldJsonName, publisherFieldJsonName, domainFieldJsonName}
	}

	return []string{}
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
