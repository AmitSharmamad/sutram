# Deployment API Documentation

This document provides detailed information about the Deployment APIs in the Sutram platform. All endpoints require authentication.

## Authentication

All API endpoints require authentication via JWT token. The token should be provided in the Authorization header or as a cookie.

```bash
# Token in Authorization header (Bearer token)
-H "Authorization: Bearer YOUR_JWT_TOKEN"

# Or token in cookie
-b "token=YOUR_JWT_TOKEN"
```

## Deployment Management

### Create Deployment

Creates a new deployment for a strategy and instrument.

**Endpoint**: `POST /deployment`

**Request**:
```bash
curl -X POST "http://localhost:8080/api/deployment" \
  -H "Content-Type: application/json" \

  -d '{
    "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
    "instrument_id": 12345,
    "name": "My BTC Strategy Deployment",
    "description": "Bitcoin strategy deployment for daily trading"
  }'
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
    "instrument_id": 12345,
    "name": "My BTC Strategy Deployment",
    "description": "Bitcoin strategy deployment for daily trading",
    "is_active": false,
    "created_at": "2025-05-23T10:30:45Z",
    "updated_at": "2025-05-23T10:30:45Z"
  }
}
```

### Get All User Deployments

Retrieves all deployments for the authenticated user.

**Endpoint**: `GET /deployment`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployment" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
      "instrument_id": 12345,
      "name": "My BTC Strategy Deployment",
      "description": "Bitcoin strategy deployment for daily trading",
      "is_active": false,
      "created_at": "2025-05-23T10:30:45Z",
      "updated_at": "2025-05-23T10:30:45Z"
    },
    {
      "id": "b8c23f65-7c0e-5c9f-0b3d-90e5d24f8d10",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "strategy_id": "661f9511-f30c-52e5-b827-557766551111",
      "instrument_id": 67890,
      "name": "ETH Trading Strategy",
      "description": "Ethereum trading strategy",
      "is_active": true,
      "created_at": "2025-05-22T14:20:30Z",
      "updated_at": "2025-05-22T14:20:30Z"
    }
  ]
}
```

### Search Deployments

Searches deployments based on filter criteria.

**Endpoint**: `GET /deployment/search`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployment/search?strategy_id=550e8400-e29b-41d4-a716-446655440000&is_active=true" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
      "instrument_id": 12345,
      "name": "My BTC Strategy Deployment",
      "description": "Bitcoin strategy deployment for daily trading",
      "is_active": true,
      "created_at": "2025-05-23T10:30:45Z",
      "updated_at": "2025-05-23T10:30:45Z"
    }
  ]
}
```

### Get Deployment By ID

Retrieves a specific deployment by its ID.

