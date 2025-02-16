package bcguid

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gofrs/uuid"
)

func New() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func NewFrom(seed string) string {
	return uuid.NewV5(uuid.NamespaceDNS, strings.ToLower(seed)).String()
}

func NewFromCaseSensitive(seed string) string {
	return uuid.NewV5(uuid.NamespaceDNS, seed).String()
}

func NewFromf(seed ...interface{}) string {
	return uuid.NewV5(uuid.NamespaceDNS, strings.Replace(strings.ToLower(fmt.Sprint(seed...)), " ", "", -1)).String()
}

func NewFromSortedStrings(seed ...string) string {
	sort.Strings(seed)

	return NewFrom(strings.Join(seed, ""))
}
