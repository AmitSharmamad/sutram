USE sutram;

CREATE TABLE IF NOT EXISTS instrument_master (
    exchange_security_id Int32,
    exchange LowCardinality(String),
    segment LowCardinality(String),
    symbol String,
    instrument String,
    start_date DateTime64(6),
    expiry DateTime64(6),
    exm_exch_id Int32,
    smst_security_id Int32,
    instrument_name String,
    trading_symbol String,
    lot_units Int32,
    strike_price Float64,
    option_type String,
    tick_size Float64,
    expiry_flag String,
    exch_instrument_type String,
    symbol_name String
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
