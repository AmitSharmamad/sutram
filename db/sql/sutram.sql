CREATE DATABASE IF NOT EXISTS sutram;

USE sutram;

CREATE USER IF NOT EXISTS test IDENTIFIED WITH plaintext_password BY 'test';

GRANT ALL ON sutram.* TO test;


-- User Management Tables
CREATE TABLE IF NOT EXISTS user (
    id UUID DEFAULT generateUUIDv4(),
    name String,
    username String,
    email String,
    phone_number String,
    password String,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    is_active Boolean DEFAULT true,
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS portfolio (
    id UUID DEFAULT generateUUIDv4(),
    user_id UUID,
    capital Decimal64(2),
    returns Decimal64(2),
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS setting (
    id UUID DEFAULT generateUUIDv4(),
    user_id UUID,
    config String, -- JSON encoded string
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS notification (
    id UUID DEFAULT generateUUIDv4(),
    user_id UUID,
    type Enum('BT_DONE' = 1, 'LT_TRADE_ENTERED' = 2, 'LT_TRADE_EXITED' = 3),
    status Enum('DELIVER' = 1, 'DELIVERED' = 2),
    receiver Enum('MOBILE' = 1, 'EMAIL' = 2),
    content String,
    PRIMARY KEY (id)
) ENGINE = MergeTree;

-- Stock Management Tables
CREATE TABLE IF NOT EXISTS stock_master (
    id UUID DEFAULT generateUUIDv4(),
    exchange Enum('NSE' = 1, 'BSE' = 2),
    name String,
    symbol String,
    is_active Boolean DEFAULT true,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

-- Strategy Management Tables
CREATE TABLE IF NOT EXISTS strategy (
    id UUID DEFAULT generateUUIDv4(),
    user_id UUID,
    name String,
    description String,
    position_type Enum8('BUY' = 1, 'SELL' = 2), 
    instruments Array(String),
    chart_type Enum8('CANDLESTICK' = 1, 'HEIKEN-ASHI' = 2, 'RENKO' = 3),
    chart_timeframe Enum8('1m' = 1, '3m' = 2, '5m' = 3, '10m' = 4, '15m' = 5, '30m' = 6, '1h' = 7, '1d' = 8),
    chart_unknown String,
    qty Int32,
    entry_expression String,
    exit_expression String,
    target_profit Decimal64(2),
    stop_loss Decimal64(2),
    trailing_stop_loss Decimal64(2),
    exit_type Enum8('PERCENTAGE' = 1, 'ABSOLUTE' = 2),
    capital Decimal64(2),
    max_trades Int32,
    entry_start_time DateTime64(3),
    entry_end_time DateTime64(3),
    start_date Date,
    end_date Date,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    is_active Boolean DEFAULT true,
    PRIMARY KEY (id)
) ENGINE = MergeTree
ORDER BY id;

CREATE TABLE IF NOT EXISTS strategy_history (
    id UUID DEFAULT generateUUIDv4(),
    strategy_id UUID,
    change String, -- JSON encoded strategy
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

-- Backtest Engine Tables
CREATE TABLE IF NOT EXISTS bt_run (
    id UUID DEFAULT generateUUIDv4(),
    strategy_id UUID,
    state Enum('start' = 1, 'in_progress' = 2, 'ended' = 3),
    qty UInt32,
    period_start_time DateTime,
    period_end_time DateTime,
    started_time DateTime,
    ended_time DateTime,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS bt_position (
    id UUID DEFAULT generateUUIDv4(),
    run_id UUID,
    state Enum('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    entry_price Decimal64(2),
    entry_time DateTime,
    exit_price Decimal64(2),
    exit_time DateTime,
    exit_reason String,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS bt_position_log (
    id UUID DEFAULT generateUUIDv4(),
    position_id UUID,
    state Enum('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    created_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

-- Paper Trading Tables
CREATE TABLE IF NOT EXISTS pt_run (
    id UUID DEFAULT generateUUIDv4(),
    strategy_id UUID,
    state Enum('start' = 1, 'in_progress' = 2, 'ended' = 3),
    qty UInt32,
    period_start_time DateTime,
    period_end_time DateTime,
    started_time DateTime,
    ended_time DateTime,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS pt_position (
    id UUID DEFAULT generateUUIDv4(),
    run_id UUID,
    state Enum('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    entry_price Decimal64(2),
    entry_time DateTime,
    exit_price Decimal64(2),
    exit_time DateTime,
    exit_reason String,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS pt_position_log (
    id UUID DEFAULT generateUUIDv4(),
    position_id UUID,
    state Enum('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    created_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

-- Live Trading Tables
CREATE TABLE IF NOT EXISTS lt_run (
    id UUID DEFAULT generateUUIDv4(),
    strategy_id UUID,
    state Enum('start' = 1, 'in_progress' = 2, 'ended' = 3),
    qty UInt32,
    period_start_time DateTime,
    period_end_time DateTime,
    started_time DateTime,
    ended_time DateTime,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS lt_position (
    id UUID DEFAULT generateUUIDv4(),
    run_id UUID,
    state Enum('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    entry_price Decimal64(2),
    entry_time DateTime,
    exit_price Decimal64(2),
    exit_time DateTime,
    exit_reason String,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS lt_position_log (
    id UUID DEFAULT generateUUIDv4(),
    position_id UUID,
    state Enum('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    created_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

-- Broker Management Tables
CREATE TABLE IF NOT EXISTS broker (
    id UUID DEFAULT generateUUIDv4(),
    name String,
    client_id String,
    client_secret String,
    auth_type Enum('oauth2.0' = 1),
    PRIMARY KEY (id)
) ENGINE = MergeTree;

CREATE TABLE IF NOT EXISTS broker_user (
    user_id UUID,
    broker_id UUID,
    token String,
    expiry DateTime,
    PRIMARY KEY (user_id, broker_id)
) ENGINE = MergeTree;

-- Template for stock tables (replace SYMBOL with actual stock symbol)
CREATE TABLE IF NOT EXISTS stock_TESLA (
    timestamp DateTime,
    open Decimal64(2),
    high Decimal64(2),
    close Decimal64(2),
    low Decimal64(2),
    volume UInt64,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now(),
    PRIMARY KEY (timestamp)
) ENGINE = MergeTree;