package helpers

import (
	"sort"
	"strings"

	"github.com/volatiletech/null/v8"
)

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

func SortBy[T any](slice []T, less func(i, j T) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return less(slice[i], slice[j])
	})
}

func JoinStrings(wrappedStrings []string) string {
	return strings.Join(wrappedStrings, ",")
}
