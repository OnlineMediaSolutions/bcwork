-- general user table

create type integration_type as enum (
    'JS Tags (Compass)',
    'JS Tags (NP)',
    'Prebid.js',
    'Prebid Server',
    'oRTB EP'
    );

create table publisher
(
    publisher_id         varchar(36)          not null primary key,
    created_at           timestamp            not null,
    name                 varchar(1024) unique not null,
    account_manager_id   varchar(36),
    media_buyer_id       varchar(36),
    campaign_manager_id  varchar(36),
    office_location      varchar(36),
    pause_timestamp      int8,
    start_timestamp      int8,
    reactivate_timestamp int8,
    integration_type     integration_type[],
    status               varchar(36)
);


create table publisher_domain
(
    domain           varchar(256) not null,
    publisher_id     varchar(36)  not null references publisher (publisher_id),
    automation       bool         not null default false,
    gpp_target       double precision      default null,
    integration_type integration_type[],
    created_at       timestamp    not null,
    updated_at       timestamp,
    primary key (domain, publisher_id)
);


create table confiant
(
    confiant_key varchar(256)     not null,
    publisher_id varchar(36)      not null references publisher (publisher_id),
    domain       varchar(256),
    rate         double precision not null default 0,
    created_at   timestamp        not null,
    updated_at   timestamp,
    constraint PK_confiant_1 primary key (domain, publisher_id)
);

create table pixalate
(
    id           varchar(256)     not null,
    publisher_id varchar(36)      not null references publisher (publisher_id),
    domain       varchar(256),
    rate         double precision not null default 0,
    active       bool             not null,
    created_at   timestamp        not null,
    updated_at   timestamp,
    constraint PK_pixalate_1 primary key (domain, publisher_id)
);




