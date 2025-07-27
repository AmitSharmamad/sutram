-- Comprehensive Options Chain and Market Data Schema for ClickHouse
-- This schema provides complete options trading data storage with optimized performance

-- Main options chain table for storing complete market data
CREATE TABLE IF NOT EXISTS options_chain (
    -- Instrument identification
    exchange_security_id Int32,
    underlying_symbol LowCardinality(String),
    underlying_price Float64,              -- Current price of underlying asset
    strike_price Float64,
    option_type LowCardinality(String),    -- 'CE' for Call, 'PE' for Put
    expiry_date Date,
    days_to_expiry Int16,
    
    -- Price data
    ltp Float64,                           -- Last Traded Price
    change_amount Float64,                 -- Absolute change from previous close
    change_percent Float64,                -- Percentage change
    
    -- OHLC for the session
    open Float64,
    high Float64,
    low Float64,
    close Float64,
    previous_close Float64,                -- Previous day's closing price
    
    -- Volume and Open Interest
    volume Float64,                        -- Current session volume
    total_traded_value Float64,            -- Total value traded
    open_interest Float64,                 -- Total open interest
    oi_change Float64,                     -- Change in open interest
    oi_change_percent Float64,             -- Percentage change in OI
    
    -- Bid/Ask data (Level 1)
    bid_price Float64,
    bid_quantity Float64,
    bid_orders Nullable(Int32),            -- Number of bid orders
    ask_price Float64,
    ask_quantity Float64,
    ask_orders Nullable(Int32),            -- Number of ask orders
    spread Float64,                        -- ask_price - bid_price
    mid_price Float64,                     -- (bid_price + ask_price) / 2
    
    -- Market depth indicators
    total_buy_quantity Float64,            -- Total quantity at all bid levels
    total_sell_quantity Float64,           -- Total quantity at all ask levels
    
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
    moneyness Float64,                     -- strike_price / underlying_price
    intrinsic_value Float64,               -- Max(0, underlying_price - strike_price) for calls
    time_value Float64,                    -- ltp - intrinsic_value
    
    -- Trading status and flags
    is_active UInt8 DEFAULT 1,             -- 1 if actively trading, 0 if suspended
    circuit_limit_upper Nullable(Float64), -- Upper circuit limit
    circuit_limit_lower Nullable(Float64), -- Lower circuit limit
    
    -- Timestamps
    market_timestamp DateTime64(6, 'Asia/Kolkata'),  -- Market data timestamp
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = ReplacingMergeTree(updated_at)  -- Latest record wins on duplicates
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, strike_price, option_type, market_timestamp)
SETTINGS index_granularity = 8192;

-- Dedicated table for detailed options Greeks and risk metrics
CREATE TABLE IF NOT EXISTS options_greeks (
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
    charm Nullable(Float64),               -- Delta decay (dDelta/dTime)
    vanna Nullable(Float64),               -- Vega/Delta cross sensitivity
    volga Nullable(Float64),               -- Vega convexity
    
    -- Volatility metrics
    implied_volatility Nullable(Float64),
    vega_weighted_iv Nullable(Float64),    -- Vega-weighted implied volatility
    
    -- Time and moneyness metrics
    time_to_expiry Float64,                -- Days to expiry (can be fractional)
    moneyness Float64,                     -- Strike/Spot ratio
    log_moneyness Nullable(Float64),       -- ln(Strike/Spot)
    
    -- Value calculations
    intrinsic_value Float64,
    time_value Float64,
    theoretical_price Nullable(Float64),   -- Black-Scholes theoretical price
    
    -- Risk metrics
    risk_free_rate Nullable(Float64),      -- Risk-free interest rate used
    dividend_yield Nullable(Float64),      -- Dividend yield used
    
    -- Timestamps
    calculation_timestamp DateTime64(6, 'Asia/Kolkata'),
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = MergeTree()
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, strike_price, option_type, calculation_timestamp)
SETTINGS index_granularity = 8192;

-- High-frequency quotes table for bid/ask snapshots
CREATE TABLE IF NOT EXISTS options_quotes (
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
    spread Nullable(Float64),              -- ask_price - bid_price
    spread_percent Nullable(Float64),      -- spread / mid_price * 100
    mid_price Nullable(Float64),           -- (bid_price + ask_price) / 2
    
    -- Market depth (Level 2 data aggregated)
    total_bid_quantity Nullable(Float64), -- Sum of all bid quantities
    total_ask_quantity Nullable(Float64), -- Sum of all ask quantities
    weighted_bid_price Nullable(Float64), -- Quantity-weighted bid price
    weighted_ask_price Nullable(Float64), -- Quantity-weighted ask price
    
    -- Quote quality indicators
    quote_timestamp DateTime64(6, 'Asia/Kolkata'),
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = MergeTree()
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, strike_price, option_type, quote_timestamp)
SETTINGS index_granularity = 8192;

