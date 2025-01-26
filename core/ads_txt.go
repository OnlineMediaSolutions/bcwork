package core

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
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

// TODO:
func (a *AdsTxtService) GetMainAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) ([]*dto.AdsTxt, error) {
	query := `
		select 
			at2.id,
			at2.publisher_id,
			p."name" as publisher_name,
			p.account_manager_id,
			p.campaign_manager_id,
			at2."domain",
			at2.domain_status,
			d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name,
			d.manager_id as demand_manager_id,
			at2.demand_status,
			at2.status,
			dpc.is_required_for_ads_txt as is_required,
			d.dp_domain || 
				', ' || 
				replace(dpc.publisher_account, '%s', at2.publisher_id)  || -- pattern for subsidiary companies
				', ' || 
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
		join publisher p on p.publisher_id = at2.publisher_id
		join demand_partner_connection dpc on at2.demand_partner_connection_id = dpc.id 
		join dpo d on d.demand_partner_id = dpc.demand_partner_id 
		union 
		select 
			at2.id,
			at2.publisher_id,
			p."name" as publisher_name,
			p.account_manager_id,
			p.campaign_manager_id,
			at2."domain",
			at2.domain_status,
			d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name,
			d.manager_id as demand_manager_id,
			at2.demand_status,
			at2.status,
			dpc.is_required_for_ads_txt as is_required,
			dpc.dp_child_domain || 
				', ' || 
				replace(dpc.publisher_account, '%s', at2.publisher_id)  || -- pattern for subsidiary companies
				', ' || 
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
		join publisher p on p.publisher_id = at2.publisher_id
		join demand_partner_child dpc on at2.demand_partner_child_id = dpc.id 
		join dpo d on d.demand_partner_id = dpc.dp_parent_id
		union 
		select 
			at2.id,
			at2.publisher_id,
			p."name" as publisher_name,
			p.account_manager_id,
			p.campaign_manager_id,
			at2."domain",
			at2.domain_status,
			so.seat_owner_name || ' - Direct' as demand_partner_name,
			null as demand_manager_id,
			at2.demand_status,
			at2.status,
			true as is_required,
			so.seat_owner_domain || 
				', ' || 
				replace(so.publisher_account, '%s', at2.publisher_id)  || -- pattern for subsidiary companies
				', ' || 
				'DIRECT' || 
				case 
					when so.certification_authority_id is not null 
					then ', ' || so.certification_authority_id 
				else '' 
				end as ads_txt_line,
			at2.last_scanned_at,
			at2.error_message
		from ads_txt at2
		join publisher p on p.publisher_id = at2.publisher_id
		join seat_owner so on at2.seat_owner_id = so.id;
	`

	var mainTable []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &mainTable)
	if err != nil {
		return nil, err
	}

	users, err := models.Users().All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	usersMap := make(map[string]string, len(users))
	for _, user := range users {
		userID := strconv.Itoa(user.ID)
		usersMap[userID] = user.FirstName + " " + user.LastName
	}

	for i := range mainTable {
		mainTable[i].AccountManagerFullName = usersMap[mainTable[i].AccountManagerID.String]
		mainTable[i].CampaignManagerFullName = usersMap[mainTable[i].CampaignManagerID.String]
		mainTable[i].DemandManagerFullName = usersMap[mainTable[i].DemandManagerID.String]
	}

	return mainTable, nil
}

// TODO:
func (a *AdsTxtService) GetGroupByDPAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) ([]*dto.AdsTxt, error) {
	return nil, nil
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
			demand_partner_name,
			seat_owner_name,
			score,
			ads_txt_line
		from (
			-- seat_owners
			select 
				so.seat_owner_name || ' - Direct' as demand_partner_name,
				so.seat_owner_name,
				case 
					when so.seat_owner_name in ('OMS', 'Brightcom') then 0
					else min(score)
				end as score,
				so.seat_owner_domain || 
					', ' || 
					replace(so.publisher_account, '%s', 'XXXXX')  ||
					', ' || 
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
				d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name,
				coalesce(so.seat_owner_name, d.demand_partner_name),
				d.score,
				d.dp_domain || 
					', ' || 
					dpc.publisher_account ||
					', ' || 
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
				d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name,
				coalesce(so.seat_owner_name, dpc.dp_child_name),
				d.score,
				dpc.dp_child_domain || 
					', ' || 
					dpc.publisher_account ||
					', ' || 
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
		order by score, is_seat_owner, demand_partner_name;
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

// TODO:
func (a *AdsTxtService) UpdateAdsTxt(ctx context.Context, data *dto.AdsTxt) error {
	return nil
}
