-- +goose Up
-- +goose StatementBegin
create table if not exists report_update
(
    report    varchar(128) not null primary key,
    update_at timestamp    not null
);

create table if not exists demand_hourly
(
    time         timestamp   not null,
    demand_partner_id varchar(36) not null,
    publisher_id varchar(36) not null,
    domain varchar(256) not null default '-',
    bid_request int8 not null default 0,
    bid_response int8 not null default 0,
    bid_price  double precision not null default 0,
    impression  int8 not null default 0,
    revenue double precision not null default 0,
    demand_partner_fee double precision not null default 0,
    data_fee double precision not null default 0,
    primary key (time, publisher_id,demand_partner_id,domain)
);

create table if not exists demand_daily
(
    time         timestamp   not null,
    demand_partner_id varchar(36) not null,
    publisher_id varchar(36) not null,
    domain varchar(256) not null default '-',
    bid_request int8 not null default 0,
    bid_response int8 not null default 0,
    bid_price  double precision not null default 0,
    impression  int8 not null default 0,
    revenue double precision not null default 0,
    demand_partner_fee double precision not null default 0,
    data_fee double precision not null default 0,
    primary key (time, publisher_id,demand_partner_id,domain)
);

create table if not exists publisher_hourly
(
    time         timestamp   not null,
    publisher_id varchar(36) not null,
    domain varchar(256) not null default '-',
    os varchar(64) not null default '-',
    country varchar(64) not null default '-',
    device_type varchar(64) not null default '-',
    bid_requests int8 not null default 0,
    bid_responses int8 not null default 0,
    bid_price_count int8 not null default 0,
    bid_price_total double precision not null default 0,
    publisher_impressions  int8 not null default 0,
    demand_impressions  int8 not null default 0,
    missed_opportunities  int8 not null default 0,
    supply_total double precision not null default 0,
    demand_total double precision not null default 0,
    demand_partner_fee double precision not null default 0,
    primary key (time, publisher_id,domain,os,country,device_type)
);

create table if not exists publisher_daily
(
    time         timestamp   not null,
    publisher_id varchar(36) not null,
    domain varchar(256) not null default '-',
    os varchar(64) not null default '-',
    country varchar(64) not null default '-',
    device_type varchar(64) not null default '-',
    bid_requests int8 not null default 0,
    bid_responses int8 not null default 0,
    bid_price_count int8 not null default 0,
    bid_price_total double precision not null default 0,
    publisher_impressions  int8 not null default 0,
    demand_impressions  int8 not null default 0,
    missed_opportunities  int8 not null default 0,
    supply_total double precision not null default 0,
    demand_total double precision not null default 0,
    demand_partner_fee double precision not null default 0,
    primary key (time, publisher_id,domain,os,country,device_type)
);

create table if not exists compass_publisher_tag
(
    id serial primary key,
    publisher_id varchar(36) not null,
    device_type varchar(64) not null,
    domain varchar(256) not null,
    constraint compass_publisher_tag_pk unique (publisher_id, device_type, domain)
);

create table if not exists iiq_testing
(
    time         timestamp   not null,
    demand_partner_id varchar(36)  not null,
    iiq_requests int8 not null,
    non_iiq_requests int8 not null,
    iiq_impressions int8 not null,
    non_iiq_impressions int8 not null,
    primary key (time,demand_partner_id)
);

create table if not exists id5_testing
(
    time         timestamp   not null,
    demand_partner_id varchar(36)  not null,
    id5_requests int8 not null,
    non_id5_requests int8 not null,
    id5_impressions int8 not null,
    non_id5_impressions int8 not null,
    primary key (time,demand_partner_id)
);

create table if not exists nb_supply_hourly
(
    time         timestamp   not null,
    publisher_id varchar(36) not null,
    domain varchar(256) not null default '-',
    os varchar(64) not null default '-',
    country varchar(64) not null default '-',
    device_type varchar(64) not null default '-',
    placement_type varchar(16) not null default '-',
    size varchar(16) not null default '-',
    request_type varchar(16) not null default '-',
    payment_type varchar(16) not null default '-',
    datacenter varchar(16) not null default '-',
    bid_requests  int8 not null default 0,
    bid_responses  int8 not null default 0,
    sold_impressions  int8 not null default 0,
    publisher_impressions  int8 not null default 0,
    cost double precision not null default 0,
    revenue double precision not null default 0,
    avg_bid_price double precision not null default 0,
    missed_opportunities int8 not null default 0,
    demand_partner_fee double precision not null default 0,
    data_impressions         int8 not null default 0,
    data_fee                  double precision not null default 0,
    primary key (time, publisher_id,domain,os,country,device_type,placement_type,size,request_type,payment_type,datacenter)
);

