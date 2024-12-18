-- +goose Up
-- +goose StatementBegin
create table if not exists history
(
    id serial primary key,
    user_id int not null,
    subject varchar(64) not null,
    item text not null,
    publisher_id varchar(64),
    domain varchar(64),
    entity_id varchar(64),
    action varchar(64) not null,
    old_value jsonb,
    new_value jsonb,
    changes jsonb,
    date timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists history;
-- +goose StatementEnd