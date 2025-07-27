-- Create options greeks table

CREATE TABLE IF NOT EXISTS sutram.options_greeks (
    exchange_security_id Int32,
    underlying_symbol LowCardinality(String),
    strike_price Float64,
    option_type LowCardinality(String),
    expiry_date Date,
    
    -- Greeks
    delta Nullable(Float64),
    gamma Nullable(Float64),
    theta Nullable(Float64),
    vega Nullable(Float64),
    rho Nullable(Float64),
    
    -- Advanced Greeks
    charm Nullable(Float64),
    vanna Nullable(Float64),
    volga Nullable(Float64),
    
    -- Volatility metrics
    implied_volatility Nullable(Float64),
    vega_weighted_iv Nullable(Float64),
    
    -- Time and moneyness metrics
    time_to_expiry Float64,
    moneyness Float64,
    log_moneyness Nullable(Float64),
    
    -- Value calculations
    intrinsic_value Float64,
    time_value Float64,
    theoretical_price Nullable(Float64),
    
    -- Risk metrics
    risk_free_rate Nullable(Float64),
    dividend_yield Nullable(Float64),
    
    -- Timestamps
    calculation_timestamp DateTime64(6, 'Asia/Kolkata'),
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = MergeTree()
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, strike_price, option_type, calculation_timestamp)
SETTINGS index_granularity = 8192;
