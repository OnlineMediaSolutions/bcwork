-- +goose Up
-- +goose StatementBegin
alter table if exists publisher
add column if not exists is_direct bool not null default false;
alter table if exists publisher_domain
add column if not exists is_direct bool not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists publisher
drop column if exists is_direct;
alter table if exists publisher_domain
drop column if exists is_direct;
-- +goose StatementEnd
