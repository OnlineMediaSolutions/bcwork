
create table factor
(
    publisher varchar(64) references publisher(publisher_id) not null,
    domain varchar(256) not null,
    country varchar(64) not null,
    device varchar(64) not null,
    factor float8 not null default 0,
    created_at timestamp not null,
    updated_at timestamp,
    rule_id varchar(36) not null,
    demand_partner_id varchar(64) not null default '',
    browser varchar(64) not null,
    os varchar(64) not null,
    placement_type varchar(64) not null,
    primary key (rule_id)
);
