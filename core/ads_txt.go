package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type AdsTxtService struct {
	historyModule       history.HistoryModule
	compassModule       compass.CompassModule
	lowPerformanceCache *LowPerformanceCache
}

func NewAdsTxtService(historyModule history.HistoryModule, compassModule compass.CompassModule) *AdsTxtService {
	return &AdsTxtService{
		historyModule:       historyModule,
		compassModule:       compassModule,
		lowPerformanceCache: &LowPerformanceCache{},
	}
}

type LowPerformanceCache struct {
	partitions     string
	lowPerformance map[string]bool
	revenuePerDP   map[string]float64
	totalRevenue   map[string]float64
}

type ReportResult struct {
	Data struct {
		Result []struct {
			Domain        string  `json:"Domain"`
			DemandPartner string  `json:"DemandPartner"`
			Revenue       float64 `json:"Revenue"`
		} `json:"result"`
	} `json:"data"`
}

type AdsTxtOptions struct {
}

func (a *AdsTxtService) GetMainAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) ([]*dto.AdsTxt, error) {
	query := `
		select
			t.*,
			p."name" as publisher_name,
			p.account_manager_id,
			p.campaign_manager_id
		from (
			select
				at2.id,
				at2.publisher_id,
				at2."domain",
				at2.domain_status,
				d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name_extended,
				d.manager_id as demand_manager_id,
				at2.demand_status,
				at2.status,
				dpc.is_required_for_ads_txt as is_required,
				d.dp_domain || ', ' ||
					dpc.publisher_account || ', ' ||
					case
						when dpc.is_direct then 'DIRECT'
						else 'RESELLER'
					end ||
					case
						when d.certification_authority_id is not null
						then ', ' || d.certification_authority_id
					else ''
					end as ads_txt_line,
				at2.last_scanned_at,
				at2.error_message
			from ads_txt at2
			join demand_partner_connection dpc on at2.demand_partner_connection_id = dpc.id
			join dpo d on d.demand_partner_id = dpc.demand_partner_id
			union
			select
				at2.id,
				at2.publisher_id,
				at2."domain",
				at2.domain_status,
				d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name_extended,
				d.manager_id as demand_manager_id,
				at2.demand_status,
				at2.status,
				dpc.is_required_for_ads_txt as is_required,
				dpc.dp_child_domain || ', ' ||
					dpc.publisher_account || ', ' ||
					case
						when dpc.is_direct then 'DIRECT'
						else 'RESELLER'
					end ||
					case
						when dpc.certification_authority_id is not null
						then ', ' || dpc.certification_authority_id
					else ''
					end as ads_txt_line,
				at2.last_scanned_at,
				at2.error_message
			from ads_txt at2
			join demand_partner_child dpc on at2.demand_partner_child_id = dpc.id
			join demand_partner_connection dpc2 on dpc2.id = dpc.dp_connection_id
			join dpo d on d.demand_partner_id = dpc2.demand_partner_id
			union
			select
				at2.id,
				at2.publisher_id,
				at2."domain",
				at2.domain_status,
				so.seat_owner_name || ' - Direct' as demand_partner_name_extended,
				null as demand_manager_id,
				at2.demand_status,
				at2.status,
				true as is_required,
				so.seat_owner_domain || ', ' ||
					replace(so.publisher_account, '%s', at2.publisher_id) || ', ' ||
					'DIRECT' ||
					case
						when so.certification_authority_id is not null
						then ', ' || so.certification_authority_id
					else ''
					end as ads_txt_line,
				at2.last_scanned_at,
				at2.error_message
			from ads_txt at2
			join seat_owner so on at2.seat_owner_id = so.id
		) as t
		join publisher p on p.publisher_id = t.publisher_id
		order by t.id;
	`

	mainTable, err := a.getAdsTxtTableByQueryWithUsersFullNames(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve main table: %w", err)
	}

	return mainTable, nil
}

