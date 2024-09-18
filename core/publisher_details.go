package core

import (
	"context"
	"database/sql"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
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
	PublisherID filter.StringArrayFilter `json:"publisher_id,omitempty"`
	Domain      filter.StringArrayFilter `json:"domain,omitempty"`
	Automation  filter.StringArrayFilter `json:"automation,omitempty"`
	GppTarget   filter.StringArrayFilter `json:"gpp_target,omitempty"`
}

func GetPublisherDetails(ctx context.Context, ops *GetPublisherDetailsOptions) (PublisherDetailsSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.PublisherColumns.PublisherID).
		AddArray(ops.Pagination.Do())
	qmods = qmods.Add(
		qm.Load(models.PublisherRels.PublisherDomains),
	)

	var mods models.PublisherSlice
	mods, err := models.Publishers(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve publisher, domains and factor values")
	}

	res := make(PublisherDetailsSlice, 0, len(mods))
	res.FromModel(mods)

	return res, nil
}

func (filter *PublisherDetailsFilter) QueryMod() qmods.QueryModsSlice {
	const amountOfFields = 6

	mods := make(qmods.QueryModsSlice, 0, amountOfFields)
	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.PublisherColumns.PublisherID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.PublisherDomainColumns.Domain))
	}

	if len(filter.GppTarget) > 0 {
		mods = append(mods, filter.GppTarget.AndIn(models.PublisherDomainColumns.GPPTarget))
	}

	if len(filter.Automation) > 0 {
		mods = append(mods, filter.Automation.AndIn(models.PublisherDomainColumns.Automation))
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
}

func (pd *PublisherDetail) FromModel(mod *models.Publisher, domain *models.PublisherDomain) error {
	pd.Name = mod.Name
	pd.PublisherID = mod.PublisherID
	pd.Domain = domain.Domain
	pd.AccountManagerID = mod.AccountManagerID.String
	pd.Automation = domain.Automation
	pd.GPPTarget = domain.GPPTarget.Float64

	return nil
}

type PublisherDetailsSlice []*PublisherDetail

func (pds *PublisherDetailsSlice) FromModel(mods models.PublisherSlice) error {
	for _, mod := range mods {
		for _, domain := range mod.R.GetPublisherDomains() {
			pd := PublisherDetail{}
			err := pd.FromModel(mod, domain)
			if err != nil {
				return eris.Cause(err)
			}
			*pds = append(*pds, &pd)
		}
	}

	return nil
}
