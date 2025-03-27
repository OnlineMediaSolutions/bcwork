package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/m6yf/bcwork/utils/helpers"
)

type FilterService struct{}

func NewFilterService() *FilterService {
	return &FilterService{}
}

func (f *FilterService) GetFilterFields(ctx context.Context, filterName string) ([]string, error) {
	filter, err := getFilter(filterName)
	if err != nil {
		return nil, err
	}

	filterFields, err := extractFilterFieldsNames(filter.Type())
	if err != nil {
		return nil, err
	}

	return filterFields, nil
}

func getFilter(filterName string) (reflect.Value, error) {
	var filter any
	switch filterName {
	case "ads_txt_main":
		filter = new(AdsTxtGetMainFilter)
	case "ads_txt_group_by_dp":
		filter = new(AdsTxtGetGroupByDPFilter)
	default:
		return reflect.Value{}, fmt.Errorf("unknown filter name [%v]", filterName)
	}

	value, err := helpers.GetStructReflectValue(filter)
	if err != nil {
		return reflect.Value{}, err
	}

	return value, nil
}

func extractFilterFieldsNames(t reflect.Type) ([]string, error) {
	var filters []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			innerFilters, err := extractFilterFieldsNames(field.Type)
			if err != nil {
				return nil, err
			}
			filters = append(filters, innerFilters...)
		} else {
			filters = append(filters, strings.Split(field.Tag.Get("json"), ",")[0])
		}
	}

	return filters, nil
}