func (a *AdsTxtService) GetGroupByDPAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) (map[string]*dto.AdsTxtGroupedByDPData, error) {
	query := `
		select
			t.*,
			p."name" as publisher_name,
			p.account_manager_id,
			p.campaign_manager_id,
			sum(case
				when t.status = 'added' then 1
				else 0
			end) over (partition by t.publisher_id, t."domain", t.demand_partner_name, t."media_type") as added,
			count(t.status) over (partition by t.publisher_id, t."domain", t.demand_partner_name, t."media_type") as total,
			bool_and(case
				when t.status = 'added' AND t.is_required and t.demand_status = 'approved' then true
				when not t.is_required then true
				else false
			end) over (partition by t.publisher_id, t."domain", t.demand_partner_name, t."media_type") as is_ready_to_go_live
		from (
			select
				at2.id,
				at2.publisher_id,
				at2."domain",
				at2.domain_status,
				d.demand_partner_name,
				dpc."media_type",
				d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name_extended,
				d.manager_id as demand_manager_id,
				at2.demand_status,
				at2.status,
				dpc.is_required_for_ads_txt as is_required,
				d.active,
				d.dp_domain || ', ' ||
					dpc.publisher_account || ', ' ||
					case
						when dpc.is_direct then 'DIRECT'
						else 'RESELLER'
					end ||
					case
						when d.certification_authority_id is not null
						then ', ' || d.certification_authority_id
					else ''
					end as ads_txt_line,
				at2.last_scanned_at
			from ads_txt at2
			join demand_partner_connection dpc on at2.demand_partner_connection_id = dpc.id
			join dpo d on d.demand_partner_id = dpc.demand_partner_id
			where dpc.is_required_for_ads_txt
			union
			select
				at2.id,
				at2.publisher_id,
				at2."domain",
				at2.domain_status,
				d.demand_partner_name,
				dpc2."media_type",
				d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name_extended,
				d.manager_id as demand_manager_id,
				at2.demand_status,
				at2.status,
				dpc.is_required_for_ads_txt as is_required,
				d.active,
				dpc.dp_child_domain || ', ' ||
					dpc.publisher_account || ', ' ||
					case
						when dpc.is_direct then 'DIRECT'
						else 'RESELLER'
					end ||
					case
						when dpc.certification_authority_id is not null
						then ', ' || dpc.certification_authority_id
					else ''
					end as ads_txt_line,
				at2.last_scanned_at
			from ads_txt at2
			join demand_partner_child dpc on at2.demand_partner_child_id = dpc.id
			join demand_partner_connection dpc2 on dpc2.id = dpc.dp_connection_id
			join dpo d on d.demand_partner_id = dpc2.demand_partner_id
			union all
			select
				at2.id,
				at2.publisher_id,
				at2."domain",
				at2.domain_status,
				d.demand_partner_name,
				dpc."media_type",
				so.seat_owner_name || ' - Direct' as demand_partner_name_extended,
				null as demand_manager_id,
				at2.demand_status,
				at2.status,
				true as is_required,
				d.active,
				so.seat_owner_domain || ', ' ||
					replace(so.publisher_account, '%s', at2.publisher_id) || ', ' ||
					'DIRECT' ||
					case
						when so.certification_authority_id is not null
						then ', ' || so.certification_authority_id
					else ''
					end as ads_txt_line,
				at2.last_scanned_at
			from ads_txt at2
			join seat_owner so on at2.seat_owner_id = so.id
			join dpo d on d.seat_owner_id = so.id
			join demand_partner_connection dpc on d.demand_partner_id = dpc.demand_partner_id
		) as t
		join publisher p on t.publisher_id = p.publisher_id
		where t.active
		order by t.publisher_id, t."domain", t.demand_partner_name, t."media_type";
	`

	var rawTable []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &rawTable)
	if err != nil {
		return nil, err
	}

	//usersMap, err := getUsersMap(ctx)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to retrieve users: %w", err)
	//}
	//
	//groupByDpTable := make(map[string]*dto.AdsTxtGroupedByDPData)
	//for _, row := range rawTable {
	//name := fmt.Sprintf("%v:%v:%v:%v", row.PublisherID, row.Domain, row.DemandPartnerName, strings.Join(row.MediaType, ","))

	//row.AccountManagerFullName = usersMap[row.AccountManagerID.String]
	//row.CampaignManagerFullName = usersMap[row.CampaignManagerID.String]
	//row.DemandManagerFullName = usersMap[row.DemandManagerID.String]

	//dpData, ok := groupByDpTable[name]
	//if !ok {
	//	groupByDpTable[name] = &dto.AdsTxtGroupedByDPData{
	//		Parent:   row,
	//		Children: []*dto.AdsTxt{row},
	//	}
	//} else {
	//	dpData.Children = append(dpData.Children, row)
	//}
	//}

	return nil, nil
}

