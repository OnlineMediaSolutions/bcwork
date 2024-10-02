CREATE TABLE ads_txt
(
    publisher_id        varchar(64) references publisher (publisher_id) not null,
    domain              varchar(256),
    demand_partner_name varchar(128),
    active              bool                                            not null default false,
    created_at          timestamp                                       not null,
    updated_at          timestamp,
    primary key (publisher_id, domain, demand_partner_name)
);