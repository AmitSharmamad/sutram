package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"time"
)

type EMA struct {
	period int
	offset int
}

func (indicator *EMA) SliceData(prices Prices) Prices {
	// We need enough data to calculate EMA properly
	// For EMA calculation we need at least 2*period data points
	// This ensures we have enough data for initial SMA and subsequent EMA calculation
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

func (indicator *EMA) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	field := config[0].(string)
	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.period = period
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, field)

	// Check if we have enough data points
	if len(prices) < period+offset {
		return 0.0
	}

	prices = indicator.SliceData(prices)

	// Need at least period data points to calculate
	if len(prices) < period {
		return 0.0
	}

	// Create unique indicator key
	indicatorKey := fmt.Sprintf("EMA_%s_%d", field, period)

	// Get current timestamp (using the last data point's timestamp)
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()
	// timeFrame := ctx.Strategy.ChartTimeframe

	// Try to get previous EMA value from Results
	if prevValue, exists := getPreviousValueFromResults(ctx, indicatorKey, currentTimestamp-3600); exists {
		if prevEMA, ok := prevValue.(float64); ok {
			// Calculate multiplier (alpha)
			multiplier := 2.0 / float64(period+1)

			// Use the most recent price for EMA calculation
			currentPrice := prices[len(prices)-1]

			// Calculate EMA using previous value: EMA = prevEMA + multiplier * (currentPrice - prevEMA)
			ema := prevEMA + multiplier*(currentPrice-prevEMA)

			// Store the result for future calculations
			insertValueToResults(ctx, indicatorKey, currentTimestamp, ema)

			return util.Round(ema)
		}
	}

	// If no previous value exists, calculate from scratch
	// Calculate multiplier (alpha)
	multiplier := 2.0 / float64(period+1)

	// Calculate initial SMA for the first EMA value
	var sma float64
	for i := 0; i < period; i++ {
		sma += prices[i]
	}
	sma /= float64(period)

	// Calculate EMA recursively
	// Start with the initial SMA value
	ema := sma

	// Only process remaining data points after the initial period
	for i := period; i < len(prices); i++ {
		// EMA formula: EMA(today) = (Price(today) * multiplier) + (EMA(yesterday) * (1 - multiplier))
		// Alternative form: EMA = Price * multiplier + EMA * (1 - multiplier)
		// Simplified: EMA = EMA + multiplier * (Price - EMA)
		ema = ema + multiplier*(prices[i]-ema)
	}

	// Store the result for future calculations
	insertValueToResults(ctx, indicatorKey, currentTimestamp, ema)

	return util.Round(ema)
}

func (indicator *EMA) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 20 // default period from indicators.json
	}

	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}

	return int32(indicator.period + indicator.offset)
}

func (indicator *EMA) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is field (string)
	// config[1] is period (float64)
	// config[2] is offset (float64)

	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	// For EMA calculation, we need sufficient historical data
	// Using 2*period to ensure we have enough data for proper EMA calculation
	// as used in SliceData method, plus the offset
	return int64(5*period + offset)
}

// CalculateAll calculates EMA values for all data points
func (indicator *EMA) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	field := config[0].(string)
	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	indicator.period = period
	indicator.offset = offset

	stockData := *ctx.Data
	var results []IndicatorValues

	// Get prices for the entire dataset
	prices := GetPrices(stockData, field)

	// We need at least enough data to calculate initial EMA
	// For EMA, we typically need at least 2*period data points for proper calculation
	minDataPoints := 2*period + offset
	if len(stockData) < minDataPoints {
		return []IndicatorValues{}
	}

	// Calculate multiplier (alpha)
	multiplier := 2.0 / float64(period+1)

	// Calculate initial SMA for the first EMA value
	startIndex := period + offset
	var sma float64
	smaStart := startIndex - period - offset
	smaEnd := startIndex - offset

	for i := smaStart; i < smaEnd; i++ {
		sma += prices[i]
	}
	sma /= float64(period)

	// Start with the initial SMA value as the first EMA
	ema := sma

	// Calculate EMA for each data point
	for i := startIndex; i < len(stockData); i++ {
		// For the first point, use the SMA calculated above
		if i == startIndex {
			results = append(results, IndicatorValues{
				Timestamp: stockData[i].Timestamp,
				Value:     util.Round(ema),
			})
			continue
		}

		// Calculate EMA using the previous EMA and current price
		currentPrice := prices[i-offset]
		ema = ema + multiplier*(currentPrice-ema)

		results = append(results, IndicatorValues{
			Timestamp: stockData[i].Timestamp,
			Value:     util.Round(ema),
		})
	}

	// Iterate through results to filter by startDate
	var filteredResults []IndicatorValues
	for _, result := range results {
		if !result.Timestamp.Before(startDate) {
			filteredResults = append(filteredResults, result)
		}
	}

	return filteredResults
}
