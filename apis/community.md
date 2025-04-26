# Marketplace API Endpoints

This document outlines the API endpoints required for the marketplace feature of the algo trading platform.

## Authentication

All endpoints require authentication via Bearer token unless specified otherwise.

```
Authorization: Bearer <token>
```

## Strategy Marketplace Endpoints

### 1. List Marketplace Strategies

Returns a list of all strategies available in the marketplace.

```bash
curl -X GET \
  'https://api.algotrading.com/v1/marketplace/strategies?page=1&limit=10' \
  -H 'Authorization: Bearer <token>'
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Number of items per page (default: 10)
- `sort`: Sort field (options: popularity, newest, price)
- `category`: Filter by category
- `author`: Filter by author ID

**Response:**
```json
{
  "total": 120,
  "page": 1,
  "limit": 10,
  "strategies": [
    {
      "id": "strat-123",
      "name": "Golden Cross Strategy",
      "description": "A trend-following strategy based on moving average crossovers",
      "author": {
        "id": "user-456",
        "name": "John Doe"
      },
      "price": 99.99,
      "rating": 4.5,
      "reviews": 27,
      "category": "Technical Analysis",
      "created_at": "2025-03-15T08:23:45Z",
      "updated_at": "2025-04-10T14:12:30Z",
      "preview_image": "https://storage.algotrading.com/strategy-previews/strat-123.png"
    },
    // ... other strategies
  ]
}
```

### 2. Get Strategy Details

Returns detailed information about a specific marketplace strategy.

```bash
curl -X GET \
  'https://api.algotrading.com/v1/marketplace/strategies/{strategy_id}' \
  -H 'Authorization: Bearer <token>'
```

**Response:**
```json
{
  "id": "strat-123",
  "name": "Golden Cross Strategy",
  "description": "A trend-following strategy based on moving average crossovers",
  "detailed_description": "This strategy uses the crossover of the 50-day and 200-day moving averages to identify long-term trend changes in the market...",
  "author": {
    "id": "user-456",
    "name": "John Doe",
    "profile_image": "https://storage.algotrading.com/profiles/user-456.jpg",
    "bio": "Professional trader with 10 years of experience",
    "total_followers": 1250,
    "total_strategies": 8
  },
  "price": 99.99,
  "rating": 4.5,
  "reviews": [
    {
      "id": "rev-789",
      "user": {
        "id": "user-321",
        "name": "Jane Smith"
      },
      "rating": 5,
      "comment": "Great strategy, consistent returns!",
      "created_at": "2025-04-01T09:45:22Z"
    },
    // ... other reviews
  ],
  "performance": {
    "annual_return": 22.5,
    "max_drawdown": 8.2,
    "sharpe_ratio": 1.8,
    "win_rate": 62.3,
    "backtest_period": {
      "start": "2023-01-01",
      "end": "2025-01-01"
    }
  },
  "category": "Technical Analysis",
  "tags": ["Moving Average", "Crossover", "Trend Following"],
  "requirements": {
    "min_capital": 5000,
    "supported_instruments": ["stocks", "etfs"],
    "timeframes": ["1d", "4h"]
  },
  "preview_images": [
    "https://storage.algotrading.com/strategy-previews/strat-123-1.png",
    "https://storage.algotrading.com/strategy-previews/strat-123-2.png"
  ],
  "created_at": "2025-03-15T08:23:45Z",
  "updated_at": "2025-04-10T14:12:30Z"
}
```

### 3. Purchase Strategy

Purchases a strategy from the marketplace.

```bash
curl -X POST \
  'https://api.algotrading.com/v1/marketplace/strategies/{strategy_id}/purchase' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "payment_method_id": "pm-876"
  }'
```

**Response:**
```json
{
  "id": "purchase-567",
  "strategy_id": "strat-123",
  "user_id": "user-321",
  "amount": 99.99,
  "status": "completed",
  "purchase_date": "2025-04-25T10:30:15Z",
  "receipt_url": "https://api.algotrading.com/receipts/purchase-567"
}
```

### 4. List User Purchased Strategies

Returns all strategies purchased by the authenticated user.

```bash
curl -X GET \
  'https://api.algotrading.com/v1/marketplace/purchases' \
  -H 'Authorization: Bearer <token>'
```

**Response:**
```json
{
  "purchases": [
    {
      "purchase_id": "purchase-567",
      "purchase_date": "2025-04-25T10:30:15Z",
      "strategy": {
        "id": "strat-123",
        "name": "Golden Cross Strategy",
        "description": "A trend-following strategy based on moving average crossovers",
        "author": {
          "id": "user-456",
          "name": "John Doe"
        },
        "price": 99.99
      }
    },
    // ... other purchases
  ]
}
```

### 5. Publish Strategy to Marketplace

Publishes a user's strategy to the marketplace.

```bash
curl -X POST \
  'https://api.algotrading.com/v1/marketplace/strategies' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "strategy_id": "user-strat-789",
    "price": 129.99,
    "detailed_description": "Comprehensive description of how the strategy works...",
    "category": "Mean Reversion",
    "tags": ["RSI", "Bollinger Bands", "Oversold"],
    "preview_images": [
      "data:image/png;base64,..."
    ]
  }'
