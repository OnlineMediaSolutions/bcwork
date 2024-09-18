create table competitors
(
    competitor_name varchar(100) not null,
    url TEXT,
    primary key (competitor_name)
);


ALTER TABLE competitors ADD CONSTRAINT unique_url UNIQUE (url);

CREATE TABLE sellers_json_history (
  competitor_name VARCHAR(100) references competitors(competitor_name),
  url TEXT references competitors(url),
  added_domains TEXT,
  added_publishers TEXT,
  backup_today JSONB,
  backup_yesterday JSONB,
  created_at timestamp,
  updated_at timestamp,
  primary key (competitor_name,url)

);
