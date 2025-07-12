package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
)

// Dema represents the Double Exponential Moving Average indicator.
type Dema struct {
	period int
	offset int
}

func (indicator *Dema) SliceData(prices Prices) Prices {
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset
	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}
	return prices[start:end]
}

func (indicator *Dema) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	if len(config) != 3 {
		return "", fmt.Errorf("DEMA requires 3 parameters: field, period, offset")
	}
	field, ok := config[0].(string)
	if !ok {
		return "", fmt.Errorf("DEMA field parameter must be a string")
	}
	periodFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("DEMA period parameter must be a number")
	}
	period := int(periodFloat)
	offsetFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("DEMA offset parameter must be a number")
	}
	offset := int(offsetFloat)
	if offset != 0 {
		return fmt.Sprintf("DEMA_%s_%d_%d", strings.ToUpper(field), period, offset), nil
	}
	return fmt.Sprintf("DEMA_%s_%d", strings.ToUpper(field), period), nil
}

func (indicator *Dema) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	field := config[0].(string)
	period := int(config[1].(float64))
	offset := int(config[2].(float64))
	indicator.period = period
	indicator.offset = offset
	prices := GetPrices(*ctx.Data, field)
	if len(prices) < period+offset {
		return 0.0
	}
	prices = indicator.SliceData(prices)
	if len(prices) < period {
		return 0.0
	}
	demaValues := talib.Dema(prices, period)
	if len(demaValues) == 0 {
		return 0.0
	}
	dema := demaValues[len(demaValues)-1]
	indicatorKey := fmt.Sprintf("DEMA_%s_%d", field, period)
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()
	insertValueToResults(ctx, indicatorKey, currentTimestamp, dema)
	return util.Round(dema)
}

func (indicator *Dema) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 20 // default period for DEMA
	}
	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}
	return int32(indicator.period + indicator.offset)
}

func (indicator *Dema) GetLookBack(config []any) int64 {
	period := int(config[1].(float64))
	offset := int(config[2].(float64))
	return int64(5*period + offset)
}

func (indicator *Dema) CalculateAll(
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
	if len(stockData) < period+offset {
		return []IndicatorValues{}
	}
	prices := GetPrices(stockData, field)
	demaValues := talib.Dema(prices, period)
	if len(demaValues) == 0 {
		return []IndicatorValues{}
	}
	var filteredResults []IndicatorValues
	startIndex := period - 1 + offset
	for i, demaValue := range demaValues {
		dataIndex := i + offset
		if dataIndex >= len(stockData) || demaValue != demaValue {
			continue
		}
		if !stockData[dataIndex].Timestamp.Before(startDate) && dataIndex >= startIndex {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     util.Round(demaValue),
			})
		}
	}
	return filteredResults
}
