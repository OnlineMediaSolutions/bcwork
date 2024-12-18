package utils

import (
	"fmt"
	"strings"
)

func CreateDeleteQuery(ids []string, deleteQuery string) string {
	var wrappedStrings []string
	for _, ruleId := range ids {
		wrappedStrings = append(wrappedStrings, fmt.Sprintf(`'%s'`, ruleId))
	}

	return fmt.Sprintf(deleteQuery, strings.Join(wrappedStrings, ","))
}
