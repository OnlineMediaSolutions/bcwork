-- +goose Up
-- +goose StatementBegin
create table if not exists dpo_automation_log
(
    time        timestamp        not null,
    eval_time   timestamp        not null,
    domain      varchar(255)     not null,
    publisher   varchar(255)     not null,
    os          varchar(50)      not null,
    country     varchar(2)       not null,
    dp          varchar(255)     not null,
    bid_request integer          not null,
    revenue     double precision not null,
    erpm        double precision not null,
    old_factor  double precision not null,
    new_factor  double precision not null,
    resp_status integer          not null,
    constraint dpo_automation_log_pk_1
        primary key (time, dp, country, publisher, domain, os)
);

CREATE INDEX IF NOT EXISTS idx_dpo_automation_log_time ON dpo_automation_log(time);
CREATE INDEX IF NOT EXISTS idx_dpo_automation_log_domain ON dpo_automation_log(domain);
CREATE INDEX IF NOT EXISTS idx_dpo_automation_log_dp ON dpo_automation_log(dp);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists dpo_automation_log;
drop index if exists idx_dpo_automation_log_time;
drop index if exists idx_dpo_automation_log_domain;
drop index if exists idx_dpo_automation_log_dp;
-- +goose StatementEnd