func (a *AdsTxtService) GetAMAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) ([]*dto.AdsTxt, error) {
	type seatOwnersWithActiveDP struct {
		SeatOwnerName string
		HasActiveDP   bool
	}

	var seatOwners []*seatOwnersWithActiveDP
	err := queries.Raw(`
		select
			so.seat_owner_name,
			bool_or(d.active) as has_active_dp
		from seat_owner so
		join dpo d on d.seat_owner_id = so.id
		group by so.seat_owner_name
	`).Bind(ctx, bcdb.DB(), &seatOwners)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve seat owner map: %w", err)
	}

	seatOwnerMap := make(map[string]bool, len(seatOwners))
	for _, seatOwner := range seatOwners {
		seatOwnerMap[seatOwner.SeatOwnerName] = seatOwner.HasActiveDP
	}

	query := `
		with am as (
			select
				t.*,
				p."name" as publisher_name,
				p.account_manager_id,
				p.campaign_manager_id
			from (
				select
					at2.id,
					at2.publisher_id,
					at2."domain",
					at2.domain_status,
					d.demand_partner_name,
					d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name_extended,
					d.manager_id as demand_manager_id,
					at2.demand_status,
					at2.status,
					dpc.is_required_for_ads_txt as is_required,
					d.dp_domain || ', ' ||
						dpc.publisher_account || ', ' ||
						case
							when dpc.is_direct then 'DIRECT'
							else 'RESELLER'
						end ||
						case
							when d.certification_authority_id is not null
							then ', ' || d.certification_authority_id
						else ''
						end as ads_txt_line,
					at2.last_scanned_at,
					at2.error_message,
					d.active as is_demand_partner_active
				from ads_txt at2
				join demand_partner_connection dpc on at2.demand_partner_connection_id = dpc.id
				join dpo d on d.demand_partner_id = dpc.demand_partner_id
				union
				select
					at2.id,
					at2.publisher_id,
					at2."domain",
					at2.domain_status,
					dpc.dp_child_name as demand_partner_name,
					d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name_extended,
					d.manager_id as demand_manager_id,
					at2.demand_status,
					at2.status,
					dpc.is_required_for_ads_txt as is_required,
					dpc.dp_child_domain || ', ' ||
						dpc.publisher_account || ', ' ||
						case
							when dpc.is_direct then 'DIRECT'
							else 'RESELLER'
						end ||
						case
							when dpc.certification_authority_id is not null
							then ', ' || dpc.certification_authority_id
						else ''
						end as ads_txt_line,
					at2.last_scanned_at,
					at2.error_message,
					d.active as is_demand_partner_active
				from ads_txt at2
				join demand_partner_child dpc on at2.demand_partner_child_id = dpc.id
				join demand_partner_connection dpc2 on dpc2.id = dpc.dp_connection_id
				join dpo d on d.demand_partner_id = dpc2.demand_partner_id
				union all
				select
					at2.id,
					at2.publisher_id,
					at2."domain",
					at2.domain_status,
					so.seat_owner_name as demand_partner_name,
					so.seat_owner_name || ' - Direct' as demand_partner_name_extended,
					null as demand_manager_id,
					at2.demand_status,
					at2.status,
					true as is_required,
					so.seat_owner_domain || ', ' ||
						replace(so.publisher_account, '%s', at2.publisher_id) || ', ' ||
						'DIRECT' ||
						case
							when so.certification_authority_id is not null
							then ', ' || so.certification_authority_id
						else ''
						end as ads_txt_line,
					at2.last_scanned_at,
					at2.error_message,
					true as is_demand_partner_active
				from ads_txt at2
				join seat_owner so on at2.seat_owner_id = so.id
			) as t
			join publisher p on p.publisher_id = t.publisher_id
			where t.status != 'not_scanned' and t.domain_status != 'paused'
		)
		select
			*
		from am a
		order by a.id;
	`

	var (
		errGroup errgroup.Group
		amTable  []*dto.AdsTxt
	)

	errGroup.Go(func() error {
		amTable, err = a.getAdsTxtTableByQueryWithUsersFullNames(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to retrieve am table: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		err = a.updatePerfomanceData()
		if err != nil {
			return fmt.Errorf("failed to update perfomance data: %w", err)
		}

		return nil
	})

	err = errGroup.Wait()
	if err != nil {
		return nil, err
	}

	for _, row := range amTable {
		key := fmt.Sprintf("%s_%s_%s", strings.ToLower(row.Domain), strings.ToLower(row.DemandPartnerName), "nb")
		action := dto.AdsTxtActionNone
		hasActiveDP, isSeatOwner := seatOwnerMap[row.DemandPartnerName]

		if slices.Contains(
			[]string{
				dto.DPStatusRejected,
				dto.DPStatusRejectedTQ,
				dto.DPStatusWillNotBeSent,
				dto.DPStatusDisabledSPO,
				dto.DPStatusApprovedPaused,
				dto.DPStatusDisabledNoImps,
			},
			row.DemandStatus,
		) || (!isSeatOwner && !row.IsDemandPartnerActive) {
			action = dto.AdsTxtActionNone
			if row.Status == dto.AdsTxtStatusAdded {
				action = dto.AdsTxtActionRemove
			}
			//} else if row.Status == dto.AdsTxtStatusAdded {
			//	action = dto.AdsTxtActionKeep
			//	if row.ErrorMessage.Valid {
			//		action = dto.AdsTxtActionFix
			//	}
		} else {
			action = dto.AdsTxtActionAdd
		}

		if !isSeatOwner &&
			row.Status == dto.AdsTxtStatusAdded &&
			!slices.Contains([]string{dto.AdsTxtActionFix, dto.AdsTxtActionRemove}, action) &&
			a.lowPerformanceCache.lowPerformance[key] {
			action = dto.AdsTxtActionLowPerfomance
		}

		if isSeatOwner && !hasActiveDP {
			if action == dto.AdsTxtActionAdd {
				action = dto.AdsTxtActionNone
			} else if slices.Contains([]string{dto.AdsTxtActionKeep, dto.AdsTxtActionFix}, action) {
				action = dto.AdsTxtActionRemove
			}
		}

		row.Action = action
	}

	return amTable, nil
}

