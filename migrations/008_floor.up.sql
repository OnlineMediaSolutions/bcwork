
create table floor
(
    publisher varchar(64) references publisher(publisher_id) not null,
    domain varchar(256) not null,
    country varchar(64) not null default '',
    device varchar(64) not null default '',
    floor float8 not null default 0,
    created_at timestamp not null,
    updated_at timestamp,
    rule_id varchar(36) not null default '',
    demand_partner_id varchar(64) not null default '',
    browser varchar(64) not null default '',
    os varchar(64) not null default '',
    placement_type varchar(64) not null default '',
    primary key (rule_id)
);
