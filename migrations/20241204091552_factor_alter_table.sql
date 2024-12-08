-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS factor
DROP CONSTRAINT factor_pkey;

ALTER TABLE IF EXISTS factor
ALTER COLUMN country DROP NOT NULL,
ALTER COLUMN device DROP NOT NULL;

ALTER TABLE IF EXISTS factor
ADD CONSTRAINT fk_factor_publisher
FOREIGN KEY (publisher) REFERENCES publisher(publisher_id);

ALTER TABLE IF EXISTS factor
ADD PRIMARY KEY (rule_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS factor
DROP CONSTRAINT factor_pkey;

ALTER TABLE IF EXISTS factor
ALTER COLUMN country SET NOT NULL,
ALTER COLUMN device SET NOT NULL;

ALTER TABLE IF EXISTS factor
DROP CONSTRAINT fk_factor_publisher; 

ALTER TABLE IF EXISTS factor
ADD PRIMARY KEY (publisher, domain, device, country); 
-- +goose StatementEnd