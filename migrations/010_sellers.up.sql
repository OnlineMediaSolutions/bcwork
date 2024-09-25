CREATE TABLE competitors (
                             name VARCHAR(100) NOT NULL,
                             url TEXT NOT NULL,
                             PRIMARY KEY (name),
                             CONSTRAINT unique_url UNIQUE (url)
);

CREATE TABLE sellers_json_history (
                                      competitor_name VARCHAR(100) NOT NULL,
                                      added_domains TEXT NOT NULL,
                                      added_publishers TEXT NOT NULL,
                                      backup_today JSONB NOT NULL,
                                      backup_yesterday JSONB NOT NULL,
                                      backup_before_yesterday JSONB NOT NULL,
                                      created_at TIMESTAMP,
                                      updated_at TIMESTAMP,
                                      PRIMARY KEY (competitor_name),
                                      FOREIGN KEY (competitor_name) REFERENCES competitors (name)
);
