-- Insert default system user and trading strategies

-- Insert system user
INSERT INTO "user" (
  id, name, username, email, phone_number, password, is_active
) VALUES (
  'e65c67f1-2b43-4581-b33b-ec26295e741b',
  'system',
  'system',
  'system@sutram.com',
  '-',
  'system',
  true
);

-- Top Trading Strategies with Entry Expressions, Target Profits, and Stop Losses

INSERT INTO strategy (
    user_id, name, description, position_type, instruments, chart_type, chart_timeframe, 
    trading_days, trade_type, qty, entry_expression, exit_expression, target_profit, target_profit_type, 
    stop_loss, stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, capital, max_allocation, 
    max_trades, segment, strategy_type
) VALUES (
    'e65c67f1-2b43-4581-b33b-ec26295e741b', -- system user ID
    'Moving Average Crossover',
    'Classic strategy that uses the crossing of a fast moving average over a slow moving average as a buy signal, and the opposite for a sell signal.',
    'LONG',
    '{}', -- empty array for instruments
    'CANDLESTICK',
    '15m',
    '{MON,TUE,WED,THU,FRI}',
    'MIS',
    100,
    'CLOSE(0) > EMA(CLOSE, 9, 0) AND EMA(CLOSE, 9, 0) CrossesAbove EMA(CLOSE, 21, 0)',
    'EMA(CLOSE, 9, 0) CrossesBelow EMA(CLOSE, 21, 0)',
    3.5,
    'PERCENTAGE',
    2.0,
    'PERCENTAGE',
    1.5,
    'PERCENTAGE',
    100000,
    10000,
    5,
    'EQ',
    'CONDITIONAL'
);

-- 2. RSI Oversold/Overbought Strategy
INSERT INTO strategy (
    user_id, name, description, position_type, instruments, chart_type, chart_timeframe, 
    trading_days, trade_type, qty, entry_expression, exit_expression, target_profit, target_profit_type, 
    stop_loss, stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, capital, max_allocation, 
    max_trades, segment, strategy_type
) VALUES (
    'e65c67f1-2b43-4581-b33b-ec26295e741b', -- system user ID
    'RSI Reversal Strategy',
    'This strategy uses RSI to identify oversold conditions to enter long positions and overbought conditions to exit.',
    'LONG',
    '{}', -- empty array for instruments
    'CANDLESTICK',
    '5m',
    '{MON,TUE,WED,THU,FRI}',
    'MIS',
    100,
    'RSI(CLOSE, 14, 0) CrossesAbove NUMBER(30)',
    'RSI(CLOSE, 14, 0) CrossesAbove NUMBER(70)',
    4.0,
    'PERCENTAGE',
    2.5,
    'PERCENTAGE',
    1.8,
    'PERCENTAGE',
    100000,
    10000,
    5,
    'EQ',
    'CONDITIONAL'
);

-- 4. MACD Bullish Crossover
INSERT INTO strategy (
    user_id, name, description, position_type, instruments, chart_type, chart_timeframe, 
    trading_days, trade_type, qty, entry_expression, exit_expression, target_profit, target_profit_type, 
    stop_loss, stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, capital, max_allocation, 
    max_trades, segment, strategy_type
) VALUES (
    'e65c67f1-2b43-4581-b33b-ec26295e741b', -- system user ID
    'MACD Bullish Crossover',
    'Enters long when MACD line crosses above the signal line with positive momentum.',
    'LONG',
    '{}', -- empty array for instruments
    'CANDLESTICK',
    '30m',
    '{MON,TUE,WED,THU,FRI}',
    'MIS',
    100,
    'MACD(12, 26, 9, STANDARD, 0) CrossesAbove MACD(12, 26, 9, SIGNAL, 0) AND MACD(12, 26, 9, STANDARD, 0) < NUMBER(0)',
    'MACD(12, 26, 9, STANDARD, 0) CrossesBelow MACD(12, 26, 9, SIGNAL, 0)',
    4.2,
    'PERCENTAGE',
    2.8,
    'PERCENTAGE',
    2.0,
    'PERCENTAGE',
    100000,
    10000,
    3,
    'EQ',
    'CONDITIONAL'
);

-- 6. Triple EMA Crossover Strategy
INSERT INTO strategy (
    user_id, name, description, position_type, instruments, chart_type, chart_timeframe, 
    trading_days, trade_type, qty, entry_expression, exit_expression, target_profit, target_profit_type, 
    stop_loss, stop_loss_type, trailing_stop_loss, trailing_stop_loss_type, capital, max_allocation, 
    max_trades, segment, strategy_type
) VALUES (
    'e65c67f1-2b43-4581-b33b-ec26295e741b', -- system user ID
    'Triple EMA Strategy',
    'Uses three exponential moving averages of different periods to confirm trend direction.',
    'LONG',
    '{}', -- empty array for instruments
    'CANDLESTICK',
    '15m',
    '{MON,TUE,WED,THU,FRI}',
    'MIS',
    100,
    'EMA(CLOSE, 5, 0) > EMA(CLOSE, 10, 0) AND EMA(CLOSE, 10, 0) > EMA(CLOSE, 20, 0) AND CLOSE(0) > EMA(CLOSE, 5, 0)',
    'EMA(CLOSE, 5, 0) < EMA(CLOSE, 10, 0)',
    3.5,
    'PERCENTAGE',
    2.2,
    'PERCENTAGE',
    1.5,
    'PERCENTAGE',
    100000,
    10000,
    4,
    'EQ',
    'CONDITIONAL'
);