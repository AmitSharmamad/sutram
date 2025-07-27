-- Create options quotes table

CREATE TABLE IF NOT EXISTS sutram.options_quotes (
    exchange_security_id Int32,
    underlying_symbol LowCardinality(String),
    strike_price Float64,
    option_type LowCardinality(String),
    expiry_date Date,
    
    -- Level 1 data
    bid_price Nullable(Float64),
    ask_price Nullable(Float64),
    bid_quantity Nullable(Float64),
    ask_quantity Nullable(Float64),
    bid_orders Nullable(Int32),
    ask_orders Nullable(Int32),
    
    -- Derived metrics
    spread Nullable(Float64),
    spread_percent Nullable(Float64),
    mid_price Nullable(Float64),
    
    -- Market depth (Level 2 data aggregated)
    total_bid_quantity Nullable(Float64),
    total_ask_quantity Nullable(Float64),
    weighted_bid_price Nullable(Float64),
    weighted_ask_price Nullable(Float64),
    
    -- Quote quality indicators
    quote_timestamp DateTime64(6, 'Asia/Kolkata'),
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = MergeTree()
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, strike_price, option_type, quote_timestamp)
SETTINGS index_granularity = 8192;
