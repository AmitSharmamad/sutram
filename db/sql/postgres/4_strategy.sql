-- Strategy Management Tables
CREATE TABLE IF NOT EXISTS strategy (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    description VARCHAR,
    position_type direction,
    instruments INT[],
    chart_type chart_type,
    chart_timeframe chart_timeframe,
    trading_days trading_day[],
    trade_type trade_type,
    qty INT,
    entry_expression TEXT,
    exit_expression TEXT,
    target_profit DECIMAL(18, 8) NOT NULL,
    target_profit_type value_type NOT NULL,
    stop_loss DECIMAL(18, 8) NOT NULL,
    stop_loss_type value_type NOT NULL,
    trailing_stop_loss DECIMAL(18, 8),
    trailing_stop_loss_type value_type,
    exit_type VARCHAR,
    capital DECIMAL(18, 8),
    max_allocation DECIMAL(18, 8),
    max_trades INT,
    brokerage INT,
    slippage DECIMAL(5, 2),
    entry_start_time TIME,
    entry_end_time TIME,
    start_date DATE,
    end_date DATE,
    segment segment NOT NULL,
    strategy_type strategy_type NOT NULL,
    entry_time TIME,
    exit_time TIME,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT true,
    PRIMARY KEY (id)
);


-- constraints for strategy table
-- when strategy_type is 'CONDITIONAL', entry_expression must be provided
-- ALTER TABLE strategy ADD CONSTRAINT entry_expression_required CHECK (
--     (strategy_type != 'CONDITIONAL') OR (entry_expression IS NOT NULL AND entry_expression != '')
-- );
-- when strategy_type is 'TIME_BASED', entry_expression must not be provided and entry_time and exit_time must be provided
-- ALTER TABLE strategy ADD CONSTRAINT time_based_strategy CHECK (
--     (strategy_type != 'TIME_BASED') OR (entry_time IS NOT NULL AND exit_time IS NOT NULL)
-- );

--  when segment is 'EQ', position_type can not be null and trade_type can not be null and qty cannot be null and max_allocation cannot be null
-- ALTER TABLE strategy ADD CONSTRAINT segment_eq_constraints CHECK (
--     (segment != 'EQ') OR (
--         position_type IS NOT NULL AND
--         trade_type IS NOT NULL AND
--         qty IS NOT NULL AND
--         max_allocation IS NOT NULL
--     )
-- );
-- ALTER TABLE strategy ADD CONSTRAINT position_type_required CHECK (
--     (segment != 'EQ') OR (position_type IS NOT NULL)
-- );
