package fields

import (
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Registry map[string][]qm.QueryMod

func (f Registry) Select(selector string) []qm.QueryMod {
	if len(f) == 0 {
		return nil
	}

	return f[selector]
}

func (f Registry) Validate(selector string) error {
	if len(f) == 0 {
		return nil
	}
	if _, found := f[selector]; !found {
		return fmt.Errorf("scheme not found for selector %s", selector)
	}

	return nil
}
