-- +goose Up
-- +goose StatementBegin
create table if not exists metadata_queue
(
    transaction_id varchar(36) primary key not null,
    key varchar(256) not null,
    version varchar(16),
    value jsonb not null,
    commited_instances int8 not null,
    created_at timestamp not null,
    updated_at timestamp
);

create table if not exists metadata_queue_temp
(
    transaction_id varchar(36) primary key not null,
    key varchar(256) not null,
    version varchar(16),
    value jsonb not null,
    commited_instances int8 not null,
    created_at timestamp not null,
    updated_at timestamp
);

create table if not exists metadata_instance
(
    instance_id varchar(64) primary key not null,
    bitwise int8 not null,
    type varchar(16) not null,
    config jsonb
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists metadata_queue;
drop table if exists metadata_queue_temp;
drop table if exists metadata_instance;
-- +goose StatementEnd
