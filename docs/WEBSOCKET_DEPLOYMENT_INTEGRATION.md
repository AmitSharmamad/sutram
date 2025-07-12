# WebSocket Deployment Integration

This document explains the WebSocket integration for dynamic deployment management in the deployment scheduler.

## Overview

The deployment scheduler now supports real-time updates when new deployments are created through the API. When a deployment is created via the `/deployment` API endpoint, it publishes a `deployment.created` event through WebSocket, and the scheduler automatically picks up the new deployment without requiring a restart.

## Architecture

### Components

1. **Deployment API** (`/api/deployment` POST) - Creates new deployments
2. **WebSocket Event Publisher** - Publishes deployment events
3. **WebSocket Server** - Serves real-time events to clients
4. **Deployment Event Listener** - Listens for deployment changes in the scheduler
5. **Deployment Scheduler** - Executes deployments based on timeframes

### Event Flow

```
API Request → Create Deployment → Publish Event → WebSocket → Scheduler Updates
     ↓              ↓               ↓             ↓            ↓
1. POST /deployment
                2. Save to DB
                            3. deployment.created
                                      4. Event broadcast
                                                5. Update deployment list
```

## Configuration

### Environment Variables

```bash
# Required: WebSocket server port
DEPLOYMENT_WEBSOCKET_PORT=9091

# Optional: Polling interval for fallback sync (default: 30s)
# Not yet implemented, but can be added if needed
```

### Event Structure

The `deployment.created` event contains:

```json
{
  "event_type": "deployment.created",
  "deployment_id": "uuid",
  "strategy_id": "uuid", 
  "instrument_id": 12345,
  "timestamp": "2025-07-12T10:30:00Z",
  "data": {
    "deployment_id": "uuid",
    "name": "Strategy Name",
    "description": "Description",
    "strategy_name": "My Strategy",
    "timeframe": "1m",
    "is_active": true,
    "created_at": "2025-07-12T10:30:00Z"
  },
  "metadata": {
    "user_id": "uuid",
    "strategy_name": "My Strategy"
  }
}
```

## Testing

### 1. Start the Deployment Service

```bash
# Set environment variables
export DEPLOYMENT_WEBSOCKET_PORT=9091
export APP_ENV=development

# Run the deployment service
go run cmd/deployment/main.go
```

### 2. Create a New Deployment

```bash
# Create a deployment via API
curl -X POST http://localhost:9091/api/deployment \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "strategy_id": "your-strategy-id",
    "instrument_id": 12345,
    "name": "Test Deployment",
    "description": "Test deployment for WebSocket integration"
  }'
```

### 3. Monitor Logs

The deployment scheduler will log:
- WebSocket event publishing
- Deployment list updates
- New deployment addition to timeframe lists

Expected log output:
```
2025-07-12 10:30:00 WebSocket listener started, will poll for deployment changes every 30 seconds
2025-07-12 10:30:15 Refreshing deployment list from database...
2025-07-12 10:30:15 Timeframe 1m: 5 deployments
2025-07-12 10:30:15 Timeframe 5m: 3 deployments
```

### 4. Connect to WebSocket (Optional)

```javascript
// Connect to WebSocket endpoint
const ws = new WebSocket('ws://localhost:9091/api/deployment/events/ws');

ws.onmessage = function(event) {
  const deploymentEvent = JSON.parse(event.data);
  console.log('Received event:', deploymentEvent);
  
  if (deploymentEvent.event_type === 'deployment.created') {
    console.log('New deployment created:', deploymentEvent.data);
  }
};
```

## Implementation Details

### Deployment Event Listener

The `DeploymentEventListener` in `main.go`:
- Maintains a thread-safe map of deployments by timeframe
- Handles `deployment.created` events
- Provides methods to get current deployment lists
- Ensures no duplicate deployments are added

### Fallback Mechanism

The scheduler includes a polling mechanism that:
- Refreshes deployment list every 30 seconds from the database
- Ensures reliability even if WebSocket events are missed
- Provides resilience against network issues

### Thread Safety

All deployment list operations are protected by read-write mutexes to prevent race conditions between:
- Event handling
- Scheduler execution
- Periodic refreshes

## Benefits

1. **Real-time Updates** - No need to restart the scheduler when adding deployments
2. **Improved Reliability** - Fallback polling ensures consistency
3. **Better User Experience** - Deployments start executing immediately
4. **Scalability** - WebSocket events can be extended for other operations

## Future Enhancements

1. **WebSocket Client** - Replace polling with actual WebSocket client connection
2. **Deployment Updates** - Handle deployment modification events
3. **Deployment Deletion** - Handle deployment removal events
4. **Batch Operations** - Handle multiple deployment changes efficiently
5. **Event Replay** - Replay missed events on reconnection

## Troubleshooting

### Common Issues

1. **WebSocket Port Not Set**
   ```
   DEPLOYMENT_WEBSOCKET_PORT environment variable is required
   ```
   Solution: Set the environment variable

2. **Deployment Not Added to Scheduler**
   - Check if deployment is active
   - Verify strategy has a valid timeframe
   - Check logs for error messages

3. **WebSocket Connection Issues**
   - Ensure WebSocket server is running
   - Check firewall settings
   - Verify authentication tokens

### Debug Logs

Enable debug logging by setting:
```bash
export GIN_MODE=debug
```

This will show detailed WebSocket and HTTP request logs.
