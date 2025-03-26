-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS missing_sellers
ADD column IF NOT EXISTS yesterdaybackup TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
