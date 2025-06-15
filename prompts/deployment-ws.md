# Deployment WebSocket Events API

This document describes the real-time WebSocket API for receiving deployment events from the Sutram trading platform. The WebSocket connection provides live updates about strategy executions, position changes, and market activity.

## Connection

### WebSocket Endpoint
```
ws://localhost:9090/api/deployment/events/ws
```

### Authentication
The WebSocket connection requires authentication via JWT token. Include the token in one of the following ways:

1. **Cookie Authentication** (Recommended):
   ```javascript
   // Token should be set as an HTTP-only cookie named 'token'
   const ws = new WebSocket('ws://localhost:9090/api/deployment/events/ws');
   ```

2. **Header Authentication**:
   ```javascript
   const ws = new WebSocket('ws://localhost:9090/api/deployment/events/ws', [], {
     headers: {
       'Authorization': 'Bearer YOUR_JWT_TOKEN'
     }
   });
   ```

### Connection Example (JavaScript)
```javascript
const ws = new WebSocket('ws://localhost:9090/api/deployment/events/ws');

ws.onopen = function(event) {
    console.log('Connected to deployment events WebSocket');
};

ws.onmessage = function(event) {
    const deploymentEvent = JSON.parse(event.data);
    console.log('Received event:', deploymentEvent);
    handleDeploymentEvent(deploymentEvent);
};

ws.onclose = function(event) {
    console.log('WebSocket connection closed:', event.code, event.reason);
};

ws.onerror = function(error) {
    console.error('WebSocket error:', error);
};
```

## Event Structure

All events follow this base structure:

```json
{
  "event_type": "position_opened",
  "deployment_id": "123e4567-e89b-12d3-a456-426614174000",
  "strategy_id": "987fcdeb-51a2-43d1-9f12-123456789abc",
  "instrument_id": 12345,
  "timestamp": "2025-06-15T14:30:00Z",
  "data": {
    // Event-specific data (varies by event type)
  },
  "metadata": {
    "user_id": "456e7890-e12b-34d5-a678-901234567def",
    "strategy_name": "My Trading Strategy",
    "instrument_symbol": "RELIANCE",
    "market_price": 2450.75
  }
}
```

## Event Types

### 1. Strategy Evaluation Events

#### `evaluation_started`
Fired when strategy evaluation begins for a market data point.

```json
{
  "event_type": "evaluation_started",
  "data": {
    "market_price": 2450.75,
    "timestamp": "2025-06-15T14:30:00Z",
    "timeframe": "1m"
  }
}
```

#### `entry_signal_detected`
Fired when entry conditions are met.

```json
{
  "event_type": "entry_signal_detected",
  "data": {
    "market_price": 2450.75,
    "signal_time": "2025-06-15T14:30:00Z",
    "expression": "CLOSE > SMA(20)",
    "signal_value": true
  }
}
```

#### `exit_signal_detected`
Fired when exit conditions are met.

```json
{
  "event_type": "exit_signal_detected",
  "data": {
    "market_price": 2445.50,
    "signal_time": "2025-06-15T14:35:00Z",
    "expression": "CLOSE < SMA(20)",
    "signal_value": true
  }
}
```

#### `no_signal`
Fired when strategy evaluation completes but no action is taken.

```json
{
  "event_type": "no_signal",
  "data": {
    "price": 2450.75,
    "volume": 150000,
    "timestamp": "2025-06-15T14:30:00Z",
    "ohlc": {
      "open": 2448.00,
      "high": 2452.25,
      "low": 2447.50,
      "close": 2450.75
    }
  }
}
```

### 2. Position Management Events

#### `position_opened`
Fired when a new position is created.

```json
{
  "event_type": "position_opened",
  "data": {
    "position_id": "pos_123456789",
    "entry_price": 2450.75,
    "quantity": 10.0,
    "direction": "Long",
    "entry_time": "2025-06-15T14:30:00Z",
    "stop_loss_price": 2425.75,
    "target_profit_price": 2475.75
  }
}
```

#### `position_closed`
Fired when a position is closed for any reason.

```json
{
  "event_type": "position_closed",
  "data": {
    "position_id": "pos_123456789",
    "exit_price": 2465.25,
    "exit_time": "2025-06-15T15:15:00Z",
    "exit_reason": "Target Profit",
    "pnl": 145.00,
    "duration": "45m0s",
    "entry_price": 2450.75,
    "quantity": 10.0
  }
}
```

