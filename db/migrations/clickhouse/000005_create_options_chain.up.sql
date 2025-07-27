-- Create options chain main table

CREATE TABLE IF NOT EXISTS sutram.options_chain (
    -- Instrument identification
    exchange_security_id Int32,
    underlying_symbol LowCardinality(String),
    underlying_price Float64,
    strike_price Float64,
    option_type LowCardinality(String),
    expiry_date Date,
    days_to_expiry Int16,
    
    -- Price data
    ltp Float64,
    change_amount Float64,
    change_percent Float64,
    
    -- OHLC for the session
    open Float64,
    high Float64,
    low Float64,
    close Float64,
    previous_close Float64,
    
    -- Volume and Open Interest
    volume Float64,
    total_traded_value Float64,
    open_interest Float64,
    oi_change Float64,
    oi_change_percent Float64,
    
    -- Bid/Ask data (Level 1)
    bid_price Float64,
    bid_quantity Float64,
    bid_orders Nullable(Int32),
    ask_price Float64,
    ask_quantity Float64,
    ask_orders Nullable(Int32),
    spread Float64,
    mid_price Float64,
    
    -- Market depth indicators
    total_buy_quantity Float64,
    total_sell_quantity Float64,
    
    -- Greeks
    delta Nullable(Float64),
    gamma Nullable(Float64),
    theta Nullable(Float64),
    vega Nullable(Float64),
    rho Nullable(Float64),
    
    -- Volatility data
    implied_volatility Nullable(Float64),
    historical_volatility Nullable(Float64),
    
    -- Calculated metrics
    moneyness Float64,
    intrinsic_value Float64,
    time_value Float64,
    
    -- Trading status and flags
    is_active UInt8 DEFAULT 1,
    circuit_limit_upper Nullable(Float64),
    circuit_limit_lower Nullable(Float64),
    
    -- Timestamps
    market_timestamp DateTime64(6, 'Asia/Kolkata'),
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, strike_price, option_type, market_timestamp)
SETTINGS index_granularity = 8192;
