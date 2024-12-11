-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS floor
ALTER COLUMN country DROP NOT NULL,
ALTER COLUMN device DROP NOT NULL,
ALTER COLUMN placement_type DROP NOT NULL,
ALTER COLUMN os DROP NOT NULL,
ALTER COLUMN browser DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS floor
ALTER COLUMN country SET NOT NULL,
ALTER COLUMN device SET NOT NULL,
ALTER COLUMN placement_type SET NOT NULL,
ALTER COLUMN os SET NOT NULL,
ALTER COLUMN browser SET NOT NULL;
-- +goose StatementEnd