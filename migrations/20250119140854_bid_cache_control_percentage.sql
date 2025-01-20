-- +goose Up
-- +goose StatementBegin
alter table if exists bid_caching
add column if not exists control_percentage float8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists bid_caching
drop column if exists control_percentage;
-- +goose StatementEnd
