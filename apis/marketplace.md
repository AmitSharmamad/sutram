# Marketplace API Documentation

This document outlines all the API endpoints required for the marketplace feature in the Algo Trading Platform.

## Marketplace Listings APIs

### 1. Get All Marketplace Listings

Retrieves a paginated list of all available strategy listings in the marketplace.

```bash
curl -X GET "http://localhost:9090/api/marketplace/listings" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 2. Get Filtered Marketplace Listings

Retrieves a filtered list of marketplace listings based on category, sorting preference, and search terms.

```bash
curl -X GET "http://localhost:9090/api/marketplace/listings?category=trend-following&sort=rating&search=momentum" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 3. Get Specific Marketplace Listing

Retrieves detailed information about a specific marketplace listing.

```bash
curl -X GET "http://localhost:9090/api/marketplace/listings/{listing_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 4. Create New Marketplace Listing

Creates a new marketplace listing for a strategy.

```bash
curl -X POST "http://localhost:9090/api/marketplace/listings" \
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
curl -X PUT "http://localhost:9090/api/marketplace/listings/{listing_id}" \
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
curl -X DELETE "http://localhost:9090/api/marketplace/listings/{listing_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## Categories APIs

### 7. Get All Categories

Retrieves all available strategy categories for the marketplace.

```bash
curl -X GET "http://localhost:9090/api/marketplace/categories" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## Performance Metrics APIs

### 8. Update Strategy Performance Metrics

Updates performance metrics for a specific marketplace listing.

```bash
curl -X PUT "http://localhost:9090/api/marketplace/listings/{listing_id}/performance" \
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
curl -X POST "http://localhost:9090/api/marketplace/subscriptions" \
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
curl -X GET "http://localhost:9090/api/marketplace/subscriptions/user" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 11. Cancel Subscription

Cancels an active subscription.

```bash
curl -X PUT "http://localhost:9090/api/marketplace/subscriptions/{subscription_id}" \
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
curl -X GET "http://localhost:9090/api/marketplace/listings/{listing_id}/reviews" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 13. Submit a Review

Submits a new review for a marketplace listing.

```bash
curl -X POST "http://localhost:9090/api/marketplace/listings/{listing_id}/reviews" \
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
curl -X PUT "http://localhost:9090/api/marketplace/reviews/{review_id}" \
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
curl -X DELETE "http://localhost:9090/api/marketplace/reviews/{review_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## Transaction APIs

### 16. Get Transaction History

Retrieves the transaction history for the authenticated user.

```bash
curl -X GET "http://localhost:9090/api/marketplace/transactions" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 17. Get Transaction Details

Retrieves details for a specific transaction.

```bash
curl -X GET "http://localhost:9090/api/marketplace/transactions/{transaction_id}" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## Seller APIs

### 18. Get Seller Earnings/Analytics

Retrieves earnings and analytics data for a seller.

```bash
curl -X GET "http://localhost:9090/api/marketplace/seller/analytics" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### 19. Get Seller Listings

Retrieves all listings created by the authenticated seller.

```bash
curl -X GET "http://localhost:9090/api/marketplace/seller/listings" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```