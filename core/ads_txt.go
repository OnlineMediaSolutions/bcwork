package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"golang.org/x/net/context"
)

type AdsTxtService struct {
	historyModule history.HistoryModule
}

func NewAdsTxtService(historyModule history.HistoryModule) *AdsTxtService {
	return &AdsTxtService{
		historyModule: historyModule,
	}
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
			join dpo d on d.demand_partner_id = dpc.dp_parent_id
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

	var mainTable []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &mainTable)
	if err != nil {
		return nil, err
	}

	usersMap, err := getUsersMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	for i := range mainTable {
		mainTable[i].AccountManagerFullName = usersMap[mainTable[i].AccountManagerID.String]
		mainTable[i].CampaignManagerFullName = usersMap[mainTable[i].CampaignManagerID.String]
		mainTable[i].DemandManagerFullName = usersMap[mainTable[i].DemandManagerID.String]
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
			end) over (partition by t.publisher_id, t."domain", t.demand_partner_name) as added, 
			count(t.status) over (partition by t.publisher_id, t."domain", t.demand_partner_name) as total,
			bool_and(case 
				when t.status = 'added' AND t.is_required and t.demand_status = 'approved' then true
				when not t.is_required then true
				else false
			end) over (partition by t.publisher_id, t."domain", t.demand_partner_name) as is_ready_to_go_live
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
			join dpo d on d.demand_partner_id = dpc.dp_parent_id
			union all
			select 
				at2.id,
				at2.publisher_id,
				at2."domain",
				at2.domain_status,
				d.demand_partner_name,
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
		) as t
		join publisher p on t.publisher_id = p.publisher_id 
		where t.active
		order by t.publisher_id, t."domain", t.demand_partner_name;
	`

	var rawTable []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &rawTable)
	if err != nil {
		return nil, err
	}

	usersMap, err := getUsersMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	groupByDpTable := make(map[string]*dto.AdsTxtGroupedByDPData)
	for _, row := range rawTable {
		dpData, ok := groupByDpTable[row.DemandPartnerName]
		if !ok {
			row.AccountManagerFullName = usersMap[row.AccountManagerID.String]
			row.CampaignManagerFullName = usersMap[row.CampaignManagerID.String]
			row.DemandManagerFullName = usersMap[row.DemandManagerID.String]

			groupByDpTable[row.DemandPartnerName] = &dto.AdsTxtGroupedByDPData{
				Parent:   row,
				Children: []*dto.AdsTxt{row},
			}
		} else {
			dpData.Children = append(dpData.Children, row)
		}
	}

	return groupByDpTable, nil
}

// TODO:
func (a *AdsTxtService) GetAMAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) ([]*dto.AdsTxt, error) {
	return nil, nil
}

// TODO:
func (a *AdsTxtService) GetCMAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) ([]*dto.AdsTxt, error) {
	return nil, nil
}

// GetMBAdsTxtTable
// Main goal of this table: we show ALL active DP required lines and active - Direct lines (it means this direct has active DP)
// Order by score.
// Where - Direct lines has score of their “best” performed DP and follow exactly after it.
// Also, OMS and Brightcom lines always with score 0 - because its our main lines and we add them always.
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
					'DIRECT'
				as ads_txt_line,
				true as active,
				true as is_seat_owner
			from seat_owner so
			join dpo d on so.id = d.seat_owner_id
			where d.active
			group by so.seat_owner_name, so.seat_owner_domain, so.publisher_account
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
				dpc.active,
				false as is_seat_owner
			from dpo d 
			join demand_partner_child dpc on d.demand_partner_id = dpc.dp_parent_id
			left join seat_owner so on d.seat_owner_id = so.id
			where dpc.is_required_for_ads_txt
		)
		where active
		order by score, is_seat_owner, demand_partner_name_extended;
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
