package dto

import (
	"time"

	"github.com/google/uuid"
)

// DashboardStatisticsResponse represents the complete dashboard statistics
type DashboardStatisticsResponse struct {
	TradingPerformance  TradingPerformanceStats   `json:"trading_performance"`
	LiveTradingStatus   LiveTradingStatusStats    `json:"live_trading_status"`
	StrategyDevelopment StrategyDevelopmentStats  `json:"strategy_development"`
	MarketplaceActivity *MarketplaceActivityStats `json:"marketplace_activity,omitempty"`
	CommunityEngagement CommunityEngagementStats  `json:"community_engagement"`
	QuickAnalytics      QuickAnalyticsStats       `json:"quick_analytics"`
	RecentActivity      RecentActivityStats       `json:"recent_activity"`
	GeneratedAt         time.Time                 `json:"generated_at"`
}

// TradingPerformanceStats represents portfolio and trading performance metrics
type TradingPerformanceStats struct {
	TotalPnL               float64           `json:"total_pnl"`
	CurrentPortfolioValue  float64           `json:"current_portfolio_value"`
	ActiveDeploymentsCount int               `json:"active_deployments_count"`
	TodaysPnL              float64           `json:"todays_pnl"`
	OverallWinRate         float64           `json:"overall_win_rate"`
	BestPerformingStrategy *BestStrategyInfo `json:"best_performing_strategy,omitempty"`
	PortfolioSharpeRatio   *float64          `json:"portfolio_sharpe_ratio,omitempty"`
	MaxDrawdown            float64           `json:"max_drawdown"`
	TotalTradesExecuted    int               `json:"total_trades_executed"`
}

// BestStrategyInfo represents information about the best performing strategy
type BestStrategyInfo struct {
	StrategyID   uuid.UUID `json:"strategy_id"`
	StrategyName string    `json:"strategy_name"`
	ReturnPct    float64   `json:"return_pct"`
}

// LiveTradingStatusStats represents current live trading status
type LiveTradingStatusStats struct {
	ActivePositions  int                  `json:"active_positions"`
	PendingOrders    int                  `json:"pending_orders"`
	RecentTrades     []RecentTradeInfo    `json:"recent_trades"`
	StrategyStatuses []StrategyStatusInfo `json:"strategy_statuses"`
}

// RecentTradeInfo represents a recent trade summary
type RecentTradeInfo struct {
	TradeID        uuid.UUID `json:"trade_id"`
	StrategyName   string    `json:"strategy_name"`
	InstrumentName string    `json:"instrument_name"`
	Direction      string    `json:"direction"`
	PnL            float64   `json:"pnl"`
	ExecutedAt     time.Time `json:"executed_at"`
}

// StrategyStatusInfo represents strategy deployment status
type StrategyStatusInfo struct {
	DeploymentID   uuid.UUID `json:"deployment_id"`
	StrategyName   string    `json:"strategy_name"`
	InstrumentName string    `json:"instrument_name"`
	IsActive       bool      `json:"is_active"`
	CurrentPnL     float64   `json:"current_pnl"`
}

// StrategyDevelopmentStats represents strategy creation and backtesting metrics
type StrategyDevelopmentStats struct {
	TotalStrategiesCreated     int      `json:"total_strategies_created"`
	BacktestsCompleted         int      `json:"backtests_completed"`
	MostProfitableBacktestPct  *float64 `json:"most_profitable_backtest_pct,omitempty"`
	AverageBacktestPerformance *float64 `json:"average_backtest_performance,omitempty"`
}

// MarketplaceActivityStats represents marketplace seller statistics
type MarketplaceActivityStats struct {
	TotalRevenueEarned    float64 `json:"total_revenue_earned"`
	ActiveListings        int     `json:"active_listings"`
	TotalSubscribers      int     `json:"total_subscribers"`
	AverageStrategyRating float64 `json:"average_strategy_rating"`
}

// CommunityEngagementStats represents community participation metrics
type CommunityEngagementStats struct {
	ReputationScore int               `json:"reputation_score"`
	BadgesEarned    int               `json:"badges_earned"`
	TopicCount      int               `json:"topic_count"`
	PostCount       int               `json:"post_count"`
	FollowersCount  int               `json:"followers_count"`
	RecentBadges    []RecentBadgeInfo `json:"recent_badges"`
}

// RecentBadgeInfo represents a recently earned badge
type RecentBadgeInfo struct {
	BadgeID     int       `json:"badge_id"`
	BadgeName   string    `json:"badge_name"`
	Description string    `json:"description"`
	EarnedAt    time.Time `json:"earned_at"`
}

// QuickAnalyticsStats represents analytical summaries
type QuickAnalyticsStats struct {
	MonthlyPerformance []MonthlyPerformanceData `json:"monthly_performance"`
	AssetAllocation    []AssetAllocationData    `json:"asset_allocation"`
	TradingFrequency   TradingFrequencyData     `json:"trading_frequency"`
	RiskMetrics        RiskMetricsData          `json:"risk_metrics"`
}

// MonthlyPerformanceData represents monthly P&L data
type MonthlyPerformanceData struct {
	Month string  `json:"month"`
	PnL   float64 `json:"pnl"`
}

// AssetAllocationData represents allocation by instrument
type AssetAllocationData struct {
	InstrumentName string  `json:"instrument_name"`
	AllocationPct  float64 `json:"allocation_pct"`
	CurrentValue   float64 `json:"current_value"`
}

// TradingFrequencyData represents trading frequency metrics
type TradingFrequencyData struct {
	TradesPerDay     float64 `json:"trades_per_day"`
	TradesPerWeek    float64 `json:"trades_per_week"`
	AvgTradeDuration float64 `json:"avg_trade_duration_hours"`
}

// RiskMetricsData represents risk assessment metrics
type RiskMetricsData struct {
	ValueAtRisk       *float64 `json:"value_at_risk,omitempty"`
	Volatility        float64  `json:"volatility"`
	DownsideDeviation *float64 `json:"downside_deviation,omitempty"`
}

// RecentActivityStats represents recent activity feed
type RecentActivityStats struct {
	DeploymentEvents    []ActivityEvent `json:"deployment_events"`
	SystemNotifications []ActivityEvent `json:"system_notifications"`
	CommunityActivity   []ActivityEvent `json:"community_activity"`
	MarketplaceUpdates  []ActivityEvent `json:"marketplace_updates"`
}

// ActivityEvent represents a single activity event
type ActivityEvent struct {
	EventID     uuid.UUID   `json:"event_id"`
	EventType   string      `json:"event_type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Timestamp   time.Time   `json:"timestamp"`
	Metadata    interface{} `json:"metadata,omitempty"`
}
