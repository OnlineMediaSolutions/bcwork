package core

import (
	"context"
	"database/sql"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type GlobalFactorService struct {
	historyModule history.HistoryModule
}

func NewGlobalFactorService(historyModule history.HistoryModule) *GlobalFactorService {
	return &GlobalFactorService{
		historyModule: historyModule,
	}
}

type GetGlobalFactorOptions struct {
	Filter     GlobalFactorFilter     `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type GlobalFactorFilter struct {
	Key       filter.StringArrayFilter `json:"key"`
	Publisher filter.StringArrayFilter `json:"publisher_id"`
	Value     filter.StringArrayFilter `json:"value"`
}

type GlobalFactorRequest struct {
	Key       string  `json:"key" validate:"globalFactorKey"`
	Publisher string  `json:"publisher_id"`
	Value     float64 `json:"value" validate:"gte=0"`
}

type GlobalFactor struct {
	Key         string     `boil:"key" json:"key" toml:"key" yaml:"key"`
	PublisherID string     `boil:"publisher_id" json:"publisher_id,omitempty" toml:"publisher_id" yaml:"publisher_id"`
	Value       float64    `boil:"value" json:"value" toml:"value" yaml:"value"`
	CreatedAt   *time.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

type GlobalFactorSlice []*GlobalFactor

func (g *GlobalFactorService) GetGlobalFactor(ctx context.Context, ops *GetGlobalFactorOptions) (GlobalFactorSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.GlobalFactorColumns.Key).AddArray(ops.Pagination.Do())

	if ops.Selector == "id" {
		qmods = qmods.Add(qm.Select("DISTINCT " + models.GlobalFactorColumns.Key))
	} else {
		qmods = qmods.Add(qm.Select("DISTINCT *"))

	}
	mods, err := models.GlobalFactors(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve Global Factors")
	}

	res := make(GlobalFactorSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (cs *GlobalFactorSlice) FromModel(slice models.GlobalFactorSlice) error {

	for _, mod := range slice {
		c := GlobalFactor{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (globalFactor *GlobalFactor) FromModel(mod *models.GlobalFactor) error {
	globalFactor.PublisherID = mod.PublisherID
	globalFactor.CreatedAt = &mod.CreatedAt.Time
	globalFactor.UpdatedAt = mod.UpdatedAt.Ptr()
	globalFactor.Key = mod.Key
	globalFactor.Value = mod.Value.Float64

	return nil
}

func (g *GlobalFactorService) UpdateGlobalFactor(ctx context.Context, data *GlobalFactorRequest) error {
	var oldModPointer any
	mod, err := models.GlobalFactors(
		models.GlobalFactorWhere.PublisherID.EQ(data.Publisher),
		models.GlobalFactorWhere.Key.EQ(data.Key),
	).One(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if mod == nil {
		mod = &models.GlobalFactor{
			Key:         data.Key,
			PublisherID: data.Publisher,
			Value:       null.Float64From(data.Value),
			CreatedAt:   null.TimeFrom(time.Now().UTC()),
		}

		err := mod.Insert(ctx, bcdb.DB(), boil.Infer())
		if err != nil {
			return err
		}
	} else {
		oldMod := *mod
		oldModPointer = &oldMod

		mod.Value = null.Float64From(data.Value)
		mod.UpdatedAt = null.TimeFrom(time.Now().UTC())

		_, err := mod.Update(ctx, bcdb.DB(), boil.Infer())
		if err != nil {
			return err
		}
	}

	g.historyModule.SaveAction(ctx, oldModPointer, mod, nil)

	return nil
}

func (filter *GlobalFactorFilter) QueryMod() qmods.QueryModsSlice {

	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.GlobalFactorColumns.PublisherID))
	}

	if len(filter.Key) > 0 {
		mods = append(mods, filter.Key.AndIn(models.GlobalFactorColumns.Key))
	}

	if len(filter.Value) > 0 {
		mods = append(mods, filter.Value.AndIn(models.GlobalFactorColumns.Value))
	}

	return mods
}
