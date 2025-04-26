-- Run Group Tables
CREATE TABLE IF NOT EXISTS run_group (
    id UUID DEFAULT gen_random_uuid(),
    strategy_id UUID,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (strategy_id) REFERENCES strategy(id) ON DELETE CASCADE
);

-- Backtest Engine Tables
CREATE TABLE IF NOT EXISTS run (
    id UUID DEFAULT gen_random_uuid(),
    run_group_id UUID,
    strategy_id UUID,
    -- Instrument Details
    exchange VARCHAR,
    instrument_id INT,
    type run_type,
    state run_state,
    qty INT,
    period_start_time TIMESTAMP,
    period_end_time TIMESTAMP,
    started_time TIMESTAMP,
    ended_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (run_group_id) REFERENCES run_group(id) ON DELETE CASCADE,
    FOREIGN KEY (strategy_id) REFERENCES strategy(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS position (
    id UUID DEFAULT gen_random_uuid(),
    run_id UUID,
    state position_state,
    direction direction,
    -- Entry Details
    entry_price_signal DECIMAL(18, 8),
    entry_time_signal TIMESTAMP,
    entry_price_order DECIMAL(18, 8),
    entry_time_order TIMESTAMP,
    entry_signal VARCHAR,
    -- Exit Details
    exit_price_signal DECIMAL(18, 8),
    exit_time_signal TIMESTAMP,
    exit_price_order DECIMAL(18, 8),
    exit_time_order TIMESTAMP,
    exit_signal VARCHAR,
    exit_reason VARCHAR,
    -- Profit/Loss Details
    profit_loss_absolute DECIMAL(18, 8),
    profit_loss_percentage DECIMAL(18, 8),
    -- Transaction Costs
    transaction_cost DECIMAL(18, 8), -- Commission, fees, slippage
    entry_fees DECIMAL(18, 8),
    exit_fees DECIMAL(18, 8),
    net_profit_loss DECIMAL(18, 8), -- P/L after fees
    -- Position Details
    position_size DECIMAL(18, 8), -- Quantity of the instrument
    -- Order Details
    entry_order_id UUID,
    exit_order_id UUID,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (run_id) REFERENCES run(id) ON DELETE CASCADE
);

-- Metrics Table
CREATE TABLE IF NOT EXISTS run_metrics (
    id UUID DEFAULT gen_random_uuid(),
    run_id UUID,
    -- Returns Metrics
    total_return_pct DECIMAL(10, 4),
    annualized_return DECIMAL(10, 4),
    risk_adjusted_return DECIMAL(10, 4),
    daily_returns DECIMAL(10, 4)[],
    monthly_returns DECIMAL(10, 4)[],

    -- Risk Metrics
    sharpe_ratio DECIMAL(10, 4),
    sortino_ratio DECIMAL(10, 4),
    max_drawdown DECIMAL(10, 4),
    max_drawdown_duration DECIMAL(10, 4),
    value_at_risk DECIMAL(10, 4),
    std_deviation DECIMAL(10, 4),
    downside_deviation DECIMAL(10, 4),

    -- Trade Metrics
    total_trades INT,
    win_rate DECIMAL(10, 4),
    avg_win_size DECIMAL(18, 8),
    avg_loss_size DECIMAL(18, 8),
    avg_trade_duration DECIMAL(10, 4),
    profit_factor DECIMAL(10, 4),

    -- Position Metrics
    avg_position_size DECIMAL(18, 8),
    max_position_size DECIMAL(18, 8),
    position_turnover DECIMAL(10, 4),
    avg_exposure DECIMAL(18, 8),

    -- Cost Metrics
    total_commission DECIMAL(18, 8),
    total_slippage DECIMAL(18, 8),
    transaction_costs_pct DECIMAL(10, 4),

    -- json to save all the metric details
    metric_json TEXT,

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),

    PRIMARY KEY (id),

    FOREIGN KEY (run_id) REFERENCES run(id) ON DELETE CASCADE,
    UNIQUE (run_id)
);