create table impression_log_hourly
(
    time         timestamp   not null,
    publisher_id varchar(36) not null,
    demand_partner_id varchar(36) not null,
    domain varchar(256) not null default '-',
    os varchar(64) not null default '-',
    country varchar(64) not null default '-',
    device_type varchar(64) not null default '-',
    size  varchar(36) not null default '-',
    is_first bool not null,
    had_followup bool not null,
    sold_impressions  int8 not null default 0,
    pub_impressions  int8 not null default 0,
    cost double precision not null default 0,
    revenue double precision not null default 0,
    demand_partner_fees double precision not null default 0,
    primary key (time, publisher_id,demand_partner_id,domain,os,country,device_type,size,is_first,had_followup)
);


create table impression_log_daily
(
    time         timestamp   not null,
    publisher_id varchar(36) not null,
    demand_partner_id varchar(36) not null,
    domain varchar(256) not null default '-',
    os varchar(64) not null default '-',
    country varchar(64) not null default '-',
    device_type varchar(64) not null default '-',
    size  varchar(36) not null default '-',
    is_first bool not null,
    had_followup bool not null,
    sold_impressions  int8 not null default 0,
    pub_impressions  int8 not null default 0,
    cost double precision not null default 0,
    revenue double precision not null default 0,
    demand_partner_fees double precision not null default 0,
    primary key (time, publisher_id,demand_partner_id,domain,os,country,device_type,size,is_first,had_followup)
);

