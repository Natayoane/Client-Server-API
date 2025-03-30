CREATE TABLE IF NOT EXISTS exchange_rate (
    id          CHAR(36) PRIMARY KEY,
    code        VARCHAR(10) NOT NULL,
    code_in     VARCHAR(10) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    high        VARCHAR(20) NOT NULL,
    low         VARCHAR(20) NOT NULL,
    var_bid     VARCHAR(20) NOT NULL,
    pct_change  VARCHAR(20) NOT NULL,
    bid         VARCHAR(20) NOT NULL,
    ask         VARCHAR(20) NOT NULL,
    timestamp   VARCHAR(20) NOT NULL,
    create_date VARCHAR(20) NOT NULL
);