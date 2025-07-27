-- User Management Tables
CREATE TABLE IF NOT EXISTS "user" (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    password VARCHAR(255),
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT true,
    PRIMARY KEY (id)
);

-- User Roles and Permissions Tables

-- Roles table for defining platform roles
CREATE TABLE IF NOT EXISTS role (
    id SERIAL PRIMARY KEY,
    name role_type UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- User-Role assignment (many-to-many relationship)
CREATE TABLE IF NOT EXISTS user_role (
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES role(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (user_id, role_id)
);

-- Permissions table
CREATE TABLE IF NOT EXISTS permission (
    id SERIAL PRIMARY KEY,
    name permission_type UNIQUE NOT NULL,
    description TEXT,
    resource_type VARCHAR(50) NOT NULL, -- 'strategy', 'backtest', 'marketplace', etc.
    action VARCHAR(50) NOT NULL, -- 'create', 'read', 'update', 'delete', etc.
    created_at TIMESTAMP DEFAULT now()
);

-- Role-Permission assignment (many-to-many relationship)
CREATE TABLE IF NOT EXISTS role_permission (
    role_id INT NOT NULL REFERENCES role(id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permission(id) ON DELETE CASCADE,
    granted_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (role_id, permission_id)
);

-- User-specific permissions (override role permissions)
CREATE TABLE IF NOT EXISTS user_permission (
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permission(id) ON DELETE CASCADE,
    is_granted BOOLEAN NOT NULL DEFAULT true, -- true=grant, false=deny (override)
    granted_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (user_id, permission_id)
);

-- Create indexes for performance
CREATE INDEX idx_user_role_user ON user_role(user_id);
CREATE INDEX idx_user_role_role ON user_role(role_id);
CREATE INDEX idx_role_permission_role ON role_permission(role_id);
CREATE INDEX idx_user_permission_user ON user_permission(user_id);
