-- Community Tables for Algo Trading Platform
-- This defines tables for managing community features like discussions, topics, replies, likes, follows, badges, etc.

BEGIN;

-- User Reputation table to track user standing in the community
CREATE TABLE user_reputation (
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    reputation_score INT NOT NULL DEFAULT 0,
    contribution_count INT NOT NULL DEFAULT 0,
    joined_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id)
);

-- Categories for organizing discussions and content
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(20),
    parent_id INT REFERENCES categories(id) ON DELETE SET NULL,
    display_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Topics/Discussions
CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    is_pinned BOOLEAN NOT NULL DEFAULT false,
    is_closed BOOLEAN NOT NULL DEFAULT false,
    is_hidden BOOLEAN NOT NULL DEFAULT false,
    view_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index for faster topic searches
CREATE INDEX idx_topics_category ON topics(category_id);
CREATE INDEX idx_topics_author ON topics(author_id);
CREATE INDEX idx_topics_created_at ON topics(created_at);
CREATE INDEX idx_topics_last_activity ON topics(last_activity_at);

-- Topics metadata and statistics
CREATE TABLE topic_stats (
    topic_id INT PRIMARY KEY REFERENCES topics(id) ON DELETE CASCADE,
    reply_count INT NOT NULL DEFAULT 0,
    like_count INT NOT NULL DEFAULT 0,
    bookmark_count INT NOT NULL DEFAULT 0,
    last_reply_id INT,
    last_reply_at TIMESTAMPTZ,
    last_reply_user_id UUID REFERENCES "user"(id) ON DELETE SET NULL
);

-- Posts/Replies to topics
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    topic_id INT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_solution BOOLEAN NOT NULL DEFAULT false,
    parent_post_id INT REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    edited_by UUID REFERENCES "user"(id) ON DELETE SET NULL,
    edit_count INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_posts_topic ON posts(topic_id);
CREATE INDEX idx_posts_author ON posts(author_id);
CREATE INDEX idx_posts_parent ON posts(parent_post_id);

-- Topic Tags for categorizing discussions
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    color VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Many-to-many relationship between topics and tags
CREATE TABLE topic_tags (
    topic_id INT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (topic_id, tag_id)
);

-- Likes for topics and posts
CREATE TABLE likes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    likeable_type VARCHAR(20) NOT NULL, -- 'topic' or 'post'
    likeable_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, likeable_type, likeable_id)
);

CREATE INDEX idx_likes_user ON likes(user_id);
CREATE INDEX idx_likes_likeable ON likes(likeable_type, likeable_id);

-- Bookmarks for topics
CREATE TABLE bookmarks (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    topic_id INT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, topic_id)
);

-- User follows (user following other user)
CREATE TABLE user_follows (
    follower_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id) -- User cannot follow themselves
);

-- Category follows (user following categories)
CREATE TABLE category_follows (
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, category_id)
);

-- Topic follows (user following specific topics)
CREATE TABLE topic_follows (
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    topic_id INT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, topic_id)
);

-- Notifications system
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    notification_type VARCHAR(50) NOT NULL, -- 'mention', 'reply', 'like', etc.
    actor_id UUID REFERENCES "user"(id) ON DELETE SET NULL, -- User who triggered notification
    reference_type VARCHAR(30) NOT NULL, -- 'topic', 'post', 'badge', etc.
    reference_id INT NOT NULL, -- ID of the referenced entity
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(is_read);

-- Badges system
CREATE TABLE badge_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Badge definitions
CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    type_id INT NOT NULL REFERENCES badge_types(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    requirement TEXT NOT NULL, -- Description of how to earn
    tier INT NOT NULL DEFAULT 1, -- 1=bronze, 2=silver, 3=gold, etc.
    points INT NOT NULL DEFAULT 0, -- Reputation points earned
    icon VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- User earned badges
CREATE TABLE user_badges (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    badge_id INT NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    awarded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, badge_id)
);

