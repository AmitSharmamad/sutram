package service

import (
	"math"
	"time"

	"algo/internal/modules/backtest/domain/model"
	mathUtil "algo/internal/util/math"
)

// MetricsCalculationService defines the interface for metrics calculation operations
type MetricsCalculationService interface {
	CalculateRunMetrics(positions []*model.Position, initialCapital, finalCapital float64, run *model.Run) *model.RunMetrics
	CalculateDrawdown(positions []*model.Position, initialCapital float64) model.DrawdownResult
	CalculateReturns(positions []*model.Position) model.ReturnsData
	CalculateBasicStats(positions []*model.Position) BasicStats
	CalculatePositionMetrics(positions []*model.Position) PositionMetrics
	CalculateCapitalProgression(positions []*model.Position) []interface{}
}

// BasicStats holds basic trading statistics
type BasicStats struct {
	TotalTrades     int
	WinningTrades   int
	LosingTrades    int
	BreakEvenTrades int
	TotalProfitLoss float64
	TotalWinAmount  float64
	TotalLossAmount float64
	WinRate         float64
	ProfitFactor    float64
	AvgWinSize      float64
	AvgLossSize     float64
}

// PositionMetrics holds position-related metrics
type PositionMetrics struct {
	AvgPositionSize  float64
	MaxPositionSize  float64
	AvgExposure      float64
	AvgTradeDuration float64
}

// metricsCalculationService implements MetricsCalculationService
type metricsCalculationService struct{}

// NewMetricsCalculationService creates a new metrics calculation service
func NewMetricsCalculationService() MetricsCalculationService {
	return &metricsCalculationService{}
}

// CalculateRunMetrics calculates complete run metrics from positions
func (s *metricsCalculationService) CalculateRunMetrics(positions []*model.Position, initialCapital, finalCapital float64, run *model.Run) *model.RunMetrics {
	metrics := model.NewRunMetrics(run.ID)

	// Calculate basic trading statistics
	basicStats := s.CalculateBasicStats(positions)

	// Set basic metrics
	metrics.TotalTrades = basicStats.TotalTrades
	metrics.WinRate = basicStats.WinRate
	metrics.ProfitFactor = &basicStats.ProfitFactor
	metrics.AvgWinSize = &basicStats.AvgWinSize
	metrics.AvgLossSize = &basicStats.AvgLossSize

	// Calculate position metrics
	positionMetrics := s.CalculatePositionMetrics(positions)
	metrics.AvgPositionSize = &positionMetrics.AvgPositionSize
	metrics.MaxPositionSize = &positionMetrics.MaxPositionSize
	metrics.AvgExposure = &positionMetrics.AvgExposure
	metrics.AvgTradeDuration = &positionMetrics.AvgTradeDuration

	// Set transaction costs to 0 for now, can be addressed later
	metrics.TotalSlippage = 0.0
	metrics.TransactionCostsPct = 0.0

	// Calculate return percentage
	if initialCapital > 0 {
		totalReturnPct := (finalCapital - initialCapital) / initialCapital * 100
		metrics.TotalReturnPct = mathUtil.ClampFloat(totalReturnPct, -9999.9999, 9999.9999)
	}

	// Calculate more advanced metrics
	returnsData := s.CalculateReturns(positions)
	metrics.DailyReturns = returnsData.DailyReturns
	metrics.MonthlyReturns = returnsData.MonthlyReturns
	metrics.AnnualizedReturn = returnsData.AnnualizedReturn

	metrics.SharpeRatio = 0.0  // Would require daily returns calculation
	metrics.SortinoRatio = 0.0 // Would require daily returns calculation

	// Calculate max drawdown and duration
	drawdownResult := s.CalculateDrawdown(positions, initialCapital)
	metrics.MaxDrawdown = drawdownResult.MaxDrawdown
	metrics.MaxDrawdownDuration = drawdownResult.MaxDrawdownDuration

	// Calculate capital progression
	capitalProgression := s.CalculateCapitalProgression(positions)
	metrics.SetCapitalProgression(capitalProgression)

	// Store additional metrics as JSON
	additionalMetrics := map[string]interface{}{
		"initialCapital":       initialCapital,
		"finalCapital":         finalCapital,
		"netProfit":            finalCapital - initialCapital,
		"winningTrades":        basicStats.WinningTrades,
		"losingTrades":         basicStats.LosingTrades,
		"breakEvenTrades":      basicStats.BreakEvenTrades,
		"capital_progression":  capitalProgression,
		"final_capital":        s.calculateFinalCapital(positions),
		"total_capital_change": s.calculateTotalCapitalChange(positions),
		"positions_count":      len(capitalProgression),
	}

	metrics.SetMetricData(additionalMetrics)

	// Update metrics timestamp
	metrics.UpdatedAt = time.Now()

	return metrics
}