func (a *AdsTxtService) GetCMAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) ([]*dto.AdsTxt, error) {
	query := `
		with cm as (
			select
				t.*,
				p."name" as publisher_name,
				p.account_manager_id,
				p.campaign_manager_id,
				case
					when t.is_approval_needed and t.is_required then 1
					when t.is_approval_needed and not t.is_required then 2
					when not t.is_approval_needed and t.is_required then 3
					else 0
				end as approval_group
			from (
				select
					at2.id,
					at2.publisher_id,
					at2."domain",
					at2.domain_status,
					d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name_extended,
					d.manager_id as demand_manager_id,
					at2.demand_status,
					at2.status,
					d.is_approval_needed,
					dpc.is_required_for_ads_txt as is_required,
					d.active,
					d.dp_domain || ', ' ||
						dpc.publisher_account || ', ' ||
						case
							when dpc.is_direct then 'DIRECT'
							else 'RESELLER'
						end ||
						case
							when d.certification_authority_id is not null
							then ', ' || d.certification_authority_id
						else ''
						end as ads_txt_line,
					at2.last_scanned_at,
					at2.error_message
				from ads_txt at2
				join demand_partner_connection dpc on at2.demand_partner_connection_id = dpc.id
				join dpo d on d.demand_partner_id = dpc.demand_partner_id
				union
				select
					at2.id,
					at2.publisher_id,
					at2."domain",
					at2.domain_status,
					d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name_extended,
					d.manager_id as demand_manager_id,
					at2.demand_status,
					at2.status,
					d.is_approval_needed,
					dpc.is_required_for_ads_txt as is_required,
					d.active,
					dpc.dp_child_domain || ', ' ||
						dpc.publisher_account || ', ' ||
						case
							when dpc.is_direct then 'DIRECT'
							else 'RESELLER'
						end ||
						case
							when dpc.certification_authority_id is not null
							then ', ' || dpc.certification_authority_id
						else ''
						end as ads_txt_line,
					at2.last_scanned_at,
					at2.error_message
				from ads_txt at2
				join demand_partner_child dpc on at2.demand_partner_child_id = dpc.id
				join demand_partner_connection dpc2 on dpc2.id = dpc.dp_connection_id
				join dpo d on d.demand_partner_id = dpc2.demand_partner_id
			) as t
			join publisher p on p.publisher_id = t.publisher_id
			where t.active and t.demand_status in ('pending', 'not_sent')
		)
		select
			*
		from cm c
		where
			(c.approval_group = 1 and c.domain_status in ('new', 'active') and c.status = 'added')
			or (c.approval_group = 2 and c.domain_status in ('new', 'active') and c.status in ('added', 'no', 'not_scanned'))
			or (c.approval_group = 3 and c.domain_status = 'active' and c.status = 'added')
		order by c.id;
	`

	cmTable, err := a.getAdsTxtTableByQueryWithUsersFullNames(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cm table: %w", err)
	}

	return cmTable, nil
}

