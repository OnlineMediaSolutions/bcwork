package filter

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type BoolFilter bool

func (filter BoolFilter) Where(column string) qm.QueryMod {
	if filter {
		return qm.Where(column + " = TRUE")
	}
	return qm.Where(column + " = FALSE")

}

func NewBoolFilter(value bool) *BoolFilter {
	bf := BoolFilter(value)
	return &bf
}