// CalculateBasicStats calculates basic trading statistics
func (s *metricsCalculationService) CalculateBasicStats(positions []*model.Position) BasicStats {
	totalTrades := len(positions)
	winningTrades := 0
	losingTrades := 0
	breakEvenTrades := 0
	totalProfitLoss := 0.0
	totalWinAmount := 0.0
	totalLossAmount := 0.0

	// Process all positions
	for _, position := range positions {
		if position.NetProfitLoss == nil {
			continue
		}

		profitLoss := *position.NetProfitLoss

		if profitLoss > 0 {
			winningTrades++
			totalWinAmount += profitLoss
		} else if profitLoss < 0 {
			losingTrades++
			totalLossAmount += profitLoss // This is negative
		} else {
			breakEvenTrades++
		}

		totalProfitLoss += profitLoss
	}

	// Calculate win rate
	winRate := 0.0
	if totalTrades > 0 {
		winRate = float64(winningTrades) / float64(totalTrades)
		winRate = mathUtil.ClampFloat(winRate, 0.0, 1.0)
	}

	// Calculate profit factor
	profitFactor := 0.0
	if totalLossAmount != 0 {
		profitFactor = totalWinAmount / math.Abs(totalLossAmount)
		profitFactor = mathUtil.ClampFloat(profitFactor, -9999.9999, 9999.9999)
	}

	// Calculate average win and loss sizes
	avgWinSize := 0.0
	if winningTrades > 0 {
		avgWinSize = totalWinAmount / float64(winningTrades)
		avgWinSize = mathUtil.ClampFloat(avgWinSize, -9999999999.99999999, 9999999999.99999999)
	}

	avgLossSize := 0.0
	if losingTrades > 0 {
		avgLossSize = totalLossAmount / float64(losingTrades)
		avgLossSize = mathUtil.ClampFloat(avgLossSize, -9999999999.99999999, 9999999999.99999999)
	}

	return BasicStats{
		TotalTrades:     totalTrades,
		WinningTrades:   winningTrades,
		LosingTrades:    losingTrades,
		BreakEvenTrades: breakEvenTrades,
		TotalProfitLoss: totalProfitLoss,
		TotalWinAmount:  totalWinAmount,
		TotalLossAmount: totalLossAmount,
		WinRate:         winRate,
		ProfitFactor:    profitFactor,
		AvgWinSize:      avgWinSize,
		AvgLossSize:     avgLossSize,
	}
}

// CalculatePositionMetrics calculates position-related metrics
func (s *metricsCalculationService) CalculatePositionMetrics(positions []*model.Position) PositionMetrics {
	avgPositionSize := s.calculateAvgPositionSize(positions)
	maxPositionSize := s.calculateMaxPositionSize(positions)
	avgExposure := s.calculateAvgExposure(positions)
	avgTradeDuration := s.calculateAvgTradeDuration(positions)

	return PositionMetrics{
		AvgPositionSize:  avgPositionSize,
		MaxPositionSize:  maxPositionSize,
		AvgExposure:      avgExposure,
		AvgTradeDuration: avgTradeDuration,
	}
}