#### `position_update`
Fired for real-time position updates (P&L changes, etc.).

```json
{
  "event_type": "position_update",
  "data": {
    "position_id": "pos_123456789",
    "current_price": 2455.25,
    "unrealized_pnl": 45.00,
    "stop_loss_price": 2430.75,
    "target_profit_price": 2475.75
  }
}
```

### 3. Risk Management Events

#### `stop_loss_triggered`
Fired when a stop loss is hit.

```json
{
  "event_type": "stop_loss_triggered",
  "data": {
    "position_id": "pos_123456789",
    "trigger_price": 2425.75,
    "current_price": 2425.50,
    "trigger_reason": "Stop Loss"
  }
}
```

#### `target_profit_triggered`
Fired when target profit is reached.

```json
{
  "event_type": "target_profit_triggered",
  "data": {
    "position_id": "pos_123456789",
    "trigger_price": 2475.75,
    "current_price": 2476.00,
    "trigger_reason": "Target Profit"
  }
}
```

#### `trailing_stop_updated`
Fired when trailing stop loss is updated.

```json
{
  "event_type": "trailing_stop_updated",
  "data": {
    "position_id": "pos_123456789",
    "trigger_price": 2455.25,
    "current_price": 2460.75,
    "trigger_reason": "Trailing Stop Loss",
    "previous_stop_loss_price": 2425.75,
    "new_stop_loss_price": 2435.75
  }
}
```

#### `eod_square_off`
Fired when MIS positions are squared off at end of day.

```json
{
  "event_type": "eod_square_off",
  "data": {
    "position_id": "pos_123456789",
    "exit_price": 2455.00,
    "exit_time": "2025-06-15T15:30:00Z",
    "exit_reason": "EOD MIS Square-off",
    "pnl": 42.50,
    "duration": "6h30m0s",
    "entry_price": 2450.75,
    "quantity": 10.0
  }
}
```

### 4. Deployment Status Events

#### `deployment_error`
Fired when errors occur during strategy execution.

```json
{
  "event_type": "deployment_error",
  "data": {
    "error_message": "Failed to calculate indicators",
    "error_code": "INDICATOR_ERROR",
    "context": "Entry condition evaluation"
  }
}
```

#### `entry_window_closed`
Fired when the entry time window is outside allowed hours.

```json
{
  "event_type": "entry_window_closed",
  "data": {
    "window_type": "entry",
    "window_status": "closed",
    "start_time": "2025-06-15T09:15:00Z",
    "end_time": "2025-06-15T15:30:00Z",
    "current_time": "2025-06-15T16:00:00Z"
  }
}
```

## Frontend Integration Examples

### React Hook for WebSocket Events

```javascript
import { useState, useEffect, useCallback } from 'react';

export const useDeploymentEvents = () => {
  const [events, setEvents] = useState([]);
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [ws, setWs] = useState(null);

  const connect = useCallback(() => {
    const websocket = new WebSocket('ws://localhost:9090/api/deployment/events/ws');
    
    websocket.onopen = () => {
      setConnectionStatus('connected');
      console.log('Connected to deployment events');
    };
    
    websocket.onmessage = (event) => {
      const deploymentEvent = JSON.parse(event.data);
      setEvents(prev => [deploymentEvent, ...prev.slice(0, 99)]); // Keep last 100 events
    };
    
    websocket.onclose = () => {
      setConnectionStatus('disconnected');
      console.log('Disconnected from deployment events');
    };
    
    websocket.onerror = (error) => {
      setConnectionStatus('error');
      console.error('WebSocket error:', error);
    };
    
    setWs(websocket);
  }, []);

  const disconnect = useCallback(() => {
    if (ws) {
      ws.close();
      setWs(null);
    }
  }, [ws]);

  useEffect(() => {
    connect();
    return () => disconnect();
  }, [connect, disconnect]);

  return { events, connectionStatus, connect, disconnect };
};
```

### Event Handler Example

