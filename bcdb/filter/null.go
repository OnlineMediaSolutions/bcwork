package filter

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type NullFilter bool

func (filter NullFilter) And(column string) qm.QueryMod {

	if filter {
		return qm.And(column + " IS NULL")
	} else {
		return qm.And(column + " IS NOT NULL")
	}
}

func (filter NullFilter) Where(column string) qm.QueryMod {
	if filter {
		return qm.Where(column + " IS NULL")
	} else {
		return qm.Where(column + " IS NOT NULL")
	}
}

func (filter NullFilter) Or(column string) qm.QueryMod {

	if filter {
		return qm.Or(column + " IS NULL")
	} else {
		return qm.Or(column + " IS NOT NULL")
	}
}

func NullFilterPtr(f bool) *NullFilter {
	res := NullFilter(f)
	return &res
}

type NotNullFilter bool

func (filter NotNullFilter) And(column string) qm.QueryMod {

	if filter {
		return qm.And(column + " IS NOT NULL")
	} else {
		return qm.And(column + " IS NULL")
	}
}

func (filter NotNullFilter) Where(column string) qm.QueryMod {
	if filter {
		return qm.Where(column + " IS NOT NULL")
	} else {
		return qm.Where(column + " IS NULL")
	}
}

func (filter NotNullFilter) Or(column string) qm.QueryMod {

	if filter {
		return qm.Or(column + " IS NOT NULL")
	} else {
		return qm.Or(column + " IS NULL")
	}
}

func NotNullFilterPtr(f bool) *NotNullFilter {
	res := NotNullFilter(f)
	return &res
}
