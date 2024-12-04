-- +goose Up
-- +goose StatementBegin
create table if not exists "user"
(
    user_id varchar(36) not null primary key,
    created_at timestamp not null,
    email varchar(1024) unique not null,
    first_name varchar(256),
    last_name varchar(256),
    last_activity_at timestamp,
    invited_at timestamp,
    signedup_at timestamp,
    invited_by varchar(36)
);

create table if not exists "user_platform_role"
(
    user_id varchar(36) not null references "user"(user_id) on delete cascade,
    role_id varchar(36) not null,
    created_at timestamp not null,
    primary key(user_id,role_id)
);

create table if not exists "auth" (
    auth_user_id varchar(36) not null primary key,
    user_id varchar(36) not null references "user"(user_id),
    impersonate_as_id varchar(36) references "user"(user_id),
    created_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "user";
drop table if exists "user_platform_role";
drop table if exists "auth";
-- +goose StatementEnd
