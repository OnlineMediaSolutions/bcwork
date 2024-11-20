create materialized view search_view as 
select 
	'Publisher list' as section_type, 
	p.publisher_id, 
	p."name" as publisher_name, 
	null as "domain", 
	null as demand_partner_name, 
	coalesce(p.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(null, '') || ':' || coalesce(null, '') as query 
from publisher p 
union 
select 
	'Publisher / domain list' as section_type, 
	pd.publisher_id, 
	p."name" as publisher_name, 
	pd."domain", 
	null as demand_partner_name, 
	coalesce(pd.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(pd."domain", '') || ':' || coalesce(null, '') as query 
from publisher_domain pd 
join publisher p on p.publisher_id = pd.publisher_id 
union 
select 
	'Publisher / domain - Dashboard' as section_type, 
	pd.publisher_id, 
	p."name" as publisher_name, 
	pd."domain", 
	null as demand_partner_name, 
	coalesce(pd.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(pd."domain", '') || ':' || coalesce(null, '') as query 
from publisher_domain pd 
join publisher p on p.publisher_id = pd.publisher_id 
union 
select 
	'Targeting - Bidder' as section_type, 
	f.publisher as publisher_id, 
	p."name" as publisher_name, 
	f."domain", 
	null as demand_partner_name, 
	coalesce(f.publisher, '') || ':' || coalesce(p."name", '') || ':' || coalesce(f."domain", '') || ':' || coalesce(null, '') as query 
from factor f 
join publisher p on p.publisher_id = f.publisher 
union 
select 
	'Targeting - JS' as section_type, 
	t.publisher_id, 
	p."name" as publisher_name, 
	t."domain", 
	null as demand_partner_name, 
	coalesce(t.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(t."domain", '') || ':' || coalesce(null, '') as query 
from targeting t 
join publisher p on p.publisher_id = t.publisher_id 
union 
select 
	'Floors' as section_type, 
	f.publisher as publisher_id, 
	p."name" as publisher_name, 
	f."domain", 
	null as demand_partner_name, 
	coalesce(f.publisher, '') || ':' || coalesce(p."name", '') || ':' || coalesce(f."domain", '') || ':' || coalesce(null, '') as query 
from floor f 
join publisher p on p.publisher_id = f.publisher 
union 
select 
	'Publisher / domain - Demand' as section_type, 
	pd.publisher_id, 
	p."name" as publisher_name, 
	pd."domain", 
	d.demand_partner_name, 
	coalesce(pd.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(pd."domain", '') || ':' || coalesce(d.demand_partner_name, '') as query 
from publisher_demand pd 
join publisher p on p.publisher_id = pd.publisher_id 
join dpo d on pd.demand_partner_id = d.demand_partner_id 
union 
select 
	'DPO Rule' as section_type, 
	dr.publisher as publisher_id, 
	p."name" as publisher_name, 
	dr."domain", 
	d.demand_partner_name, 
	coalesce(dr.publisher, '') || ':' || coalesce(p."name", '') || ':' || coalesce(dr."domain", '') || ':' || coalesce(d.demand_partner_name, '') as query 
from dpo_rule dr 
join dpo d on dr.demand_partner_id = d.demand_partner_id 
left join publisher p on dr.publisher = p.publisher_id 
union 
select 
	'Demand - Demand' as section_type, 
	null as publisher_id, 
	null as publisher_name, 
	null as "domain", 
	d.demand_partner_name, 
	coalesce(null, '') || ':' || coalesce(null, '') || ':' || coalesce(null, '') || ':' || coalesce(d.demand_partner_name, '') as query 
from dpo d;