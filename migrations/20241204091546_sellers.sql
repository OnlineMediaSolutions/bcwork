-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS competitors (
    name VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    PRIMARY KEY (name),
    CONSTRAINT unique_url UNIQUE (url)
);

CREATE TABLE IF NOT EXISTS sellers_json_history (
    competitor_name VARCHAR(100) NOT NULL,
    added_domains TEXT NOT NULL,
    added_publishers TEXT NOT NULL,
    backup_today JSONB NOT NULL,
    backup_yesterday JSONB NOT NULL,
    backup_before_yesterday JSONB NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (competitor_name),
    FOREIGN KEY (competitor_name) REFERENCES competitors (name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists competitors;
drop table if exists sellers_json_history;
-- +goose StatementEnd