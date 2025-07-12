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

-- Insert default limits for existing users based on their roles
-- Backtest limits
INSERT INTO user_usage_limits (user_id, limit_type, limit_value, current_usage, reset_period)
SELECT DISTINCT u.id, 'backtest_per_month'::usage_limit_type, 
       CASE 
           WHEN EXISTS (SELECT 1 FROM user_role ur JOIN role r ON ur.role_id = r.id WHERE ur.user_id = u.id AND r.name = 'admin') THEN 999999
           WHEN EXISTS (SELECT 1 FROM user_role ur JOIN role r ON ur.role_id = r.id WHERE ur.user_id = u.id AND r.name = 'paid_user') THEN 100
           WHEN EXISTS (SELECT 1 FROM user_role ur JOIN role r ON ur.role_id = r.id WHERE ur.user_id = u.id AND r.name = 'free_user') THEN 10
           ELSE 5
       END,
       0, 'monthly'::reset_period_type
FROM "user" u
ON CONFLICT (user_id, limit_type) DO NOTHING;

-- Deployment limits
INSERT INTO user_usage_limits (user_id, limit_type, limit_value, current_usage, reset_period)
SELECT DISTINCT u.id, 'deployment_count'::usage_limit_type,
       CASE 
           WHEN EXISTS (SELECT 1 FROM user_role ur JOIN role r ON ur.role_id = r.id WHERE ur.user_id = u.id AND r.name = 'admin') THEN 999999
           WHEN EXISTS (SELECT 1 FROM user_role ur JOIN role r ON ur.role_id = r.id WHERE ur.user_id = u.id AND r.name = 'paid_user') THEN 20
           WHEN EXISTS (SELECT 1 FROM user_role ur JOIN role r ON ur.role_id = r.id WHERE ur.user_id = u.id AND r.name = 'free_user') THEN 3
           ELSE 1
       END,
       0, 'never'::reset_period_type
FROM "user" u
ON CONFLICT (user_id, limit_type) DO NOTHING;

-- Live trades per day limits
INSERT INTO user_usage_limits (user_id, limit_type, limit_value, current_usage, reset_period)
SELECT DISTINCT u.id, 'live_trades_per_day'::usage_limit_type,
       CASE 
           WHEN EXISTS (SELECT 1 FROM user_role ur JOIN role r ON ur.role_id = r.id WHERE ur.user_id = u.id AND r.name = 'admin') THEN 999999
           WHEN EXISTS (SELECT 1 FROM user_role ur JOIN role r ON ur.role_id = r.id WHERE ur.user_id = u.id AND r.name = 'paid_user') THEN 500
           WHEN EXISTS (SELECT 1 FROM user_role ur JOIN role r ON ur.role_id = r.id WHERE ur.user_id = u.id AND r.name = 'free_user') THEN 50
           ELSE 10
       END,
       0, 'daily'::reset_period_type
FROM "user" u
ON CONFLICT (user_id, limit_type) DO NOTHING;