**Endpoint**: `GET /deployment/:deployment_id`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployment/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
    "instrument_id": 12345,
    "name": "My BTC Strategy Deployment",
    "description": "Bitcoin strategy deployment for daily trading",
    "is_active": false,
    "created_at": "2025-05-23T10:30:45Z",
    "updated_at": "2025-05-23T10:30:45Z"
  }
}
```

### Update Deployment

Updates an existing deployment.

**Endpoint**: `PUT /deployment/:deployment_id`

**Request**:
```bash
curl -X PUT "http://localhost:8080/api/deployment/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09" \
  -H "Content-Type: application/json" \

  -d '{
    "id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "name": "Updated BTC Strategy Deployment",
    "description": "Updated Bitcoin strategy deployment description"
  }'
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
    "instrument_id": 12345,
    "name": "Updated BTC Strategy Deployment",
    "description": "Updated Bitcoin strategy deployment description",
    "is_active": false,
    "created_at": "2025-05-23T10:30:45Z",
    "updated_at": "2025-05-23T11:15:20Z"
  }
}
```

### Delete Deployment

Deletes a deployment.

**Endpoint**: `DELETE /deployment/:deployment_id`

**Request**:
```bash
curl -X DELETE "http://localhost:8080/api/deployment/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "message": "Deployment deleted successfully"
}
```

### Activate Deployment

Activates a deployment to start running.

**Endpoint**: `PUT /deployment/:deployment_id/activate`

**Request**:
```bash
curl -X PUT "http://localhost:8080/api/deployment/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/activate"
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
    "instrument_id": 12345,
    "name": "My BTC Strategy Deployment",
    "description": "Bitcoin strategy deployment for daily trading",
    "is_active": true,
    "created_at": "2025-05-23T10:30:45Z",
    "updated_at": "2025-05-23T11:45:30Z"
  }
}
```

### Deactivate Deployment

Deactivates a running deployment.

**Endpoint**: `PUT /deployment/:deployment_id/deactivate`

**Request**:
```bash
curl -X PUT "http://localhost:8080/api/deployment/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/deactivate" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
    "instrument_id": 12345,
    "name": "My BTC Strategy Deployment",
    "description": "Bitcoin strategy deployment for daily trading",
    "is_active": false,
    "created_at": "2025-05-23T10:30:45Z",
    "updated_at": "2025-05-23T12:15:10Z"
  }
}
```

### Stop All Deployments

Stops all active deployments for the current user.

**Endpoint**: `POST /deployment/stop-all`

**Request**:
```bash
curl -X POST "http://localhost:8080/api/deployment/stop-all" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "message": "All deployments stopped successfully",
  "data": {
    "stopped_count": 3
  }
}
```

### Start All Deployments

Starts all inactive deployments for the current user.

**Endpoint**: `POST /deployment/start-all`

**Request**:
```bash
curl -X POST "http://localhost:8080/api/deployment/start-all" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "message": "All deployments started successfully",
  "data": {
    "started_count": 2
  }
}
```

### Get Deployments By Strategy

Retrieves all deployments for a specific strategy.

**Endpoint**: `GET /deployment/strategy/:strategy_id`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployment/strategy/550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
      "instrument_id": 12345,
      "name": "My BTC Strategy Deployment",
      "description": "Bitcoin strategy deployment for daily trading",
      "is_active": false,
      "created_at": "2025-05-23T10:30:45Z",
      "updated_at": "2025-05-23T12:15:10Z"
    },
    {
      "id": "c9d34g76-8d1f-6e0g-1c4f-01f6e35g9e21",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
      "instrument_id": 67890,
      "name": "My ETH Strategy Deployment",
      "description": "Ethereum strategy deployment using the same strategy",
      "is_active": true,
      "created_at": "2025-05-23T11:45:30Z",
      "updated_at": "2025-05-23T11:45:30Z"
    }
  ]
}
```

### Deploy All Instruments

Deploys a strategy across all supported instruments.

**Endpoint**: `POST /deployment/all`

**Request**:
```bash
curl -X POST "http://localhost:8080/api/deployment/all" \
  -H "Content-Type: application/json" \

  -d '{
    "strategy_id": "550e8400-e29b-41d4-a716-446655440000",
    "name_prefix": "Crypto Strategy",
    "description": "Automated deployment across all crypto instruments"
  }'
```

**Response**:
```json
{
  "success": true,
  "message": "Successfully deployed to all instruments",
  "data": {
    "deployment_count": 5,
    "deployments": [
      {
        "id": "d0e45h87-9e2g-7f1h-2d5g-12g7f46h0f32",
        "instrument_id": 12345,
        "name": "Crypto Strategy - BTC"
      },
      {
        "id": "e1f56i98-0f3h-8g2i-3e6h-23h8g57i1g43",
        "instrument_id": 67890,
        "name": "Crypto Strategy - ETH"
      }
    ]
  }
}
```

## Deployment Logs

### Create Deployment Log

Creates a new log entry for a deployment.

**Endpoint**: `POST /deployments/:deploymentID/logs`

**Request**:
```bash
curl -X POST "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/logs" \
  -H "Content-Type: application/json" \

  -d '{
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "status": "STARTED",
    "metadata": {
      "environment": "production",
      "server_id": "srv-123456",
      "memory_usage_mb": 256
    }
  }'
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 123,
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "status": "STARTED",
    "metadata": {},
    "action_at": "2025-05-23T14:30:45Z"
  }
}
```

