-- +goose Up
-- +goose StatementBegin
create table if not exists factor
(
    publisher varchar(64),
    domain varchar(256),
    country varchar(64),
    device varchar(64),
    factor float8 not null default 0,
    created_at timestamp not null,
    updated_at timestamp,
    rule_id varchar(36) not null default '',
    demand_partner_id varchar(64) not null default '',
    browser varchar(64),
    os varchar(64),
    placement_type varchar(64),
    primary key (publisher, domain, device, country)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists factor;
-- +goose StatementEnd