package dto

import "time"

// PageType represents the type of page the user is viewing
type PageType string

const (
	PageTypeDashboard        PageType = "dashboard"
	PageTypeDeploymentList   PageType = "deployment_list"
	PageTypeDeploymentDetail PageType = "deployment_detail"
	PageTypePortfolio        PageType = "portfolio"
	PageTypeSettings         PageType = "settings"
)

// PageContext represents the context of the user's current page
type PageContext struct {
	PageType   PageType               `json:"page_type"`
	PageParams map[string]interface{} `json:"page_params,omitempty"`
}

// WebSocketConnectionContext represents the context for a WebSocket connection
type WebSocketConnectionContext struct {
	ConnectionID        string      `json:"connection_id"`
	UserID              string      `json:"user_id"`
	PageContext         PageContext `json:"page_context"`
	ConnectedAt         time.Time   `json:"connected_at"`
	LastActiveAt        time.Time   `json:"last_active_at"`
	SubscriptionFilters []string    `json:"subscription_filters,omitempty"`
}

// WebEvent represents a page-specific event sent to the frontend
type WebEvent struct {
	EventType string                 `json:"event_type"`
	PageType  PageType               `json:"page_type"`
	Priority  EventPriority          `json:"priority"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
}

// EventPriority represents the priority of an event
type EventPriority string

const (
	EventPriorityHigh   EventPriority = "high"
	EventPriorityMedium EventPriority = "medium"
	EventPriorityLow    EventPriority = "low"
)

// Dashboard Events
type DashboardSummaryData struct {
	TotalDeployments  int       `json:"total_deployments"`
	ActiveDeployments int       `json:"active_deployments"`
	TotalPnL          float64   `json:"total_pnl"`
	TodayPnL          float64   `json:"today_pnl"`
	OpenPositions     int       `json:"open_positions"`
	WinRate           float64   `json:"win_rate"`
	LastUpdated       time.Time `json:"last_updated"`
}

type DashboardRecentActivityData struct {
	Activity       string                 `json:"activity"`
	DeploymentID   string                 `json:"deployment_id"`
	DeploymentName string                 `json:"deployment_name"`
	Timestamp      time.Time              `json:"timestamp"`
	Type           string                 `json:"type"` // position_opened, position_closed, signal_detected, etc.
	Details        map[string]interface{} `json:"details,omitempty"`
}

type DashboardAlertsData struct {
	CriticalCount int `json:"critical_count"`
	WarningCount  int `json:"warning_count"`
	InfoCount     int `json:"info_count"`
}

// Deployment List Events
type DeploymentListItemData struct {
	DeploymentID  string    `json:"deployment_id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	PnL           float64   `json:"pnl"`
	OpenPositions int       `json:"open_positions"`
	LastSignal    string    `json:"last_signal,omitempty"`
	LastActive    time.Time `json:"last_active"`
}

// Deployment Detail Events
type DeploymentDetailPositionData struct {
	PositionID string    `json:"position_id"`
	EntryPrice float64   `json:"entry_price"`
	ExitPrice  *float64  `json:"exit_price,omitempty"`
	Quantity   int       `json:"quantity"`
	Direction  string    `json:"direction"` // long, short
	Status     string    `json:"status"`    // open, closed, cancelled
	PnL        *float64  `json:"pnl,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

type DeploymentDetailSignalData struct {
	SignalType string                 `json:"signal_type"` // entry, exit
	Direction  string                 `json:"direction"`   // buy, sell
	Price      float64                `json:"price"`
	Confidence float64                `json:"confidence"`
	Reason     string                 `json:"reason"`
	Indicators map[string]interface{} `json:"indicators,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
}

type DeploymentDetailMetricsData struct {
	TotalTrades   int       `json:"total_trades"`
	WinningTrades int       `json:"winning_trades"`
	LosingTrades  int       `json:"losing_trades"`
	WinRate       float64   `json:"win_rate"`
	TotalPnL      float64   `json:"total_pnl"`
	AveragePnL    float64   `json:"average_pnl"`
	MaxDrawdown   float64   `json:"max_drawdown"`
	Sharpe        float64   `json:"sharpe,omitempty"`
	LastUpdated   time.Time `json:"last_updated"`
}

type DeploymentDetailChartData struct {
	Timestamp  time.Time              `json:"timestamp"`
	Price      float64                `json:"price"`
	Volume     int64                  `json:"volume"`
	Indicators map[string]interface{} `json:"indicators,omitempty"`
}

type DeploymentDetailLogData struct {
	LogID     string                 `json:"log_id"`
	Level     string                 `json:"level"` // info, warning, error
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Portfolio Events
type PortfolioBalanceData struct {
	TotalBalance     float64   `json:"total_balance"`
	AvailableBalance float64   `json:"available_balance"`
	UsedMargin       float64   `json:"used_margin"`
	UnrealizedPnL    float64   `json:"unrealized_pnl"`
	RealizedPnL      float64   `json:"realized_pnl"`
	LastUpdated      time.Time `json:"last_updated"`
}

type PortfolioPositionsSummaryData struct {
	TotalPositions     int     `json:"total_positions"`
	LongPositions      int     `json:"long_positions"`
	ShortPositions     int     `json:"short_positions"`
	TotalUnrealizedPnL float64 `json:"total_unrealized_pnl"`
	TotalRealizedPnL   float64 `json:"total_realized_pnl"`
}

// Alert/Notification Events
type AlertData struct {
	AlertID        string                 `json:"alert_id"`
	Type           string                 `json:"type"`     // risk_warning, system_status, trade_alert
	Severity       string                 `json:"severity"` // critical, warning, info
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	ActionURL      string                 `json:"action_url,omitempty"`
	Context        map[string]interface{} `json:"context,omitempty"`
	Timestamp      time.Time              `json:"timestamp"`
	RequiresAction bool                   `json:"requires_action"`
}

// Context Update Message (for changing page context on same connection)
type ContextUpdateMessage struct {
	Type        string      `json:"type"` // should be "update_context"
	PageContext PageContext `json:"page_context"`
}