### Get Deployment Logs

Retrieves all logs for a specific deployment.

**Endpoint**: `GET /deployments/:deploymentID/logs`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/logs" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 121,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "status": "CREATED",
      "metadata": {},
      "action_at": "2025-05-23T14:15:30Z"
    },
    {
      "id": 122,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "status": "STARTED",
      "metadata": {},
      "action_at": "2025-05-23T14:30:45Z"
    },
    {
      "id": 123,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "status": "STOPPED",
      "metadata": {},
      "action_at": "2025-05-23T16:45:20Z"
    }
  ]
}
```

### Get Latest Deployment Log

Retrieves the most recent log entry for a deployment.

**Endpoint**: `GET /deployments/:deploymentID/logs/latest`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/logs/latest" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 123,
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "status": "STOPPED",
    "metadata": {},
    "action_at": "2025-05-23T16:45:20Z"
  }
}
```

## Deployment Positions

### Create Position

Creates a new trading position for a deployment.

**Endpoint**: `POST /deployments/:deploymentID/positions`

**Request**:
```bash
curl -X POST "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/positions" \
  -H "Content-Type: application/json" \

  -d '{
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "instrument_id": 12345,
    "quantity": 0.25,
    "entry_price": 45000.50,
    "entry_time": "2025-05-23T15:30:45Z",
    "metadata": {
      "strategy_signal": "MACD crossover",
      "confidence": 0.85
    }
  }'
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 456,
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "instrument_id": 12345,
    "quantity": 0.25,
    "entry_price": 45000.50,
    "entry_time": "2025-05-23T15:30:45Z",
    "status": "OPEN",
    "metadata": {},
    "created_at": "2025-05-23T15:30:46Z",
    "updated_at": "2025-05-23T15:30:46Z"
  }
}
```

### Get All Positions for a Deployment

Retrieves all positions for a specific deployment.

**Endpoint**: `GET /deployments/:deploymentID/positions`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/positions" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 455,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "instrument_id": 12345,
      "quantity": 0.15,
      "entry_price": 44800.25,
      "entry_time": "2025-05-23T10:45:30Z",
      "exit_price": 45200.75,
      "exit_time": "2025-05-23T14:20:15Z",
      "pnl": 60.075,
      "status": "CLOSED",
      "metadata": {},
      "created_at": "2025-05-23T10:45:31Z",
      "updated_at": "2025-05-23T14:20:16Z"
    },
    {
      "id": 456,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "instrument_id": 12345,
      "quantity": 0.25,
      "entry_price": 45000.50,
      "entry_time": "2025-05-23T15:30:45Z",
      "status": "OPEN",
      "metadata": {},
      "created_at": "2025-05-23T15:30:46Z",
      "updated_at": "2025-05-23T15:30:46Z"
    }
  ]
}
```

### Get Open Positions for a Deployment

Retrieves all open positions for a specific deployment.

**Endpoint**: `GET /deployments/:deploymentID/positions/open`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/positions/open" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 456,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "instrument_id": 12345,
      "quantity": 0.25,
      "entry_price": 45000.50,
      "entry_time": "2025-05-23T15:30:45Z",
      "status": "OPEN",
      "metadata": {},
      "created_at": "2025-05-23T15:30:46Z",
      "updated_at": "2025-05-23T15:30:46Z"
    }
  ]
}
```

### Get Position by ID

Retrieves a specific position by its ID.

