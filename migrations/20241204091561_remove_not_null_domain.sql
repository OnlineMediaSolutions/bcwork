-- +goose Up
-- +goose StatementBegin
alter table if exists bid_caching
alter column domain drop not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table if exists bid_caching
alter column domain set not null;
-- +goose StatementEnd