# Dashboard Statistics API

## Overview
A comprehensive API endpoint that returns all relevant user dashboard statistics for the Sutram algorithmic trading platform in a single response.

## Endpoint
```
GET /api/user/dashboard/statistics
```

## Authentication
- Requires valid JWT token in Authorization header or cookie
- Returns 401 Unauthorized if token is missing or invalid

## Response Structure

```json
{
  "success": true,
  "data": {
    "generated_at": "2025-06-20T10:30:00Z",
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
          "id": "uuid",
          "instrument_id": 1234,
          "instrument_name": "RELIANCE",
          "quantity": 100,
          "average_price": 2450.50,
          "current_price": 2475.25,
          "pnl": 2475.00,
          "direction": "LONG",
          "created_at": "2025-06-20T09:15:00Z"
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
          {"date": "2025-06-19", "pnl": 1200.50},
          {"date": "2025-06-20", "pnl": 1250.75}
        ],
        "strategy_performance": [
          {"strategy_id": "uuid1", "name": "RSI Strategy", "pnl": 5000.00},
          {"strategy_id": "uuid2", "name": "Moving Average", "pnl": 3500.00}
        ]
      }
    },
    "recent_activity": {
      "recent_trades": [
        {
          "id": "uuid",
          "instrument_name": "RELIANCE",
          "quantity": 50,
          "price": 2450.50,
          "direction": "BUY",
          "timestamp": "2025-06-20T09:30:00Z"
        }
      ],
      "recent_deployments": [
        {
          "id": "uuid",
          "strategy_name": "RSI Strategy",
          "status": "RUNNING",
          "created_at": "2025-06-20T08:00:00Z"
        }
      ]
    }
  }
}
```

## Error Response
```json
{
  "success": false,
  "error": "Error message"
}
```

## Implementation Details

### Architecture
- **Dashboard Service**: Aggregates data from multiple modules
- **Module Integration**: Uses services from deployment, strategy, backtest, marketplace, and community modules
- **Authentication**: Integrates with existing JWT authentication system

### Data Sources
1. **Portfolio**: Aggregated from deployment positions and account data
2. **Trading Performance**: Calculated from historical trades and positions
3. **Live Trading**: Real-time data from active deployments
4. **Strategy Development**: User's strategy creation and management data
5. **Marketplace Activity**: User's selling activity and performance metrics
6. **Community Engagement**: User's forum participation and reputation
7. **Analytics**: Chart data for visualization
8. **Recent Activity**: Latest trades and deployments

### Performance Considerations
- Single API call reduces client-server round trips
- Efficient data aggregation across modules
- Structured response for easy frontend consumption
- Authentication handled once per request

## Usage Example

```javascript
// Frontend usage
const response = await fetch('/api/user/dashboard/statistics', {
  headers: {
    'Authorization': 'Bearer ' + token,
    'Content-Type': 'application/json'
  }
});

const dashboardData = await response.json();
if (dashboardData.success) {
  // Use dashboardData.data to populate dashboard
  updatePortfolioSection(dashboardData.data.portfolio);
  updateTradingPerformance(dashboardData.data.trading_performance);
  // ... etc
}
```

## Files Modified/Created

1. **Created**: `internal/modules/user/application/dto/dashboard_dto.go` - Response DTOs
2. **Created**: `internal/modules/user/application/service/dashboard_service.go` - Business logic
3. **Modified**: `internal/modules/user/interfaces/http/user_controller.go` - API endpoint
4. **Modified**: `internal/modules/user/interfaces/http/router.go` - Route registration
5. **Modified**: `internal/modules/user/module.go` - Module integration
6. **Modified**: `cmd/web/main.go` - Service initialization
7. **Modified**: Multiple module files to expose services

The dashboard statistics API is now fully functional and provides comprehensive user data in a single, efficient response.
