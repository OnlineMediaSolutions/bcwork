package core

import (
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"golang.org/x/net/context"
)

const (
	cursorIDOrderByDefault = "id"
	cursorIDColumnName     = "cursor_id"
	groupByDPIDColumnName  = "group_by_dp_id"
)

func (a *AdsTxtService) GetAdsTxtDataForFilters(ctx context.Context) (map[string][]interface{}, error) {
	mods, err := models.AdsTXTMainViews().All(ctx, bcdb.DB())
	if err != nil {
		return nil, err
	}

	// Initialize the map for filtering data
	adsTxtFilterMap := make(map[string]map[interface{}]struct{})

	// Initialize a set for each filter key
	adsTxtFilterMap["publisher_id"] = make(map[interface{}]struct{})
	adsTxtFilterMap["publisher_name"] = make(map[interface{}]struct{})
	adsTxtFilterMap["domain"] = make(map[interface{}]struct{})
	adsTxtFilterMap["domain_status"] = make(map[interface{}]struct{})
	adsTxtFilterMap["demand_partner_name"] = make(map[interface{}]struct{})
	adsTxtFilterMap["demand_partner_name_extended"] = make(map[interface{}]struct{})
	adsTxtFilterMap["demand_status"] = make(map[interface{}]struct{})

	for _, mod := range mods {
		adsTxtFilterMap["publisher_id"][mod.PublisherID] = struct{}{}
		adsTxtFilterMap["publisher_name"][mod.PublisherName] = struct{}{}
		adsTxtFilterMap["domain"][mod.Domain] = struct{}{}
		adsTxtFilterMap["domain_status"][mod.DomainStatus] = struct{}{}
		if mod.DemandPartnerName.String != "" {
			adsTxtFilterMap["demand_partner_name"][mod.DemandPartnerName] = struct{}{}
		}
		adsTxtFilterMap["demand_partner_name_extended"][mod.DemandPartnerNameExtended] = struct{}{}
		adsTxtFilterMap["demand_status"][mod.DemandStatus] = struct{}{}
	}

	result := make(map[string][]interface{})
	for key, valueSet := range adsTxtFilterMap {
		var uniqueValues []interface{}
		for value := range valueSet {
			uniqueValues = append(uniqueValues, value)
		}
		result[key] = uniqueValues
	}

	return result, nil
}

type AdsTxtFilterMap struct {
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
	IsReadyToGoLive   *filter.BoolFilter       `json:"is_ready_to_go_live,omitempty"`
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

	if filter.IsReadyToGoLive != nil {
		mods = append(mods, filter.IsReadyToGoLive.Where(models.AdsTXTGroupByDPViewColumns.IsReadyToGoLive))
	}

	return mods
}
