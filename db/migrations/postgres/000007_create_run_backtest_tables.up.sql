-- Run Group Tables
CREATE TABLE IF NOT EXISTS run_group (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    strategy_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (strategy_id) REFERENCES strategy(id) ON DELETE CASCADE
);

-- Backtest Engine Tables
CREATE TABLE IF NOT EXISTS run (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    run_group_id UUID NOT NULL,
    strategy_id UUID NOT NULL,
    -- Instrument Details
    exchange VARCHAR NOT NULL,
    instrument_id INT NOT NULL,
    type run_type NOT NULL,
    state run_state NOT NULL,
    qty INT NOT NULL,
    period_start_time TIMESTAMP NOT NULL,
    period_end_time TIMESTAMP NOT NULL,
    started_time TIMESTAMP,
    ended_time TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (run_group_id) REFERENCES run_group(id) ON DELETE CASCADE,
    FOREIGN KEY (strategy_id) REFERENCES strategy(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS position (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL,
    state position_state NOT NULL,
    direction direction NOT NULL,
    -- Entry Details
    entry_price_signal DECIMAL(18, 8) NOT NULL,
    entry_time_signal TIMESTAMP NOT NULL,
    entry_price_order DECIMAL(18, 8),
    entry_time_order TIMESTAMP,
    entry_signal VARCHAR NOT NULL,
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
    position_size DECIMAL(18, 8) NOT NULL, -- Quantity of the instrument
    -- Order Details
    entry_order_id UUID,
    exit_order_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    FOREIGN KEY (run_id) REFERENCES run(id) ON DELETE CASCADE,
    FOREIGN KEY (entry_order_id) REFERENCES "order"(id) ON DELETE SET NULL,
    FOREIGN KEY (exit_order_id) REFERENCES "order"(id) ON DELETE SET NULL
);

-- Metrics Table
CREATE TABLE IF NOT EXISTS run_metrics (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL,
    -- Returns Metrics
    total_return_pct DECIMAL(10, 4) NOT NULL,
    annualized_return DECIMAL(10, 4) NOT NULL,
    risk_adjusted_return DECIMAL(10, 4),
    daily_returns DECIMAL(10, 4)[],
    monthly_returns DECIMAL(10, 4)[],

    -- Risk Metrics
    sharpe_ratio DECIMAL(10, 4) NOT NULL,
    sortino_ratio DECIMAL(10, 4) NOT NULL,
    max_drawdown DECIMAL(10, 4) NOT NULL,
    max_drawdown_duration DECIMAL(10, 4) NOT NULL,
    value_at_risk DECIMAL(10, 4),
    std_deviation DECIMAL(10, 4) NOT NULL,
    downside_deviation DECIMAL(10, 4),

    -- Trade Metrics
    total_trades INT NOT NULL,
    win_rate DECIMAL(10, 4) NOT NULL,
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
    total_commission DECIMAL(18, 8) NOT NULL,
    total_slippage DECIMAL(18, 8) NOT NULL,
    transaction_costs_pct DECIMAL(10, 4) NOT NULL,

    -- json to save all the metric details
    metric_json TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    PRIMARY KEY (id),

    FOREIGN KEY (run_id) REFERENCES run(id) ON DELETE CASCADE,
    UNIQUE (run_id)
);
