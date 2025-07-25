package service

import (
	"context"
	"fmt"
	"math"
	"time"

	backtestService "algo/internal/modules/backtest/application/service"
	communityService "algo/internal/modules/community/application/service"
	deploymentDto "algo/internal/modules/deployment/application/dto"
	deploymentService "algo/internal/modules/deployment/application/service"
	instrumentService "algo/internal/modules/instrument/application/service"
	marketplaceService "algo/internal/modules/marketplace/application/service"
	strategyService "algo/internal/modules/strategy/application/service"
	"algo/internal/modules/user/application/dto"

	"github.com/google/uuid"
)

// DashboardService provides dashboard statistics aggregation
type DashboardService struct {
	deploymentService  deploymentService.DeploymentService
	strategyService    strategyService.StrategyService
	backtestService    backtestService.BacktestService
	marketplaceService *marketplaceService.SellerService
	communityService   *communityService.UserService
	instrumentService  *instrumentService.InstrumentService
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(
	deploymentService deploymentService.DeploymentService,
	strategyService strategyService.StrategyService,
	backtestService backtestService.BacktestService,
	marketplaceService *marketplaceService.SellerService,
	communityService *communityService.UserService,
	instrumentService *instrumentService.InstrumentService,
) *DashboardService {
	return &DashboardService{
		deploymentService:  deploymentService,
		strategyService:    strategyService,
		backtestService:    backtestService,
		marketplaceService: marketplaceService,
		communityService:   communityService,
		instrumentService:  instrumentService,
	}
}

// getInstrumentName retrieves the actual instrument name from the instrument master table
func (s *DashboardService) getInstrumentName(ctx context.Context, instrumentID int) string {
	instrumentMaster, err := s.instrumentService.GetInstrumentMasterByExchangeSecurityID(ctx, int64(instrumentID))
	if err != nil {
		// Fallback to placeholder if lookup fails
		return fmt.Sprintf("Instrument-%d", instrumentID)
	}

	// Prefer TradingSymbol for display, fallback to InstrumentName or Name
	if instrumentMaster.TradingSymbol != "" {
		return instrumentMaster.TradingSymbol
	}
	if instrumentMaster.InstrumentName != "" {
		return instrumentMaster.InstrumentName
	}
	if instrumentMaster.Name != "" {
		return instrumentMaster.Name
	}

	// Final fallback
	return fmt.Sprintf("Instrument-%d", instrumentID)
}

// GetDashboardStatistics aggregates and returns all dashboard statistics for a user
func (s *DashboardService) GetDashboardStatistics(ctx context.Context, userID uuid.UUID) (*dto.DashboardStatisticsResponse, error) {
	// Initialize response
	response := &dto.DashboardStatisticsResponse{
		GeneratedAt: time.Now(),
	}

	// Aggregate trading performance stats
	tradingStats, err := s.getTradingPerformanceStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trading performance stats: %w", err)
	}
	response.TradingPerformance = *tradingStats

	// Aggregate live trading status
	liveStats, err := s.getLiveTradingStatusStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get live trading status: %w", err)
	}
	response.LiveTradingStatus = *liveStats

	// Aggregate strategy development stats
	strategyStats, err := s.getStrategyDevelopmentStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy development stats: %w", err)
	}
	response.StrategyDevelopment = *strategyStats

	// Aggregate marketplace activity (optional - only if user has marketplace activity)
	marketplaceStats, err := s.getMarketplaceActivityStats(ctx, userID)
	if err == nil && marketplaceStats != nil {
		response.MarketplaceActivity = marketplaceStats
	}

	// Aggregate community engagement stats
	communityStats, err := s.getCommunityEngagementStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community engagement stats: %w", err)
	}
	response.CommunityEngagement = *communityStats

	// Aggregate quick analytics
	analyticsStats, err := s.getQuickAnalyticsStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quick analytics: %w", err)
	}
	response.QuickAnalytics = *analyticsStats

	// Aggregate recent activity
	activityStats, err := s.getRecentActivityStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	response.RecentActivity = *activityStats

	return response, nil
}

