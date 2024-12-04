-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS global_factor 
DROP COLUMN IF EXISTS created_by_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS global_factor 
ADD COLUMN IF NOT EXISTS created_by_id VARCHAR(36);
-- +goose StatementEnd