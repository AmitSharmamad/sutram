package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"time"

	"github.com/markcheno/go-talib"
)

// Cci represents the Commodity Channel Index indicator.
type Cci struct {
	period int
	offset int
}

func (indicator *Cci) SliceData(prices Prices) Prices {
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset
	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}
	return prices[start:end]
}

func (indicator *Cci) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	if len(config) != 2 {
		return "", fmt.Errorf("CCI requires 2 parameters: period, offset")
	}
	periodFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("CCI period parameter must be a number")
	}
	period := int(periodFloat)
	offsetFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("CCI offset parameter must be a number")
	}
	offset := int(offsetFloat)
	return fmt.Sprintf("CCI_%d_%d", period, offset), nil
}

func (indicator *Cci) Calculate(
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

	cciValues := talib.Cci(high, low, close, period)
	if len(cciValues) == 0 {
		return 0.0
	}
	cci := cciValues[len(cciValues)-1]
	indicatorKey := fmt.Sprintf("CCI_%d", period)
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()
	insertValueToResults(ctx, indicatorKey, currentTimestamp, cci)
	return util.Round(cci)
}

func (indicator *Cci) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 20 // default period for CCI
	}
	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}
	return int32(indicator.period + indicator.offset)
}

func (indicator *Cci) GetLookBack(config []any) int64 {
	period := int(config[0].(float64))
	offset := int(config[1].(float64))
	return int64(5*period + offset)
}

func (indicator *Cci) CalculateAll(
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
	cciValues := talib.Cci(high, low, close, period)
	if len(cciValues) == 0 {
		return []IndicatorValues{}
	}
	var filteredResults []IndicatorValues
	startIndex := period - 1 + offset
	for i, cciValue := range cciValues {
		dataIndex := i + offset
		if dataIndex >= len(stockData) || cciValue != cciValue {
			continue
		}
		if !stockData[dataIndex].Timestamp.Before(startDate) && dataIndex >= startIndex {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     util.Round(cciValue),
			})
		}
	}
	return filteredResults
}
