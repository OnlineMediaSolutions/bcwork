package pagination

import (
	"fmt"

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

func (pg *Pagination) DoV2(columnName string, amountOfArgs int) []qm.QueryMod {
	if pg == nil || pg.PageSize == 0 || columnName == "" {
		return []qm.QueryMod{}
	}

	if pg.Page < 1 {
		pg.Page = 1
	}

	start := (pg.Page - 1) * pg.PageSize
	end := start + pg.PageSize

	return []qm.QueryMod{qm.Where(
		fmt.Sprintf(`%v between $%v and $%v`, columnName, amountOfArgs+1, amountOfArgs+2),
		start+1, end,
	)}
}
