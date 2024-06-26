create type device_type AS ENUM ('all', 'desktop', 'mobile', 'tablet');


create table factor
(
    publisher varchar(64) not null,
    domain varchar(255) not null,
    device device_type,
    factor NUMERIC(4, 2),
    country varchar(100),
    primary key (publisher, domain)
);


