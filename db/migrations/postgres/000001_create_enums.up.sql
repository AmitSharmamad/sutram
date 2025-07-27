-- Create an ENUM type
CREATE TYPE exchange AS ENUM ('NSE', 'BSE', 'MCX');
CREATE TYPE segment AS ENUM ('F&O', 'EQ');
CREATE TYPE direction AS ENUM ('LONG', 'SHORT');
CREATE TYPE transaction_type AS ENUM ('BUY', 'SELL');
CREATE TYPE product_type AS ENUM('CNC', 'MIS', 'NRML');
CREATE TYPE order_type AS ENUM('MARKET', 'LIMIT');
CREATE TYPE order_state AS ENUM('SENT', 'EXECUTED', 'CANCELLED');
CREATE TYPE value_type AS ENUM('ABSOLUTE', 'PERCENTAGE');
CREATE TYPE chart_type AS ENUM('CANDLESTICK');
CREATE TYPE chart_timeframe AS ENUM('1m', '3m', '5m', '10m', '15m', '30m', '1h', '1d');
CREATE TYPE trade_type AS ENUM ('MIS', 'CNC/NRML');

CREATE TYPE run_type AS ENUM ('BACKTEST', 'PAPER_TRADE', 'LIVE_TRADE');
CREATE TYPE run_state AS ENUM ('NOT_STARTED', 'PAUSED', 'STARTED', 'ENDED');
CREATE TYPE position_state AS ENUM ('NOT_ENTERED', 'ENTERED', 'EXITED');

-- Enum for roles
CREATE TYPE role_type AS ENUM ('admin', 'paid_user', 'free_user', 'moderator', 'analyst');

-- Enum for permissions
CREATE TYPE permission_type AS ENUM (
    -- Strategy permissions
    'create_strategy', 'view_any_strategy', 'edit_own_strategy', 'delete_own_strategy',
    -- Backtest permissions
    'run_backtest', 'view_own_backtest', 'view_any_backtest',
    -- Marketplace permissions
    'publish_strategy', 'purchase_strategy',
    -- Paper trading
    'run_paper_trade',
    -- Live trading
    'run_live_trade',
    -- Community permissions
    'create_topic', 'moderate_content'
);

-- Enum for usage limit types
CREATE TYPE usage_limit_type AS ENUM (
    'backtest_per_month',
    'deployment_count'
);

-- Enum for reset period types
CREATE TYPE reset_period_type AS ENUM (
    'hourly', 'daily', 'weekly', 'monthly', 'yearly', 'never'
);

-- Enums for options trading support
CREATE TYPE strike_type AS ENUM (
    'OTM20', 'OTM19', 'OTM18', 'OTM17', 'OTM16', 'OTM15', 'OTM14', 'OTM13', 'OTM12', 'OTM11',
    'OTM10', 'OTM9', 'OTM8', 'OTM7', 'OTM6', 'OTM5', 'OTM4', 'OTM3', 'OTM2', 'OTM1',
    'ATM',
    'ITM1', 'ITM2', 'ITM3', 'ITM4', 'ITM5', 'ITM6', 'ITM7', 'ITM8', 'ITM9', 'ITM10',
    'ITM11', 'ITM12', 'ITM13', 'ITM14', 'ITM15', 'ITM16', 'ITM17', 'ITM18', 'ITM19', 'ITM20'
);

CREATE TYPE expiry_cycle AS ENUM ('CURRENT', 'NEXT', 'NEXT2');
CREATE TYPE expiry_type AS ENUM ('WEEKLY', 'MONTHLY', 'QUARTERLY');
CREATE TYPE instrument_type AS ENUM ('EQUITY', 'FUTURES', 'OPTIONS');
CREATE TYPE option_type AS ENUM ('CALL', 'PUT');
CREATE TYPE asset AS ENUM ('INDEX', 'STOCK');
CREATE TYPE strategy_type AS ENUM ('CONDITIONAL', 'TIME_BASED');

CREATE TYPE trading_day AS ENUM ('MON', 'TUE', 'WED', 'THU', 'FRI');
