-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS missing_sellers (
                                               name varchar(100) NOT NULL,
    url varchar(100) NOT NULL,
    sellers text,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (name)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS missing_sellers;
-- +goose StatementEnd