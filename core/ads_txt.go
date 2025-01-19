package core

import (
	"log"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/modules/history"
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
	return nil, nil
}

// TODO:
func (a *AdsTxtService) GetMBAdsTxtTable(ctx context.Context, ops *AdsTxtOptions) ([]*dto.AdsTxt, error) {
	query := `
		select 
			demand_partner_name,
			seat_owner_name,
			score,
			ads_txt_line
		from (
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
			group by so.seat_owner_name, so.seat_owner_domain, so.publisher_account
			union
			select 
				d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name,
				coalesce(so.seat_owner_name, d.demand_partner_name),
				d.score,
				d.dp_domain || 
				', ' || 
				dpc.publisher_account ||
				', ' || 
				case 
					when d.is_direct then 'DIRECT' 
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
			join demand_partner_connection dpc ON d.demand_partner_id = dpc.demand_partner_id
			left join seat_owner so on d.seat_owner_id = so.id
		)
		where active
		order by score, is_seat_owner, demand_partner_name;
	`

	log.Println(query)

	return nil, nil
}

// TODO:
func (a *AdsTxtService) UpdateAdsTxt(ctx context.Context, data *dto.AdsTxt) error {
	return nil
}
