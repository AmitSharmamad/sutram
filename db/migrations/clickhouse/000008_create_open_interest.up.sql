-- Create open interest table

CREATE TABLE IF NOT EXISTS sutram.open_interest (
    exchange_security_id Int32,
    underlying_symbol LowCardinality(String),
    strike_price Float64,
    option_type LowCardinality(String),
    expiry_date Date,
    
    -- Open Interest data
    open_interest Float64,
    change_in_oi Float64,
    change_in_oi_percent Float64,
    
    -- Additional OI metrics
    oi_volume_ratio Nullable(Float64),
    oi_concentration Nullable(Float64),
    
    -- PCR (Put-Call Ratio) data at strike level
    pcr_volume Nullable(Float64),
    pcr_oi Nullable(Float64),
    
    -- Timestamps
    trading_date Date,
    data_timestamp DateTime64(6, 'Asia/Kolkata'),
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = MergeTree()
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, strike_price, option_type, trading_date)
SETTINGS index_granularity = 8192;
