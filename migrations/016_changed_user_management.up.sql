drop table "user_platform_role";
drop table "auth";
drop table "user";

create table "user"
(
    id serial primary key,
    user_id varchar(256) not null,
    email varchar(256) unique not null,
    first_name varchar(256) not null,
    last_name varchar(256) not null,
    role varchar(64) not null,
    organization_name varchar(128) not null,
    address varchar(128),
    phone varchar(32),
    enabled bool not null default true,
    password_changed bool not null default false,
    created_at timestamp not null,
    disabled_at timestamp
);