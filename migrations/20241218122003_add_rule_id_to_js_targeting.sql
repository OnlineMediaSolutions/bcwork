-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS targeting
ADD COLUMN IF NOT EXISTS rule_id varchar(36) not null default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS targeting
DROP COLUMN IF EXISTS rule_id;
-- +goose StatementEnd