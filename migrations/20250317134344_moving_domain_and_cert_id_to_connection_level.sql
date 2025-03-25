-- +goose Up
-- +goose StatementBegin
alter table if exists dpo
drop column if exists dp_domain,
drop column if exists certification_authority_id;

alter table if exists demand_partner_connection
add column if not exists dp_domain varchar(128) not null default '',
add column if not exists certification_authority_id varchar(256);

alter table if exists demand_partner_child
rename column dp_child_domain to dp_domain;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists demand_partner_child
rename column dp_domain to dp_child_domain;

alter table if exists demand_partner_connection
drop column if exists dp_domain,
drop column if exists certification_authority_id;

alter table if exists dpo
add column if not exists dp_domain varchar(128) not null default '',
add column if not exists certification_authority_id varchar(256);
-- +goose StatementEnd
