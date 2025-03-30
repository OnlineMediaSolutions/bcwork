package core

import (
	"context"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	cursorIDOrderByDefault = "id"
	cursorIDColumnName     = "cursor_id"
	groupByDPIDColumnName  = "group_by_dp_id"
)

type filterResponse struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func (a *AdsTxtService) GetAdsTxtDataForFilters(ctx context.Context, filterName string) ([]*filterResponse, error) {
	var (
		filters []*filterResponse
		query   string
		table   string
	)

	switch filterName {
	case models.AdsTXTMainViewColumns.PublisherID:
		query = `DISTINCT publisher_name || '(' || publisher_id || ')' as label, publisher_id AS value`
		table = models.ViewNames.AdsTXTMainView
	case models.AdsTXTMainViewColumns.Domain:
		query = `DISTINCT "domain" as label, "domain" AS value`
		table = models.ViewNames.AdsTXTMainView
	case models.AdsTXTMainViewColumns.DomainStatus:
		for value, label := range dto.DomainStatusMap {
			filters = append(filters, &filterResponse{Label: label, Value: value})
		}

		return filters, nil
	case models.AdsTXTMainViewColumns.DemandPartnerNameExtended:
		query = `DISTINCT demand_partner_name_extended as label, demand_partner_name_extended AS value`
		table = models.ViewNames.AdsTXTMainView
	case models.AdsTXTMainViewColumns.AccountManagerID:
		query = `DISTINCT account_manager_full_name as label, account_manager_id AS value`
		table = models.ViewNames.AdsTXTMainView
	case models.AdsTXTMainViewColumns.CampaignManagerID:
		query = `DISTINCT campaign_manager_full_name as label, campaign_manager_id AS value`
		table = models.ViewNames.AdsTXTMainView
	case models.AdsTXTMainViewColumns.DemandManagerID:
		query = `DISTINCT demand_manager_full_name as label, demand_manager_id AS value`
		table = models.ViewNames.AdsTXTMainView
	case models.AdsTXTMainViewColumns.Status:
		for value, label := range dto.StatusMap {
			filters = append(filters, &filterResponse{Label: label, Value: value})
		}

		return filters, nil
	case models.AdsTXTMainViewColumns.MediaType:
		return []*filterResponse{
			{Label: dto.WebBannersMediaType, Value: dto.WebBannersMediaType},
			{Label: dto.VideoMediaType, Value: dto.VideoMediaType},
			{Label: dto.InAppMediaType, Value: dto.InAppMediaType},
		}, nil
	case models.AdsTXTGroupByDPViewColumns.DemandPartnerName:
		query = `DISTINCT demand_partner_name as label, demand_partner_name AS value`
		table = models.ViewNames.AdsTXTGroupByDPView
	case models.AdsTXTMainViewColumns.DemandStatus:
		for value, label := range dto.DPStatusMap {
			filters = append(filters, &filterResponse{Label: label, Value: value})
		}

		return filters, nil
	}

	err := models.NewQuery(
		qm.Select(query),
		qm.From(table),
	).Bind(ctx, bcdb.DB(), &filters)
	if err != nil {
		return nil, err
	}

	return filters, nil
}

type AdsTxtGetBaseOptions struct {
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type AdsTxtGetBaseFilter struct {
	PublisherID       filter.StringArrayFilter   `json:"publisher_id,omitempty"`
	AccountManagerID  filter.StringArrayFilter   `json:"account_manager_id,omitempty"`
	CampaignManagerID filter.StringArrayFilter   `json:"campaign_manager_id,omitempty"`
	Domain            filter.StringArrayFilter   `json:"domain,omitempty"`
	MediaType         filter.String2DArrayFilter `json:"media_type,omitempty"`
	DemandStatus      filter.StringArrayFilter   `json:"demand_status,omitempty"`
	DomainStatus      filter.StringArrayFilter   `json:"domain_status,omitempty"`
	DemandManagerID   filter.IntArrayFilter      `json:"demand_manager_id,omitempty"`
}

func (filter *AdsTxtGetBaseFilter) queryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.AdsTXTMainViewColumns.PublisherID))
	}

	if len(filter.AccountManagerID) > 0 {
		mods = append(mods, filter.AccountManagerID.AndIn(models.AdsTXTMainViewColumns.AccountManagerID))
	}

	if len(filter.CampaignManagerID) > 0 {
		mods = append(mods, filter.CampaignManagerID.AndIn(models.AdsTXTMainViewColumns.CampaignManagerID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.AdsTXTMainViewColumns.Domain))
	}

	if len(filter.MediaType) > 0 {
		mods = append(mods, filter.MediaType.AndIn(models.AdsTXTMainViewColumns.MediaType))
	}

	if len(filter.DemandStatus) > 0 {
		mods = append(mods, filter.DemandStatus.AndIn(models.AdsTXTMainViewColumns.DemandStatus))
	}

	if len(filter.DomainStatus) > 0 {
		mods = append(mods, filter.DomainStatus.AndIn(models.AdsTXTMainViewColumns.DomainStatus))
	}

	if len(filter.DemandManagerID) > 0 {
		mods = append(mods, filter.DemandManagerID.AndIn(models.AdsTXTMainViewColumns.DemandManagerID))
	}

	return mods
}

type AdsTxtGetMainOptions struct {
	AdsTxtGetBaseOptions
	Filter AdsTxtGetMainFilter `json:"filter"`
}

type AdsTxtGetMainFilter struct {
	AdsTxtGetBaseFilter
	DemandPartnerNameExtended filter.StringArrayFilter `json:"demand_partner_name_extended"`
	Status                    filter.StringArrayFilter `json:"status,omitempty"`
	IsRequired                *filter.BoolFilter       `json:"is_required,omitempty"`
}

func (filter *AdsTxtGetMainFilter) queryModMain() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	mods = filter.queryMod()

	if len(filter.DemandPartnerNameExtended) > 0 {
		mods = append(mods, filter.DemandPartnerNameExtended.AndIn(models.AdsTXTMainViewColumns.DemandPartnerNameExtended))
	}

	if len(filter.Status) > 0 {
		mods = append(mods, filter.Status.AndIn(models.AdsTXTMainViewColumns.Status))
	}

	if filter.IsRequired != nil {
		mods = append(mods, filter.IsRequired.Where(models.AdsTXTMainViewColumns.IsRequired))
	}

	return mods
}

type AdsTxtGetGroupByDPOptions struct {
	AdsTxtGetBaseOptions
	Filter AdsTxtGetGroupByDPFilter `json:"filter"`
}

type AdsTxtGetGroupByDPFilter struct {
	AdsTxtGetBaseFilter
	DemandPartnerName filter.StringArrayFilter `json:"demand_partner_name,omitempty"`
	DPEnabled         *filter.BoolFilter       `json:"dp_enabled,omitempty"`
}

func (filter *AdsTxtGetGroupByDPFilter) queryModGroupByDP() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	mods = filter.queryMod()

	if len(filter.DemandPartnerName) > 0 {
		mods = append(mods, filter.DemandPartnerName.AndIn(models.AdsTXTMainViewColumns.DemandPartnerName))
	}

	if filter.DPEnabled != nil {
		mods = append(mods, filter.DPEnabled.Where(models.AdsTXTGroupByDPViewColumns.DPEnabled))
	}

	return mods
}
