USE sutram;

CREATE TABLE IF NOT EXISTS instrument_master (
    id UInt64,
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

 INSERT INTO instrument_master (
    id, exchange_security_id, exchange, segment, symbol, instrument, start_date, expiry, 
    exm_exch_id, smst_security_id, instrument_name, trading_symbol, lot_units, strike_price, 
    option_type, tick_size, expiry_flag, exch_instrument_type, symbol_name
) VALUES
    (1, 12018, 'NSE', 'EQUITY', 'SUZLON', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 12018, 'Suzlon Energy Ltd', 'SUZLON', 1, 0, '', 0.05, 'N', 'EQ', 'SUZLON'),
    (2, 14366, 'NSE', 'EQUITY', 'IDEA', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 14366, 'Vodafone Idea Ltd', 'IDEA', 1, 0, '', 0.05, 'N', 'EQ', 'IDEA'),
    (3, 11915, 'NSE', 'EQUITY', 'YESBANK', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 11915, 'Yes Bank Ltd', 'YESBANK', 1, 0, '', 0.05, 'N', 'EQ', 'YESBANK'),
    (4, 2029, 'NSE', 'EQUITY', 'IRFC', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 2029, 'Indian Railway Finance Corporation Ltd', 'IRFC', 1, 0, '', 0.05, 'N', 'EQ', 'IRFC'),
    (5, 1491, 'NSE', 'EQUITY', 'IFCI', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 1491, 'IFCI Ltd', 'IFCI', 1, 0, '', 0.05, 'N', 'EQ', 'IFCI'),
    (6, 383, 'NSE', 'EQUITY', 'BEL', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 383, 'Bharat Electronics Ltd', 'BEL', 1, 0, '', 0.05, 'N', 'EQ', 'BEL'),
    (7, 15259, 'NSE', 'EQUITY', 'RPOWER', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 15259, 'Reliance Power Ltd', 'RPOWER', 1, 0, '', 0.05, 'N', 'EQ', 'RPOWER'),
    (8, 9552, 'NSE', 'EQUITY', 'RVNL', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 9552, 'Rail Vikas Nigam Ltd', 'RVNL', 1, 0, '', 0.05, 'N', 'EQ', 'RVNL'),
    (9, 17400, 'NSE', 'EQUITY', 'NHPC', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 17400, 'NHPC Ltd', 'NHPC', 1, 0, '', 0.05, 'N', 'EQ', 'NHPC'),
    (10, 10604, 'NSE', 'EQUITY', 'BHARTIARTL', 'STOCK', '2025-01-01 09:00:00.000000', '2025-12-31 17:00:00.000000', 1, 10604, 'Bharti Airtel Ltd', 'BHARTIARTL', 1, 0, '', 0.05, 'N', 'EQ', 'BHARTIARTL');
