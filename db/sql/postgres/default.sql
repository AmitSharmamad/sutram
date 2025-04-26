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

INSERT INTO strategy (
  id,
  user_id,
  name,
  description,
  position_type,
  instruments,
  chart_type,
  chart_timeframe,
  trade_type,
  qty,
  entry_expression,
  exit_expression,
  target_profit,
  stop_loss,
  trailing_stop_loss,
  exit_type,
  capital,
  max_trades,
  entry_start_time,
  entry_end_time,
  start_date,
  end_date,
  is_active
) 

VALUES (
  'c524747e-cb18-4f2e-86df-77b19dc6d23a',
  'e65c67f1-2b43-4581-b33b-ec26295e741b',
  '50 SMA Crosses Over 200 SMA',
  '50 SMA Crosses Over 200 SMA',
  'LONG',
  ARRAY[1],
  'CANDLESTICK',  -- Example chart type
  '3m',  -- Daily timeframe
  'MIS',
  10,  -- Quantity per trade
  'SMA(CLOSE, 14, 0) CrossesAbove SMA(CLOSE, 14, 0)',  -- Empty entry expression
  'SMA(CLOSE, 14, 0) CrossesAbove SMA(CLOSE, 14, 0)',  -- Empty exit expression
  0.05,  -- 5% target profit
  0.02,  -- 2% stop loss
  0.03,  -- 3% trailing stop loss
  'PERCENTAGE',  -- Exit type
  100000.00,  -- Capital allocation
  5,  -- Maximum trades
  CURRENT_TIMESTAMP,  -- Entry start time
  CURRENT_TIMESTAMP + INTERVAL '1 day',  -- Entry end time
  CURRENT_DATE,  -- Start date
  CURRENT_DATE + INTERVAL '30 days',  -- End date
  true  -- Is active
);
