-- +goose Up
-- +goose StatementBegin
alter table if exists publisher
add column if not exists media_type varchar(64)[];

alter table if exists publisher
alter column integration_type type varchar(64)[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists publisher
alter column integration_type type integration_type[] using (integration_type::integration_type[]); 

alter table if exists publisher
drop column if exists media_type;
-- +goose StatementEnd
