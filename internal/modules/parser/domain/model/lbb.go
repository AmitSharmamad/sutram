package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"math"
	"strings"
	"time"
)

// MovingAverageType represents the type of moving average to use
type MovingAverageType string

const (
	SIMPLE      MovingAverageType = "SIMPLE"
	EXPONENTIAL MovingAverageType = "EXPONENTIAL"
	WEIGHTED    MovingAverageType = "WEIGHTED"
	TRIPLE      MovingAverageType = "TRIPLE"
	TRIANGULAR  MovingAverageType = "TRIANGULAR"
)

// LowerBollingerBands represents the Lower Bollinger Bands indicator
type LowerBollingerBands struct {
	period      float64
	standardDev float64
	maType      MovingAverageType
	offset      int
}

func (indicator *LowerBollingerBands) SliceData(prices Prices) Prices {
	// We need enough data to calculate Lower Bollinger Bands properly
	period := int(indicator.period)
	start := max(0, len(prices)-(indicator.offset+2*period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// generateIndicatorKey generates a unique key for LowerBollingerBands indicator based on its configuration
func (indicator *LowerBollingerBands) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 4 {
		return "", fmt.Errorf("LowerBollingerBands requires 4 parameters: point, period, standardDev, offset")
	}

	// Extract point parameter
	point, ok := config[0].(string)
	if !ok {
		return "", fmt.Errorf("LowerBollingerBands point parameter must be a string")
	}

	// Extract period parameter
	period, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("LowerBollingerBands period parameter must be a number")
	}

	// Extract standardDev parameter
	standardDev, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("LowerBollingerBands standardDev parameter must be a number")
	}

	// Extract offset parameter
	offsetFloat, ok := config[3].(float64)
	if !ok {
		return "", fmt.Errorf("LowerBollingerBands offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	baseKey := fmt.Sprintf("LBB_%s_%.0f_%.1f", strings.ToUpper(point), period, standardDev)

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

// Calculate the appropriate moving average based on the type
func (indicator *LowerBollingerBands) calculateMA(prices Prices, period int, maType MovingAverageType) float64 {
	if len(prices) < period {
		return 0.0
	}

	switch maType {
	case SIMPLE:
		// Simple Moving Average
		var sum float64
		for i := len(prices) - period; i < len(prices); i++ {
			sum += prices[i]
		}
		return sum / float64(period)

	case EXPONENTIAL:
		// Exponential Moving Average
		var ema float64
		multiplier := 2.0 / float64(period+1)

		// Initialize with SMA
		var sum float64
		for i := len(prices) - period; i < len(prices)-(period/2); i++ {
			sum += prices[i]
		}
		ema = sum / float64(period/2)

		// Calculate EMA
		for i := len(prices) - (period / 2); i < len(prices); i++ {
			ema = (prices[i]-ema)*multiplier + ema
		}
		return ema

	case WEIGHTED:
		// Weighted Moving Average
		var sum, weightSum float64
		for i := 0; i < period; i++ {
			weight := float64(i + 1)
			sum += prices[len(prices)-period+i] * weight
			weightSum += weight
		}
		return sum / weightSum

	case TRIPLE:
		// Triple Exponential Moving Average (TEMA)
		// TEMA = 3*EMA - 3*EMA(EMA) + EMA(EMA(EMA))
		if len(prices) < 3*period {
			return 0.0
		}

		// Calculate first EMA
		multiplier := 2.0 / float64(period+1)
		var ema1 float64

		// Initialize with SMA
		var sum float64
		for i := len(prices) - period; i < len(prices); i++ {
			sum += prices[i]
		}
		ema1 = sum / float64(period)

		// Create EMA prices for second EMA
		ema1Prices := make(Prices, len(prices))
		for i := 0; i < len(prices)-period; i++ {
			ema1Prices[i] = 0
		}

		for i := len(prices) - period; i < len(prices); i++ {
			if i == len(prices)-period {
				ema1Prices[i] = ema1
			} else {
				ema1 = (prices[i]-ema1)*multiplier + ema1
				ema1Prices[i] = ema1
			}
		}

		// Calculate second EMA (EMA of EMA)
		var ema2 float64
		sum = 0
		for i := len(ema1Prices) - period; i < len(ema1Prices); i++ {
			sum += ema1Prices[i]
		}
		ema2 = sum / float64(period)

		// Create EMA2 prices for third EMA
		ema2Prices := make(Prices, len(prices))
		for i := 0; i < len(prices)-period; i++ {
			ema2Prices[i] = 0
		}

		for i := len(ema1Prices) - period; i < len(ema1Prices); i++ {
			if i == len(ema1Prices)-period {
				ema2Prices[i] = ema2
			} else {
				ema2 = (ema1Prices[i]-ema2)*multiplier + ema2
				ema2Prices[i] = ema2
			}
		}

		// Calculate third EMA (EMA of EMA of EMA)
		var ema3 float64
		sum = 0
		for i := len(ema2Prices) - period; i < len(ema2Prices); i++ {
			sum += ema2Prices[i]
		}
		ema3 = sum / float64(period)

		for i := len(ema2Prices) - period; i < len(ema2Prices); i++ {
			if i == len(ema2Prices)-period {
				continue
			} else {
				ema3 = (ema2Prices[i]-ema3)*multiplier + ema3
			}
		}

		// TEMA = 3*EMA1 - 3*EMA2 + EMA3
		return 3*ema1 - 3*ema2 + ema3

	case TRIANGULAR:
		// Triangular Moving Average
		// First, calculate SMA with period
		if len(prices) < 2*period {
			return 0.0
		}

		n := period
		smaValues := make(Prices, period)

		for i := 0; i < period; i++ {
			var sum float64
			for j := 0; j < n; j++ {
				idx := len(prices) - period - n + i + j
				if idx >= 0 && idx < len(prices) {
					sum += prices[idx]
				}
			}
			smaValues[i] = sum / float64(n)
		}

		// Then, calculate SMA of those SMAs
		var sum float64
		for _, sma := range smaValues {
			sum += sma
		}

		return sum / float64(len(smaValues))

	default:
		// Default to Simple Moving Average
		var sum float64
		for i := len(prices) - period; i < len(prices); i++ {
			sum += prices[i]
		}
		return sum / float64(period)
	}
}

// LOWER_BOLLINGER_BANDS(CLOSE, 20, 2, TRIPLE, 0)
// Fields: OHLCV
// Period: float
// Standard Deviations: float
// Type: enum -> SIMPLE, EXPONENTIAL, WEIGHTED, TRIPLE, TRIANGULAR
// Offset: int
func (indicator *LowerBollingerBands) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	point := config[0].(string)
	period := config[1].(float64)
	standardDev := config[2].(float64)
	offset := int(config[3].(float64))

	indicator.period = period
	indicator.standardDev = standardDev
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, point)

	// Check if we have enough data points
	if len(prices) < int(period)+offset {
		return 0.0
	}

	prices = indicator.SliceData(prices)

	// Need at least period data points to calculate
	if len(prices) < int(period) {
		return 0.0
	}

	// Calculate SMA (Simple Moving Average)
	var sum float64
	periodInt := int(period)
	for i := len(prices) - periodInt; i < len(prices); i++ {
		sum += prices[i]
	}
	sma := sum / period

	// Calculate standard deviation
	var variance float64
	for i := len(prices) - periodInt; i < len(prices); i++ {
		diff := prices[i] - sma
		variance += diff * diff
	}
	stdDev := math.Sqrt(variance / period)

	// Calculate Lower Bollinger Band
	lowerBand := sma - (standardDev * stdDev)

	return util.Round(lowerBand)
}

func (indicator *LowerBollingerBands) GetDataRange(config map[string]any) int32 {
	period := indicator.period
	if period == 0 {
		if p, ok := config["period"].(float64); ok {
			period = p
		} else {
			period = 20 // default period for Bollinger Bands
		}
	}

	offset := indicator.offset
	if offset == 0 {
		if o, ok := config["offset"].(float64); ok {
			offset = int(o)
		} else {
			offset = 0 // default offset
		}
	}

	// For complex indicators like TEMA or Triangular, we need more data
	maType := indicator.maType
	if maType == "" {
		if t, ok := config["maType"].(string); ok {
			maType = MovingAverageType(t)
		} else {
			maType = SIMPLE // default to simple MA
		}
	}

	// Additional data needed for more complex MA types
	var additionalData float64
	switch maType {
	case TRIPLE:
		additionalData = period * 3
	case TRIANGULAR:
		additionalData = period * 2
	case EXPONENTIAL, WEIGHTED:
		additionalData = period * 1.5
	default:
		additionalData = period
	}

	return int32(additionalData + float64(offset))
}

// GetLookBack returns the number of historical data points needed for Lower Bollinger Bands calculation
// config[0] = point (string), config[1] = period (float64), config[2] = standardDev (float64)
// config[3] = maType (string), config[4] = offset (float64)
func (indicator *LowerBollingerBands) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is point (string) - not needed for lookback calculation
	period := config[1].(float64) // default: 20
	// config[2] is standardDev (float64) - not needed for lookback calculation
	maType := MovingAverageType(config[3].(string)) // default: SIMPLE
	offset := int(config[4].(float64))              // default: 0

	// Calculate additional data needed based on MA type
	var multiplier float64
	switch maType {
	case TRIPLE:
		multiplier = 3.0 // Triple EMA requires more data
	case TRIANGULAR:
		multiplier = 2.0 // Triangular MA requires double period
	case EXPONENTIAL, WEIGHTED:
		multiplier = 1.5 // EMA and WMA need some extra data for stabilization
	default:
		multiplier = 1.0 // Simple MA
	}

	// For Lower Bollinger Bands calculation, we need:
	// - Period data points adjusted by MA type multiplier
	// - Additional buffer for stability (add period again)
	// - Plus any offset
	return int64(period*multiplier + period + float64(offset))
}