// calculateAvgPositionSize calculates the average position size from all positions
func (s *metricsCalculationService) calculateAvgPositionSize(positions []*model.Position) float64 {
	totalPositionSize := 0.0
	positionCount := 0

	for _, position := range positions {
		if position.PositionSize > 0 {
			totalPositionSize += position.PositionSize
			positionCount++
		}
	}

	avgPositionSize := 0.0
	if positionCount > 0 {
		avgPositionSize = totalPositionSize / float64(positionCount)
		avgPositionSize = mathUtil.ClampFloat(avgPositionSize, -9999999999.99999999, 9999999999.99999999)
	}

	return avgPositionSize
}

// calculateMaxPositionSize calculates the maximum position size from all positions
func (s *metricsCalculationService) calculateMaxPositionSize(positions []*model.Position) float64 {
	maxPositionSize := 0.0

	for _, position := range positions {
		if position.PositionSize > maxPositionSize {
			maxPositionSize = position.PositionSize
		}
	}

	// Clamp to fit in database DECIMAL(18, 8)
	maxPositionSize = mathUtil.ClampFloat(maxPositionSize, -9999999999.99999999, 9999999999.99999999)

	return maxPositionSize
}

// calculateAvgExposure calculates the average dollar exposure from all positions
// Formula: Average Exposure = Σ(entry_price × position_size) / number_of_positions
func (s *metricsCalculationService) calculateAvgExposure(positions []*model.Position) float64 {
	totalExposure := 0.0
	exposureCount := 0

	for _, position := range positions {
		if position.PositionSize > 0 {
			exposure := position.EntryPriceSignal * position.PositionSize
			totalExposure += exposure
			exposureCount++
		}
	}

	avgExposure := 0.0
	if exposureCount > 0 {
		avgExposure = totalExposure / float64(exposureCount)
		// Clamp to fit in database DECIMAL(18, 8)
		avgExposure = mathUtil.ClampFloat(avgExposure, -9999999999.99999999, 9999999999.99999999)
	}

	return avgExposure
}

// calculateAvgTradeDuration calculates the average trade duration in hours
func (s *metricsCalculationService) calculateAvgTradeDuration(positions []*model.Position) float64 {
	totalDurationSeconds := 0.0
	count := 0

	for _, position := range positions {
		if position.ExitTimeSignal != nil {
			durationSeconds := position.ExitTimeSignal.Sub(position.EntryTimeSignal).Seconds()
			totalDurationSeconds += durationSeconds
			count++
		}
	}

	avgTradeDuration := 0.0
	if count > 0 {
		avgTradeDuration = totalDurationSeconds / float64(count) / 3600 // Convert to hours
		avgTradeDuration = mathUtil.ClampFloat(avgTradeDuration, -9999.9999, 9999.9999)
	}

	return avgTradeDuration
}

