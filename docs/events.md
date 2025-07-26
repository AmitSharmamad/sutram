# WebSocket Events Documentation

## Overview

The Sutram application provides real-time updates through page-specific WebSocket events. The frontend connects to the WebSocket endpoint and receives events filtered by the current page context. This ensures that each page only receives relevant events, improving performance and user experience.

## Connection

### WebSocket Endpoint
```
ws://localhost:9090/api/events/ws
```

### Connection Parameters
When connecting to the WebSocket, you must provide the page context in the query parameters:

```javascript
const pageContext = {
  page_type: "dashboard", // or deployment_list, deployment_detail, portfolio, settings
  page_params: {
    deployment_id: "uuid-here" // only required for deployment_detail page
  }
};

const ws = new WebSocket(`ws://localhost:9090/api/events/ws?page_type=${pageContext.page_type}&page_params=${encodeURIComponent(JSON.stringify(pageContext.page_params))}`);
```

### Authentication
Include the JWT token in the Authorization header:
```javascript
const ws = new WebSocket(url, [], {
  headers: {
    'Authorization': `Bearer ${jwtToken}`
  }
});
```

## Event Structure

All events follow this base structure:

```typescript
interface WebEvent {
  event_type: string;           // Specific event identifier
  page_type: PageType;          // Target page type
  priority: EventPriority;      // high, medium, low
  data: any;                   // Event-specific data
  metadata?: Record<string, any>; // Additional context
  timestamp: string;           // ISO timestamp
  expires_at?: string;         // Optional expiration time
}