// CalculateAll calculates Lower Bollinger Bands values for all data points
func (indicator *LowerBollingerBands) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	// point := config[0].(string) - handled by Calculate method
	period := config[1].(float64)
	// standardDev := config[2].(float64) - handled by Calculate method
	maType := MovingAverageType(config[3].(string))
	offset := int(config[4].(float64))

	indicator.period = period
	indicator.maType = maType
	indicator.offset = offset

	stockData := *ctx.Data
	var results []IndicatorValues

	// Calculate additional data needed based on MA type
	var multiplier float64
	switch maType {
	case TRIPLE:
		multiplier = 3.0
	case TRIANGULAR:
		multiplier = 2.0
	case EXPONENTIAL, WEIGHTED:
		multiplier = 1.5
	default:
		multiplier = 1.0
	}

	minDataPoints := int(period*multiplier) + offset
	if len(stockData) < minDataPoints {
		return []IndicatorValues{}
	}

	// Calculate for each valid data point
	for i := minDataPoints - 1; i < len(stockData); i++ {
		// Create a temporary context with data up to this point
		tempData := stockData[:i+1]
		tempCtx := &EvaluationContext{
			Data:     &tempData,
			Strategy: ctx.Strategy,
			Results:  make(map[int64]map[string]interface{}),
		}

		// Calculate Lower Bollinger Bands value for this point
		value := indicator.Calculate(config, tempCtx)

		results = append(results, IndicatorValues{
			Timestamp: stockData[i].Timestamp,
			Value:     value,
		})
	}

	// Filter results by startDate
	var filteredResults []IndicatorValues
	for _, result := range results {
		if !result.Timestamp.Before(startDate) {
			filteredResults = append(filteredResults, result)
		}
	}

	return filteredResults
}
