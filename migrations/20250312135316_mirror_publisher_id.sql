-- +goose Up
-- +goose StatementBegin
alter table if exists publisher_domain
add column if not exists mirror_publisher_id varchar(36) references publisher (publisher_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists publisher_domain
drop column if exists mirror_publisher_id;
-- +goose StatementEnd
