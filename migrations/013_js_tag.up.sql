create table targeting
(
    id serial primary key,
    hash varchar(36) not null,
    rule_id varchar(36) not null,
    publisher varchar(64) references publisher(publisher_id),
    domain varchar(256),
    unit_size varchar	(64),
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
    status  varchar(64) not null
);
