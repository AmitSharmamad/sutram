<!-- # Strategy API Documentation -->

This document provides cURL examples for all strategy API endpoints.

## Base URL
```
/api/strategy
```

## Authentication
All endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### 1. Create Strategy

Create a new trading strategy.

**Endpoint:** `POST /strategy`

**cURL Example:**
```bash
curl -X POST "/api/strategy" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "RSI Momentum Strategy",
    "description": "A momentum strategy based on RSI indicator",
    "position_type": "LONG",
    "instruments": [1, 2, 3],
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "5m",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "MIS",
    "qty": 100,
    "entry_expression": "RSI(14) < 30",
    "exit_expression": "RSI(14) > 70",
    "target_profit": 5.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 2.0,
    "stop_loss_type": "PERCENTAGE",
    "trailing_stop_loss": 1.0,
    "trailing_stop_loss_type": "PERCENTAGE",
    "capital": 100000.0,
    "max_allocation": 50000.0,
    "max_trades": 5,
    "brokerage": 20,
    "slippage": 0.1,
    "entry_start_time": "09:15:00",
    "entry_end_time": "15:00:00",
    "start_date": "2025-01-01T00:00:00Z",
    "end_date": "2025-12-31T00:00:00Z",
    "segment": "F&O",
    "strategy_type": "CONDITIONAL"
  }'
```

**Time-Based Strategy Example:**
```bash
curl -X POST "/api/strategy" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Morning Buy Strategy",
    "description": "Buy at market open, sell at specific time",
    "position_type": "LONG",
    "instruments": [1],
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "1d",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "MIS",
    "qty": 50,
    "target_profit": 3.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 1.5,
    "stop_loss_type": "PERCENTAGE",
    "capital": 50000.0,
    "max_allocation": 25000.0,
    "segment": "EQ",
    "strategy_type": "TIME_BASED",
    "entry_time": "09:15:00",
    "exit_time": "15:00:00"
  }'
```

### 2. Get All User Strategies

Retrieve all strategies for the authenticated user.

**Endpoint:** `GET /strategy`

**cURL Example:**
```bash
curl -X GET "/api/strategy" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. Get Strategy by ID

Retrieve a specific strategy by its ID.

**Endpoint:** `GET /strategy/{id}`

**cURL Example:**
```bash
curl -X GET "/api/strategy/123e4567-e89b-12d3-a456-426614174000" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Update Strategy

Update an existing strategy.

**Endpoint:** `PUT /strategy/{id}`

**cURL Example:**
```bash
curl -X PUT "/api/strategy/123e4567-e89b-12d3-a456-426614174000" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Updated RSI Momentum Strategy",
    "description": "Updated momentum strategy with new parameters",
    "position_type": "LONG",
    "instruments": [1, 2, 3, 4],
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "15m",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "MIS",
    "qty": 150,
    "entry_expression": "RSI(14) < 25",
    "exit_expression": "RSI(14) > 75",
    "target_profit": 7.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 3.0,
    "stop_loss_type": "PERCENTAGE",
    "trailing_stop_loss": 1.5,
    "trailing_stop_loss_type": "PERCENTAGE",
    "capital": 150000.0,
    "max_allocation": 75000.0,
    "max_trades": 10,
    "brokerage": 20,
    "slippage": 0.1,
    "entry_start_time": "09:15:00",
    "entry_end_time": "14:30:00",
    "start_date": "2025-01-01T00:00:00Z",
    "end_date": "2025-12-31T00:00:00Z",
    "segment": "F&O",
    "strategy_type": "CONDITIONAL"
  }'
```

### 5. Delete Strategy

Delete a strategy by its ID.

**Endpoint:** `DELETE /strategy/{id}`

**cURL Example:**
```bash
curl -X DELETE "/api/strategy/123e4567-e89b-12d3-a456-426614174000" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6. Activate Strategy

Activate a strategy by its ID.

**Endpoint:** `POST /strategy/{id}/activate`

**cURL Example:**
```bash
curl -X POST "/api/strategy/123e4567-e89b-12d3-a456-426614174000/activate" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 7. Deactivate Strategy

Deactivate a strategy by its ID.

**Endpoint:** `POST /strategy/{id}/deactivate`

**cURL Example:**
```bash
curl -X POST "/api/strategy/123e4567-e89b-12d3-a456-426614174000/deactivate" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 8. Clone Strategy

Create a copy of an existing strategy.

**Endpoint:** `POST /strategy/{id}/clone`

**cURL Example:**
```bash
curl -X POST "/api/strategy/123e4567-e89b-12d3-a456-426614174000/clone" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 9. Search Strategies

Search strategies with filter criteria.

**Endpoint:** `POST /strategy/search`

**cURL Example:**
```bash
curl -X POST "/api/strategy/search" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "RSI",
    "position_type": "LONG",
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "5m",
    "segment": "F&O",
    "strategy_type": "CONDITIONAL",
    "is_active": true,
    "created_after": "2025-01-01T00:00:00Z",
    "created_before": "2025-12-31T00:00:00Z"
  }'
```

