-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS dpo
ADD COLUMN IF NOT EXISTS automation_name varchar(64),
ADD COLUMN IF NOT EXISTS threshold float,
ADD COLUMN IF NOT EXISTS automation boolean default false not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS dpo
DROP COLUMN IF EXISTS automation_name,
DROP COLUMN IF EXISTS threshold,
DROP COLUMN IF EXISTS automation
-- +goose StatementEnd