// CalculateReturns calculates daily, monthly, and annualized returns
func (s *metricsCalculationService) CalculateReturns(positions []*model.Position) model.ReturnsData {
	emptyReturns := model.ReturnsData{
		DailyReturns:     []float64{},
		MonthlyReturns:   []float64{},
		AnnualizedReturn: 0.0,
	}

	if len(positions) == 0 {
		return emptyReturns
	}

	// Filter only closed positions with valid P&L data
	var closedPositions []*model.Position
	for _, position := range positions {
		if position.State == model.Exited && position.NetProfitLoss != nil &&
			position.ExitTimeSignal != nil {
			closedPositions = append(closedPositions, position)
		}
	}

	if len(closedPositions) == 0 {
		return emptyReturns
	}

	// Group positions by day and month to calculate returns
	dailyReturnsMap := make(map[string]float64)
	monthlyReturnsMap := make(map[string]float64)
	finalReturns := float64(0.0)

	// Calculate daily and monthly returns based on P&L
	for _, position := range closedPositions {
		if position.NetProfitLoss == nil || position.ExitTimeSignal == nil {
			continue
		}

		exitDate := position.ExitTimeSignal.Format("2006-01-02")
		exitMonth := position.ExitTimeSignal.Format("2006-01")

		// Accumulate daily returns
		dailyReturnsMap[exitDate] += *position.NetProfitLoss

		// Accumulate monthly returns
		monthlyReturnsMap[exitMonth] += *position.NetProfitLoss

		finalReturns += *position.NetProfitLoss
	}

	// Convert maps to sorted arrays
	var dailyReturns []float64
	var monthlyReturns []float64

	// Extract daily returns values
	for _, dailyReturn := range dailyReturnsMap {
		dailyReturns = append(dailyReturns, dailyReturn)
	}

	// Extract monthly returns values
	for _, monthlyReturn := range monthlyReturnsMap {
		monthlyReturns = append(monthlyReturns, monthlyReturn)
	}

	// Calculate annualized return based on monthly returns
	// NOTE: For now, annualized return is used to show the total pnl
	annualizedReturn := finalReturns

	// if len(monthlyReturns) > 0 {
	// 	// Calculate total return from monthly data
	// 	totalReturn := 0.0
	// 	for _, monthlyReturn := range monthlyReturns {
	// 		totalReturn += monthlyReturn
	// 	}

	// 	// Calculate average monthly return
	// 	avgMonthlyReturn := totalReturn / float64(len(monthlyReturns))

	// 	// Annualized return = average monthly return * 12
	// 	// This is a simplified calculation - in practice you'd want compound returns
	// 	annualizedReturn = avgMonthlyReturn * 12
	// }

	return model.ReturnsData{
		DailyReturns:     dailyReturns,
		MonthlyReturns:   monthlyReturns,
		AnnualizedReturn: annualizedReturn,
	}
}

// CalculateDrawdown calculates maximum drawdown and its duration
func (s *metricsCalculationService) CalculateDrawdown(positions []*model.Position, initialCapital float64) model.DrawdownResult {
	if len(positions) == 0 {
		return model.DrawdownResult{
			MaxDrawdown:         0.0,
			MaxDrawdownDuration: 0.0,
		}
	}

	// Filter and sort positions that have been exited with valid data
	var exitedPositions []*model.Position
	for _, position := range positions {
		if position.State == model.Exited && position.NetProfitLoss != nil &&
			position.ExitTimeSignal != nil {
			exitedPositions = append(exitedPositions, position)
		}
	}

	if len(exitedPositions) == 0 {
		return model.DrawdownResult{
			MaxDrawdown:         0.0,
			MaxDrawdownDuration: 0.0,
		}
	}

	// Sort positions chronologically by exit time
	exitedPositions = s.sortPositionsByExitTime(exitedPositions)

	// Track capital progression and drawdowns
	runningPeak := initialCapital
	maxDrawdown := 0.0
	maxDrawdownDuration := 0.0
	var drawdownStartTime *time.Time
	cumulativeEquity := initialCapital

	for _, position := range exitedPositions {
		if position.NetProfitLoss != nil {
			// Update cumulative equity
			cumulativeEquity += *position.NetProfitLoss

			// Check if we're at a new peak
			if cumulativeEquity > runningPeak {
				// If we were in a drawdown period, calculate its duration
				if drawdownStartTime != nil {
					currentDrawdownDuration := position.ExitTimeSignal.Sub(*drawdownStartTime).Hours()
					if currentDrawdownDuration > maxDrawdownDuration {
						maxDrawdownDuration = currentDrawdownDuration
					}
					// Reset drawdown tracking since we reached a new peak
					drawdownStartTime = nil
				}
				runningPeak = cumulativeEquity
			} else if cumulativeEquity < runningPeak {
				// We're in a drawdown
				if drawdownStartTime == nil {
					// This is the start of a new drawdown period
					drawdownStartTime = position.ExitTimeSignal
				}
				// Calculate current drawdown percentage
				if runningPeak > 0 {
					drawdown := (runningPeak - cumulativeEquity) / runningPeak
					if drawdown > maxDrawdown {
						maxDrawdown = drawdown
					}
				}
			}
		}
	}

	// Handle case where the backtest ends in a drawdown
	if drawdownStartTime != nil && len(exitedPositions) > 0 {
		lastPosition := exitedPositions[len(exitedPositions)-1]
		currentDrawdownDuration := lastPosition.ExitTimeSignal.Sub(*drawdownStartTime).Hours()
		if currentDrawdownDuration > maxDrawdownDuration {
			maxDrawdownDuration = currentDrawdownDuration
		}
	}

	return model.DrawdownResult{
		MaxDrawdown:         maxDrawdown * 100.0, // Return as percentage
		MaxDrawdownDuration: maxDrawdownDuration,
	}
}

