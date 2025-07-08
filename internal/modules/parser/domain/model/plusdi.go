package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"time"

	"github.com/markcheno/go-talib"
)

// PlusDI represents the Plus Directional Indicator.
type PlusDI struct {
	period int
	offset int
}

// SliceData slices the price data for the PlusDI indicator.
func (indicator *PlusDI) SliceData(prices Prices) Prices {
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset
	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}
	return prices[start:end]
}

// GenerateIndicatorKey generates the key for the PlusDI indicator.
func (indicator *PlusDI) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	if len(config) != 2 {
		return "", fmt.Errorf("PlusDI requires 2 parameters: period, offset")
	}
	periodFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("PlusDI period parameter must be a number")
	}
	period := int(periodFloat)
	offsetFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("PlusDI offset parameter must be a number")
	}
	offset := int(offsetFloat)
	return fmt.Sprintf("PLUSDI_%d_%d", period, offset), nil
}

// Calculate calculates the PlusDI indicator.
func (indicator *PlusDI) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	period := int(config[0].(float64))
	offset := int(config[1].(float64))
	indicator.period = period
	indicator.offset = offset
	data := *ctx.Data
	if len(data) < period+offset {
		return 0.0
	}
	high := GetPrices(data, "HIGH")
	low := GetPrices(data, "LOW")
	close := GetPrices(data, "CLOSE")

	high = indicator.SliceData(high)
	low = indicator.SliceData(low)
	close = indicator.SliceData(close)

	if len(high) < period || len(low) < period || len(close) < period {
		return 0.0
	}

	plusDIValues := talib.PlusDI(high, low, close, period)
	if len(plusDIValues) == 0 {
		return 0.0
	}
	plusDI := plusDIValues[len(plusDIValues)-1]
	indicatorKey := fmt.Sprintf("PLUSDI_%d", period)
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()
	insertValueToResults(ctx, indicatorKey, currentTimestamp, plusDI)
	return util.Round(plusDI)
}

// GetDataRange gets the data range for the PlusDI indicator.
func (indicator *PlusDI) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 14 // default period for PlusDI
	}
	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}
	return int32(indicator.period + indicator.offset)
}

// GetLookBack gets the lookback period for the PlusDI indicator.
func (indicator *PlusDI) GetLookBack(config []any) int64 {
	period := int(config[0].(float64))
	offset := int(config[1].(float64))
	return int64(5*period + offset)
}

// CalculateAll calculates the PlusDI indicator for all applicable data points.
func (indicator *PlusDI) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	period := int(config[0].(float64))
	offset := int(config[1].(float64))
	indicator.period = period
	indicator.offset = offset
	stockData := *ctx.Data
	if len(stockData) < period+offset {
		return []IndicatorValues{}
	}
	high := GetPrices(stockData, "HIGH")
	low := GetPrices(stockData, "LOW")
	close := GetPrices(stockData, "CLOSE")
	plusDIValues := talib.PlusDI(high, low, close, period)
	if len(plusDIValues) == 0 {
		return []IndicatorValues{}
	}
	var filteredResults []IndicatorValues
	startIndex := period - 1 + offset
	for i, plusDIValue := range plusDIValues {
		dataIndex := i + offset
		if dataIndex >= len(stockData) || plusDIValue != plusDIValue {
			continue
		}
		if !stockData[dataIndex].Timestamp.Before(startDate) && dataIndex >= startIndex {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     util.Round(plusDIValue),
			})
		}
	}
	return filteredResults
}