// getTradingPerformanceStats aggregates trading performance metrics
func (s *DashboardService) getTradingPerformanceStats(ctx context.Context, userID uuid.UUID) (*dto.TradingPerformanceStats, error) {
	// Get all user deployments
	deployments, err := s.deploymentService.GetUserDeployments(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user deployments: %w", err)
	}

	stats := &dto.TradingPerformanceStats{}

	// Calculate aggregated metrics
	var totalPnL float64
	var totalTrades int
	var winningTrades int
	var maxDrawdown float64
	var bestStrategyReturn float64
	var bestStrategyInfo *dto.BestStrategyInfo
	activeCount := 0

	for _, deployment := range deployments {
		if deployment.IsActive {
			activeCount++
		}

		// Get latest metrics for this deployment
		latestMetrics, err := s.deploymentService.GetLatestMetrics(ctx, deployment.ID, userID)
		if err != nil {
			continue // Skip if no metrics available
		}

		totalPnL += latestMetrics.TotalPNL
		totalTrades += latestMetrics.TotalTrades
		winningTrades += latestMetrics.WinTrades

		if latestMetrics.MaxDrawdown > maxDrawdown {
			maxDrawdown = latestMetrics.MaxDrawdown
		}

		// Calculate return percentage (simplified)
		baseCapital := 100000.0 // Default base capital
		returnPct := (latestMetrics.TotalPNL / baseCapital) * 100

		if returnPct > bestStrategyReturn {
			bestStrategyReturn = returnPct
			bestStrategyInfo = &dto.BestStrategyInfo{
				StrategyID:   deployment.StrategyID,
				StrategyName: deployment.Name,
				ReturnPct:    returnPct,
			}
		}
	}

	// Calculate today's P&L (simplified - would need daily metrics in real implementation)
	todaysPnL := totalPnL * 0.02 // Placeholder: assume 2% of total was today

	// Calculate overall win rate
	var overallWinRate float64
	if totalTrades > 0 {
		overallWinRate = float64(winningTrades) / float64(totalTrades)
	}

	// Portfolio value calculation (simplified)
	currentPortfolioValue := 100000.0 + totalPnL // Base capital + total P&L

	stats.TotalPnL = totalPnL
	stats.CurrentPortfolioValue = currentPortfolioValue
	stats.ActiveDeploymentsCount = activeCount
	stats.TodaysPnL = todaysPnL
	stats.OverallWinRate = overallWinRate
	stats.BestPerformingStrategy = bestStrategyInfo
	stats.MaxDrawdown = maxDrawdown
	stats.TotalTradesExecuted = totalTrades

	return stats, nil
}

