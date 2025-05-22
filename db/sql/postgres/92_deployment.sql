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
