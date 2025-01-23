-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists real_time_report
(
    time                 varchar(64) not null,
    publisher            varchar(64) not null,
    publisher_id         varchar(64) not null,
    domain               varchar(64) not null,
    bid_requests         float8 not null,
    device               varchar(64) not null,
    country              varchar(64) not null,
    revenue              float8 not null,
    cost                 float8 not null,
    sold_impressions     float8 not null,
    publisher_impressions float8 not null,
    pub_fill_rate        float8 not null,
    cpm                  float8 not null,
    rpm                  float8 not null,
    dp_rpm               float8 not null,
    gp                   float8 not null,
    gpp                  float8 not null,
    consultant_fee       float8 not null,
    tam_fee              float8 not null,
    tech_fee             float8 not null,
    demand_partner_fee   float8 not null,
    data_fee             float8 not null,
    constraint real_time_report_log_pk_1
    primary key (publisher_id, time, domain, device, country)
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists real_time_report;
-- +goose StatementEnd