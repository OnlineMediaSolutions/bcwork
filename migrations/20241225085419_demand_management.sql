-- +goose Up
-- +goose StatementBegin
alter table if exists dpo
add column if not exists dp_domain varchar(128) not null default '',
add column if not exists certification_authority_id varchar(256),
add column if not exists seat_owner_id int,
add column if not exists manager_id int references "user"(id),
add column if not exists poc_name varchar(128) not null default '',
add column if not exists poc_email varchar(128) not null default '',
add column if not exists is_direct bool not null default false,
add column if not exists is_approval_needed bool not null default false,	
add column if not exists approval_before_going_live bool not null default false,
add column if not exists approval_process varchar(64) not null default 'Other',
add column if not exists is_required_for_ads_txt bool not null default false,
add column if not exists dp_blocks varchar(64) not null default 'Other',
add column if not exists score int not null default 1000,
add column if not exists "comments" text;

create table if not exists seat_owner
(
	id serial primary key,
	seat_owner_name varchar(128) not null default '',
	seat_owner_domain varchar(128) not null default '',
	publisher_account varchar(256) not null default '%s',
	certification_authority_id varchar(256),
    created_at timestamp not null,
    updated_at timestamp
);

create table if not exists demand_partner_child
(
	id serial primary key,
	dp_parent_id varchar(64) not null references dpo(demand_partner_id),
	dp_child_name varchar(128) not null default '',
	dp_child_domain varchar(128) not null default '',
	publisher_account varchar(256) not null default '',
	certification_authority_id varchar(256),
	active bool not null default true,
	is_required_for_ads_txt bool not null default false,
	created_at timestamp not null,
    updated_at timestamp
);

create table if not exists demand_partner_connection
(
	id serial primary key,
	demand_partner_id varchar(64) not null references dpo(demand_partner_id),
	publisher_account varchar(256) not null default '',
	integration_type varchar(64)[],
	active bool not null default true,
    created_at timestamp not null,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists demand_partner_connection;
drop table if exists demand_partner_child;
drop table if exists seat_owner;

alter table if exists dpo
drop column if exists dp_domain,
drop column if exists certification_authority_id,
drop column if exists seat_owner_id,
drop column if exists manager_id,
drop column if exists poc_name,
drop column if exists poc_email,
drop column if exists is_direct,
drop column if exists is_approval_needed,	
drop column if exists approval_before_going_live,
drop column if exists approval_process,
drop column if exists is_required_for_ads_txt,
drop column if exists dp_blocks,
drop column if exists score,
drop column if exists "comments";
-- +goose StatementEnd
