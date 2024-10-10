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
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Competitor struct {
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`
	Url  string `boil:"url" json:"url" toml:"url" yaml:"url"`
}

type CompetitorSlice []*Competitor

type GetCompetitorOptions struct {
	Filter     CompetitorFilter       `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type CompetitorFilter struct {
	Name filter.StringArrayFilter `json:"name,omitempty"`
	Url  filter.StringArrayFilter `json:"url,omitempty"`
}

func (competitor *Competitor) FromModel(mod *models.Competitor) error {

	competitor.Name = mod.Name
	competitor.Url = mod.URL
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
	if len(filter.Url) > 0 {
		mods = append(mods, filter.Url.AndIn(models.CompetitorColumns.URL))
	}
	return mods
}

func UpdateCompetitor(c *fiber.Ctx, data []constant.CompetitorUpdateRequest) error {

	for _, competitor := range data {
		modConf := models.Competitor{
			Name: competitor.Name,
			URL:  competitor.URL,
		}
		err := modConf.Upsert(c.Context(), bcdb.DB(), true, []string{models.CompetitorColumns.Name}, boil.Infer(), boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
}
