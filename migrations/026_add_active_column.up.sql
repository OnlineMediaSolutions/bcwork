ALTER TABLE factor
    ADD COLUMN active bool not null default true;

ALTER TABLE floor
    ADD COLUMN active bool not null default true;