-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS missing_sellers
ADD column IF NOT EXISTS yesterdaybackup TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS missing_sellers
DROP column IF EXISTS yesterdaybackup;
-- +goose StatementEnd
