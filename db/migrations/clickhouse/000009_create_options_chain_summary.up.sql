-- Create options chain summary table

CREATE TABLE IF NOT EXISTS sutram.options_chain_summary (
    underlying_symbol LowCardinality(String),
    expiry_date Date,
    trading_date Date,
    
    -- Summary statistics
    total_call_volume Float64,
    total_put_volume Float64,
    total_call_oi Float64,
    total_put_oi Float64,
    
    -- PCR ratios
    pcr_volume Float64,
    pcr_oi Float64,
    
    -- ATM data
    atm_strike Float64,
    atm_call_iv Nullable(Float64),
    atm_put_iv Nullable(Float64),
    
    -- IV statistics
    avg_call_iv Nullable(Float64),
    avg_put_iv Nullable(Float64),
    iv_skew Nullable(Float64),
    
    -- Max Pain and key levels
    max_pain_strike Nullable(Float64),
    highest_oi_call_strike Nullable(Float64),
    highest_oi_put_strike Nullable(Float64),
    
    -- Market indicators
    vix_equivalent Nullable(Float64),
    fear_greed_index Nullable(Float64),
    
    -- Timestamps
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, trading_date)
SETTINGS index_granularity = 8192;
