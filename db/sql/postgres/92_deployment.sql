-- Deployment Management Tables
CREATE TABLE IF NOT EXISTS deployment (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    strategy_id UUID NOT NULL REFERENCES strategy(id) ON DELETE CASCADE,
    instrument_id INT NOT NULL,
    name VARCHAR NOT NULL,
    description VARCHAR,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id),
    UNIQUE (user_id, strategy_id, instrument_id)
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_deployment_user_id ON deployment(user_id);
CREATE INDEX IF NOT EXISTS idx_deployment_strategy_id ON deployment(strategy_id);
CREATE INDEX IF NOT EXISTS idx_deployment_instrument_id ON deployment(instrument_id);

-- Deployment Log Table - Tracks progression of a deployment (immutable)
CREATE TABLE IF NOT EXISTS deployment_log (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    deployment_id UUID NOT NULL REFERENCES deployment(id) ON DELETE CASCADE,
    status VARCHAR NOT NULL,
    metadata JSONB,
    action_at TIMESTAMP DEFAULT now()
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_deployment_log_deployment_id ON deployment_log(deployment_id);
CREATE INDEX IF NOT EXISTS idx_deployment_log_status ON deployment_log(status);

-- Deployment Position Table - Stores deployment positions
CREATE TABLE IF NOT EXISTS deployment_position (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    deployment_id UUID NOT NULL REFERENCES deployment(id) ON DELETE CASCADE,
    instrument_id INT NOT NULL,
    quantity DECIMAL NOT NULL,
    entry_price DECIMAL NOT NULL,
    entry_time TIMESTAMP NOT NULL,
    exit_price DECIMAL,
    exit_time TIMESTAMP,
    pnl DECIMAL,
    status VARCHAR NOT NULL,
    direction VARCHAR,
    stop_loss_price DECIMAL,
    initial_stop_loss_price DECIMAL,
    target_profit_price DECIMAL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_deployment_position_deployment_id ON deployment_position(deployment_id);
CREATE INDEX IF NOT EXISTS idx_deployment_position_instrument_id ON deployment_position(instrument_id);
CREATE INDEX IF NOT EXISTS idx_deployment_position_status ON deployment_position(status);

-- Deployment Metrics Table - Stores deployment metrics
CREATE TABLE IF NOT EXISTS deployment_metrics (
    id UUID DEFAULT gen_random_uuid()  PRIMARY KEY,
    deployment_id UUID NOT NULL REFERENCES deployment(id) ON DELETE CASCADE,
    metric_date DATE NOT NULL,
    total_trades INT DEFAULT 0,
    win_trades INT DEFAULT 0,
    loss_trades INT DEFAULT 0,
    total_pnl DECIMAL DEFAULT 0,
    max_drawdown DECIMAL DEFAULT 0,
    sharpe_ratio DECIMAL,
    sortino_ratio DECIMAL,
    win_rate DECIMAL,
    avg_win DECIMAL,
    avg_loss DECIMAL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (id, deployment_id)
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_deployment_metrics_deployment_id ON deployment_metrics(deployment_id);
CREATE INDEX IF NOT EXISTS idx_deployment_metrics_metric_date ON deployment_metrics(metric_date);
