package core

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DomainService struct {
	historyModule history.HistoryModule
}

func NewDomainService(historyModule history.HistoryModule) *DomainService {
	return &DomainService{
		historyModule: historyModule,
	}
}

type GetPublisherDomainOptions struct {
	Filter     PublisherDomainFilter  `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type PublisherDomainFilter struct {
	Domain      filter.StringArrayFilter `json:"domain,omitempty"`
	PublisherID filter.StringArrayFilter `json:"publisher_id,omitempty"`
	Automation  *filter.BoolFilter       `json:"automation,omitempty"`
	GppTarget   filter.StringArrayFilter `json:"gpp_target,omitempty"`
}

func (d *DomainService) GetPublisherDomain(ctx context.Context, ops *GetPublisherDomainOptions) (dto.PublisherDomainSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.PublisherDomainColumns.PublisherID).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))
	qmods = qmods.Add(qm.Load(models.PublisherDomainRels.Publisher))

	mods, err := models.PublisherDomains(qmods...).All(ctx, bcdb.DB())
	if err != nil {
		return nil, eris.Wrap(err, "Failed to retrieve publisher domains values")
	}

	if len(mods) == 0 {
		return dto.PublisherDomainSlice{}, nil
	}

	confiantMap, err := LoadConfiantByPublisherAndDomain(ctx, mods)
	if err != nil {
		return nil, eris.Wrap(err, "Error while retreving confiant data for publisher domains values")
	}

	pixalateMap, err := LoadPixalateByPublisherAndDomain(ctx, mods)
	if err != nil {
		return nil, eris.Wrap(err, "Error while retreving pixalate data for publisher domains values")
	}

	bidCachingMap, err := LoadBidCacheByPublisherAndDomain(ctx, mods)
	if err != nil {
		return nil, eris.Wrap(err, "Error while retreving bid cache data for publisher domains values")
	}

	refreshCacheMap, err := LoadRefreshCacheByPublisherAndDomain(ctx, mods)
	if err != nil {
		return nil, eris.Wrap(err, "Error while retreving refresh cache data for publisher domains values")
	}

	res := make(dto.PublisherDomainSlice, 0, len(mods))
	err = res.FromModel(mods, confiantMap, pixalateMap, bidCachingMap, refreshCacheMap)
	if err != nil {
		return nil, eris.Wrap(err, "failed to map publisher domain")
	}

	return res, nil
}

func (filter *PublisherDomainFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.PublisherDomainColumns.PublisherID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.PublisherDomainColumns.Domain))
	}

	if len(filter.GppTarget) > 0 {
		mods = append(mods, filter.GppTarget.AndIn(models.PublisherDomainColumns.GPPTarget))
	}

	if filter.Automation != nil {
		mods = append(mods, filter.Automation.Where(models.PublisherDomainColumns.Automation))
	}

	return mods
}

func (d *DomainService) UpdatePublisherDomain(ctx context.Context, data *dto.PublisherDomainUpdateRequest) error {
	var oldModPointer any
	mod, err := models.PublisherDomains(
		models.PublisherDomainWhere.PublisherID.EQ(data.PublisherID),
		models.PublisherDomainWhere.Domain.EQ(data.Domain),
	).One(ctx, bcdb.DB())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if mod == nil {
		mod = &models.PublisherDomain{
			Domain:            data.Domain,
			PublisherID:       data.PublisherID,
			Automation:        data.Automation,
			GPPTarget:         null.Float64FromPtr(data.GppTarget),
			IntegrationType:   data.IntegrationType,
			MirrorPublisherID: null.StringFromPtr(data.MirrorPublisherID),
			IsDirect:          null.BoolFromPtr(data.IsDirect),
		}

		err := mod.Insert(ctx, bcdb.DB(), boil.Infer())
		if err != nil {
			return err
		}
	} else {
		oldMod := *mod
		oldModPointer = &oldMod

		mod.Automation = data.Automation
		mod.GPPTarget = null.Float64FromPtr(data.GppTarget)
		mod.IntegrationType = data.IntegrationType
		mod.UpdatedAt = null.TimeFrom(time.Now().UTC())
		mod.MirrorPublisherID = null.StringFromPtr(data.MirrorPublisherID)
		mod.IsDirect = null.BoolFromPtr(data.IsDirect)

		_, err := mod.Update(ctx, bcdb.DB(), boil.Infer())
		if err != nil {
			return err
		}
	}

	d.historyModule.SaveAction(ctx, oldModPointer, mod, nil)

	return nil
}
