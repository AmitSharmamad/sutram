# Strategy API Documentation

This document provides sample cURL commands for the Strategy API endpoints.

## Authentication

All strategy endpoints require authentication. Include the Authorization header in your requests:

```bash
Authorization: Bearer <your-jwt-token>
```

## Create Strategy

### Create EQ Strategy

```bash
curl -X POST http://localhost:8080/api/v1/strategy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "name": "Simple EQ Strategy",
    "description": "A basic equity strategy for long positions",
    "position_type": "LONG",
    "instruments": [1, 2, 3],
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "5m",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "MIS",
    "qty": 100,
    "entry_expression": "close > sma(20)",
    "exit_expression": "close < sma(20)",
    "target_profit": 5.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 2.0,
    "stop_loss_type": "PERCENTAGE",
    "capital": 100000.0,
    "max_allocation": 50000.0,
    "max_trades": 5,
    "brokerage": 20,
    "slippage": 0.1,
    "entry_start_time": "09:15:00",
    "entry_end_time": "15:00:00",
    "segment": "EQ",
    "strategy_type": "CONDITIONAL"
  }'
```

### Create F&O Strategy with Option Legs

```bash
curl -X POST http://localhost:8080/api/v1/strategy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "name": "Iron Condor Options Strategy",
    "description": "A market neutral options strategy using iron condor",
    "position_type": "BOTH",
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "15m",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "NRML",
    "entry_expression": "volatility < 20",
    "exit_expression": "volatility > 30 OR days_to_expiry < 5",
    "target_profit": 10.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 5.0,
    "stop_loss_type": "PERCENTAGE",
    "capital": 500000.0,
    "max_allocation": 100000.0,
    "max_trades": 3,
    "brokerage": 40,
    "slippage": 0.2,
    "entry_start_time": "09:30:00",
    "entry_end_time": "14:30:00",
    "segment": "F&O",
    "strategy_type": "CONDITIONAL",
    "option_legs": [
      {
        "leg_sequence": 1,
        "select_index": "NIFTY",
        "segment": "F&O",
        "option_type": "CALL",
        "action_type": "SELL",
        "strike_type": "OTM1",
        "strike_offset": 100,
        "total_lot": 1,
        "order_type": "MARKET",
        "expiry_type": "Monthly",
        "expiry_cycle": "CURRENT"
      },
      {
        "leg_sequence": 2,
        "select_index": "NIFTY",
        "segment": "F&O",
        "option_type": "PUT",
        "action_type": "SELL",
        "strike_type": "OTM1",
        "strike_offset": -100,
        "total_lot": 1,
        "order_type": "MARKET",
        "expiry_type": "Monthly",
        "expiry_cycle": "CURRENT"
      },
      {
        "leg_sequence": 3,
        "select_index": "NIFTY",
        "segment": "F&O",
        "option_type": "CALL",
        "action_type": "BUY",
        "strike_type": "OTM2",
        "strike_offset": 200,
        "total_lot": 1,
        "order_type": "MARKET",
        "expiry_type": "Monthly",
        "expiry_cycle": "CURRENT"
      },
      {
        "leg_sequence": 4,
        "select_index": "NIFTY",
        "segment": "F&O",
        "option_type": "PUT",
        "action_type": "BUY",
        "strike_type": "OTM2",
        "strike_offset": -200,
        "total_lot": 1,
        "order_type": "MARKET",
        "expiry_type": "Monthly",
        "expiry_cycle": "CURRENT"
      }
    ]
  }'
```

### Create Time-Based Strategy

```bash
curl -X POST http://localhost:8080/api/v1/strategy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "name": "Time Based Entry Strategy",
    "description": "Strategy that enters trades at specific times",
    "position_type": "LONG",
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "1m",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "MIS",
    "qty": 50,
    "entry_expression": "time_based_entry",
    "exit_expression": "time_based_exit",
    "target_profit": 3.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 1.5,
    "stop_loss_type": "PERCENTAGE",
    "capital": 200000.0,
    "segment": "EQ",
    "strategy_type": "TIME_BASED",
    "entry_time": "10:00:00",
    "exit_time": "14:30:00"
  }'
```

## Get Strategy by ID

```bash
curl -X GET http://localhost:8080/api/v1/strategy/{strategy_id} \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Get User Strategies

```bash
curl -X GET http://localhost:8080/api/v1/strategy \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Update Strategy

### Update EQ Strategy

```bash
curl -X PUT http://localhost:8080/api/v1/strategy/{strategy_id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "name": "Updated EQ Strategy",
    "description": "Updated description for equity strategy",
    "position_type": "LONG",
    "instruments": [1, 2, 3, 4],
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "10m",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "MIS",
    "qty": 150,
    "entry_expression": "close > sma(50)",
    "exit_expression": "close < sma(50)",
    "target_profit": 7.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 3.0,
    "stop_loss_type": "PERCENTAGE",
    "capital": 150000.0,
    "max_allocation": 75000.0,
    "segment": "EQ",
    "strategy_type": "CONDITIONAL"
  }'
```

### Update F&O Strategy with New Option Legs

