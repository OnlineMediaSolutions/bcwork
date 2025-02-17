package qmods

import (
	"github.com/m6yf/bcwork/bcdb/fields"
	"github.com/m6yf/bcwork/bcdb/join"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type QueryModsSlice []qm.QueryMod

func (q QueryModsSlice) Add(mods ...qm.QueryMod) QueryModsSlice {
	if q == nil {
		q = QueryModsSlice{}
	}
	for _, mod := range mods {
		q = append(q, mod)
	}

	return q
}

func (q QueryModsSlice) AddArray(mods []qm.QueryMod) QueryModsSlice {
	return q.Add(mods...)
}

func (q QueryModsSlice) Paginate(pg *pagination.Pagination) QueryModsSlice {
	return q.AddArray(pg.Do())
}

func (q QueryModsSlice) Fields(reg fields.Registry, selector string) QueryModsSlice {
	return q.AddArray(reg.Select(selector))
}

func (q QueryModsSlice) Order(sort order.Sort, custom order.CustomSort, primaryKey string) QueryModsSlice {
	return q.Add(sort.Do(custom, primaryKey))
}

func (q QueryModsSlice) Join(left string, right string, column string) QueryModsSlice {
	return q.Add(join.Inner(left, right, column))
}
