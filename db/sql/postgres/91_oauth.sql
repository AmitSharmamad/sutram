-- Add OAuth-related tables

-- Table to store OAuth provider information
CREATE TABLE IF NOT EXISTS oauth_provider (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, -- 'google', 'github', etc.
    display_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Table to store user OAuth connection
CREATE TABLE IF NOT EXISTS user_oauth (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    provider_id INT NOT NULL REFERENCES oauth_provider(id) ON DELETE CASCADE,
    provider_user_id VARCHAR(255) NOT NULL, -- ID from the provider
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP,
    user_data JSONB, -- Store additional user data from provider
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(provider_id, provider_user_id)
);

-- Insert default providers
INSERT INTO oauth_provider (name, display_name) VALUES
('google', 'Google'),
('github', 'GitHub')
ON CONFLICT (name) DO NOTHING;

-- Create indexes
CREATE INDEX idx_user_oauth_user_id ON user_oauth(user_id);
CREATE INDEX idx_user_oauth_provider ON user_oauth(provider_id);
