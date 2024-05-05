package filter

import (
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ContainsFilter string

func (filter ContainsFilter) And(column string) qm.QueryMod {

	return qm.And("LOWER(CAST (" + column + " AS TEXT)) LIKE '%" + strings.ToLower(string(filter)) + "%'")
}

func (filter ContainsFilter) Where(column string) qm.QueryMod {

	return qm.Where("LOWER(CAST (" + column + " AS TEXT))  LIKE '%" + strings.ToLower(string(filter)) + "%'")
}

func (filter ContainsFilter) Or(column string) qm.QueryMod {

	return qm.Or("LOWER(CAST (" + column + " AS TEXT))  LIKE '%" + strings.ToLower(string(filter)) + "%'")
}

type NotContainsFilter string

func (filter NotContainsFilter) And(column string) qm.QueryMod {

	return qm.And("LOWER(CAST (" + column + " AS TEXT)) NOT LIKE '%" + strings.ToLower(string(filter)) + "%'")
}

func (filter NotContainsFilter) Where(column string) qm.QueryMod {

	return qm.Where("LOWER(CAST (" + column + " AS TEXT))  NOT LIKE '%" + strings.ToLower(string(filter)) + "%'")
}

func (filter NotContainsFilter) Or(column string) qm.QueryMod {

	return qm.Or("LOWER(CAST (" + column + " AS TEXT))  NOT LIKE '%" + strings.ToLower(string(filter)) + "%'")
}
