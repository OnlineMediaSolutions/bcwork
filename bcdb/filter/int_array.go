package filter

import (
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type IntArrayFilter []int

func (filter IntArrayFilter) AndIn(column string) qm.QueryMod {
	return qm.AndIn(column+" = ANY (?)", pq.Array([]int(filter)))
}

func (filter IntArrayFilter) WhereIn(column string) qm.QueryMod {
	return qm.WhereIn(column+" = ANY (?)", pq.Array([]int(filter)))
}

func (filter IntArrayFilter) OrIn(column string) qm.QueryMod {
	return qm.OrIn(column+" = ANY (?)", pq.Array([]int(filter)))
}
