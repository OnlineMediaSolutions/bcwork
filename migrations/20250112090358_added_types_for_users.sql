-- +goose Up
-- +goose StatementBegin
alter table if exists "user"
add column if not exists types varchar(64)[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists "user"
drop column if exists types;
-- +goose StatementEnd
