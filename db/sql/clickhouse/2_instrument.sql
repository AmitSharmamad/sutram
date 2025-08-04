USE sutram;

CREATE TABLE IF NOT EXISTS sutram.instrument_master (
    exchange_security_id Int32,
    exchange LowCardinality(String),
    segment Enum8('Equities' = 1, 'Derivatives' = 2, 'Commodities' = 3, 'Currencies' = 4, 'Indices' = 5),
    symbol String,
    instrument Enum8('Stocks (Cash Market)' = 1, 'Stock Futures' = 2, 'Stock Options' = 3, 'Index Spot' = 4, 'Index Futures' = 5, 'Index Options' = 6, 'Commodity Futures' = 7, 'Commodity Options' = 8),
    start_date DateTime64(6),
    expiry Nullable(DateTime64(6)),
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
ORDER BY (exchange, segment);

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
