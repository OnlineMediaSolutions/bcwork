
create table bid_caching
(
    publisher         varchar(64)                                    not null
    references publisher,
    domain            varchar(256)                                   not null,
    country           varchar(64),
    device            varchar(64),
    bid_caching       SMALLINT                  not null,
    created_at        timestamp                                      not null,
    updated_at        timestamp,
    rule_id           varchar(36)                                    not null
    primary key,
    demand_partner_id varchar(64)      default ''::character varying not null,
    browser           varchar(64),
    os                varchar(64),
    placement_type    varchar(64),
    active bool not null default true

);

alter table bid_caching
owner to postgres;
