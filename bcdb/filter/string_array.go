package filter

import (
	"github.com/lib/pq"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type StringArrayFilter []string

func (filter StringArrayFilter) AndIn(column string) qm.QueryMod {

	return qm.And(column+" = ANY(?)", pq.Array([]string(filter)))
}

func (filter StringArrayFilter) WhereIn(column string) qm.QueryMod {

	return qm.WhereIn(column+" = ANY(?)", pq.Array([]string(filter)))
}

func (filter StringArrayFilter) OrIn(column string) qm.QueryMod {

	return qm.OrIn(column+" = ANY(?)", pq.Array([]string(filter)))
}

func (filter StringArrayFilter) AndNotIn(column string) qm.QueryMod {

	return qm.AndIn("NOT "+column+" = ANY(?)", []string(filter))
}

func (filter StringArrayFilter) WhereNotIn(column string) qm.QueryMod {

	return qm.WhereIn("NOT "+column+" = ANY(?)", []string(filter))
}

func (filter StringArrayFilter) OrNotIn(column string) qm.QueryMod {

	return qm.WhereIn("NOT "+column+" = ANY(?)", []string(filter))
}

type String2DArrayFilter []string

func (filter String2DArrayFilter) AndIn(column string) qm.QueryMod {
	mods := make([]qm.QueryMod, 0)

	for _, f := range filter {
		mods = append(mods, qm.Or("?  =  ANY("+column+")", f))
	}

	return qm.Expr(mods...)
}
