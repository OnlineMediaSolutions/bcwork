-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS price_factor_log
(
    time            TIMESTAMP    NOT NULL,
    eval_time       TIMESTAMP    NOT NULL,
    pubimps         INTEGER      NOT NULL,
    soldimps        INTEGER      NOT NULL,
    cost            FLOAT        NOT NULL,
    revenue         FLOAT        NOT NULL,
    gp              FLOAT        NOT NULL,
    gpp             FLOAT        NOT NULL,
    publisher       VARCHAR(10)  NOT NULL,
    domain          VARCHAR(255) NOT NULL,
    country         CHAR(2)      NOT NULL,
    device          VARCHAR(50)  NOT NULL,
    old_factor      FLOAT        NOT NULL,
    new_factor      FLOAT        NOT NULL,
    response_status INTEGER      NOT NULL,
    increase        FLOAT        NOT NULL,
    source          varchar(20)  NOT NULL,
    PRIMARY KEY (publisher, domain, country, device, time)
);

CREATE TABLE IF NOT EXISTS configuration
(
    key         VARCHAR(36) UNIQUE not null,
    value       TEXT               not null,
    description TEXT,
    updated_at  TIMESTAMP,
    created_at  TIMESTAMP,
    primary key (key)
);

CREATE TABLE IF NOT EXISTS global_factor
(
    key          VARCHAR(36) not null,
    publisher_id VARCHAR(36),
    value  FLOAT,
    created_by_id VARCHAR(36),
    updated_at   TIMESTAMP,
    created_at   TIMESTAMP,
    primary key (key, publisher_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists price_factor_log;
drop table if exists configuration;
drop table if exists global_factor;
-- +goose StatementEnd