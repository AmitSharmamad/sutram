USE sutram;

CREATE TABLE IF NOT EXISTS instrument_master (
    id UInt64,
    exchange_security_id Int32,
    exchange LowCardinality(String),
    segment LowCardinality(String),
    symbol String,
    instrument String,
    start_date DateTime64(6),
    expiry DateTime64(6)
) ENGINE = MergeTree()
ORDER BY (
    exchange,
    segment,
);

CREATE TABLE IF NOT EXISTS instrument (
    exchange_security_id Int32,
    open Float32,
    high Float32,
    low Float32,
    close Float32,
    volume Float32,
    `timestamp` DateTime64(6, 'Asia/Kolkata'),
    `created_at` DateTime64(6, 'Asia/Kolkata'),
    `updated_at` DateTime64(6, 'Asia/Kolkata')
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(`timestamp`) -- Partition by year and month
ORDER BY (
    exchange_security_id,
    `timestamp`
);

CREATE INDEX IF NOT EXISTS idx_instrument_exchange_security_id
ON instrument (
    exchange_security_id
) TYPE bloom_filter GRANULARITY 4;

INSERT INTO instrument_master (
    id,
    exchange_security_id,
    exchange,
    segment,
    symbol,
    instrument,
    start_date,
    expiry
) VALUES (
    1,
    12345,
    'NSE',
    'F&O',
    'Tata Motors',
    'Tata Motors',
    '2023-01-01 00:00:00.000000',
    '2024-01-01 00:00:00.000000'
);