-- +goose Up
-- +goose StatementBegin
create table if not exists dpo
(
    demand_partner_id varchar(64) not null primary key,
    is_include bool not null default false,
    created_at timestamp not null,
    updated_at timestamp,
    demand_partner_name varchar(128),
    active bool not null default true
);

create table if not exists dpo_rule
(
    rule_id varchar(36) not null primary key,
    demand_partner_id varchar(64) not null references dpo(demand_partner_id),
    publisher varchar(64) references publisher(publisher_id),
    domain varchar(256),
    country varchar(64),
    browser varchar(64),
    os varchar(64),
    device_type varchar(64),
    placement_type varchar(64),
    factor float8 not null default 0,
    created_at timestamp not null,
    updated_at timestamp,
    active bool not null default true
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists dpo_rule;
drop table if exists dpo;
-- +goose StatementEnd