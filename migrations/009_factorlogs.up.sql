CREATE TABLE price_factor_log (
                                time TIMESTAMP NOT NULL,
                                eval_time TIMESTAMP NOT NULL,
                                pubimps INTEGER NOT NULL,
                                soldimps INTEGER NOT NULL,
                                cost FLOAT NOT NULL,
                                revenue FLOAT NOT NULL,
                                gp FLOAT NOT NULL,
                                gpp FLOAT NOT NULL,
                                publisher VARCHAR(10) NOT NULL,
                                domain VARCHAR(255) NOT NULL,
                                country CHAR(2) NOT NULL,
                                device VARCHAR(50) NOT NULL,
                                old_factor FLOAT NOT NULL,
                                new_factor FLOAT NOT NULL,
                                response_status INTEGER NOT NULL,
                                increase FLOAT NOT NULL,
                                PRIMARY KEY (publisher, domain, country, device, time)
);