-- Insert default roles
INSERT INTO role (name, description) VALUES
('admin', 'Full system access with all permissions'),
('paid_user', 'Access to all trading features'),
('free_user', 'Limited access to basic platform features'),
('moderator', 'Community moderation capabilities'),
('analyst', 'Can publish verified market analysis');

-- Insert core permissions
INSERT INTO permission (name, resource_type, action, description) VALUES
-- Strategy permissions
('create_strategy', 'strategy', 'create', 'Can create strategies'),
('view_any_strategy', 'strategy', 'read', 'Can view any strategy'),
('edit_own_strategy', 'strategy', 'update', 'Can edit own strategies'),
('delete_own_strategy', 'strategy', 'delete', 'Can delete own strategies'),

-- Backtest permissions
('run_backtest', 'backtest', 'create', 'Can run backtests'),
('view_own_backtest', 'backtest', 'read', 'Can view own backtest results'),
('view_any_backtest', 'backtest', 'read', 'Can view any backtest results'),

-- Marketplace permissions
('publish_strategy', 'marketplace', 'create', 'Can publish strategies to marketplace'),
('purchase_strategy', 'marketplace', 'create', 'Can purchase strategies'),

-- Paper trading
('run_paper_trade', 'paper_trade', 'create', 'Can run paper trading'),

-- Live trading
('run_live_trade', 'live_trade', 'create', 'Can run live trading'),

-- Community permissions
('create_topic', 'community', 'create', 'Can create community topics'),
('moderate_content', 'community', 'update', 'Can moderate community content');

-- Assign permissions to roles
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r, permission p
WHERE r.name = 'admin'; -- Admin gets all permissions

-- Basic permissions for free users
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r, permission p
WHERE r.name = 'free_user' AND p.name IN (
    'create_strategy',
    'edit_own_strategy',
    'delete_own_strategy',
    'view_own_backtest',
    'run_backtest',
    'run_paper_trade',
    'create_topic'
);

-- Additional permissions for paid users (includes all free user permissions)
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r, permission p
WHERE r.name = 'paid_user' AND p.name IN (
    'create_strategy',
    'edit_own_strategy',
    'delete_own_strategy',
    'view_own_backtest',
    'run_backtest',
    'run_paper_trade',
    'create_topic',
    'publish_strategy',
    'purchase_strategy',
    'run_live_trade',
    'view_any_strategy'
);

-- Moderator permissions
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r, permission p
WHERE r.name = 'moderator' AND p.name IN (
    'moderate_content',
    'view_any_strategy'
);

-- Analyst permissions
INSERT INTO role_permission (role_id, permission_id)
SELECT r.id, p.id
FROM role r, permission p
WHERE r.name = 'analyst' AND p.name IN (
    'create_strategy',
    'edit_own_strategy',
    'delete_own_strategy',
    'view_own_backtest',
    'run_backtest',
    'publish_strategy'
);
