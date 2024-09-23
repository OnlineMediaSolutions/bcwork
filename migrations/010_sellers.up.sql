create table competitors
(
    name varchar(100) not null,
    url TEXT,
    primary key (name)
);


ALTER TABLE competitors ADD CONSTRAINT unique_url UNIQUE (url);

CREATE TABLE sellers_json_history (
  competitor_name VARCHAR(100) references competitors(name),
  added_domains TEXT not null,
  added_publishers TEXT not null,
  backup_today JSONB,
  backup_yesterday JSONB,
  created_at timestamp,
  updated_at timestamp,
  primary key (competitor_name)

);
