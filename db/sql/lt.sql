USE sutram;

-- Live Trading Tables
CREATE TABLE IF NOT EXISTS lt_run (
    id UUID DEFAULT generateUUIDv4(),
    strategy_id UUID,
    state String,
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
    state String,
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
    state String,
    created_at DateTime DEFAULT now(),
    PRIMARY KEY (id)
) ENGINE = MergeTree;
