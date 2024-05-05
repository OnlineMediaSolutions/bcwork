package filter

import (
	"github.com/lib/pq"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type StringArrayDynFilter map[string][]string

func (filter StringArrayDynFilter) AndIn() qmods.QueryModsSlice {

	res := make(qmods.QueryModsSlice, 0)
	for col, vals := range filter {
		res = append(res, qm.And(col+" = ANY(?)", pq.Array(vals)))

	}

	return res
}
