package filter

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type BoolFilter string

func (filter BoolFilter) And(column string) qm.QueryMod {

	if filter == "true" {
		return qm.And(column + " = TRUE")
	} else {
		return qm.And(column + " = FALSE OR " + column + " IS NULL")
	}
}

func (filter BoolFilter) Where(column string) qm.QueryMod {

	if filter == "true" {
		return qm.Where(column + " = TRUE")
	} else {
		return qm.Where(column + " = FALSE OR " + column + " IS NULL")
	}
}

func (filter BoolFilter) Or(column string) qm.QueryMod {

	if filter == "true" {
		return qm.Or(column + " = TRUE")
	} else {
		return qm.Or(column + " = FALSE OR " + column + " IS NULL")
	}
}
