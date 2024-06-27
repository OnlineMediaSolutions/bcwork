create table factor
(
    publisher varchar(64) not null,
    domain varchar(255) not null,
    device varchar(64) not null default '-',
    factor float8 not null default 0,
    country varchar(64) not null default '-',
    primary key (publisher, domain)
);