type PageType = "dashboard" | "deployment_list" | "deployment_detail" | "portfolio" | "settings";
type EventPriority = "high" | "medium" | "low";
```

## Page-Specific Events

### Dashboard Events

#### 1. Dashboard Summary Updated
```typescript
{
  event_type: "dashboard.summary_updated",
  page_type: "dashboard",
  priority: "medium",
  data: {
    total_deployments: 15,
    active_deployments: 8,
    total_pnl: 12500.75,
    today_pnl: 450.25,
    open_positions: 12,
    win_rate: 0.68,
    last_updated: "2025-07-26T10:15:00Z"
  }
}
```

#### 2. Recent Activity
```typescript
{
  event_type: "dashboard.recent_activity",
  page_type: "dashboard",
  priority: "medium",
  data: {
    activity: "Position opened: LONG 100 @ 1520.30",
    deployment_id: "559ce6bf-3b5e-41a2-b51a-b5a44cab8929",
    deployment_name: "RSI Momentum Strategy",
    timestamp: "2025-07-26T10:15:00Z",
    type: "position_opened",
    details: {
      direction: "LONG",
      quantity: 100,
      entry_price: 1520.30,
      position_id: "f6f0619a-0f28-4b58-8f39-f53867ca3fec"
    }
  }
}
```

#### 3. Alerts Updated
```typescript
{
  event_type: "dashboard.alerts_updated",
  page_type: "dashboard",
  priority: "high",
  data: {
    critical_count: 2,
    warning_count: 5,
    info_count: 3
  }
}
```

### Deployment List Events

#### 1. Deployment Item Updated
```typescript
{
  event_type: "deployment_list.item_updated",
  page_type: "deployment_list",
  priority: "medium",
  data: {
    deployment_id: "559ce6bf-3b5e-41a2-b51a-b5a44cab8929",
    name: "RSI Momentum Strategy",
    status: "active",
    pnl: 1250.50,
    open_positions: 3,
    last_signal: "position_opened",
    last_active: "2025-07-26T10:15:00Z"
  }
}
```

### Deployment Detail Events

#### 1. Position Opened
```typescript
{
  event_type: "deployment_detail.position_opened",
  page_type: "deployment_detail",
  priority: "high",
  data: {
    position_id: "f6f0619a-0f28-4b58-8f39-f53867ca3fec",
    entry_price: 1520.30,
    quantity: 100,
    direction: "LONG",
    status: "open",
    timestamp: "2025-07-26T10:15:00Z"
  }
}
```

#### 2. Position Closed
```typescript
{
  event_type: "deployment_detail.position_closed",
  page_type: "deployment_detail",
  priority: "high",
  data: {
    position_id: "f6f0619a-0f28-4b58-8f39-f53867ca3fec",
    exit_price: 1545.80,
    pnl: 2550.00,
    status: "closed",
    timestamp: "2025-07-26T10:20:00Z"
  }
}
```

#### 3. Signal Detected
```typescript
{
  event_type: "deployment_detail.signal_detected",
  page_type: "deployment_detail",
  priority: "high",
  data: {
    signal_type: "entry",
    price: 1520.30,
    confidence: 0.85,
    reason: "RSI oversold condition with bullish divergence",
    timestamp: "2025-07-26T10:15:00Z"
  }
}
```

#### 4. Metrics Updated
```typescript
{
  event_type: "deployment_detail.metrics_updated",
  page_type: "deployment_detail",
  priority: "medium",
  data: {
    total_trades: 47,
    winning_trades: 32,
    losing_trades: 15,
    win_rate: 0.68,
    total_pnl: 12500.75,
    average_pnl: 265.96,
    max_drawdown: -850.50,
    sharpe: 1.45,
    last_updated: "2025-07-26T10:15:00Z"
  }
}
```

#### 5. Chart Data Updated
```typescript
{
  event_type: "deployment_detail.chart_updated",
  page_type: "deployment_detail",
  priority: "low",
  data: {
    timestamp: "2025-07-26T10:15:00Z",
    price: 1520.30,
    volume: 125000,
    indicators: {
      rsi: 32.5,
      sma_20: 1518.75,
      bollinger_upper: 1555.20,
      bollinger_lower: 1485.40
    }
  }
}
```

#### 6. Log Entry
```typescript
{
  event_type: "deployment_detail.log_entry",
  page_type: "deployment_detail",
  priority: "low",
  data: {
    log_id: "log_123456",
    level: "info",
    message: "Entry signal detected",
    context: {
      rsi_value: 32.5,
      price: 1520.30,
      volume: 125000
    },
    timestamp: "2025-07-26T10:15:00Z"
  }
}
```

### Portfolio Events

#### 1. Balance Updated
```typescript
{
  event_type: "portfolio.balance_updated",
  page_type: "portfolio",
  priority: "medium",
  data: {
    total_balance: 105250.75,
    available_balance: 89750.25,
    used_margin: 15500.50,
    unrealized_pnl: 2250.30,
    realized_pnl: 450.25,
    last_updated: "2025-07-26T10:15:00Z"
  }
}
```

#### 2. Positions Summary Updated
```typescript
{
  event_type: "portfolio.positions_summary_updated",
  page_type: "portfolio",
  priority: "medium",
  data: {
    total_positions: 12,
    long_positions: 8,
    short_positions: 4,
    total_unrealized_pnl: 2250.30,
    total_realized_pnl: 12500.75
  }
}
```

### Global Events (All Pages)

#### 1. Alert Notifications
```typescript
{
  event_type: "alerts.new_notification",
  page_type: "", // Empty means all pages
  priority: "high",
  data: {
    alert_id: "alert_123456",
    type: "deployment_error",
    severity: "critical",
    title: "Deployment Error",
    message: "Error in RSI Momentum Strategy: Connection timeout",
    requires_action: true,
    timestamp: "2025-07-26T10:15:00Z",
    context: {
      deployment_id: "559ce6bf-3b5e-41a2-b51a-b5a44cab8929",
      error_code: "CONN_TIMEOUT"
    }
  },
  expires_at: "2025-07-27T10:15:00Z"
}
```

#### 2. Stop Loss Triggered
```typescript
{
  event_type: "alerts.new_notification",
  page_type: "",
  priority: "high",
  data: {
    alert_id: "stoploss_123456",
    type: "stop_loss",
    severity: "warning",
    title: "Stop Loss Triggered",
    message: "Stop loss hit at 1485.50 for RSI Momentum Strategy",
    timestamp: "2025-07-26T10:15:00Z",
    context: {
      position_id: "f6f0619a-0f28-4b58-8f39-f53867ca3fec",
      trigger_price: 1485.50,
      deployment_id: "559ce6bf-3b5e-41a2-b51a-b5a44cab8929"
    }
  }
}
```

#### 3. Target Profit Reached
```typescript
{
  event_type: "alerts.new_notification",
  page_type: "",
  priority: "medium",
  data: {
    alert_id: "target_123456",
    type: "target_profit",
    severity: "info",
    title: "Target Reached",
    message: "Target profit hit at 1565.80 for RSI Momentum Strategy",
    timestamp: "2025-07-26T10:15:00Z",
    context: {
      position_id: "f6f0619a-0f28-4b58-8f39-f53867ca3fec",
      trigger_price: 1565.80,
      deployment_id: "559ce6bf-3b5e-41a2-b51a-b5a44cab8929"
    }
  }
}
```

## Frontend Implementation

### 1. Basic WebSocket Connection

```javascript
class SutramWebSocket {
  constructor(pageType, pageParams = {}, token) {
    this.pageType = pageType;
    this.pageParams = pageParams;
    this.token = token;
    this.ws = null;
    this.eventHandlers = new Map();
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
  }

