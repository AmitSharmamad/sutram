-- Create instrument master table

CREATE TABLE IF NOT EXISTS sutram.instrument_master (
    exchange_security_id Int32,
    exchange LowCardinality(String),
    segment Enum8('Equities' = 1, 'Indices' = 2, 'Commodities' = 3, 'Currencies' = 4, 'Derivatives' = 5),
    symbol String,
    instrument Enum8('Stocks (Cash Market)' = 1, 'Stock Futures' = 2, 'Stock Options' = 3, 'Index Spot' = 4, 'Index Futures' = 5, 'Index Options' = 6, 'Commodity Futures' = 7, 'Commodity Options' = 8),
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
ORDER BY (exchange, segment);
