-- Tables to store chatbot conversations

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- AI Chat table to store conversation metadata
CREATE TABLE IF NOT EXISTS ai_chat (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- AI Message table to store individual messages in a conversation
CREATE TABLE IF NOT EXISTS ai_message (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID NOT NULL REFERENCES ai_chat(id) ON DELETE CASCADE,
    query TEXT NOT NULL,
    response TEXT NOT NULL,
    feedback VARCHAR(10) CHECK (feedback IN ('like', 'dislike', NULL)),
    feedback_comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- AI Actions table to store actions associated with messages
CREATE TABLE IF NOT EXISTS ai_action (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id UUID NOT NULL REFERENCES ai_message(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_ai_message_chat_id ON ai_message(chat_id);
CREATE INDEX IF NOT EXISTS idx_ai_action_message_id ON ai_action(message_id);
CREATE INDEX IF NOT EXISTS idx_ai_chat_user_id ON ai_chat(user_id);