**Endpoint**: `GET /deployments/positions/:positionID`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployments/positions/456" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 456,
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "instrument_id": 12345,
    "quantity": 0.25,
    "entry_price": 45000.50,
    "entry_time": "2025-05-23T15:30:45Z",
    "status": "OPEN",
    "metadata": {},
    "created_at": "2025-05-23T15:30:46Z",
    "updated_at": "2025-05-23T15:30:46Z"
  }
}
```

### Close Position

Closes an open position.

**Endpoint**: `PUT /deployments/positions/:positionID/close`

**Request**:
```bash
curl -X PUT "http://localhost:8080/api/deployments/positions/456/close" \
  -H "Content-Type: application/json" \

  -d '{
    "exit_price": 45500.75,
    "exit_time": "2025-05-23T18:45:30Z"
  }'
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 456,
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "instrument_id": 12345,
    "quantity": 0.25,
    "entry_price": 45000.50,
    "entry_time": "2025-05-23T15:30:45Z",
    "exit_price": 45500.75,
    "exit_time": "2025-05-23T18:45:30Z",
    "pnl": 125.0625,
    "status": "CLOSED",
    "metadata": {},
    "created_at": "2025-05-23T15:30:46Z",
    "updated_at": "2025-05-23T18:45:31Z"
  }
}
```

### Cancel Position

Cancels an open position without executing it.

**Endpoint**: `PUT /deployments/positions/:positionID/cancel`

**Request**:
```bash
curl -X PUT "http://localhost:8080/api/deployments/positions/457/cancel" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "message": "Position cancelled successfully"
}
```

## Deployment Metrics

### Create or Update Metrics

Creates or updates metrics for a deployment on a specific date.

**Endpoint**: `POST /deployments/:deploymentID/metrics`

**Request**:
```bash
curl -X POST "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/metrics" \
  -H "Content-Type: application/json" \

  -d '{
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "metric_date": "2025-05-23T00:00:00Z",
    "total_trades": 15,
    "win_trades": 10,
    "loss_trades": 5,
    "total_pnl": 250.75,
    "max_drawdown": 45.5,
    "sharpe_ratio": 1.8,
    "sortino_ratio": 2.2,
    "metadata": {
      "avg_trade_duration_minutes": 120,
      "max_concurrent_positions": 3
    }
  }'
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 789,
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "metric_date": "2025-05-23T00:00:00Z",
    "total_trades": 15,
    "win_trades": 10,
    "loss_trades": 5,
    "total_pnl": 250.75,
    "max_drawdown": 45.5,
    "sharpe_ratio": 1.8,
    "sortino_ratio": 2.2,
    "win_rate": 0.6667,
    "avg_win": 35.075,
    "avg_loss": -12.0,
    "metadata": {},
    "created_at": "2025-05-23T23:55:00Z",
    "updated_at": "2025-05-23T23:55:00Z"
  }
}
```

### Get All Metrics for a Deployment

Retrieves all metrics for a specific deployment.

**Endpoint**: `GET /deployments/:deploymentID/metrics`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/metrics" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 787,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "metric_date": "2025-05-21T00:00:00Z",
      "total_trades": 12,
      "win_trades": 7,
      "loss_trades": 5,
      "total_pnl": 180.25,
      "max_drawdown": 35.0,
      "sharpe_ratio": 1.5,
      "sortino_ratio": 1.8,
      "win_rate": 0.5833,
      "avg_win": 32.5,
      "avg_loss": -11.0,
      "metadata": {},
      "created_at": "2025-05-21T23:55:00Z",
      "updated_at": "2025-05-21T23:55:00Z"
    },
    {
      "id": 788,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "metric_date": "2025-05-22T00:00:00Z",
      "total_trades": 14,
      "win_trades": 9,
      "loss_trades": 5,
      "total_pnl": 220.50,
      "max_drawdown": 40.0,
      "sharpe_ratio": 1.7,
      "sortino_ratio": 2.0,
      "win_rate": 0.6429,
      "avg_win": 33.5,
      "avg_loss": -11.5,
      "metadata": {},
      "created_at": "2025-05-22T23:55:00Z",
      "updated_at": "2025-05-22T23:55:00Z"
    },
    {
      "id": 789,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "metric_date": "2025-05-23T00:00:00Z",
      "total_trades": 15,
      "win_trades": 10,
      "loss_trades": 5,
      "total_pnl": 250.75,
      "max_drawdown": 45.5,
      "sharpe_ratio": 1.8,
      "sortino_ratio": 2.2,
      "win_rate": 0.6667,
      "avg_win": 35.075,
      "avg_loss": -12.0,
      "metadata": {},
      "created_at": "2025-05-23T23:55:00Z",
      "updated_at": "2025-05-23T23:55:00Z"
    }
  ]
}
```

