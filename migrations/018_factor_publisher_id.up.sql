ALTER TABLE factor
DROP CONSTRAINT factor_pkey;

ALTER TABLE factor
    ALTER COLUMN country DROP NOT NULL,
    ALTER COLUMN device DROP NOT NULL;

ALTER TABLE factor
ADD CONSTRAINT fk_factor_publisher
FOREIGN KEY (publisher) REFERENCES publisher(publisher_id);



ALTER TABLE factor
ADD PRIMARY KEY (rule_id);