// CalculateCapitalProgression calculates capital progression over time
func (s *metricsCalculationService) CalculateCapitalProgression(positions []*model.Position) []interface{} {
	currentCapital := 0.0
	var capitalProgression []interface{}

	// Sort positions by exit time to ensure chronological order
	sortedPositions := s.sortPositionsByExitTime(positions)

	// Calculate capital progression for each position
	for i, position := range sortedPositions {
		if position.NetProfitLoss != nil {
			// Add the position's P&L to current capital
			currentCapital += *position.NetProfitLoss

			// Create capital progression entry
			progressionEntry := map[string]interface{}{
				"position_index":  i + 1,
				"position_id":     position.ID.String(),
				"entry_time":      position.EntryTimeSignal,
				"exit_time":       position.ExitTimeSignal,
				"net_profit_loss": *position.NetProfitLoss,
				"capital_before":  currentCapital - *position.NetProfitLoss,
				"capital_after":   currentCapital,
				"capital_change":  *position.NetProfitLoss,
				"direction":       string(position.Direction),
				"position_size":   position.PositionSize,
			}

			// Add entry price
			progressionEntry["entry_price"] = position.EntryPriceSignal

			// Add exit price if available
			if position.ExitPriceSignal != nil {
				progressionEntry["exit_price"] = *position.ExitPriceSignal
			}

			// Add exit reason if available
			if position.ExitReason != nil {
				progressionEntry["exit_reason"] = *position.ExitReason
			}

			capitalProgression = append(capitalProgression, progressionEntry)
		}
	}

	return capitalProgression
}

// sortPositionsByExitTime sorts positions chronologically by their exit time
// Positions with nil exit times are placed at the end
func (s *metricsCalculationService) sortPositionsByExitTime(positions []*model.Position) []*model.Position {
	// Create a copy to avoid modifying the original slice
	sortedPositions := make([]*model.Position, len(positions))
	copy(sortedPositions, positions)

	// Simple bubble sort implementation for clarity
	n := len(sortedPositions)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			// Handle nil exit times - put them at the end
			if sortedPositions[j].ExitTimeSignal == nil {
				if sortedPositions[j+1].ExitTimeSignal != nil {
					sortedPositions[j], sortedPositions[j+1] = sortedPositions[j+1], sortedPositions[j]
				}
				continue
			}
			if sortedPositions[j+1].ExitTimeSignal == nil {
				continue
			}

			// Sort by exit time (ascending)
			if sortedPositions[j].ExitTimeSignal.After(*sortedPositions[j+1].ExitTimeSignal) {
				sortedPositions[j], sortedPositions[j+1] = sortedPositions[j+1], sortedPositions[j]
			}
		}
	}

	return sortedPositions
}

// calculateFinalCapital calculates the final capital from all positions
func (s *metricsCalculationService) calculateFinalCapital(positions []*model.Position) float64 {
	finalCapital := 0.0
	for _, position := range positions {
		if position.NetProfitLoss != nil {
			finalCapital += *position.NetProfitLoss
		}
	}
	return finalCapital
}

// calculateTotalCapitalChange calculates the total capital change
func (s *metricsCalculationService) calculateTotalCapitalChange(positions []*model.Position) float64 {
	return s.calculateFinalCapital(positions) // Same as final capital in this context
}
