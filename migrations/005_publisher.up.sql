




-- general user table
create table publisher
(
    publisher_id varchar(36) not null primary key,
    created_at timestamp not null,
    name varchar(1024) unique not null,
    account_manager_id varchar(36),
    media_buyer_id varchar(36),
    campaign_manager_id varchar(36),
    office_location varchar(36),
    pause_timestamp int8,
    start_timestamp int8,
    reactivate_timestamp int8,
    status varchar(36)
);

create table publisher_domain
(
    name         varchar(256) not null,
    publisher_id varchar(36)  not null references publisher (publisher_id),
    created_at timestamp not null,
    primary key (name, publisher_id)
);


create table confiant
(
    confiant_key  varchar(256) not null primary key,
    publisher_id varchar(36) not null references publisher(publisher_id),
    domain varchar(256),
    rate double precision not null default 1,
    created_at timestamp not null,
    updated_at timestamp
);


