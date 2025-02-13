-- +goose Up
-- +goose StatementBegin
alter table if exists dpo
drop column if exists is_required_for_ads_txt,
drop column if exists is_direct,
add column if not exists integration_type varchar(64)[],
add constraint dpo_seat_owner_id_fkey foreign key (seat_owner_id) references seat_owner(id);

alter table if exists demand_partner_connection
add column if not exists is_required_for_ads_txt bool not null default false,
add column if not exists is_direct bool not null default false,
drop column if exists active;

alter table if exists demand_partner_connection
rename column "integration_type" to media_type;

alter table if exists demand_partner_child
drop column if exists dp_parent_id,
add column if not exists dp_connection_id int not null references demand_partner_connection(id),
drop column if exists active;

create table if not exists ads_txt 
(
	id serial primary key, 
    demand_partner_connection_id int references demand_partner_connection(id), 
    demand_partner_child_id int references demand_partner_child(id), 
    seat_owner_id int references seat_owner(id), 
    publisher_id varchar(64) not null references publisher(publisher_id), 
    domain varchar(256) not null, 
    status varchar(64) not null default 'not_scanned', 
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
add column if not exists dp_parent_id varchar(64) not null references dpo(demand_partner_id),
add column if not exists active bool not null default true;

alter table if exists demand_partner_connection
drop column if exists is_required_for_ads_txt,
drop column if exists is_direct,
add column if not exists active bool not null default true;

alter table if exists demand_partner_connection
rename column media_type to "integration_type";

alter table if exists dpo
add column if not exists is_required_for_ads_txt bool not null default false,
add column if not exists is_direct bool not null default false,
drop column if exists integration_type,
drop constraint dpo_seat_owner_id_fkey;
-- +goose StatementEnd
