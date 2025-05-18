USE sutram;

CREATE TABLE IF NOT EXISTS instrument_master (
    id UInt64,
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
    id UInt64,
    instrument String,
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
    id,
    `timestamp`
);

CREATE INDEX IF NOT EXISTS idx_instrument_id
ON instrument (
    id
) TYPE bloom_filter GRANULARITY 4;

INSERT INTO instrument_master (
    id,
    exchange,
    segment,
    symbol,
    instrument,
    start_date,
    expiry
) VALUES (
    1,
    'NSE',
    'F&O',
    'Tata Motors',
    'Tata Motors',
    '2023-01-01 00:00:00.000000',
    '2024-01-01 00:00:00.000000'
);