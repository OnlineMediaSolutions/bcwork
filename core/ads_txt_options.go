package core

import (
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
)

const cursorIDColumnName = "cursor_id"

type AdsTxtGetBaseOptions struct {
	Filter     AdsTxtGetBaseFilter    `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type AdsTxtGetBaseFilter struct {
	PublisherID               filter.StringArrayFilter   `json:"publisher_id,omitempty"`
	AccountManagerID          filter.StringArrayFilter   `json:"account_manager_id,omitempty"`
	CampaignManagerID         filter.StringArrayFilter   `json:"campaign_manager_id,omitempty"`
	Domain                    filter.StringArrayFilter   `json:"domain,omitempty"`
	DomainStatus              filter.StringArrayFilter   `json:"domain_status,omitempty"`
	DemandPartnerName         filter.StringArrayFilter   `json:"demand_partner_name,omitempty"`
	DemandPartnerNameExtended filter.StringArrayFilter   `json:"demand_partner_name_extended"`
	MediaType                 filter.String2DArrayFilter `json:"media_type,omitempty"`
	DemandManagerID           filter.IntArrayFilter      `json:"demand_manager_id,omitempty"`
	DemandStatus              filter.StringArrayFilter   `json:"demand_status,omitempty"`
	Status                    filter.StringArrayFilter   `json:"status,omitempty"`
	IsRequired                *filter.BoolFilter         `json:"is_required,omitempty"`
	// TODO: add mirror filtering
}

// TODO: return total amount of rows
// TODO: return filters

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

	if len(filter.DemandPartnerName) > 0 {
		mods = append(mods, filter.DemandPartnerName.AndIn(models.AdsTXTMainViewColumns.DemandPartnerName))
	}

	if len(filter.DemandPartnerNameExtended) > 0 {
		mods = append(mods, filter.DemandPartnerNameExtended.AndIn(models.AdsTXTMainViewColumns.DemandPartnerNameExtended))
	}

	if len(filter.MediaType) > 0 {
		mods = append(mods, filter.MediaType.AndIn(models.AdsTXTMainViewColumns.MediaType))
	}

	if len(filter.DemandManagerID) > 0 {
		mods = append(mods, filter.DemandManagerID.AndIn(models.AdsTXTMainViewColumns.DemandManagerID))
	}

	if len(filter.DemandStatus) > 0 {
		mods = append(mods, filter.DemandStatus.AndIn(models.AdsTXTMainViewColumns.DemandStatus))
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
	Filter     AdsTxtGetGroupByDPFilter `json:"filter"`
	Pagination *pagination.Pagination   `json:"pagination"`
	Order      order.Sort               `json:"order"`
	Selector   string                   `json:"selector"`
}

type AdsTxtGetGroupByDPFilter struct {
	AdsTxtGetBaseFilter
	IsReadyToGoLive *filter.BoolFilter `json:"is_ready_to_go_live,omitempty"`
}

func (filter *AdsTxtGetGroupByDPFilter) queryModGroupByDP() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	mods = filter.queryMod()

	if filter.IsReadyToGoLive != nil {
		mods = append(mods, filter.IsReadyToGoLive.Where(models.AdsTXTGroupByDPViewColumns.IsReadyToGoLive))
	}

	return mods
}
