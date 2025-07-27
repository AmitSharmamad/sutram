-- User Usage Limits Tables

-- User usage limits table
CREATE TABLE IF NOT EXISTS user_usage_limits (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    limit_type usage_limit_type NOT NULL,
    limit_value INT NOT NULL DEFAULT 0,
    current_usage INT NOT NULL DEFAULT 0,
    reset_period reset_period_type NOT NULL DEFAULT 'monthly',
    last_reset_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(user_id, limit_type)
);

-- Usage tracking log
CREATE TABLE IF NOT EXISTS user_usage_log (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    limit_type usage_limit_type NOT NULL,
    action_type VARCHAR(50) NOT NULL, -- 'increment', 'decrement', 'reset'
    previous_value INT NOT NULL,
    new_value INT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT now()
);

-- Create indexes for performance
CREATE INDEX idx_user_usage_limits_user ON user_usage_limits(user_id);
CREATE INDEX idx_user_usage_limits_type ON user_usage_limits(limit_type);
CREATE INDEX idx_user_usage_log_user ON user_usage_log(user_id);
CREATE INDEX idx_user_usage_log_type ON user_usage_log(limit_type);
CREATE INDEX idx_user_usage_log_created ON user_usage_log(created_at);
