-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS missing_publishers_sellers (
  name varchar(100) NOT NULL,
  url varchar(100) NOT NULL,
  sellers varchar(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS missing_publishers_sellers;
-- +goose StatementEnd