func (a *AdsTxtService) GetMBAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) ([]*dto.AdsTxt, error) {
	query := `
		select
			demand_partner_name_extended,
			seat_owner_name,
			score,
			ads_txt_line
		from (
			-- seat_owners
			select
				so.seat_owner_name || ' - Direct' as demand_partner_name_extended,
				so.seat_owner_name,
				case
					when so.seat_owner_name in ('OMS', 'Brightcom') then 0
					else min(score)
				end as score,
				so.seat_owner_domain || ', ' ||
					replace(so.publisher_account, '%s', 'XXXXX') || ', ' ||
					'DIRECT' ||
				case
					when so.certification_authority_id is not null
					then ', ' || so.certification_authority_id
				else ''
				end as ads_txt_line,
				true as active,
				true as is_seat_owner
			from seat_owner so
			join dpo d on so.id = d.seat_owner_id
			where d.active
			group by so.seat_owner_name, so.seat_owner_domain, so.publisher_account, so.certification_authority_id
			union
			-- demand_partners
			select
				d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name_extended,
				coalesce(so.seat_owner_name, d.demand_partner_name),
				d.score,
				d.dp_domain || ', ' ||
					dpc.publisher_account || ', ' ||
					case
						when dpc.is_direct then 'DIRECT'
						else 'RESELLER'
					end ||
					case
						when d.certification_authority_id is not null
						then ', ' || d.certification_authority_id
					else ''
				end as ads_txt_line,
				d.active,
				false as is_seat_owner
			from dpo d
			join demand_partner_connection dpc on d.demand_partner_id = dpc.demand_partner_id
			left join seat_owner so on d.seat_owner_id = so.id
			where dpc.is_required_for_ads_txt
			union
			-- children
			select
				d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name_extended,
				coalesce(so.seat_owner_name, dpc.dp_child_name),
				d.score,
				dpc.dp_child_domain || ', ' ||
					dpc.publisher_account || ', ' ||
					case
						when dpc.is_direct then 'DIRECT'
						else 'RESELLER'
					end ||
					case
						when dpc.certification_authority_id is not null
						then ', ' || dpc.certification_authority_id
					else ''
				end as ads_txt_line,
				true as active,
				false as is_seat_owner
			from dpo d
			join demand_partner_connection dpc2 on d.demand_partner_id = dpc2.demand_partner_id
			join demand_partner_child dpc on dpc2.id = dpc.dp_connection_id
			left join seat_owner so on d.seat_owner_id = so.id
			where dpc.is_required_for_ads_txt
		) as t
		where active
		order by t.score, t.is_seat_owner, t.demand_partner_name_extended;
	`

	var mbTable []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &mbTable)
	if err != nil {
		return nil, err
	}

	for i := range mbTable {
		mbTable[i].ID = i + 1
	}

	return mbTable, nil
}

