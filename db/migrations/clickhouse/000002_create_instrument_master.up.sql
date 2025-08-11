-- Create instrument master table with essential columns for algorithmic trading

CREATE TABLE IF NOT EXISTS sutram.instrument_master (
    -- Core Identification Columns
    exchange_security_id Int32,
    exchange Enum8('NSE' = 1, 'BSE' = 2, 'MCX' = 3),
    segment Enum8('EQUITY' = 1, 'INDEX' = 2, 'COMMODITY' = 3, 'CURRENCY' = 4, 'DERIVATIVE' = 5),
    symbol String,
    instrument Enum8('EQUITY_CASH' = 1, 'STOCK_FUTURES' = 2, 'STOCK_OPTIONS' = 3, 'INDEX' = 4, 'INDEX_FUTURES' = 5, 'INDEX_OPTIONS' = 6, 'COMMODITY_FUTURES' = 7, 'COMMODITY_OPTIONS' = 8, 'CURRENCY_FUTURES' = 9, 'CURRENCY_OPTIONS' = 10),

    -- Trading Symbols and Names
    instrument_name String, -- Maps to DISPLAY_NAME from Dhan

    -- ISIN and Security Identifiers
    isin String,
    underlying_security_id Nullable(Int32),
    underlying_symbol Nullable(String),

    -- Date and Expiry Information
    expiry DateTime64(6),
    expiry_flag Nullable(String), -- M = Monthly, W = Weekly

    -- Trading Parameters
    lot_units Int32,
    tick_size Float64,
    strike_price Nullable(Float64),
    option_type Nullable(String), -- CE = Call, PE = Put
    instrument_type String,

    -- Order Type Availability Flags
    bracket_flag Enum8('N' = 0, 'Y' = 1) DEFAULT 'N', -- Bracket order availability
    cover_flag Enum8('N' = 0, 'Y' = 1) DEFAULT 'N',   -- Cover order availability

    -- Surveillance and Trading Permissions
    asm_gsm_flag Enum8('N' = 0, 'R' = 1, 'Y' = 2) DEFAULT 'N', -- N=Not in ASM/GSM, R=Removed, Y=ASM/GSM
    buy_sell_indicator LowCardinality(String) DEFAULT 'A', -- A = Both Buy/Sell allowed

    -- Cover Order Margin Parameters
    buy_co_min_margin_per Nullable(Float64),
    sell_co_min_margin_per Nullable(Float64),
    buy_co_sl_range_max_perc Nullable(Float64),
    sell_co_sl_range_min_perc Nullable(Float64),

    -- Bracket Order Parameters
    buy_bo_min_margin_per Nullable(Float64),
    buy_bo_sl_range_max_perc Nullable(Float64),
    sell_bo_sl_min_range Nullable(Float64),
    buy_bo_profit_range_max_perc Nullable(Float64),
    sell_bo_profit_range_max_perc Nullable(Float64),

    -- Margin Trading Facility
    mtf_leverage Nullable(Float64),

    -- Metadata
    created_at DateTime64(6) DEFAULT now64(),
    updated_at DateTime64(6) DEFAULT now64()

) ENGINE = MergeTree()
ORDER BY (exchange, segment, symbol)
SETTINGS index_granularity = 8192;