```

**Response:**
```json
{
  "id": "strat-456",
  "name": "RSI Reversal Strategy",
  "status": "pending_review",
  "submission_date": "2025-04-25T11:45:30Z",
  "estimated_review_completion": "2025-04-27T11:45:30Z"
}
```

### 6. Update Marketplace Strategy

Updates the details of a strategy published by the authenticated user.

```bash
curl -X PUT \
  'https://api.algotrading.com/v1/marketplace/strategies/{strategy_id}' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "price": 149.99,
    "detailed_description": "Updated description with new information...",
    "category": "Mean Reversion",
    "tags": ["RSI", "Bollinger Bands", "Oversold", "Volatility"]
  }'
```

**Response:**
```json
{
  "id": "strat-456",
  "name": "RSI Reversal Strategy",
  "status": "pending_review",
  "submission_date": "2025-04-25T14:22:10Z",
  "estimated_review_completion": "2025-04-27T14:22:10Z"
}
```

### 7. Remove Strategy from Marketplace

Removes a strategy from the marketplace.

```bash
curl -X DELETE \
  'https://api.algotrading.com/v1/marketplace/strategies/{strategy_id}' \
  -H 'Authorization: Bearer <token>'
```

**Response:**
```json
{
  "id": "strat-456",
  "status": "removed",
  "removal_date": "2025-04-25T15:30:45Z"
}
```

### 8. Submit Strategy Review

Submits a review for a purchased strategy.

```bash
curl -X POST \
  'https://api.algotrading.com/v1/marketplace/strategies/{strategy_id}/reviews' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "rating": 4,
    "comment": "Great strategy with solid performance. Documentation could be more detailed."
  }'
```

**Response:**
```json
{
  "id": "rev-123",
  "strategy_id": "strat-123",
  "user": {
    "id": "user-321",
    "name": "Jane Smith"
  },
  "rating": 4,
  "comment": "Great strategy with solid performance. Documentation could be more detailed.",
  "created_at": "2025-04-25T16:15:30Z"
}
```

### 9. Get Marketplace Categories

Returns a list of available marketplace categories.

```bash
curl -X GET \
  'https://api.algotrading.com/v1/marketplace/categories' \
  -H 'Authorization: Bearer <token>'
```

**Response:**
```json
{
  "categories": [
    {
      "id": "cat-1",
      "name": "Technical Analysis",
      "description": "Strategies based on price action and technical indicators",
      "strategy_count": 45
    },
    {
      "id": "cat-2",
      "name": "Mean Reversion",
      "description": "Strategies that capitalize on price returning to the mean",
      "strategy_count": 32
    },
    {
      "id": "cat-3",
      "name": "Momentum",
      "description": "Strategies that follow market trends and momentum",
      "strategy_count": 28
    },
    {
      "id": "cat-4",
      "name": "Arbitrage",
      "description": "Strategies that exploit price differences between markets",
      "strategy_count": 15
    },
    {
      "id": "cat-5",
      "name": "Options",
      "description": "Strategies focused on options trading",
      "strategy_count": 22
    }
  ]
}
```

### 10. Get Author Profile

Returns detailed information about a strategy author.

```bash
curl -X GET \
  'https://api.algotrading.com/v1/marketplace/authors/{author_id}' \
  -H 'Authorization: Bearer <token>'
```

**Response:**
```json
{
  "id": "user-456",
  "name": "John Doe",
  "profile_image": "https://storage.algotrading.com/profiles/user-456.jpg",
  "bio": "Professional trader with 10 years of experience specializing in technical analysis and mean reversion strategies.",
  "joined_date": "2023-05-12T09:30:15Z",
  "total_followers": 1250,
  "total_strategies": 8,
  "average_rating": 4.7,
  "total_sales": 3450,
  "strategies": [
    {
      "id": "strat-123",
      "name": "Golden Cross Strategy",
      "description": "A trend-following strategy based on moving average crossovers",
      "price": 99.99,
      "rating": 4.5,
      "reviews": 27,
      "purchases": 320
    },
    // ... other strategies
  ]
}
```

## Marketplace Analytics Endpoints

### 1. Get User Marketplace Analytics

Returns analytics data for strategies published by the authenticated user.

```bash
curl -X GET \
  'https://api.algotrading.com/v1/marketplace/analytics' \
  -H 'Authorization: Bearer <token>'
```

**Response:**
```json
{
  "total_revenue": 12450.75,
  "total_sales": 145,
  "average_rating": 4.3,
  "total_reviews": 87,
  "strategy_analytics": [
    {
      "id": "strat-456",
      "name": "RSI Reversal Strategy",
      "revenue": 8950.25,
      "sales": 112,
      "average_rating": 4.5,
      "views": 2340,
      "conversion_rate": 4.8
    },
    // ... other strategies
  ],
  "monthly_trends": [
    {
      "month": "2025-01",
      "revenue": 1250.50,
      "sales": 15
    },
    {
      "month": "2025-02",
      "revenue": 1890.25,
      "sales": 22
    },
    {
      "month": "2025-03",
      "revenue": 2150.75,
      "sales": 25
    },
    {
      "month": "2025-04",
      "revenue": 1760.25,
      "sales": 19
    }
  ]
}
```

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
