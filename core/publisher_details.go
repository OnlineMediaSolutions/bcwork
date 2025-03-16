package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type GetPublisherDetailsOptions struct {
	Filter     PublisherDetailsFilter `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type PublisherDetailsFilter struct {
	PublisherID      filter.StringArrayFilter `json:"publisher_id,omitempty"`
	Domain           filter.StringArrayFilter `json:"domain,omitempty"`
	Automation       *filter.BoolFilter       `json:"automation,omitempty"`
	AccountManagerID filter.StringArrayFilter `json:"account_manager,omitempty"`
}

func (p *PublisherService) GetPublisherDetails(
	ctx context.Context,
	ops *GetPublisherDetailsOptions,
	activityStatus map[string]map[string]dto.ActivityStatus,
) (dto.PublisherDetailsSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(updateFieldNames(ops.Order), nil, models.TableNames.Publisher+"."+models.PublisherColumns.PublisherID).
		AddArray(ops.Pagination.Do()).
		Add(
			qm.Select(
				models.TableNames.Publisher+"."+models.PublisherColumns.Name,
				models.TableNames.Publisher+"."+models.PublisherColumns.PublisherID,
				models.TableNames.PublisherDomain+"."+models.PublisherDomainColumns.Domain,
				models.TableNames.Publisher+"."+models.PublisherColumns.AccountManagerID,
				models.TableNames.PublisherDomain+"."+models.PublisherDomainColumns.Automation,
				models.TableNames.PublisherDomain+"."+models.PublisherDomainColumns.GPPTarget,
				models.TableNames.User+"."+models.UserColumns.FirstName,
				models.TableNames.User+"."+models.UserColumns.LastName,
			),
			qm.From(models.TableNames.Publisher),
			qm.InnerJoin(
				models.TableNames.PublisherDomain+" ON "+
					models.TableNames.Publisher+"."+models.PublisherColumns.PublisherID+
					" = "+
					models.TableNames.PublisherDomain+"."+models.PublisherDomainColumns.PublisherID,
			),
			qm.LeftOuterJoin(
				`"`+models.TableNames.User+`" ON `+
					models.TableNames.Publisher+`.`+models.PublisherColumns.AccountManagerID+
					` = `+
					`"`+models.TableNames.User+`".`+models.UserColumns.ID+`::varchar`,
			),
		)

	var mods []*dto.PublisherDetailModel
	err := models.NewQuery(qmods...).Bind(ctx, bcdb.DB(), &mods)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrap(err, "Failed to retrieve publisher, domains and factor values")
	}

	res := make(dto.PublisherDetailsSlice, 0, len(mods))

	confiantMap, pixalateMap, bidCachingMap, refreshCacheMap, err := p.fetchExtraDataPerPublisherDomain(ctx, mods)
	if err != nil {
		return nil, err
	}

	err = res.FromModel(mods, activityStatus, confiantMap, pixalateMap, bidCachingMap, refreshCacheMap)
	if err != nil {
		return nil, fmt.Errorf("failed to map publisher details: %w", err)
	}

	return res, nil
}

func (p *PublisherService) fetchExtraDataPerPublisherDomain(ctx context.Context, mods []*dto.PublisherDetailModel) (map[string]models.Confiant, map[string]models.Pixalate, map[string][]models.BidCaching, map[string][]models.RefreshCache, error) {
	var pubDom models.PublisherDomainSlice
	for _, mod := range mods {
		pubDom = append(pubDom, &models.PublisherDomain{
			Domain:      mod.PublisherDomain.Domain,
			PublisherID: mod.Publisher.PublisherID,
		})
	}

	confiantMap, err := LoadConfiantByPublisherAndDomain(ctx, pubDom)
	if err != nil {
		return nil, nil, nil, nil, eris.Wrap(err, "Error while retrieving confiant data for publisher domains values")
	}

	pixalateMap, err := LoadPixalateByPublisherAndDomain(ctx, pubDom)
	if err != nil {
		return nil, nil, nil, nil, eris.Wrap(err, "Error while retrieving pixalate data for publisher domains values")
	}

	bidCachingMap, err := LoadBidCacheByPublisherAndDomain(ctx, pubDom)
	if err != nil {
		return nil, nil, nil, nil, eris.Wrap(err, "Error while retrieving bid cache data for publisher domains values")
	}

	refreshCacheMap, err := LoadRefreshCacheByPublisherAndDomain(ctx, pubDom)
	if err != nil {
		return nil, nil, nil, nil, eris.Wrap(err, "Error while retrieving refresh cache data for publisher domains values")
	}
	return confiantMap, pixalateMap, bidCachingMap, refreshCacheMap, nil
}

func generatePublisherDomainSlice(mods []*dto.PublisherDetailModel) models.PublisherDomainSlice {
	var pubDomains models.PublisherDomainSlice
	for _, mod := range mods {
		pubDomains = append(pubDomains, &models.PublisherDomain{
			Domain:      mod.PublisherDomain.Domain,
			PublisherID: mod.Publisher.PublisherID,
		})
	}
	return pubDomains
}

// updateFieldNames To solve problem of column names ambiguous
func updateFieldNames(order order.Sort) order.Sort {
	for i := range order {
		switch order[i].Name {
		case models.PublisherColumns.PublisherID:
			order[i].Name = models.TableNames.Publisher + "." + order[i].Name
		}
	}

	return order
}

func (filter *PublisherDetailsFilter) QueryMod() qmods.QueryModsSlice {
	const amountOfFields = 4

	mods := make(qmods.QueryModsSlice, 0, amountOfFields)
	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(
			models.TableNames.Publisher+"."+models.PublisherColumns.PublisherID,
		))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(
			models.TableNames.PublisherDomain+"."+models.PublisherDomainColumns.Domain,
		))
	}

	if len(filter.AccountManagerID) > 0 {
		mods = append(mods, filter.AccountManagerID.AndIn(
			models.TableNames.Publisher+"."+models.PublisherColumns.AccountManagerID,
		))
	}

	if filter.Automation != nil {
		mods = append(mods, filter.Automation.Where(
			models.TableNames.PublisherDomain+"."+models.PublisherDomainColumns.Automation,
		))
	}

	return mods
}