### 10. Get Default Strategies

Retrieve all default (system) strategies.

**Endpoint:** `GET /strategy/default`

**cURL Example:**
```bash
curl -X GET "/api/strategy/default" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 11. Clone Default Strategy

Clone a default strategy for the current user.

**Endpoint:** `POST /strategy/default/{id}/clone`

**cURL Example:**
```bash
curl -X POST "/api/strategy/default/123e4567-e89b-12d3-a456-426614174000/clone" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Field Descriptions

### Required Fields
- **name**: Strategy name (3-100 characters)
- **instruments**: Array of instrument IDs
- **chart_type**: Chart type ("CANDLESTICK")
- **chart_timeframe**: Time frame ("1m", "3m", "5m", "10m", "15m", "30m", "1h", "1d")
- **target_profit**: Target profit value
- **target_profit_type**: Type of target profit ("ABSOLUTE", "PERCENTAGE")
- **stop_loss**: Stop loss value
- **stop_loss_type**: Type of stop loss ("ABSOLUTE", "PERCENTAGE")
- **segment**: Market segment ("F&O", "EQ")
- **strategy_type**: Strategy type ("CONDITIONAL", "TIME_BASED")

### Conditional Fields
- **entry_expression**: Required when strategy_type is "CONDITIONAL"
- **entry_time**: Required when strategy_type is "TIME_BASED"
- **exit_time**: Required when strategy_type is "TIME_BASED"

### EQ Segment Requirements
When segment is "EQ", the following fields are required:
- **position_type**: Position direction ("LONG", "SHORT")
- **trade_type**: Trade type ("MIS", "CNC/NRML")
- **qty**: Quantity
- **max_allocation**: Maximum allocation amount

### Optional Fields
- **description**: Strategy description (max 500 characters)
- **position_type**: Position direction ("LONG", "SHORT") - required for EQ segment
- **trading_days**: Array of trading days ("MON", "TUE", "WED", "THU", "FRI")
- **trade_type**: Trade type ("MIS", "CNC/NRML") - required for EQ segment
- **qty**: Quantity - required for EQ segment
- **exit_expression**: Exit condition expression
- **trailing_stop_loss**: Trailing stop loss value
- **trailing_stop_loss_type**: Type of trailing stop loss ("ABSOLUTE", "PERCENTAGE")
- **exit_type**: Exit type
- **capital**: Capital amount
- **max_allocation**: Maximum allocation - required for EQ segment
- **max_trades**: Maximum number of trades
- **brokerage**: Brokerage amount
- **slippage**: Slippage percentage (0-100)
- **entry_start_time**: Entry start time (HH:MM:SS)
- **entry_end_time**: Entry end time (HH:MM:SS)
- **start_date**: Strategy start date (ISO 8601 format)
- **end_date**: Strategy end date (ISO 8601 format)

## Response Format

All successful responses follow this format:
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data here
  }
}
```

Error responses follow this format:
```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

## HTTP Status Codes

- **200**: Success
- **201**: Created
- **400**: Bad Request (validation errors)
- **401**: Unauthorized (invalid or missing token)
- **403**: Forbidden (insufficient permissions)
- **404**: Not Found
- **409**: Conflict (e.g., strategy cannot be modified due to deployments)
- **500**: Internal Server Error

## Business Rules

1. **Conditional Strategies**: Must have `entry_expression` provided
2. **Time-Based Strategies**: Must have both `entry_time` and `exit_time` provided
3. **EQ Segment**: Must have `position_type`, `trade_type`, `qty`, and `max_allocation` provided
4. **Strategy Modification**: Strategies with existing deployments become read-only and cannot be modified
5. **Default Strategies**: Only system administrators can create default strategies
6. **Validation**: All numeric values must be non-negative, and percentages must be between 0-100

## Examples for Different Use Cases

### Equity Strategy (EQ Segment)
```bash
curl -X POST "/api/strategy" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Equity Buy and Hold",
    "description": "Simple buy and hold strategy for equity",
    "position_type": "LONG",
    "instruments": [1],
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "1d",
    "trade_type": "CNC/NRML",
    "qty": 10,
    "target_profit": 10.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 5.0,
    "stop_loss_type": "PERCENTAGE",
    "max_allocation": 10000.0,
    "segment": "EQ",
    "strategy_type": "CONDITIONAL",
    "entry_expression": "PRICE > SMA(50)"
  }'
```

### Futures & Options Strategy (F&O Segment)
```bash
curl -X POST "/api/strategy" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Options Scalping",
    "description": "Quick scalping strategy for options",
    "instruments": [101, 102],
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "1m",
    "target_profit": 20.0,
    "target_profit_type": "ABSOLUTE",
    "stop_loss": 10.0,
    "stop_loss_type": "ABSOLUTE",
    "segment": "F&O",
    "strategy_type": "CONDITIONAL",
    "entry_expression": "VOLUME > AVG_VOLUME(20) * 2"
  }'
```
