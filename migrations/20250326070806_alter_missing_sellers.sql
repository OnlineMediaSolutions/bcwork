-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS missing_sellers
ADD column IF NOT EXISTS yesterdaybackup TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS missing_sellers
ALTER column yesterdaybackup DROP NOT NULL,
-- +goose StatementEnd
