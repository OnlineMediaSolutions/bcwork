package core

var (
	dynamicPublisherIDPlaceholder = "%s"

	adsTxtDemandPartnerConnectionBaseQuery = `
		select 
			at2.id,
			at2.publisher_id,
			at2."domain",
			at2.domain_status,
			d.demand_partner_id,
			d.demand_partner_name,
			dpc.id as demand_partner_connection_id,
			dpc."media_type",
			d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name_extended,
			d.manager_id as demand_manager_id,
			at2.demand_status,
			at2.status,
			d.is_approval_needed,
			dpc.is_required_for_ads_txt as is_required,
			d.active as is_demand_partner_active,
			dpc.dp_domain || ', ' || 
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
		join demand_partner_connection dpc on at2.demand_partner_connection_id = dpc.id 
		join dpo d on d.demand_partner_id = dpc.demand_partner_id 
	`

	adsTxtdemandPartnerChildrenBaseQuery = `
		select 
			at2.id,
			at2.publisher_id,
			at2."domain",
			at2.domain_status,
			d.demand_partner_id,
			%s as demand_partner_name,
			dpc2.id as demand_partner_connection_id,
			dpc2."media_type",
			d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name_extended,
			d.manager_id as demand_manager_id,
			at2.demand_status,
			at2.status,
			d.is_approval_needed,
			dpc.is_required_for_ads_txt as is_required,
			d.active as is_demand_partner_active,
			dpc.dp_domain || ', ' || 
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
	`

	adsTxtSeatOwnersBaseQuery = `
		select 
			at2.id,
			at2.publisher_id,
			at2."domain",
			at2.domain_status,
			d.demand_partner_id,
			%v as demand_partner_name,
			%v as demand_partner_connection_id,
			%v as media_type,
			so.seat_owner_name || ' - Direct' as demand_partner_name_extended,
			null as demand_manager_id,
			at2.demand_status,
			at2.status,
			d.is_approval_needed,
			true as is_required,
			d.active as is_demand_partner_active,
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
		join dpo d on d.seat_owner_id = so.id 
	`
)
