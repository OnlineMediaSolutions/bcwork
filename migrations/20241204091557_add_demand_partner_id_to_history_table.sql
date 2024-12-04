-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS history
ADD COLUMN IF NOT EXISTS demand_partner_id varchar(64)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS history
DROP COLUMN IF EXISTS demand_partner_id;
-- +goose StatementEnd