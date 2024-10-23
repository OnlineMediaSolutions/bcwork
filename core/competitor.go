package core

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type BaseCompetitor struct {
	Name     string `json:"name" validate:"required"`
	URL      string `json:"url" validate:"required,url"`
	Type     string `json:"type" validate:"required"`
	Position int8   `json:"position" validate:"required"`
}

type CompetitorUpdateRequest struct {
	BaseCompetitor
}

type Competitor struct {
	BaseCompetitor
}

type CompetitorSlice []*Competitor

type GetCompetitorOptions struct {
	Filter     CompetitorFilter       `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type CompetitorFilter struct {
	Name     filter.StringArrayFilter `json:"name,omitempty"`
	URL      filter.StringArrayFilter `json:"url,omitempty"`
	Type     filter.StringArrayFilter `json:"type,omitempty"`
	Position filter.StringArrayFilter `json:"position,omitempty"`
}

func (competitor *Competitor) FromModel(mod *models.Competitor) error {
	competitor.Name = mod.Name
	competitor.URL = mod.URL
	competitor.Type = mod.Type
	competitor.Position = int8(mod.Position)
	return nil
}

func (cs *CompetitorSlice) FromModel(slice models.CompetitorSlice) error {

	for _, mod := range slice {
		c := Competitor{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func GetCompetitors(ctx context.Context, ops *GetCompetitorOptions) (CompetitorSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.CompetitorColumns.Name).AddArray(ops.Pagination.Do())

	qmods = qmods.Add(qm.Select("DISTINCT *"))

	mods, err := models.Competitors(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve competitors")
	}

	res := make(CompetitorSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (filter *CompetitorFilter) QueryMod() qmods.QueryModsSlice {

	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Name) > 0 {
		mods = append(mods, filter.Name.AndIn(models.CompetitorColumns.Name))
	}
	if len(filter.URL) > 0 {
		mods = append(mods, filter.URL.AndIn(models.CompetitorColumns.URL))
	}
	return mods
}

func UpdateCompetitor(c *fiber.Ctx, data *CompetitorUpdateRequest) error {

	modConf := models.Competitor{
		Name:     data.Name,
		URL:      data.URL,
		Type:     data.Type,
		Position: int64(data.Position),
	}

	return modConf.Upsert(c.Context(), bcdb.DB(), true, []string{models.CompetitorColumns.Name}, boil.Infer(), boil.Infer())
}