func (a *AdsTxtService) UpdateAdsTxt(ctx context.Context, data *dto.AdsTxtUpdateRequest) error {
	mod, err := models.AdsTXTS(models.AdsTXTWhere.ID.EQ(data.ID)).One(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to get ads txt line with id [%v] to update: %w", data.ID, err)
	}

	columns := make([]string, 0, 3)
	columns = append(columns, models.AdsTXTColumns.UpdatedAt)
	mod.UpdatedAt = null.TimeFrom(time.Now())

	if mod.DomainStatus != data.DomainStatus {
		mod.DomainStatus = data.DomainStatus
		columns = append(columns, models.AdsTXTColumns.DomainStatus)
	}

	if mod.DemandStatus != data.DemandStatus {
		mod.DemandStatus = data.DemandStatus
		columns = append(columns, models.AdsTXTColumns.DemandStatus)
	}

	if len(columns) == 1 {
		return errors.New("there are no new values to update ads txt line")
	}

	_, err = mod.Update(ctx, bcdb.DB(), boil.Whitelist(columns...))
	if err != nil {
		return fmt.Errorf("failed to update ads txt line: %w", err)
	}

	return nil
}

func (a *AdsTxtService) getAdsTxtTableByQueryWithUsersFullNames(ctx context.Context, query string) ([]*dto.AdsTxt, error) {
	var table []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &table)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ads txt lines: %w", err)
	}

	usersMap, err := getUsersMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	for i := range table {
		table[i].AccountManagerFullName = usersMap[table[i].AccountManagerID]
		table[i].CampaignManagerFullName = usersMap[table[i].CampaignManagerID]
		table[i].DemandManagerFullName = usersMap[table[i].DemandManagerID]
	}

	return table, nil
}

