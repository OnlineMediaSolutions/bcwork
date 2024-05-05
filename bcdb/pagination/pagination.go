package pagination

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func (pg *Pagination) Do() []qm.QueryMod {

	if pg == nil || pg.PageSize == 0 {
		return []qm.QueryMod{}
	}

	if pg.Page < 1 {
		pg.Page = 1
	}

	return []qm.QueryMod{qm.Limit(pg.PageSize), qm.Offset((pg.Page - 1) * pg.PageSize)}

}
