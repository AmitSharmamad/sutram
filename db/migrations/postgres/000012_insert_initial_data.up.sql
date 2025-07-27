-- Insert initial data

-- Insert default OAuth providers
INSERT INTO oauth_provider (name, display_name) VALUES
('google', 'Google'),
('github', 'GitHub')
ON CONFLICT (name) DO NOTHING;

-- Insert default usage limits for existing users based on their roles
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