func (a *AdsTxtService) updatePerfomanceData() error {
	const (
		compassKey          = "compass"
		bothKey             = "both"
		newBidderKey        = "nb"
		newBidderReportPath = "/report-dashboard/report-new-bidder"
		compassReportPath   = "/report-dashboard/report-query" //nolint:gosec
	)

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth-1, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := time.Date(currentYear, currentMonth-1, 1, 23, 59, 59, 0, currentLocation).AddDate(0, 1, -1)

	lpPartitions := fmt.Sprintf(
		"%s-%s",
		firstOfMonth.Format("2006-01-02"),
		lastOfMonth.Format("2006-01-02"),
	)

	if a.lowPerformanceCache.partitions != lpPartitions {
		a.lowPerformanceCache.partitions = lpPartitions
		a.lowPerformanceCache.lowPerformance = nil
		a.lowPerformanceCache.revenuePerDP = nil
		a.lowPerformanceCache.totalRevenue = nil
	}

	if a.lowPerformanceCache.lowPerformance == nil {
		body := buildGetLowPerfomanceRequestBody(firstOfMonth, lastOfMonth)
		var (
			compassResult, nbResult ReportResult
			errGroup                errgroup.Group
		)

		errGroup.Go(func() error {
			data, err := a.compassModule.Request(compassReportPath, http.MethodPost, body, true)
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &compassResult)
			if err != nil {
				return err
			}

			return nil
		})

		errGroup.Go(func() error {
			data, err := a.compassModule.Request(newBidderReportPath, http.MethodPost, body, true)
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &nbResult)
			if err != nil {
				return err
			}

			return nil
		})

		err := errGroup.Wait()
		if err != nil {
			return err
		}

		totalDomainsRevenue := make(map[string]float64)
		totalDomainsRevenuePerDP := make(map[string]float64)
		totalRevenuePerDP := make(map[string]float64)

		// for _, row := range compassResult.Data.Result {
		// 	domain := strings.ToLower(row.Domain)
		// 	demandPartner := strings.ToLower(row.DemandPartner)
		// 	ckey := fmt.Sprintf("%s_%s_%s", domain, demandPartner, compassKey)
		// 	bkey := fmt.Sprintf("%s_%s_%s", domain, demandPartner, bothKey)

		// 	totalDomainsRevenue[domain] += row.Revenue

		// 	totalDomainsRevenuePerDP[ckey] += row.Revenue
		// 	totalDomainsRevenuePerDP[bkey] += row.Revenue

		// 	totalRevenuePerDP[fmt.Sprintf("%s_%s", demandPartner, compassKey)] += row.Revenue
		// 	totalRevenuePerDP[fmt.Sprintf("%s_%s", demandPartner, bothKey)] += row.Revenue
		// }

		for _, row := range nbResult.Data.Result {
			domain := strings.ToLower(row.Domain)
			demandPartner := strings.ToLower(row.DemandPartner)
			nbkey := fmt.Sprintf("%s_%s_%s", domain, demandPartner, newBidderKey)
			// bkey := fmt.Sprintf("%s_%s_%s", domain, demandPartner, bothKey)

			totalDomainsRevenue[domain] += row.Revenue

			totalDomainsRevenuePerDP[nbkey] += row.Revenue
			// totalDomainsRevenuePerDP[bkey] += row.Revenue

			totalRevenuePerDP[fmt.Sprintf("%s_%s", demandPartner, newBidderKey)] += row.Revenue
			// totalRevenuePerDP[fmt.Sprintf("%s_%s", demandPartner, bothKey)] += row.Revenue
		}

		a.lowPerformanceCache.totalRevenue = map[string]float64{
			compassKey:   0,
			bothKey:      0,
			newBidderKey: 0,
		}

		for key, value := range totalRevenuePerDP {
			if strings.HasSuffix(key, bothKey) {
				a.lowPerformanceCache.totalRevenue[bothKey] += value
			} else if strings.HasSuffix(key, compassKey) {
				a.lowPerformanceCache.totalRevenue[compassKey] += value
			} else {
				a.lowPerformanceCache.totalRevenue[newBidderKey] += value
			}
		}

		lowPerformance := make(map[string]bool)
		for key, value := range totalDomainsRevenuePerDP {
			domain := strings.Split(key, "_")[0]
			lowPerformance[key] = (value/totalDomainsRevenue[domain])*100 < 1
		}

		a.lowPerformanceCache.lowPerformance = lowPerformance
		a.lowPerformanceCache.revenuePerDP = totalRevenuePerDP
	}

	return nil
}

func buildGetLowPerfomanceRequestBody(firstOfMonth, lastOfMonth time.Time) []byte {
	return []byte(fmt.Sprintf(`
		{
			"data": {
				"date": {
					"range": [
						"%v",
						"%v"
					],
					"interval": "month"
				},
				"dimensions": [
					"Domain", "DemandPartner"
				],
				"metrics": [
					"Revenue"
				]
			}
		}`,
		firstOfMonth.Format("2006-01-02 00:00:00"),
		lastOfMonth.Format("2006-01-02 15:04:05"),
	))
}

func getUsersMap(ctx context.Context) (map[string]string, error) {
	users, err := models.Users().All(ctx, bcdb.DB())
	if err != nil {
		return nil, err
	}

	usersMap := make(map[string]string, len(users))
	for _, user := range users {
		userID := strconv.Itoa(user.ID)
		usersMap[userID] = user.FirstName + " " + user.LastName
	}

	return usersMap, nil
}
