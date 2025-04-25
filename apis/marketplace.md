# Marketplace API Documentation

This document outlines all the API endpoints required for the marketplace feature in the Algo Trading Platform.

## Marketplace Listings APIs

### 1. Get All Marketplace Listings

Retrieves a paginated list of all available strategy listings in the marketplace.

```bash
curl -X GET "/marketplace/listings" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 2. Get Filtered Marketplace Listings

Retrieves a filtered list of marketplace listings based on category, sorting preference, and search terms.

```bash
curl -X GET "/marketplace/listings?category=trend-following&sort=rating&search=momentum" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 3. Get Specific Marketplace Listing

Retrieves detailed information about a specific marketplace listing.

```bash
curl -X GET "/marketplace/listings/{listing_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 4. Create New Marketplace Listing

Creates a new marketplace listing for a strategy.

```bash
curl -X POST "/marketplace/listings" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "strategy_id": "46bd1258-0bed-4be8-8c8c-8d1c6ee5eab1",
    "title": "Advanced Momentum Strategy",
    "description": "A strategy that uses multiple momentum indicators to identify strong trends",
    "price": 49.99,
    "subscription_type": "monthly",
    "category_ids": ["3", "8"],
    "features": [
      "Advanced momentum detection algorithm",
      "Multi-timeframe analysis", 
      "Automatic position sizing",
      "Risk management rules",
      "Email and SMS alerts"
    ]
  }'
```

### 5. Update Marketplace Listing

Updates an existing marketplace listing.

```bash
curl -X PUT "/marketplace/listings/{listing_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Momentum Strategy",
    "description": "Updated description with improved features",
    "price": 59.99,
    "subscription_type": "yearly",
    "category_ids": ["3", "8"],
    "features": [
      "Advanced momentum detection algorithm",
      "Multi-timeframe analysis",
      "Automatic position sizing",
      "Risk management rules",
      "Email and SMS alerts",
      "New mobile notifications"
    ],
    "is_active": true
  }'
```

### 6. Delete Marketplace Listing

Removes a marketplace listing from the platform.

```bash
curl -X DELETE "/marketplace/listings/{listing_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## Categories APIs

### 7. Get All Categories

Retrieves all available strategy categories for the marketplace.

```bash
curl -X GET "/marketplace/categories" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## Performance Metrics APIs

### 8. Update Strategy Performance Metrics

Updates performance metrics for a specific marketplace listing.

```bash
curl -X PUT "/marketplace/listings/{listing_id}/performance" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "sharpe_ratio": 1.8,
    "sortino_ratio": 2.4,
    "total_return_pct": 32.5,
    "annualized_return": 24.7,
    "max_drawdown": 12.3,
    "win_rate": 68.7,
    "profit_factor": 2.1,
    "backtest_period_days": 365
  }'
```

## Subscription APIs

### 9. Purchase/Subscribe to Strategy

Creates a new subscription for the user to access a strategy.

```bash
curl -X POST "/marketplace/subscriptions" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "listing_id": "1",
    "payment_method": "credit_card",
    "payment_details": {
      "card_token": "tok_visa"
    }
  }'
```

### 10. Get User's Subscriptions

Retrieves all subscriptions for the authenticated user.

```bash
curl -X GET "/marketplace/subscriptions/user" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 11. Cancel Subscription

Cancels an active subscription.

```bash
curl -X PUT "/marketplace/subscriptions/{subscription_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "cancelled",
    "cancellation_reason": "Found a better strategy"
  }'
```

## Review APIs

### 12. Get Reviews for a Listing

Retrieves all reviews for a specific marketplace listing.

```bash
curl -X GET "/marketplace/listings/{listing_id}/reviews" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 13. Submit a Review

Submits a new review for a marketplace listing.

```bash
curl -X POST "/marketplace/listings/{listing_id}/reviews" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rating": 5,
    "comment": "Excellent strategy, consistently profitable in trending markets. The alerts are very timely."
  }'
```

### 14. Update a Review

Updates an existing review submitted by the user.

```bash
curl -X PUT "/marketplace/reviews/{review_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rating": 4,
    "comment": "Updated review: Good performance overall, but sometimes gives false signals in choppy markets."
  }'
```

### 15. Delete a Review

Removes a review submitted by the user.

```bash
curl -X DELETE "/marketplace/reviews/{review_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## Transaction APIs

### 16. Get Transaction History

Retrieves the transaction history for the authenticated user.

```bash
curl -X GET "/marketplace/transactions" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 17. Get Transaction Details

Retrieves details for a specific transaction.

```bash
curl -X GET "/marketplace/transactions/{transaction_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## Seller APIs

### 18. Get Seller Earnings/Analytics

Retrieves earnings and analytics data for a seller.

```bash
curl -X GET "/marketplace/seller/analytics" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 19. Get Seller Listings

Retrieves all listings created by the authenticated seller.

```bash
curl -X GET "/marketplace/seller/listings" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

-- Marketplace Database Schema for Algorithmic Trading Platform

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

-- Many-to-many relationship between listings and categories
CREATE TABLE IF NOT EXISTS marketplace_listing_category (
    listing_id UUID NOT NULL,
    category_id UUID NOT NULL,
    
    PRIMARY KEY (listing_id, category_id),
    FOREIGN KEY (listing_id) REFERENCES marketplace_listing(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES marketplace_category(id) ON DELETE CASCADE
);

-- Feature Toggles for the marketplace
CREATE TABLE IF NOT EXISTS marketplace_settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT now(),
    updated_by UUID REFERENCES "user"(id)
);

-- Indexes for performance
CREATE INDEX idx_marketplace_listing_user ON marketplace_listing(user_id);
CREATE INDEX idx_marketplace_listing_active ON marketplace_listing(is_active);
CREATE INDEX idx_marketplace_subscription_subscriber ON marketplace_subscription(subscriber_id);
CREATE INDEX idx_marketplace_subscription_status ON marketplace_subscription(status);
CREATE INDEX idx_marketplace_review_listing ON marketplace_review(listing_id);

-- Insert default categories
INSERT INTO marketplace_category (name, description)
VALUES 
('Trend Following', 'Strategies that follow market trends'),
('Mean Reversion', 'Strategies that capitalize on price reversals'),
('Momentum', 'Strategies based on momentum indicators'),
('Arbitrage', 'Strategies that exploit price differences'),
('Breakout', 'Strategies that trade on price breakouts'),
('Volatility', 'Strategies that trade based on market volatility'),
('Options', 'Strategies focused on options trading'),
('Multi-Factor', 'Strategies using multiple factors for decisions');