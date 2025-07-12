package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
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
	start := max(0, len(prices)-(indicator.offset+5*period))
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
	if len(config) != 5 {
		return "", fmt.Errorf("LowerBollingerBands requires 5 parameters: point, period, standardDev, maType, offset")
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

	// Extract maType parameter
	maType, ok := config[3].(string)
	if !ok {
		return "", fmt.Errorf("LowerBollingerBands maType parameter must be a string")
	}

	// Extract offset parameter
	offsetFloat, ok := config[4].(float64)
	if !ok {
		return "", fmt.Errorf("LowerBollingerBands offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	baseKey := fmt.Sprintf("LBB_%s_%.0f_%.1f_%s", strings.ToUpper(point), period, standardDev, strings.ToUpper(maType))

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

// LOWER_BOLLINGER_BANDS(CLOSE, 20, 2, SIMPLE, 0)
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
	maType := config[3].(string)
	offset := int(config[4].(float64))

	indicator.period = period
	indicator.standardDev = standardDev
	indicator.maType = MovingAverageType(strings.ToUpper(maType))
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

	// Convert maType string to talib MAType using utility function
	talibMAType := util.ConvertMATypeToTalib(maType)

	// Use go-talib BBands function
	_, _, lowerBand := talib.BBands(prices, int(period), standardDev, standardDev, talibMAType)

	if len(lowerBand) == 0 {
		return 0.0
	}

	// Return the last calculated value
	return util.Round(lowerBand[len(lowerBand)-1])
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
	// - Plus any offset
	return int64(5*period*multiplier + float64(offset))
}

// CalculateAll calculates Lower Bollinger Bands values for all data points
func (indicator *LowerBollingerBands) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	// Parse config parameters
	point := config[0].(string)
	period := config[1].(float64)
	standardDev := config[2].(float64)
	maType := config[3].(string)
	offset := int(config[4].(float64))

	indicator.period = period
	indicator.standardDev = standardDev
	indicator.maType = MovingAverageType(strings.ToUpper(maType))
	indicator.offset = offset

	// Get stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	// Check if we have enough data
	periodInt := int(period)
	if len(stockData) < periodInt {
		return []IndicatorValues{}
	}

	// Extract price array based on point type
	prices := make([]float64, len(stockData))
	for i, data := range stockData {
		switch strings.ToUpper(point) {
		case "HIGH":
			prices[i] = data.High
		case "LOW":
			prices[i] = data.Low
		case "OPEN":
			prices[i] = data.Open
		case "VOLUME":
			prices[i] = data.Volume
		default: // "CLOSE"
			prices[i] = data.Close
		}
	}

	// Convert maType string to talib MAType using utility function
	talibMAType := util.ConvertMATypeToTalib(maType)

	// Use go-talib BBands function
	_, _, lowerBand := talib.BBands(prices, periodInt, standardDev, standardDev, talibMAType)

	// Build results starting from valid calculation points
	var result []IndicatorValues

	for i := periodInt - 1; i < len(stockData); i++ {
		// Apply offset
		if i < periodInt-1+offset {
			continue
		}

		// Check start date filter
		if stockData[i].Timestamp.Before(startDate) {
			continue
		}

		// Make sure we have valid data from talib
		if i-periodInt+1 >= 0 && i-periodInt+1 < len(lowerBand) {
			result = append(result, IndicatorValues{
				Timestamp: stockData[i].Timestamp,
				Value:     util.Round(lowerBand[i-periodInt+1]),
			})
		}
	}

	return result
}
