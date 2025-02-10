CREATE TABLE IF NOT EXISTS instrument_master (
    id UInt64,
    exchange LowCardinality(String),
    segment LowCardinality(String),
    symbol String,
    instrument LowCardinality(String),
    instrument_id UInt64,
    instrument_name String,
    start_date DateTime64(6),
    expiry DateTime64(6)
) ENGINE = MergeTree()
ORDER BY (
    exchange,
    segment,
    instrument_id
);

CREATE TABLE IF NOT EXISTS instrument (
    id UInt64,
    open Float32,
    high Float32,
    low Float32,
    close Float32,
    volume Float32,
    `timestamp` DateTime64(6),
    created_at DateTime64(6),
    updated_at DateTime64(6)
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
    instrument_id,
    instrument_name,
    start_date,
    expiry
)
VALUES (
    1,
    'MCX',
    'F&O',
    '<VendorSymbol>',
    'NATURALGAS_FUT_JAN_2025',
    1,
    'Natural Gas Futures Jan',
    '09-09-2016',
    '09-09-2016'
);

INSERT INTO instrument (
    id,
    open,
    high,
    low,
    close,
    volume,
    `timestamp`,
    created_at,
    updated_at
)
VALUES (
    1,
    918,
    927,
    900,
    909,
    9000000,
    '09-09-2016',
    '09-09-2016',
    '09-09-2016'
);
