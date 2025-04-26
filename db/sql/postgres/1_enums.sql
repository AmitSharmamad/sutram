
-- Create an ENUM type
CREATE TYPE exchange AS ENUM ('NSE', 'BSE', 'MCX');
CREATE TYPE segment AS ENUM ('F&O', 'EQ');
CREATE TYPE direction AS ENUM ('LONG', 'SHORT');
CREATE TYPE transaction_type AS ENUM ('BUY', 'SELL');
CREATE TYPE product_type AS ENUM('CNC', 'MIS', 'NRML');
CREATE TYPE order_type AS ENUM('MARKET');
CREATE TYPE order_state AS ENUM('SENT', 'EXECUTED', 'CANCELLED');
CREATE TYPE value_type AS ENUM('ABSOLUTE', 'PERCENTAGE');
CREATE TYPE chart_type AS ENUM('CANDLESTICK');
CREATE TYPE chart_timeframe AS ENUM('1m', '3m', '5m', '10m', '15m', '30m', '1h', '1d');
CREATE TYPE trade_type AS ENUM ('MIS', 'CNC/NRML');

CREATE TYPE run_type AS ENUM ('BACKTEST', 'PAPER_TRADE', 'LIVE_TRADE');
CREATE TYPE run_state AS ENUM ('NOT_STARTED', 'PAUSED', 'STARTED', 'ENDED');
CREATE TYPE position_state AS ENUM ('NOT_ENTERED', 'ENTERED', 'EXITED');
