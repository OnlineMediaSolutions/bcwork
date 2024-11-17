create materialized view global_search_view as
select 
	'Publisher list' as section_type, 
	p.publisher_id, 
	p."name" as publisher_name, 
	null as "domain", 
	null as demand_partner_name
from publisher p 
union
select 
	'Publisher / domain list' as section_type, 
	pd.publisher_id, 
	p."name" as publisher_name, 
	pd."domain", 
	null as demand_partner_name
from publisher_domain pd 
join publisher p on p.publisher_id = pd.publisher_id 
union
select 
	'Publisher / domain - Dashboard' as section_type, 
	pd.publisher_id, 
	p."name" as publisher_name, 
	pd."domain", 
	null as demand_partner_name
from publisher_domain pd 
join publisher p on p.publisher_id = pd.publisher_id 
union
select 
	'Targeting - Bidder' as section_type, 
	f.publisher as publisher_id, 
	p."name" as publisher_name, 
	f."domain", 
	null as demand_partner_name
from factor f
join publisher p on p.publisher_id = f.publisher
union
select 
	'Targeting - JS' as section_type, 
	t.publisher_id, 
	p."name" as publisher_name, 
	t."domain", 
	null as demand_partner_name
from targeting t
join publisher p on p.publisher_id = t.publisher_id
union
select 
	'Floors' as section_type, 
	f.publisher as publisher_id, 
	p."name" as publisher_name, 
	f."domain", 
	null as demand_partner_name
from floor f
join publisher p on p.publisher_id = f.publisher
union
select 
	'Publisher /domain - Demand' as section_type, 
	pd.publisher_id, 
	p."name" as publisher_name, 
	pd."domain", 
	d.demand_partner_name 
from publisher_demand pd 
join publisher p on p.publisher_id = pd.publisher_id
join dpo d on pd.demand_partner_id = d.demand_partner_id
union
select 
	'Demand / Publisher / Domain - DPO' as section_type, 
	dr.publisher as publisher_id, 
	p."name" as publisher_name, 
	dr."domain", 
	d.demand_partner_name 
from dpo_rule dr 
join dpo d on dr.demand_partner_id = d.demand_partner_id 
left join publisher p on dr.publisher = p.publisher_id 
union
select 
	'Demand - Demand' as section_type, 
	null as publisher_id, 
	null as publisher_name, 
	null as "domain", 
	d.demand_partner_name 
from dpo d;