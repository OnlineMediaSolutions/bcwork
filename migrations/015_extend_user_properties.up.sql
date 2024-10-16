ALTER TABLE "user"
ADD COLUMN id serial,
ADD COLUMN role varchar(64) not null,
ADD COLUMN organization_name varchar(128) not null,
ADD COLUMN address varchar(128) not null,
ADD COLUMN phone varchar(32) not null,
ADD COLUMN enabled bool not null default true,
ADD COLUMN disabled_at timestamp;