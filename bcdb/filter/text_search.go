package filter

import (
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TextSearchFilter string

func (filter TextSearchFilter) And(column string) qm.QueryMod {
	return qm.And("ts_rank(to_tsvector("+column+"),to_tsquery(?)) > 0.01", strings.Replace(string(filter), " ", "&", -1))
}
