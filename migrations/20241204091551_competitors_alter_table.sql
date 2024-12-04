-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS competitors
ADD COLUMN IF NOT EXISTS type VARCHAR(50) NOT NULL,
ADD COLUMN IF NOT EXISTS position varchar(50) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS competitors
DROP COLUMN IF EXISTS type,
DROP COLUMN IF EXISTS position;
-- +goose StatementEnd