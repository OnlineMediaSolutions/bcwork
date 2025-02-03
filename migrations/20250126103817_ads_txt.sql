-- +goose Up
-- +goose StatementBegin
alter table if exists dpo
drop column if exists is_required_for_ads_txt,
drop column if exists is_direct;

alter table if exists demand_partner_connection
add column if not exists is_required_for_ads_txt bool not null default false,
add column if not exists is_direct bool not null default false;

alter table if exists demand_partner_child
drop column if exists dp_parent_id,
add column if not exists dp_connection_id int not null references demand_partner_connection(id);

create table if not exists ads_txt 
(
	id serial primary key, 
    demand_partner_connection_id int references demand_partner_connection(id), 
    demand_partner_child_id int references demand_partner_child(id), 
    seat_owner_id int references seat_owner(id), 
    publisher_id varchar(64) not null references publisher(publisher_id), 
    domain varchar(256) not null, 
    status varchar(64) not null default 'NOT_SCANNED', 
    demand_status varchar(64) not null default 'not_sent', 
    domain_status varchar(64) not null default 'new', 
    created_at timestamp not null, 
    updated_at timestamp, 
    status_changed_at timestamp, 
	last_scanned_at timestamp, 
	error_message varchar(256), 
	retries int, 
	valid_url varchar(128)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists ads_txt;

alter table if exists demand_partner_child
drop column if exists dp_connection_id,
add column if not exists dp_parent_id varchar(64) not null references dpo(demand_partner_id);

alter table if exists demand_partner_connection
drop column if exists is_required_for_ads_txt,
drop column if exists is_direct;

alter table if exists dpo
add column if not exists is_required_for_ads_txt bool not null default false,
add column if not exists is_direct bool not null default false;
-- +goose StatementEnd
