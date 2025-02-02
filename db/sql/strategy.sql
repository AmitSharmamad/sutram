USE sutram;

-- Strategy Management Tables
CREATE TABLE IF NOT EXISTS strategy (
    id UUID DEFAULT generateUUIDv4(),
    user_id UUID,
    name String,
    description String,
    position_type String,
    instruments Array(String),
    chart_type String,
    chart_timeframe String,
    chart_unknown String,
    qty Int32,
    entry_expression String,
    exit_expression String,
    target_profit Decimal64(2),
    stop_loss Decimal64(2),
    trailing_stop_loss Decimal64(2),
    exit_type String,
    capital Decimal64(2),
    max_trades Int32,
    entry_start_time DateTime,
    entry_end_time DateTime,
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