```bash
curl -X PUT http://localhost:8080/api/v1/strategy/{strategy_id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "name": "Updated Options Strategy",
    "description": "Updated to bull call spread",
    "position_type": "LONG",
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "5m",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "NRML",
    "entry_expression": "trend_bullish",
    "exit_expression": "trend_bearish OR days_to_expiry < 3",
    "target_profit": 15.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 8.0,
    "stop_loss_type": "PERCENTAGE",
    "capital": 300000.0,
    "segment": "F&O",
    "strategy_type": "CONDITIONAL",
    "option_legs": [
      {
        "leg_sequence": 1,
        "select_index": "BANKNIFTY",
        "segment": "F&O",
        "option_type": "CALL",
        "action_type": "BUY",
        "strike_type": "ATM",
        "strike_offset": 0,
        "total_lot": 2,
        "order_type": "MARKET",
        "expiry_type": "Weekly",
        "expiry_cycle": "CURRENT"
      },
      {
        "leg_sequence": 2,
        "select_index": "BANKNIFTY",
        "segment": "F&O",
        "option_type": "CALL",
        "action_type": "SELL",
        "strike_type": "OTM1",
        "strike_offset": 100,
        "total_lot": 2,
        "order_type": "MARKET",
        "expiry_type": "Weekly",
        "expiry_cycle": "CURRENT"
      }
    ]
  }'
```

## Delete Strategy

```bash
curl -X DELETE http://localhost:8080/api/v1/strategy/{strategy_id} \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Activate Strategy

```bash
curl -X POST http://localhost:8080/api/v1/strategy/{strategy_id}/activate \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Deactivate Strategy

```bash
curl -X POST http://localhost:8080/api/v1/strategy/{strategy_id}/deactivate \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Clone Strategy

```bash
curl -X POST http://localhost:8080/api/v1/strategy/{strategy_id}/clone \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Get Default Strategies

```bash
curl -X GET http://localhost:8080/api/v1/strategy/default \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Clone Default Strategy

```bash
curl -X POST http://localhost:8080/api/v1/strategy/default/{default_strategy_id}/clone \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Search Strategies

```bash
curl -X POST http://localhost:8080/api/v1/strategy/search \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "name": "Options",
    "position_type": "BOTH",
    "segment": "F&O",
    "is_active": true
  }'
```

## Response Examples

### Successful F&O Strategy Creation Response

```json
{
  "success": true,
  "message": "Strategy created successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "user_id": "user-uuid",
    "name": "Iron Condor Options Strategy",
    "description": "A market neutral options strategy using iron condor",
    "position_type": "BOTH",
    "instruments": [],
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "15m",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "NRML",
    "entry_expression": "volatility < 20",
    "exit_expression": "volatility > 30 OR days_to_expiry < 5",
    "target_profit": 10.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 5.0,
    "stop_loss_type": "PERCENTAGE",
    "capital": 500000.0,
    "max_allocation": 100000.0,
    "max_trades": 3,
    "brokerage": 40,
    "slippage": 0.2,
    "entry_start_time": "09:30:00",
    "entry_end_time": "14:30:00",
    "segment": "F&O",
    "strategy_type": "CONDITIONAL",
    "created_at": "2025-07-28T10:30:00Z",
    "updated_at": "2025-07-28T10:30:00Z",
    "is_active": false,
    "is_readonly": false,
    "option_legs": [
      {
        "id": "leg-uuid-1",
        "strategy_id": "123e4567-e89b-12d3-a456-426614174000",
        "leg_sequence": 1,
        "select_index": "NIFTY",
        "segment": "F&O",
        "option_type": "CALL",
        "action_type": "SELL",
        "strike_type": "OTM1",
        "strike_offset": 100,
        "total_lot": 1,
        "order_type": "MARKET",
        "expiry_type": "mONTHLY",
        "expiry_cycle": "CURRENT",
        "created_at": "2025-07-28T10:30:00Z",
        "updated_at": "2025-07-28T10:30:00Z"
      }
    ]
  }
}
```

### Successful EQ Strategy Creation Response

```json
{
  "success": true,
  "message": "Strategy created successfully",
  "data": {
    "id": "456e7890-e89b-12d3-a456-426614174111",
    "user_id": "user-uuid",
    "name": "Simple EQ Strategy",
    "description": "A basic equity strategy for long positions",
    "position_type": "LONG",
    "instruments": [1, 2, 3],
    "chart_type": "CANDLESTICK",
    "chart_timeframe": "5m",
    "trading_days": ["MON", "TUE", "WED", "THU", "FRI"],
    "trade_type": "MIS",
    "qty": 100,
    "entry_expression": "close > sma(20)",
    "exit_expression": "close < sma(20)",
    "target_profit": 5.0,
    "target_profit_type": "PERCENTAGE",
    "stop_loss": 2.0,
    "stop_loss_type": "PERCENTAGE",
    "capital": 100000.0,
    "max_allocation": 50000.0,
    "max_trades": 5,
    "brokerage": 20,
    "slippage": 0.1,
    "entry_start_time": "09:15:00",
    "entry_end_time": "15:00:00",
    "segment": "EQ",
    "strategy_type": "CONDITIONAL",
    "created_at": "2025-07-28T10:30:00Z",
    "updated_at": "2025-07-28T10:30:00Z",
    "is_active": false,
    "is_readonly": false,
    "option_legs": []
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Invalid request: validation failed"
}
```

## Notes

1. **Option Legs**: The `option_legs` field is always present in all strategy responses. For `F&O` segment strategies, it contains the actual option legs data. For non-F&O segments (like `EQ`), it returns an empty array `[]`. Option legs are only created and managed when the strategy segment is `F&O`.

2. **Leg Sequence**: If `leg_sequence` is not provided in option legs, it will be automatically assigned based on the array index (starting from 1).

3. **Strategy ID**: When creating option legs, the `strategy_id` is automatically set and should not be provided in the request.

4. **Authentication**: All endpoints require a valid JWT token in the Authorization header.

5. **Validation**: The API validates all required fields and data types. Invalid requests will return a 400 Bad Request status with error details.

6. **Readonly Strategies**: Strategies that have deployments cannot be modified and will have `is_readonly: true` in the response.

7. **Update Behavior**: When updating F&O strategies with option legs, existing option legs are deleted and replaced with the new ones provided in the request.
