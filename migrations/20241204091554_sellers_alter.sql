-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS sellers_json_history
ADD column IF NOT EXISTS deleted_domains TEXT NOT NULL,
ADD column IF NOT EXISTS deleted_publishers TEXT NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS sellers_json_history
DROP column IF EXISTS deleted_domains,
DROP column IF EXISTS deleted_publishers;
-- +goose StatementEnd