```javascript
const handleDeploymentEvent = (event) => {
  const { event_type, data, metadata } = event;
  
  switch (event_type) {
    case 'position_opened':
      showNotification('Position Opened', `${metadata.strategy_name}: ${data.direction} ${data.quantity} @ ${data.entry_price}`, 'success');
      updatePositionsList(event);
      break;
      
    case 'position_closed':
      const pnlColor = data.pnl >= 0 ? 'success' : 'error';
      showNotification('Position Closed', `PnL: ${data.pnl} (${data.exit_reason})`, pnlColor);
      removeFromPositionsList(data.position_id);
      break;
      
    case 'stop_loss_triggered':
      showNotification('Stop Loss Hit', `Position ${data.position_id} stopped out at ${data.trigger_price}`, 'warning');
      break;
      
    case 'target_profit_triggered':
      showNotification('Target Reached', `Position ${data.position_id} hit target at ${data.trigger_price}`, 'success');
      break;
      
    case 'entry_signal_detected':
      showNotification('Entry Signal', `${metadata.strategy_name}: Signal detected at ${data.market_price}`, 'info');
      break;
      
    case 'deployment_error':
      showNotification('Deployment Error', data.error_message, 'error');
      break;
      
    default:
      console.log('Unhandled event type:', event_type, data);
  }
};
```

### Real-time Dashboard Component

```javascript
import React from 'react';
import { useDeploymentEvents } from './hooks/useDeploymentEvents';

const DeploymentDashboard = () => {
  const { events, connectionStatus } = useDeploymentEvents();
  
  const getPositionEvents = () => 
    events.filter(e => e.event_type.includes('position'));
  
  const getSignalEvents = () => 
    events.filter(e => e.event_type.includes('signal'));
  
  return (
    <div className="deployment-dashboard">
      <div className="connection-status">
        Status: <span className={connectionStatus}>{connectionStatus}</span>
      </div>
      
      <div className="events-grid">
        <div className="position-events">
          <h3>Position Events</h3>
          {getPositionEvents().map((event, index) => (
            <div key={index} className={`event ${event.event_type}`}>
              <div className="event-type">{event.event_type}</div>
              <div className="event-time">{new Date(event.timestamp).toLocaleTimeString()}</div>
              <div className="event-data">{JSON.stringify(event.data, null, 2)}</div>
            </div>
          ))}
        </div>
        
        <div className="signal-events">
          <h3>Signal Events</h3>
          {getSignalEvents().map((event, index) => (
            <div key={index} className={`event ${event.event_type}`}>
              <div className="event-type">{event.event_type}</div>
              <div className="event-time">{new Date(event.timestamp).toLocaleTimeString()}</div>
              <div className="strategy">{event.metadata.strategy_name}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default DeploymentDashboard;
```

## Error Handling

### Connection Errors
- **401 Unauthorized**: Invalid or missing authentication token
- **403 Forbidden**: User doesn't have permission to access deployment events
- **1006 Abnormal Closure**: Network connection lost

### Reconnection Strategy
```javascript
const connectWithRetry = (maxRetries = 5, delay = 1000) => {
  let retryCount = 0;
  
  const connect = () => {
    const ws = new WebSocket('ws://localhost:9090/api/deployment/events/ws');
    
    ws.onopen = () => {
      retryCount = 0; // Reset retry count on successful connection
      console.log('Connected to deployment events');
    };
    
    ws.onclose = () => {
      if (retryCount < maxRetries) {
        retryCount++;
        console.log(`Reconnecting... (${retryCount}/${maxRetries})`);
        setTimeout(connect, delay * retryCount);
      } else {
        console.error('Max reconnection attempts reached');
      }
    };
    
    return ws;
  };
  
  return connect();
};
```

## Best Practices

1. **Handle Connection State**: Always monitor connection status and provide user feedback
2. **Event Filtering**: Filter events based on user's active deployments to reduce noise
3. **Rate Limiting**: Implement client-side throttling for high-frequency events
4. **Persistence**: Consider storing critical events locally for offline viewing
5. **Error Recovery**: Implement proper reconnection logic with exponential backoff
6. **Performance**: Use efficient data structures for real-time event handling
7. **User Experience**: Provide clear visual feedback for different event types

## Security Considerations

- Events are filtered by user ID - users only receive events for their own deployments
- JWT tokens must be valid and not expired
- WebSocket connections are automatically closed on authentication failure
- Consider implementing rate limiting on the server side for high-frequency users

## Testing

You can test the WebSocket connection using browser developer tools or WebSocket testing tools:

```javascript
// Open browser console and run:
const ws = new WebSocket('ws://localhost:9090/api/deployment/events/ws');
ws.onmessage = (event) => console.log(JSON.parse(event.data));
```