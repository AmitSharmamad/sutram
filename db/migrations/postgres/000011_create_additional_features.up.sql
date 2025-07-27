-- Additional feature tables

-- Create waitlist table
CREATE TABLE waitlist (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX idx_waitlist_email ON waitlist(email);

-- Create index on created_at for analytics
CREATE INDEX idx_waitlist_created_at ON waitlist(created_at);

-- OAuth provider information
CREATE TABLE IF NOT EXISTS oauth_provider (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, -- 'google', 'github', etc.
    display_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- User OAuth connection
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

-- Create indexes
CREATE INDEX idx_user_oauth_user_id ON user_oauth(user_id);
CREATE INDEX idx_user_oauth_provider ON user_oauth(provider_id);

-- AI Chat table to store conversation metadata
CREATE TABLE IF NOT EXISTS ai_chat (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- AI Message table to store individual messages in a conversation
CREATE TABLE IF NOT EXISTS ai_message (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID NOT NULL REFERENCES ai_chat(id) ON DELETE CASCADE,
    query TEXT NOT NULL,
    response TEXT NOT NULL,
    feedback VARCHAR(10) CHECK (feedback IN ('like', 'dislike', NULL)),
    feedback_comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- AI Actions table to store actions associated with messages
CREATE TABLE IF NOT EXISTS ai_action (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL REFERENCES ai_message(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_ai_message_chat_id ON ai_message(chat_id);
CREATE INDEX IF NOT EXISTS idx_ai_action_message_id ON ai_action(message_id);
CREATE INDEX IF NOT EXISTS idx_ai_chat_user_id ON ai_chat(user_id);
