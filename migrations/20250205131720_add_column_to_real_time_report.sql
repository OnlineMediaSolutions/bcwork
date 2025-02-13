-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS real_time_report
ADD COLUMN IF NOT EXISTS bid_responses float8 not null default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS real_time_report
DROP COLUMN IF EXISTS bid_responses;
-- +goose StatementEnd
