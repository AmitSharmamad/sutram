-- Create materialized view for ATM options

CREATE MATERIALIZED VIEW IF NOT EXISTS sutram.mv_atm_options TO sutram.options_chain_summary AS
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
    0 as max_pain_strike,
    0 as highest_oi_call_strike,
    0 as highest_oi_put_strike,
    0 as vix_equivalent,
    0 as fear_greed_index,
    now64(6, 'Asia/Kolkata') as created_at,
    now64(6, 'Asia/Kolkata') as updated_at
FROM sutram.options_chain
WHERE market_timestamp >= now() - INTERVAL 1 HOUR
GROUP BY underlying_symbol, expiry_date, trading_date;
