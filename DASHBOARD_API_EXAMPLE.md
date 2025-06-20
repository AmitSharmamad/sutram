# Dashboard Statistics API - cURL Examples

## Prerequisites

1. **Start the Server**
```bash
cd /Users/amit/workspace/sutram
./bin/web
```

2. **Database Setup** (if not already done)
```bash
# Make sure your database is running and migrated
cd db
make up
make migrate
```

## Authentication Required

The dashboard endpoint requires authentication. You need to either:

### Option 1: Register a New User
```bash
curl -X POST http://localhost:9090/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "phone_number": "+1234567890",
    "password": "securepassword123"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-06-21T10:30:00Z",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "John Doe",
      "username": "johndoe",
      "email": "john@example.com",
      "phone_number": "+1234567890",
      "verified": false,
      "is_active": true,
      "roles": ["user"]
    }
  }
}
```

### Option 2: Login with Existing User
```bash
curl -X POST http://localhost:9090/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

**Expected Response:** (Same as registration)

## Dashboard Statistics API Call

### Using Authorization Header
```bash
# Replace YOUR_JWT_TOKEN with the actual token from login/register response
curl -X GET http://localhost:9090/api/user/dashboard/statistics \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

### Using Cookie (if server sets cookies)
```bash
curl -X GET http://localhost:9090/api/user/dashboard/statistics \
  -H "Content-Type: application/json" \
  --cookie "token=YOUR_JWT_TOKEN"
```

## Expected Dashboard Response

```json
{
  "success": true,
  "data": {
    "generated_at": "2025-06-20T10:30:45.123Z",
    "portfolio": {
      "total_value": 125000.50,
      "total_pnl": 25000.50,
      "todays_pnl": 1250.75,
      "available_balance": 50000.00,
      "margin_used": 75000.50,
      "margin_available": 50000.00
    },
    "trading_performance": {
      "total_trades": 150,
      "winning_trades": 95,
      "losing_trades": 55,
      "win_rate": 63.33,
      "total_profit": 15000.00,
      "total_loss": -8000.00,
      "average_profit": 157.89,
      "average_loss": -145.45,
      "max_profit": 2500.00,
      "max_loss": -1800.00,
      "sharpe_ratio": 1.85,
      "volatility": 12.5
    },
    "live_trading": {
      "active_deployments": 3,
      "total_deployments": 15,
      "running_strategies": 8,
      "paused_strategies": 2,
      "positions": [
        {
          "id": "123e4567-e89b-12d3-a456-426614174001",
          "instrument_id": 1234,
          "instrument_name": "RELIANCE",
          "quantity": 100,
          "average_price": 2450.50,
          "current_price": 2475.25,
          "pnl": 2475.00,
          "direction": "LONG",
          "created_at": "2025-06-20T09:15:30.456Z"
        },
        {
          "id": "123e4567-e89b-12d3-a456-426614174002",
          "instrument_id": 5678,
          "instrument_name": "TCS",
          "quantity": -50,
          "average_price": 3200.00,
          "current_price": 3180.75,
          "pnl": 962.50,
          "direction": "SHORT",
          "created_at": "2025-06-20T08:45:15.789Z"
        }
      ]
    },
    "strategy_development": {
      "total_strategies": 25,
      "active_strategies": 18,
      "draft_strategies": 7,
      "published_strategies": 12,
      "strategies_in_marketplace": 5
    },
    "marketplace_activity": {
      "listings_count": 5,
      "total_sales": 12,
      "total_revenue": 15000.00,
      "subscribers_count": 45,
      "average_rating": 4.2
    },
    "community_engagement": {
      "topics_created": 8,
      "posts_count": 25,
      "likes_received": 120,
      "reputation_score": 750,
      "badges_earned": 3
    },
    "analytics": {
      "chart_data": {
        "daily_pnl": [
          {
            "date": "2025-06-18",
            "pnl": 980.25
          },
          {
            "date": "2025-06-19",
            "pnl": 1200.50
          },
          {
            "date": "2025-06-20",
            "pnl": 1250.75
          }
        ],
        "strategy_performance": [
          {
            "strategy_id": "123e4567-e89b-12d3-a456-426614174010",
            "name": "RSI Momentum Strategy",
            "pnl": 5000.00
          },
          {
            "strategy_id": "123e4567-e89b-12d3-a456-426614174011",
            "name": "Moving Average Crossover",
            "pnl": 3500.00
          },
          {
            "strategy_id": "123e4567-e89b-12d3-a456-426614174012",
            "name": "Bollinger Bands Mean Reversion",
            "pnl": 2800.00
          }
        ]
      }
    },
    "recent_activity": {
      "recent_trades": [
        {
          "id": "123e4567-e89b-12d3-a456-426614174020",
          "instrument_name": "RELIANCE",
          "quantity": 50,
          "price": 2450.50,
          "direction": "BUY",
          "timestamp": "2025-06-20T09:30:15.123Z"
        },
        {
          "id": "123e4567-e89b-12d3-a456-426614174021",
          "instrument_name": "TCS",
          "quantity": 25,
          "price": 3200.00,
          "direction": "SELL",
          "timestamp": "2025-06-20T08:45:30.456Z"
        }
      ],
      "recent_deployments": [
        {
          "id": "123e4567-e89b-12d3-a456-426614174030",
          "strategy_name": "RSI Momentum Strategy",
          "status": "RUNNING",
          "created_at": "2025-06-20T08:00:00.000Z"
        },
        {
          "id": "123e4567-e89b-12d3-a456-426614174031",
          "strategy_name": "Bollinger Bands Strategy",
          "status": "PAUSED",
          "created_at": "2025-06-20T07:30:00.000Z"
        }
      ]
    }
  }
}
```

