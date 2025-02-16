-- +goose Up
-- +goose StatementBegin
alter table if exists publisher_domain
alter column integration_type type varchar(64)[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists publisher_domain
alter column integration_type type integration_type[] using (integration_type::integration_type[]); 
-- +goose StatementEnd
