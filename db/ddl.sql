CREATE DATABASE IF NOT EXISTS sutram_test_db;

CREATE TABLE IF NOT EXISTS sutram_test_db.user
(
    id UInt32,
    name String,
    username String,
    email String,
    phone_number String,
    password String,
    created_at DateTime,
    updated_at DateTime,
    is_active UInt8
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.portfolio
(
    id UInt32,
    user_id UInt32,
    stock String,
    description String,
    currency String,
    invested_value Float32,
    current_value Float32,
    amount Float32,
    profit Float32,
    loss Float32,
    created_at DateTime,
    updated_at DateTime,
    is_active UInt8
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.portfolio_log
(
    id UInt32,
    portfolio_id UInt32,
    action_type Enum8('FUND_ADDED' = 1, 'FUND_WITHDRAWN' = 2, 'FUND_REFUNDED' = 3),
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.notification
(
    id UInt32,
    user_id UInt32
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.setting
(
    id UInt32,
    user_id UInt32,
    config String,
    created_at DateTime,
    updated_at DateTime,
    is_active UInt8
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.strategy
(
    id UInt32,
    name String,
    description String,
    symbol String,
    instrument_type Enum8('EQUITY' = 1, 'FUTURES' = 2, 'OPTIONS' = 3),
    position_type Enum8('BUY' = 1, 'SELL' = 2),
    entry_expression String,
    exit_expression String,
    stop_loss Float32,
    trailing_stop_loss Float32,
    max_trades UInt32,
    created_at DateTime,
    updated_at DateTime,
    is_active UInt8
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.strategy_history
(
    id UInt32,
    strategy_id UInt32,
    change String,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.stock
(
    id UInt32,
    exchange String,
    name String,
    symbol String,
    created_at DateTime,
    modified_at DateTime,
    is_active UInt8
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.stock_AAPL
(
    timestamp DateTime,
    open Float32,
    high Float32,
    close Float32,
    low Float32,
    volume UInt32,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (timestamp);

CREATE TABLE IF NOT EXISTS sutram_test_db.stock_TSLA
(
    timestamp DateTime,
    open Float32,
    high Float32,
    close Float32,
    low Float32,
    volume UInt32,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (timestamp);

CREATE TABLE IF NOT EXISTS sutram_test_db.bt_run
(
    id UInt32,
    strategy_id UInt32,
    state Enum8('not_picked' = 1, 'started' = 2, 'in_progress' = 3, 'ended' = 4),
    qty UInt32,
    period_start_time DateTime,
    period_end_time DateTime,
    run_start_time DateTime,
    run_stop_time DateTime,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.bt_position
(
    id UInt32,
    run_id UInt32,
    state Enum8('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    entry_price Float32,
    entry_time DateTime,
    exit_price Float32,
    exit_time DateTime,
    exit_reason String,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.bt_position_log
(
    id UInt32,
    position_id UInt32,
    state Enum8('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    created_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.pt_run
(
    id UInt32,
    strategy_id UInt32,
    state Enum8('not_picked' = 1, 'started' = 2, 'in_progress' = 3, 'ended' = 4),
    start_time DateTime,
    stop_time DateTime,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.pt_position
(
    id UInt32,
    run_id UInt32,
    state Enum8('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    entry_price Float32,
    entry_time DateTime,
    exit_price Float32,
    exit_time DateTime,
    exit_reason String,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.pt_position_log
(
    id UInt32,
    position_id UInt32,
    state Enum8('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    created_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.lt_run
(
    id UInt32,
    strategy_id UInt32,
    state Enum8('not_picked' = 1, 'started' = 2, 'in_progress' = 3, 'ended' = 4),
    start_time DateTime,
    stop_time DateTime,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.lt_position
(
    id UInt32,
    run_id UInt32,
    state Enum8('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    entry_price Float32,
    entry_time DateTime,
    exit_price Float32,
    exit_time DateTime,
    exit_reason String,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.lt_position_log
(
    id UInt32,
    position_id UInt32,
    state Enum8('enter' = 1, 'entered' = 2, 'hold' = 3, 'exited' = 4),
    created_at DateTime
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.broker
(
    id UInt32,
    name String,
    client_id String,
    client_secret String,
    auth_type Enum8('oauth2.0' = 1)
)
ENGINE = MergeTree()
PRIMARY KEY (id);

CREATE TABLE IF NOT EXISTS sutram_test_db.broker_user
(
    user_id UInt32,
    broker_id UInt32,
    token String,
    expiry DateTime,
    is_active UInt8
)
ENGINE = MergeTree()
PRIMARY KEY (user_id, broker_id);
