package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
)

// StochRsiK represents the Stochastic RSI %K indicator.
type StochRsiK struct {
	rsiPeriod   int
	stochPeriod int
	smoothK     int
	maType      string
	offset      int
}

func (indicator *StochRsiK) SliceData(prices Prices) Prices {
	// We need enough data to calculate StochRsiK properly
	totalPeriod := indicator.rsiPeriod + indicator.stochPeriod + indicator.smoothK
	start := max(0, len(prices)-(indicator.offset+3*totalPeriod))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// GenerateIndicatorKey generates a unique key for StochRsiK indicator based on its configuration
func (indicator *StochRsiK) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 6 {
		return "", fmt.Errorf("StochRsiK requires 6 parameters: field, rsiPeriod, stochPeriod, smoothK, maType, offset")
	}

	// Extract field parameter
	field, ok := config[0].(string)
	if !ok {
		return "", fmt.Errorf("StochRsiK field parameter must be a string")
	}

	// Extract rsiPeriod parameter
	rsiPeriod, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("StochRsiK rsiPeriod parameter must be a number")
	}

	// Extract stochPeriod parameter
	stochPeriod, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("StochRsiK stochPeriod parameter must be a number")
	}

	// Extract smoothK parameter
	smoothK, ok := config[3].(float64)
	if !ok {
		return "", fmt.Errorf("StochRsiK smoothK parameter must be a number")
	}

	// Extract maType parameter
	maType, ok := config[4].(string)
	if !ok {
		return "", fmt.Errorf("StochRsiK maType parameter must be a string")
	}

	// Extract offset parameter
	offsetFloat, ok := config[5].(float64)
	if !ok {
		return "", fmt.Errorf("StochRsiK offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	baseKey := fmt.Sprintf("STOCHRSIK_%s_%.0f_%.0f_%.0f_%s",
		strings.ToUpper(field), rsiPeriod, stochPeriod, smoothK, strings.ToUpper(maType))

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

// STOCHRSIK(CLOSE, 14, 14, 3, SIMPLE, 0)
// Fields: OHLCV
// RSI Period: int
// Stochastic Period: int
// Smooth K Period: int
// MA Type: enum -> SIMPLE, EXPONENTIAL, WEIGHTED, etc.
// Offset: int
func (indicator *StochRsiK) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	field := config[0].(string)
	rsiPeriod := int(config[1].(float64))
	stochPeriod := int(config[2].(float64))
	smoothK := int(config[3].(float64))
	maType := config[4].(string)
	offset := int(config[5].(float64))

	indicator.rsiPeriod = rsiPeriod
	indicator.stochPeriod = stochPeriod
	indicator.smoothK = smoothK
	indicator.maType = maType
	indicator.offset = offset

	prices := GetPrices(*ctx.Data, field)

	// Check if we have enough data points
	totalPeriod := rsiPeriod + stochPeriod + smoothK
	if len(prices) < totalPeriod+offset {
		return 0.0
	}

	prices = indicator.SliceData(prices)

	// Need at least enough data points to calculate
	if len(prices) < totalPeriod {
		return 0.0
	}

	// Convert maType string to talib MAType using utility function
	talibMAType := util.ConvertMATypeToTalib(maType)

	// Use go-talib StochRsi function
	// StochRsi returns (%K, %D) - we want %K
	stochRsiK, _ := talib.StochRsi(prices, rsiPeriod, stochPeriod, smoothK, talibMAType)

	if len(stochRsiK) == 0 {
		return 0.0
	}

	// Return the last calculated value
	return util.Round(stochRsiK[len(stochRsiK)-1])
}

func (indicator *StochRsiK) GetDataRange(config map[string]any) int32 {
	rsiPeriod := indicator.rsiPeriod
	if rsiPeriod == 0 {
		if p, ok := config["rsiPeriod"].(float64); ok {
			rsiPeriod = int(p)
		} else {
			rsiPeriod = 14 // default RSI period
		}
	}

	stochPeriod := indicator.stochPeriod
	if stochPeriod == 0 {
		if p, ok := config["stochPeriod"].(float64); ok {
			stochPeriod = int(p)
		} else {
			stochPeriod = 14 // default stochastic period
		}
	}

	smoothK := indicator.smoothK
	if smoothK == 0 {
		if p, ok := config["smoothK"].(float64); ok {
			smoothK = int(p)
		} else {
			smoothK = 3 // default smooth K period
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

	// StochRsiK needs data for RSI calculation plus stochastic calculation
	totalPeriod := rsiPeriod + stochPeriod + smoothK
	return int32(totalPeriod*2 + offset)
}

// GetLookBack returns the number of historical data points needed for StochRsiK calculation
func (indicator *StochRsiK) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	rsiPeriod := int(config[1].(float64))   // default: 14
	stochPeriod := int(config[2].(float64)) // default: 14
	smoothK := int(config[3].(float64))     // default: 3
	// maType := config[4].(string)          // not needed for lookback calculation
	offset := int(config[5].(float64)) // default: 0

	// For StochRsiK calculation, we need:
	// - RSI period data points
	// - Plus stochastic period for the stochastic calculation on RSI
	// - Plus smoothing period
	// - Multiplied by a factor for stabilization
	totalPeriod := rsiPeriod + stochPeriod + smoothK
	return int64(5*totalPeriod + offset)
}

// CalculateAll calculates StochRsiK values for all data points
func (indicator *StochRsiK) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	// Parse config parameters
	field := config[0].(string)
	rsiPeriod := int(config[1].(float64))
	stochPeriod := int(config[2].(float64))
	smoothK := int(config[3].(float64))
	maType := config[4].(string)
	offset := int(config[5].(float64))

	indicator.rsiPeriod = rsiPeriod
	indicator.stochPeriod = stochPeriod
	indicator.smoothK = smoothK
	indicator.maType = maType
	indicator.offset = offset

	// Get stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	// Check if we have enough data
	totalPeriod := rsiPeriod + stochPeriod + smoothK
	if len(stockData) < totalPeriod {
		return []IndicatorValues{}
	}

	// Extract price array based on field type
	prices := make([]float64, len(stockData))
	for i, data := range stockData {
		switch strings.ToUpper(field) {
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

	// Use go-talib StochRsi function
	stochRsiK, _ := talib.StochRsi(prices, rsiPeriod, stochPeriod, smoothK, talibMAType)

	// Build results starting from valid calculation points
	var result []IndicatorValues

	for i := totalPeriod - 1; i < len(stockData); i++ {
		// Apply offset
		if i < totalPeriod-1+offset {
			continue
		}

		// Check start date filter
		if stockData[i].Timestamp.Before(startDate) {
			continue
		}

		// Make sure we have valid data from talib
		if i-totalPeriod+1 >= 0 && i-totalPeriod+1 < len(stochRsiK) {
			value := stochRsiK[i-totalPeriod+1]
			// Check for NaN values
			if value == value { // NaN check (NaN != NaN)
				result = append(result, IndicatorValues{
					Timestamp: stockData[i].Timestamp,
					Value:     util.Round(value),
				})
			}
		}
	}

	return result
}
