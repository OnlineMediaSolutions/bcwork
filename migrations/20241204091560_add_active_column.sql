-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS factor
ADD COLUMN IF NOT EXISTS active bool not null default true;

ALTER TABLE IF EXISTS floor
ADD COLUMN IF NOT EXISTS active bool not null default true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS floor
DROP COLUMN IF EXISTS active;

ALTER TABLE IF EXISTS factor
DROP COLUMN IF EXISTS active;
-- +goose StatementEnd