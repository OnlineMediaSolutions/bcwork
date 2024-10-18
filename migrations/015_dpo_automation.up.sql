CREATE TABLE dpo_automation_log (
                             id SERIAL PRIMARY KEY,
                             time TIMESTAMP NOT NULL,
                             eval_time TIMESTAMP NOT NULL,
                             domain VARCHAR(255) NOT NULL,
                             publisher VARCHAR(255) NOT NULL,
                             os VARCHAR(50) NOT NULL,
                             country VARCHAR(2) NOT NULL,
                             dp VARCHAR(255) NOT NULL,
                             bid_request INTEGER NOT NULL,
                             revenue FLOAT NOT NULL,
                             erpm FLOAT NOT NULL,
                             old_factor FLOAT NOT NULL,
                             new_factor FLOAT NOT NULL,
                             resp_status INTEGER NOT NULL
);

-- Optional: Create indexes for frequently queried columns
CREATE INDEX idx_dpo_automation_log_time ON dpo_automation_log(time);
CREATE INDEX idx_dpo_automation_log_domain ON dpo_automation_log(domain);
CREATE INDEX idx_dpo_automation_log_dp ON dpo_automation_log(dp);