-- Content reporting system
CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    reporter_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    content_type VARCHAR(20) NOT NULL, -- 'topic', 'post', 'user'
    content_id INT NOT NULL,
    reason VARCHAR(100) NOT NULL,
    details TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'resolved', 'dismissed'
    reviewed_by UUID REFERENCES "user"(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- User activity tracking for algo trading specific features
CREATE TABLE activity_log (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    activity_type VARCHAR(50) NOT NULL, -- 'shared_strategy', 'published_article', etc.
    description TEXT NOT NULL,
    reference_type VARCHAR(30), -- 'strategy', 'backtest', etc.
    reference_id UUID, -- UUID of the referenced entity
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Strategy rating system (specific to algo trading)
CREATE TABLE strategy_ratings (
    id SERIAL PRIMARY KEY,
    strategy_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    review TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(strategy_id, user_id)
);

-- User expertise areas (for algo trading platform)
CREATE TABLE user_expertise (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    expertise_area VARCHAR(100) NOT NULL, -- 'machine_learning', 'mean_reversion', 'momentum', etc.
    endorsement_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, expertise_area)
);

-- Poll system for community voting
CREATE TABLE polls (
    id SERIAL PRIMARY KEY,
    topic_id INT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    is_multiple_choice BOOLEAN NOT NULL DEFAULT false,
    is_closed BOOLEAN NOT NULL DEFAULT false,
    closes_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE poll_options (
    id SERIAL PRIMARY KEY,
    poll_id INT NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    option_text VARCHAR(255) NOT NULL,
    display_order INT NOT NULL DEFAULT 0
);

CREATE TABLE poll_votes (
    id SERIAL PRIMARY KEY,
    poll_id INT NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    option_id INT NOT NULL REFERENCES poll_options(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(poll_id, user_id, option_id)
);

-- Add example categories
INSERT INTO categories (name, slug, description, display_order) VALUES
('General Discussion', 'general-discussion', 'General topics related to algo trading', 1),
('Strategy Development', 'strategy-development', 'Discuss and share algo trading strategies', 2),
('Technical Analysis', 'technical-analysis', 'Discussions about technical indicators and analysis', 3),
('Platform Support', 'platform-support', 'Get help using our platform', 4),
('Market Commentary', 'market-commentary', 'Discussions about current market conditions', 5),
('Programming Help', 'programming-help', 'Get help with coding your algorithms', 6),
('Success Stories', 'success-stories', 'Share your algo trading success stories', 7);

-- Add example badge types
INSERT INTO badge_types (name, description) VALUES
('Achievement', 'Badges earned for specific achievements'),
('Participation', 'Badges earned for participation in the community'),
('Contribution', 'Badges earned for contributing valuable content'),
('Experience', 'Badges earned for experience milestones');

-- Add example badges
INSERT INTO badges (type_id, name, description, requirement, tier, points) VALUES
(1, 'First Strategy', 'Created your first trading strategy', 'Create and save a trading strategy', 1, 5),
(1, 'Profitable Backtest', 'Achieved a profitable backtest', 'Complete a backtest with positive returns', 1, 10),
(2, 'Discussion Starter', 'Started your first discussion', 'Create a new topic in the community', 1, 5),
(2, 'Helpful Response', 'Provided a response marked as helpful', 'Have your post marked as a solution', 2, 15),
(3, 'Strategy Contributor', 'Shared a strategy with the community', 'Share a strategy publicly', 2, 20),
(3, 'Code Reviewer', 'Helped improve others code', 'Suggest improvements to another user strategy', 2, 15),
(4, '1 Year Member', 'Member for one year', 'Maintain an active account for 1 year', 1, 10),
(4, 'Trading Veteran', 'Consistent trading activity', 'Run live algorithms for at least 6 months', 3, 50);

-- Create functions for updating counts
CREATE OR REPLACE FUNCTION update_topic_stats_on_post() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE topic_stats 
        SET 
            reply_count = reply_count + 1,
            last_reply_id = NEW.id,
            last_reply_at = NOW(),
            last_reply_user_id = NEW.author_id
        WHERE topic_id = NEW.topic_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE topic_stats 
        SET reply_count = reply_count - 1
        WHERE topic_id = OLD.topic_id;
        
        -- Also update last reply info if needed
        IF OLD.id = (SELECT last_reply_id FROM topic_stats WHERE topic_id = OLD.topic_id) THEN
            WITH last_post AS (
                SELECT id, author_id, created_at
                FROM posts
                WHERE topic_id = OLD.topic_id
                ORDER BY created_at DESC
                LIMIT 1
            )
            UPDATE topic_stats
            SET 
                last_reply_id = lp.id,
                last_reply_at = lp.created_at,
                last_reply_user_id = lp.author_id
            FROM last_post lp
            WHERE topic_id = OLD.topic_id;
        END IF;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_topic_stats_trigger
AFTER INSERT OR DELETE ON posts
FOR EACH ROW EXECUTE PROCEDURE update_topic_stats_on_post();

-- Create functions for updating like counts
CREATE OR REPLACE FUNCTION update_like_counts() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.likeable_type = 'topic' THEN
        UPDATE topic_stats 
        SET like_count = like_count + 1
        WHERE topic_id = NEW.likeable_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_like_counts_trigger
AFTER INSERT ON likes
FOR EACH ROW EXECUTE PROCEDURE update_like_counts();

-- Function to initialize topic stats on topic creation
CREATE OR REPLACE FUNCTION initialize_topic_stats() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO topic_stats (topic_id, reply_count, like_count, bookmark_count)
    VALUES (NEW.id, 0, 0, 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER initialize_topic_stats_trigger
AFTER INSERT ON topics
FOR EACH ROW EXECUTE PROCEDURE initialize_topic_stats();

COMMIT;