// getLiveTradingStatusStats gets current live trading status
func (s *DashboardService) getLiveTradingStatusStats(ctx context.Context, userID uuid.UUID) (*dto.LiveTradingStatusStats, error) {
	// Get active deployments
	deployments, err := s.deploymentService.GetUserDeployments(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	stats := &dto.LiveTradingStatusStats{
		RecentTrades:     []dto.RecentTradeInfo{},
		StrategyStatuses: []dto.StrategyStatusInfo{},
	}

	activePositions := 0
	pendingOrders := 0 // Would need order service integration

	for _, deployment := range deployments {
		// Get positions for this deployment
		positions, err := s.deploymentService.GetPositionsByDeployment(ctx, deployment.ID, userID)
		if err != nil {
			continue
		}

		// Count active positions
		for _, position := range positions {
			if position.Status == "OPEN" {
				activePositions++
			}
		}

		// Add to strategy statuses
		latestMetrics, _ := s.deploymentService.GetLatestMetrics(ctx, deployment.ID, userID)
		currentPnL := 0.0
		if latestMetrics != nil {
			currentPnL = latestMetrics.TotalPNL
		}

		strategyStatus := dto.StrategyStatusInfo{
			DeploymentID:   deployment.ID,
			StrategyName:   deployment.Name,
			InstrumentName: s.getInstrumentName(ctx, deployment.InstrumentID),
			IsActive:       deployment.IsActive,
			CurrentPnL:     currentPnL,
		}
		stats.StrategyStatuses = append(stats.StrategyStatuses, strategyStatus)

		// Add recent trades (last 5 closed positions)
		recentCount := 0
		for _, position := range positions {
			if position.Status == "CLOSED" && position.PNL != nil && recentCount < 5 {
				recentTrade := dto.RecentTradeInfo{
					TradeID:        position.ID,
					StrategyName:   deployment.Name,
					InstrumentName: s.getInstrumentName(ctx, deployment.InstrumentID),
					Direction:      string(position.Direction),
					PnL:            *position.PNL,
					ExecutedAt:     position.UpdatedAt,
				}
				stats.RecentTrades = append(stats.RecentTrades, recentTrade)
				recentCount++
			}
		}
	}

	stats.ActivePositions = activePositions
	stats.PendingOrders = pendingOrders

	return stats, nil
}

// getStrategyDevelopmentStats gets strategy creation and backtesting metrics
func (s *DashboardService) getStrategyDevelopmentStats(ctx context.Context, userID uuid.UUID) (*dto.StrategyDevelopmentStats, error) {
	// Get user strategies
	strategies, err := s.strategyService.GetUserStrategies(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user strategies: %w", err)
	}

	stats := &dto.StrategyDevelopmentStats{
		TotalStrategiesCreated: len(strategies),
	}

	// Count backtests and find best performance
	var totalBacktests int
	var totalReturnSum float64
	var bestBacktestReturn float64
	var hasBacktests bool

	for _, strategy := range strategies {
		// Get runs for this strategy
		runs, err := s.backtestService.GetRunsByStrategy(ctx, strategy.ID)
		if err != nil {
			continue
		}

		totalBacktests += len(runs)

		for _, run := range runs {
			// Get metrics for this run
			metrics, err := s.backtestService.GetRunMetrics(ctx, run.ID)
			if err != nil {
				continue
			}

			hasBacktests = true
			totalReturnSum += metrics.TotalReturnPct
			if metrics.TotalReturnPct > bestBacktestReturn {
				bestBacktestReturn = metrics.TotalReturnPct
			}
		}
	}

	stats.BacktestsCompleted = totalBacktests

	if hasBacktests {
		stats.MostProfitableBacktestPct = &bestBacktestReturn
		avgReturn := totalReturnSum / float64(totalBacktests)
		stats.AverageBacktestPerformance = &avgReturn
	}

	return stats, nil
}

// getMarketplaceActivityStats gets marketplace seller statistics
func (s *DashboardService) getMarketplaceActivityStats(ctx context.Context, userID uuid.UUID) (*dto.MarketplaceActivityStats, error) {
	// Try to get seller analytics
	analytics, err := s.marketplaceService.GetSellerAnalytics(ctx, userID)
	if err != nil {
		// User might not be a seller
		return nil, nil
	}

	stats := &dto.MarketplaceActivityStats{
		TotalRevenueEarned:    analytics.TotalRevenue,
		ActiveListings:        analytics.ActiveListings,
		TotalSubscribers:      analytics.TotalSubscriptions,
		AverageStrategyRating: float64(analytics.AverageRating),
	}

	return stats, nil
}

// getCommunityEngagementStats gets community participation metrics
func (s *DashboardService) getCommunityEngagementStats(ctx context.Context, userID uuid.UUID) (*dto.CommunityEngagementStats, error) {
	// Get user profile from community service
	profile, err := s.communityService.GetUserProfile(ctx, userID)
	if err != nil {
		// Return default stats if user hasn't engaged with community
		return &dto.CommunityEngagementStats{
			RecentBadges: []dto.RecentBadgeInfo{},
		}, nil
	}

	stats := &dto.CommunityEngagementStats{
		ReputationScore: profile.ReputationScore,
		BadgesEarned:    len(profile.Badges),
		TopicCount:      profile.TopicCount,
		PostCount:       profile.PostCount,
		FollowersCount:  0, // Would need follower service
		RecentBadges:    []dto.RecentBadgeInfo{},
	}

	// Add recent badges (last 3)
	for i, badge := range profile.Badges {
		if i >= 3 {
			break
		}
		recentBadge := dto.RecentBadgeInfo{
			BadgeID:     badge.ID,
			BadgeName:   badge.Name,
			Description: badge.Description,
			EarnedAt:    badge.CreatedAt, // Assuming badge has created date
		}
		stats.RecentBadges = append(stats.RecentBadges, recentBadge)
	}

	return stats, nil
}

// getQuickAnalyticsStats gets analytical summaries
func (s *DashboardService) getQuickAnalyticsStats(ctx context.Context, userID uuid.UUID) (*dto.QuickAnalyticsStats, error) {
	stats := &dto.QuickAnalyticsStats{
		MonthlyPerformance: []dto.MonthlyPerformanceData{},
		AssetAllocation:    []dto.AssetAllocationData{},
	}

	// Get monthly performance for last 12 months
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	deployments, err := s.deploymentService.GetUserDeployments(ctx, userID)
	if err != nil {
		return stats, nil
	}

	// Calculate monthly performance (simplified)
	monthlyPnL := make(map[string]float64)
	assetAllocation := make(map[int]float64)

	for _, deployment := range deployments {
		// Get metrics time range for this deployment
		filter := &deploymentDto.MetricsFilter{
			DeploymentID: deployment.ID,
			StartDate:    &startDate,
			EndDate:      &endDate,
		}

		metrics, err := s.deploymentService.GetMetricsTimeRange(ctx, filter, userID)
		if err != nil {
			continue
		}

		for _, metric := range metrics {
			monthKey := metric.MetricDate.Format("2006-01")
			monthlyPnL[monthKey] += metric.TotalPNL
		}

		// Track asset allocation by instrument
		latestMetrics, _ := s.deploymentService.GetLatestMetrics(ctx, deployment.ID, userID)
		if latestMetrics != nil {
			assetAllocation[deployment.InstrumentID] += math.Abs(latestMetrics.TotalPNL)
		}
	}

	// Convert to response format
	for month, pnl := range monthlyPnL {
		stats.MonthlyPerformance = append(stats.MonthlyPerformance, dto.MonthlyPerformanceData{
			Month: month,
			PnL:   pnl,
		})
	}

	// Calculate asset allocation percentages
	var totalValue float64
	for _, value := range assetAllocation {
		totalValue += value
	}

	if totalValue > 0 {
		for instrumentID, value := range assetAllocation {
			allocationPct := (value / totalValue) * 100
			stats.AssetAllocation = append(stats.AssetAllocation, dto.AssetAllocationData{
				InstrumentName: s.getInstrumentName(ctx, instrumentID),
				AllocationPct:  allocationPct,
				CurrentValue:   value,
			})
		}
	}

	// Calculate trading frequency (simplified)
	totalDays := float64(endDate.Sub(startDate).Hours() / 24)
	var totalTrades int
	for _, deployment := range deployments {
		latestMetrics, _ := s.deploymentService.GetLatestMetrics(ctx, deployment.ID, userID)
		if latestMetrics != nil {
			totalTrades += latestMetrics.TotalTrades
		}
	}

	tradesPerDay := float64(totalTrades) / totalDays
	tradesPerWeek := tradesPerDay * 7

	stats.TradingFrequency = dto.TradingFrequencyData{
		TradesPerDay:     tradesPerDay,
		TradesPerWeek:    tradesPerWeek,
		AvgTradeDuration: 24.0, // Placeholder - would need position duration analysis
	}

	// Risk metrics (simplified)
	stats.RiskMetrics = dto.RiskMetricsData{
		Volatility: 0.15, // Placeholder - would need proper volatility calculation
	}

	return stats, nil
}

// getRecentActivityStats gets recent activity feed
func (s *DashboardService) getRecentActivityStats(ctx context.Context, userID uuid.UUID) (*dto.RecentActivityStats, error) {
	stats := &dto.RecentActivityStats{
		DeploymentEvents:    []dto.ActivityEvent{},
		SystemNotifications: []dto.ActivityEvent{},
		CommunityActivity:   []dto.ActivityEvent{},
		MarketplaceUpdates:  []dto.ActivityEvent{},
	}

	// Get recent deployment events (simplified)
	deployments, err := s.deploymentService.GetUserDeployments(ctx, userID)
	if err == nil {
		for i, deployment := range deployments {
			if i >= 5 { // Limit to 5 recent events
				break
			}

			event := dto.ActivityEvent{
				EventID:     uuid.New(),
				EventType:   "deployment_status",
				Title:       fmt.Sprintf("Strategy %s", deployment.Name),
				Description: fmt.Sprintf("Deployment is %s", map[bool]string{true: "active", false: "inactive"}[deployment.IsActive]),
				Timestamp:   deployment.UpdatedAt,
			}
			stats.DeploymentEvents = append(stats.DeploymentEvents, event)
		}
	}

	// Add placeholder system notifications
	systemEvent := dto.ActivityEvent{
		EventID:     uuid.New(),
		EventType:   "system_update",
		Title:       "Platform Update",
		Description: "New features available in the strategy builder",
		Timestamp:   time.Now().Add(-24 * time.Hour),
	}
	stats.SystemNotifications = append(stats.SystemNotifications, systemEvent)

	return stats, nil
}
