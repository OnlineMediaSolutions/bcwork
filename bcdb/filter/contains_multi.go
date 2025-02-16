package filter

import (
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ContainsMultiFilter string

func (filter ContainsMultiFilter) And(columns ...string) qm.QueryMod {
	return qm.And("LOWER(CONCAT(" + strings.Join(columns, ",") + "))  LIKE '%" + strings.ToLower(string(filter)) + "%'")
}

func (filter ContainsMultiFilter) Where(columns ...string) qm.QueryMod {
	return qm.Where("LOWER(CONCAT(" + strings.Join(columns, ",") + "))  LIKE '%" + strings.ToLower(string(filter)) + "%'")
}

func (filter ContainsMultiFilter) Or(columns ...string) qm.QueryMod {
	return qm.Or("LOWER(CONCAT(" + strings.Join(columns, ",") + "))  LIKE '%" + strings.ToLower(string(filter)) + "%'")
}