### Get Latest Metrics for a Deployment

Retrieves the most recent metrics for a deployment.

**Endpoint**: `GET /deployments/:deploymentID/metrics/latest`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/metrics/latest" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 789,
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "metric_date": "2025-05-23T00:00:00Z",
    "total_trades": 15,
    "win_trades": 10,
    "loss_trades": 5,
    "total_pnl": 250.75,
    "max_drawdown": 45.5,
    "sharpe_ratio": 1.8,
    "sortino_ratio": 2.2,
    "win_rate": 0.6667,
    "avg_win": 35.075,
    "avg_loss": -12.0,
    "metadata": {},
    "created_at": "2025-05-23T23:55:00Z",
    "updated_at": "2025-05-23T23:55:00Z"
  }
}
```

### Get Metrics for a Specific Date

Retrieves metrics for a deployment on a specific date.

**Endpoint**: `GET /deployments/:deploymentID/metrics/:date`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/metrics/2025-05-22" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 788,
    "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
    "metric_date": "2025-05-22T00:00:00Z",
    "total_trades": 14,
    "win_trades": 9,
    "loss_trades": 5,
    "total_pnl": 220.50,
    "max_drawdown": 40.0,
    "sharpe_ratio": 1.7,
    "sortino_ratio": 2.0,
    "win_rate": 0.6429,
    "avg_win": 33.5,
    "avg_loss": -11.5,
    "metadata": {},
    "created_at": "2025-05-22T23:55:00Z",
    "updated_at": "2025-05-22T23:55:00Z"
  }
}
```

### Get Metrics for a Time Range

Retrieves metrics for a deployment within a specific time range.

**Endpoint**: `GET /deployments/:deploymentID/metrics/range`

**Request**:
```bash
curl -X GET "http://localhost:8080/api/deployments/a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09/metrics/range?startDate=2025-05-21&endDate=2025-05-23" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 787,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "metric_date": "2025-05-21T00:00:00Z",
      "total_trades": 12,
      "win_trades": 7,
      "loss_trades": 5,
      "total_pnl": 180.25,
      "max_drawdown": 35.0,
      "sharpe_ratio": 1.5,
      "sortino_ratio": 1.8,
      "win_rate": 0.5833,
      "avg_win": 32.5,
      "avg_loss": -11.0,
      "metadata": {},
      "created_at": "2025-05-21T23:55:00Z",
      "updated_at": "2025-05-21T23:55:00Z"
    },
    {
      "id": 788,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "metric_date": "2025-05-22T00:00:00Z",
      "total_trades": 14,
      "win_trades": 9,
      "loss_trades": 5,
      "total_pnl": 220.50,
      "max_drawdown": 40.0,
      "sharpe_ratio": 1.7,
      "sortino_ratio": 2.0,
      "win_rate": 0.6429,
      "avg_win": 33.5,
      "avg_loss": -11.5,
      "metadata": {},
      "created_at": "2025-05-22T23:55:00Z",
      "updated_at": "2025-05-22T23:55:00Z"
    },
    {
      "id": 789,
      "deployment_id": "a7f12e54-6b9d-4b8e-9a2c-89d4c13f7c09",
      "metric_date": "2025-05-23T00:00:00Z",
      "total_trades": 15,
      "win_trades": 10,
      "loss_trades": 5,
      "total_pnl": 250.75,
      "max_drawdown": 45.5,
      "sharpe_ratio": 1.8,
      "sortino_ratio": 2.2,
      "win_rate": 0.6667,
      "avg_win": 35.075,
      "avg_loss": -12.0,
      "metadata": {},
      "created_at": "2025-05-23T23:55:00Z",
      "updated_at": "2025-05-23T23:55:00Z"
    }
  ]
}
```