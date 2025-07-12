package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
)

// Mom represents the Momentum indicator.
type Mom struct {
	period int
	offset int
}

func (indicator *Mom) SliceData(prices Prices) Prices {
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset
	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}
	return prices[start:end]
}

func (indicator *Mom) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	if len(config) != 3 {
		return "", fmt.Errorf("Mom requires 3 parameters: field, period, offset")
	}
	field, ok := config[0].(string)
	if !ok {
		return "", fmt.Errorf("Mom field parameter must be a string")
	}
	periodFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("Mom period parameter must be a number")
	}
	period := int(periodFloat)
	offsetFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("Mom offset parameter must be a number")
	}
	offset := int(offsetFloat)
	return fmt.Sprintf("MOM_%s_%d_%d", strings.ToUpper(field), period, offset), nil
}

func (indicator *Mom) Calculate(
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
	momValues := talib.Mom(prices, period)
	if len(momValues) == 0 {
		return 0.0
	}
	mom := momValues[len(momValues)-1]
	indicatorKey := fmt.Sprintf("MOM_%s_%d_%d", field, period, offset)
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()
	insertValueToResults(ctx, indicatorKey, currentTimestamp, mom)
	return util.Round(mom)
}

func (indicator *Mom) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 10 // default period for Mom
	}
	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}
	return int32(indicator.period + indicator.offset)
}

func (indicator *Mom) GetLookBack(config []any) int64 {
	period := int(config[1].(float64))
	offset := int(config[2].(float64))
	return int64(5*period + offset)
}

func (indicator *Mom) CalculateAll(
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
	momValues := talib.Mom(prices, period)
	if len(momValues) == 0 {
		return []IndicatorValues{}
	}
	var filteredResults []IndicatorValues
	startIndex := period - 1 + offset
	for i, momValue := range momValues {
		dataIndex := i + offset
		if dataIndex >= len(stockData) || momValue != momValue {
			continue
		}
		if !stockData[dataIndex].Timestamp.Before(startDate) && dataIndex >= startIndex {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     util.Round(momValue),
			})
		}
	}
	return filteredResults
}
