-- Marketplace Schema for Algorithmic Trading Platform

-- Categories for organizing strategies
CREATE TABLE IF NOT EXISTS marketplace_category (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id)
);

-- Strategy Marketplace Table
CREATE TABLE IF NOT EXISTS marketplace_listing (
    id UUID DEFAULT gen_random_uuid(),
    strategy_id UUID NOT NULL,
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(18, 2) NOT NULL,
    subscription_type VARCHAR(50) NOT NULL, -- 'one-time', 'monthly', 'yearly'
    is_featured BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    views INT DEFAULT 0,
    ratings DECIMAL(3, 2),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (strategy_id) REFERENCES strategy(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);

-- Strategy Performance Metrics (for marketplace display)
CREATE TABLE IF NOT EXISTS marketplace_performance (
    listing_id UUID NOT NULL,
    sharpe_ratio DECIMAL(10, 4),
    sortino_ratio DECIMAL(10, 4),
    total_return_pct DECIMAL(10, 4),
    annualized_return DECIMAL(10, 4),
    max_drawdown DECIMAL(10, 4),
    win_rate DECIMAL(10, 4),
    profit_factor DECIMAL(10, 4),
    backtest_period_days INT,
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (listing_id),
    FOREIGN KEY (listing_id) REFERENCES marketplace_listing(id) ON DELETE CASCADE
);

-- Strategy Subscriptions
CREATE TABLE IF NOT EXISTS marketplace_subscription (
    id UUID DEFAULT gen_random_uuid(),
    listing_id UUID NOT NULL,
    subscriber_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL, -- 'active', 'expired', 'cancelled'
    subscription_type VARCHAR(50) NOT NULL, -- 'one-time', 'monthly', 'yearly'
    price_paid DECIMAL(18, 2) NOT NULL,
    start_date TIMESTAMP NOT NULL DEFAULT now(),
    end_date TIMESTAMP,
    renewal_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (listing_id) REFERENCES marketplace_listing(id),
    FOREIGN KEY (subscriber_id) REFERENCES "user"(id) ON DELETE CASCADE
);

-- Reviews and Ratings
CREATE TABLE IF NOT EXISTS marketplace_review (
    id UUID DEFAULT gen_random_uuid(),
    listing_id UUID NOT NULL,
    user_id UUID NOT NULL,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (listing_id) REFERENCES marketplace_listing(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    UNIQUE (listing_id, user_id) -- One review per user per listing
);

-- Transaction History
CREATE TABLE IF NOT EXISTS marketplace_transaction (
    id UUID DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL,
    amount DECIMAL(18, 2) NOT NULL,
    status VARCHAR(50) NOT NULL, -- 'completed', 'pending', 'failed', 'refunded'
    payment_method VARCHAR(50),
    payment_id VARCHAR(255), -- External payment reference ID
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (subscription_id) REFERENCES marketplace_subscription(id) ON DELETE CASCADE
);