-- Open Interest tracking with historical data
CREATE TABLE IF NOT EXISTS open_interest (
    exchange_security_id Int32,
    underlying_symbol LowCardinality(String),
    strike_price Float64,
    option_type LowCardinality(String),
    expiry_date Date,
    
    -- Open Interest data
    open_interest Float64,
    change_in_oi Float64,                  -- Absolute change from previous session
    change_in_oi_percent Float64,          -- Percentage change from previous session
    
    -- Additional OI metrics
    oi_volume_ratio Nullable(Float64),     -- OI / Volume ratio
    oi_concentration Nullable(Float64),    -- This strike's OI / Total series OI
    
    -- PCR (Put-Call Ratio) data at strike level
    pcr_volume Nullable(Float64),          -- Put volume / Call volume at this strike
    pcr_oi Nullable(Float64),              -- Put OI / Call OI at this strike
    
    -- Timestamps
    trading_date Date,                     -- Trading session date
    data_timestamp DateTime64(6, 'Asia/Kolkata'),
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = MergeTree()
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, strike_price, option_type, trading_date)
SETTINGS index_granularity = 8192;

-- Options chain summary for quick lookups and analytics
CREATE TABLE IF NOT EXISTS options_chain_summary (
    underlying_symbol LowCardinality(String),
    expiry_date Date,
    trading_date Date,
    
    -- Summary statistics
    total_call_volume Float64,
    total_put_volume Float64,
    total_call_oi Float64,
    total_put_oi Float64,
    
    -- PCR ratios
    pcr_volume Float64,                    -- Total Put Volume / Total Call Volume
    pcr_oi Float64,                        -- Total Put OI / Total Call OI
    
    -- ATM data
    atm_strike Float64,
    atm_call_iv Nullable(Float64),
    atm_put_iv Nullable(Float64),
    
    -- IV statistics
    avg_call_iv Nullable(Float64),
    avg_put_iv Nullable(Float64),
    iv_skew Nullable(Float64),             -- Put IV - Call IV
    
    -- Max Pain and key levels
    max_pain_strike Nullable(Float64),
    highest_oi_call_strike Nullable(Float64),
    highest_oi_put_strike Nullable(Float64),
    
    -- Market indicators
    vix_equivalent Nullable(Float64),      -- VIX-like volatility index
    fear_greed_index Nullable(Float64),    -- Market sentiment indicator
    
    -- Timestamps
    created_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata'),
    updated_at DateTime64(6, 'Asia/Kolkata') DEFAULT now64(6, 'Asia/Kolkata')
    
) ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY (underlying_symbol, toYYYYMM(expiry_date))
ORDER BY (underlying_symbol, expiry_date, trading_date)
SETTINGS index_granularity = 8192;

-- Materialized views for common queries

-- Real-time ATM options view
CREATE MATERIALIZED VIEW IF NOT EXISTS mv_atm_options TO options_chain_summary AS
SELECT 
    underlying_symbol,
    expiry_date,
    toDate(market_timestamp) as trading_date,
    argMin(strike_price, abs(strike_price - underlying_price)) as atm_strike,
    avgIf(implied_volatility, option_type = 'CE' AND abs(strike_price - underlying_price) <= 50) as atm_call_iv,
    avgIf(implied_volatility, option_type = 'PE' AND abs(strike_price - underlying_price) <= 50) as atm_put_iv,
    sumIf(volume, option_type = 'CE') as total_call_volume,
    sumIf(volume, option_type = 'PE') as total_put_volume,
    sumIf(open_interest, option_type = 'CE') as total_call_oi,
    sumIf(open_interest, option_type = 'PE') as total_put_oi,
    sumIf(volume, option_type = 'PE') / nullIf(sumIf(volume, option_type = 'CE'), 0) as pcr_volume,
    sumIf(open_interest, option_type = 'PE') / nullIf(sumIf(open_interest, option_type = 'CE'), 0) as pcr_oi,
    avgIf(implied_volatility, option_type = 'CE') as avg_call_iv,
    avgIf(implied_volatility, option_type = 'PE') as avg_put_iv,
    avgIf(implied_volatility, option_type = 'PE') - avgIf(implied_volatility, option_type = 'CE') as iv_skew,
    0 as max_pain_strike,  -- Will be calculated separately
    0 as highest_oi_call_strike,  -- Will be calculated separately
    0 as highest_oi_put_strike,   -- Will be calculated separately
    0 as vix_equivalent,  -- Will be calculated separately
    0 as fear_greed_index,  -- Will be calculated separately
    now64(6, 'Asia/Kolkata') as created_at,
    now64(6, 'Asia/Kolkata') as updated_at
FROM options_chain
WHERE market_timestamp >= now() - INTERVAL 1 HOUR
GROUP BY underlying_symbol, expiry_date, trading_date;

-- Indexes for performance optimization
-- These are automatically created based on ORDER BY clauses, but can be explicitly defined for complex queries

-- Create dictionaries for fast lookups (optional, for frequently accessed reference data)
-- CREATE DICTIONARY IF NOT EXISTS option_instruments_dict (
--     exchange_security_id Int32,
--     underlying_symbol String,
--     strike_price Float64,
--     option_type String,
--     expiry_date Date,
--     lot_size Int32
-- ) PRIMARY KEY exchange_security_id
-- SOURCE(CLICKHOUSE(HOST 'localhost' PORT 9000 USER 'default' TABLE 'instrument_master' DB 'market_data'))
-- LAYOUT(HASHED())
-- LIFETIME(MIN 3600 MAX 7200);