create table if not exists nb_demand_hourly
(
    time                     timestamp        not null,
    demand_partner_id        varchar(36)      not null,
    demand_partner_placement_id        varchar(36)      not null,
    publisher_id             varchar(36)      not null,
    domain                   varchar(256)     not null default '-',
    os                       varchar(64)      not null default '-',
    country                  varchar(64)      not null default '-',
    device_type              varchar(64)      not null default '-',
    placement_type           varchar(16)      not null default '-',
    size                     varchar(16) not null default '-',
    request_type             varchar(16) not null default '-',
    payment_type varchar(16) not null default '-',
    datacenter varchar(16) not null default '-',
    bid_requests             int8             not null default 0,
    bid_responses           int8             not null default 0,
    avg_bid_price            double precision not null default 0,
    dp_fee                   double precision not null default 0,
    auction_wins             int8             not null default 0,
    auction                  double precision not null default 0,
    sold_impressions         int8             not null default 0,
    revenue                  double precision not null default 0,
    data_impressions         int8 not null default 0,
    data_fee                  double precision not null default 0,
    primary key (time, demand_partner_id, demand_partner_placement_id,publisher_id, domain, os, country, device_type,placement_type,size,request_type,payment_type,datacenter)
);

create table if not exists demand_partner
(
    demand_partner_id varchar(36) not null primary key,
    name varchar(128) not null,
    integration_type varchar(36) not null
);

create table if not exists demand_parnter_placement
(
    demand_partner_placement_id varchar(36) not null primary key,
    demand_partner_id varchar(36) not null references  demand_partner(demand_partner_id) on delete cascade ,
    name varchar(256) not null
);

create table if not exists revenue_hourly
(
    time         timestamp   primary key not null,
    publisher_impressions  int8 not null default 0,
    sold_impressions  int8 not null default 0,
    cost double precision not null default 0,
    revenue double precision not null default 0,
    demand_partner_fees double precision not null default 0,
    missed_opportunities int8 not null default 0,
    data_fee double precision not null default 0,
    dp_bid_requests int8 not null default 0
);

create table if not exists revenue_daily
(
    time         timestamp   primary key not null,
    publisher_impressions  int8 not null default 0,
    sold_impressions  int8 not null default 0,
    cost double precision not null default 0,
    revenue double precision not null default 0,
    demand_partner_fees double precision not null default 0,
    missed_opportunities int8 not null default 0,
    data_fee double precision not null default 0,
    dp_bid_requests int8 not null default 0
);

create table if not exists iiq_hourly
(
    time         timestamp   not null,
    dpid varchar(36)  not null,
    datacenter varchar(16) not null default '-',
    request int8 not null default 0,
    response int8 not null default 0,
    impression int8 not null default 0,
    revenue float8 not null default 0,
    primary key (time,dpid,datacenter)
);

create table if not exists iiq_daily
(
    time         timestamp   not null,
    dpid varchar(36)  not null,
    datacenter varchar(16) not null default '-',
    request int8 not null default 0,
    response int8 not null default 0,
    impression int8 not null default 0,
    revenue float8 not null default 0,
    primary key (time,dpid,datacenter)
);

create table if not exists demand_partner_hourly
(
    time         timestamp   not null,
    demand_partner_id varchar(36) not null,
    domain varchar(256) not null,
    impression int8 not null default 0,
    revenue float8 not null default 0,
    primary key (time,demand_partner_id,domain)
);

create table if not exists demand_partner_daily
(
    time         timestamp   not null,
    demand_partner_id varchar(36) not null,
    domain varchar(256) not null,
    impression int8 not null default 0,
    revenue float8 not null default 0,
    primary key (time,demand_partner_id,domain)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table report_update;
drop table demand_hourly;
drop table demand_daily;
drop table publisher_hourly;
drop table publisher_daily;
drop table compass_publisher_tag;
drop table iiq_testing;
drop table id5_testing;
drop table nb_supply_hourly;
drop table nb_demand_hourly;
drop table  demand_partner;
drop table demand_parnter_placement;
drop table revenue_hourly;
drop table revenue_daily;
drop table iiq_hourly;
drop table iiq_daily;
drop table demand_partner_hourly;
drop table demand_partner_daily;
-- +goose StatementEnd