  connect() {
    const url = `ws://localhost:9090/api/events/ws?page_type=${this.pageType}&page_params=${encodeURIComponent(JSON.stringify(this.pageParams))}`;
    
    this.ws = new WebSocket(url);
    
    // Set authorization header (if supported by browser)
    if (this.token) {
      this.ws.onopen = () => {
        // Send authentication message
        this.ws.send(JSON.stringify({
          type: 'auth',
          token: this.token
        }));
      };
    }

    this.ws.onmessage = (event) => {
      try {
        const webEvent = JSON.parse(event.data);
        this.handleEvent(webEvent);
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    };

    this.ws.onclose = () => {
      this.handleReconnect();
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  handleEvent(webEvent) {
    const handlers = this.eventHandlers.get(webEvent.event_type) || [];
    handlers.forEach(handler => {
      try {
        handler(webEvent);
      } catch (error) {
        console.error(`Error handling event ${webEvent.event_type}:`, error);
      }
    });

    // Handle global handlers
    const globalHandlers = this.eventHandlers.get('*') || [];
    globalHandlers.forEach(handler => {
      try {
        handler(webEvent);
      } catch (error) {
        console.error('Error in global event handler:', error);
      }
    });
  }

  on(eventType, handler) {
    if (!this.eventHandlers.has(eventType)) {
      this.eventHandlers.set(eventType, []);
    }
    this.eventHandlers.get(eventType).push(handler);
  }

  off(eventType, handler) {
    const handlers = this.eventHandlers.get(eventType);
    if (handlers) {
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
  }

  updatePageContext(newPageType, newPageParams = {}) {
    this.pageType = newPageType;
    this.pageParams = newPageParams;
    
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({
        type: 'update_context',
        page_context: {
          page_type: newPageType,
          page_params: newPageParams
        }
      }));
    }
  }

  handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
      setTimeout(() => this.connect(), delay);
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
  }
}
```

### 2. React Hook Example

```javascript
import { useEffect, useRef, useState } from 'react';

export function useWebSocket(pageType, pageParams, token) {
  const [isConnected, setIsConnected] = useState(false);
  const [events, setEvents] = useState([]);
  const wsRef = useRef(null);

  useEffect(() => {
    const ws = new SutramWebSocket(pageType, pageParams, token);
    wsRef.current = ws;

    // Global event handler to track all events
    ws.on('*', (event) => {
      setEvents(prev => [event, ...prev.slice(0, 99)]); // Keep last 100 events
    });

    ws.ws.onopen = () => setIsConnected(true);
    ws.ws.onclose = () => setIsConnected(false);

    ws.connect();

    return () => {
      ws.disconnect();
    };
  }, [pageType, JSON.stringify(pageParams), token]);

  const subscribe = (eventType, handler) => {
    if (wsRef.current) {
      wsRef.current.on(eventType, handler);
    }
  };

  const unsubscribe = (eventType, handler) => {
    if (wsRef.current) {
      wsRef.current.off(eventType, handler);
    }
  };

  const updateContext = (newPageType, newPageParams) => {
    if (wsRef.current) {
      wsRef.current.updatePageContext(newPageType, newPageParams);
    }
  };

  return {
    isConnected,
    events,
    subscribe,
    unsubscribe,
    updateContext
  };
}
```

### 3. Dashboard Component Example

```javascript
import React, { useEffect, useState } from 'react';
import { useWebSocket } from './hooks/useWebSocket';

export function Dashboard({ token }) {
  const [summary, setSummary] = useState(null);
  const [activities, setActivities] = useState([]);
  const [alerts, setAlerts] = useState({ critical_count: 0, warning_count: 0, info_count: 0 });
  
  const { isConnected, subscribe, unsubscribe } = useWebSocket('dashboard', {}, token);

  useEffect(() => {
    const handleSummaryUpdate = (event) => {
      setSummary(event.data);
    };

    const handleRecentActivity = (event) => {
      setActivities(prev => [event.data, ...prev.slice(0, 9)]); // Keep last 10
    };

    const handleAlertsUpdate = (event) => {
      setAlerts(event.data);
    };

    subscribe('dashboard.summary_updated', handleSummaryUpdate);
    subscribe('dashboard.recent_activity', handleRecentActivity);
    subscribe('dashboard.alerts_updated', handleAlertsUpdate);

    return () => {
      unsubscribe('dashboard.summary_updated', handleSummaryUpdate);
      unsubscribe('dashboard.recent_activity', handleRecentActivity);
      unsubscribe('dashboard.alerts_updated', handleAlertsUpdate);
    };
  }, [subscribe, unsubscribe]);

  return (
    <div className="dashboard">
      <div className="connection-status">
        Status: {isConnected ? 'Connected' : 'Disconnected'}
      </div>
      
      {summary && (
        <div className="summary">
          <h2>Portfolio Summary</h2>
          <div>Total P&L: ${summary.total_pnl}</div>
          <div>Today P&L: ${summary.today_pnl}</div>
          <div>Active Deployments: {summary.active_deployments}</div>
          <div>Win Rate: {(summary.win_rate * 100).toFixed(1)}%</div>
        </div>
      )}

      <div className="alerts">
        <h3>Alerts</h3>
        <div>Critical: {alerts.critical_count}</div>
        <div>Warning: {alerts.warning_count}</div>
        <div>Info: {alerts.info_count}</div>
      </div>

      <div className="recent-activity">
        <h3>Recent Activity</h3>
        {activities.map((activity, index) => (
          <div key={index} className="activity-item">
            <div>{activity.activity}</div>
            <div className="timestamp">{new Date(activity.timestamp).toLocaleString()}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
```

## Event Filtering

Events are automatically filtered based on:

1. **Page Type**: Only events matching the current page type are sent
2. **User ID**: Only events for the authenticated user are sent
3. **Deployment Context**: For deployment detail pages, only events for the specific deployment
4. **Global Events**: Some events (like alerts) are sent to all pages

## Error Handling

- **Connection Lost**: Implement automatic reconnection with exponential backoff
- **Invalid Events**: Log and ignore malformed events
- **Authentication Errors**: Redirect to login page
- **Rate Limiting**: Implement client-side rate limiting for reconnection attempts

## Best Practices

1. **Event Handlers**: Always wrap event handlers in try-catch blocks
2. **Memory Management**: Limit the number of stored events to prevent memory leaks
3. **Performance**: Use debouncing for high-frequency events like chart updates
4. **User Experience**: Show connection status to users
5. **Error Recovery**: Implement graceful degradation when WebSocket is unavailable

## Security Considerations

1. **Authentication**: Always validate JWT tokens on the server side
2. **Authorization**: Ensure users only receive events they're authorized to see
3. **Rate Limiting**: Implement server-side rate limiting to prevent abuse
4. **Input Validation**: Validate all incoming WebSocket messages
5. **CORS**: Configure proper CORS settings for production deployments
