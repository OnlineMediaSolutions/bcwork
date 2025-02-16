-- +goose Up
-- +goose StatementBegin
create table if not exists no_dp_response_report
(
    time                   varchar(64) not null,
    demand_partner_id      varchar(64) not null,
    publisher_id           varchar(64) not null,
    domain                 varchar(64) not null,
    bid_requests           float8 not null,
    constraint no_dp_response_report_pk primary key (time, demand_partner_id, publisher_id, domain)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists no_dp_response_report;
-- +goose StatementEnd
