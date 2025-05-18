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