## Error Responses

### No Authentication Token
```bash
curl -X GET http://localhost:9090/api/user/dashboard/statistics \
  -H "Content-Type: application/json"
```

**Response:**
```json
{
  "success": false,
  "error": "No token provided"
}
```

### Invalid Token
```bash
curl -X GET http://localhost:9090/api/user/dashboard/statistics \
  -H "Authorization: Bearer invalid_token_here" \
  -H "Content-Type: application/json"
```

**Response:**
```json
{
  "success": false,
  "error": "Invalid token: token is malformed"
}
```

### Server Error
```json
{
  "success": false,
  "error": "Failed to fetch dashboard statistics: database connection error"
}
```

## Complete Test Sequence

Here's a complete sequence to test the API:

```bash
#!/bin/bash

# 1. Start the server (in another terminal)
# ./bin/web

# 2. Register a new user and capture the token
echo "1. Registering new user..."
RESPONSE=$(curl -s -X POST http://localhost:9090/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "username": "testuser",
    "email": "test@example.com",
    "phone_number": "+1234567890",
    "password": "testpassword123"
  }')

echo "Registration response: $RESPONSE"

# 3. Extract token from response (requires jq)
TOKEN=$(echo $RESPONSE | jq -r '.data.token')
echo "Extracted token: $TOKEN"

# 4. Call dashboard statistics API
echo "2. Fetching dashboard statistics..."
curl -X GET http://localhost:9090/api/user/dashboard/statistics \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq .

echo "Dashboard API test completed!"
```

## Notes

- **New User Data**: For a newly registered user, most statistics will be zero or empty since they haven't traded yet
- **Real Data**: To see meaningful statistics, you would need to:
  - Create strategies
  - Deploy strategies
  - Have some trading activity
  - Participate in marketplace and community
- **Development**: The current implementation may return placeholder/calculated values for demonstration purposes
- **Database State**: Ensure your database is properly set up and contains some test data if you want to see non-zero values

## Frontend Integration Example

```javascript
// JavaScript/React example
const fetchDashboardData = async () => {
  try {
    const token = localStorage.getItem('authToken');
    const response = await fetch('/api/user/dashboard/statistics', {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    const data = await response.json();
    
    if (data.success) {
      // Update your dashboard UI
      setPortfolioData(data.data.portfolio);
      setTradingPerformance(data.data.trading_performance);
      setLiveTrading(data.data.live_trading);
      // ... etc
    } else {
      console.error('Dashboard API error:', data.error);
    }
  } catch (error) {
    console.error('Network error:', error);
  }
};
```
