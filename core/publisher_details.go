package core

import (
	"context"
	"database/sql"
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

type modelsPublisherDetail struct {
	Publisher       models.Publisher       `boil:"publisher,bind"`
	PublisherDomain models.PublisherDomain `boil:"publisher_domain,bind"`
}

func (p *PublisherService) GetPublisherDetails(ctx context.Context, ops *GetPublisherDetailsOptions, activityStatus map[string]map[string]dto.ActivityStatus) (PublisherDetailsSlice, error) {
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
			),
			qm.From(models.TableNames.Publisher),
			qm.InnerJoin(
				models.TableNames.PublisherDomain+" ON "+
					models.TableNames.Publisher+"."+models.PublisherColumns.PublisherID+
					" = "+
					models.TableNames.PublisherDomain+"."+models.PublisherDomainColumns.PublisherID,
			),
		)

	var mods []*modelsPublisherDetail
	err := models.NewQuery(qmods...).Bind(ctx, bcdb.DB(), &mods)
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve publisher, domains and factor values")
	}

	res := make(PublisherDetailsSlice, 0, len(mods))
	res.FromModel(mods, activityStatus)

	return res, nil
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

type PublisherDetail struct {
	Name             string  `boil:"name" json:"name" toml:"name" yaml:"name"`
	PublisherID      string  `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain           string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	AccountManagerID string  `boil:"account_manager_id" json:"account_manager_id,omitempty" toml:"account_manager_id" yaml:"account_manager_id,omitempty"`
	Automation       bool    `boil:"automation" json:"automation" toml:"automation" yaml:"automation"`
	GPPTarget        float64 `boil:"gpp_target" json:"gpp_target" toml:"gpp_target" yaml:"gpp_target,omitempty"`
	ActivityStatus   string  `boil:"activity_status" json:"activity_status" toml:"activity_status" yaml:"activity_status"`
}

func (pd *PublisherDetail) FromModel(mod *modelsPublisherDetail, activityStatus map[string]map[string]dto.ActivityStatus) error {
	pd.Name = mod.Publisher.Name
	pd.PublisherID = mod.Publisher.PublisherID
	pd.Domain = mod.PublisherDomain.Domain
	pd.AccountManagerID = mod.Publisher.AccountManagerID.String
	pd.Automation = mod.PublisherDomain.Automation
	pd.GPPTarget = mod.PublisherDomain.GPPTarget.Float64
	pd.ActivityStatus = activityStatus[pd.Domain][pd.PublisherID].String()
	return nil
}

type PublisherDetailsSlice []*PublisherDetail

func (pds *PublisherDetailsSlice) FromModel(mods []*modelsPublisherDetail, activityStatus map[string]map[string]dto.ActivityStatus) error {
	for _, mod := range mods {
		pd := PublisherDetail{}
		err := pd.FromModel(mod, activityStatus)
		if err != nil {
			return eris.Cause(err)
		}
		*pds = append(*pds, &pd)

	}

	return nil
}
