create table targeting
(
    id serial primary key,
    publisher_id varchar(64) not null references publisher(publisher_id),
    domain varchar(256) not null,
    unit_size varchar(64) not null,
    placement_type varchar(64),
    country text[],
    device_type text[],
    browser text[],
    os text[],
	kv jsonb,
	price_model varchar(64) not null,
	value float8 not null,
	daily_cap int,
    created_at timestamp not null,
    updated_at timestamp,
    status varchar(64) not null
);
