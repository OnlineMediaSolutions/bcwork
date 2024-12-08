-- +goose Up
-- +goose StatementBegin
create table if not exists publisher_sync
(
    key         varchar(50)     not null primary key,
    had_error   bool            not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists publisher_sync;
-- +goose StatementEnd
