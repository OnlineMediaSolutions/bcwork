CREATE TABLE publisher_demand
(
    publisher_id      varchar(64)                                    not null,
    domain            varchar(256),
    demand_partner_id varchar(64) references dpo (demand_partner_id) not null,
    ads_txt_status    bool                                           not null default false,
    active            bool                                           not null default true,
    created_at        timestamp                                      not null,
    updated_at        timestamp,
    primary key (publisher_id, domain, demand_partner_id)
);