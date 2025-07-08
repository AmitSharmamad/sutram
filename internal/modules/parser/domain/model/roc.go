package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
)

// Roc represents the Rate of Change indicator.
type Roc struct {
	period int
	offset int
}

func (indicator *Roc) SliceData(prices Prices) Prices {
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset
	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}
	return prices[start:end]
}

func (indicator *Roc) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	if len(config) != 3 {
		return "", fmt.Errorf("ROC requires 3 parameters: field, period, offset")
	}
	field, ok := config[0].(string)
	if !ok {
		return "", fmt.Errorf("ROC field parameter must be a string")
	}
	periodFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("ROC period parameter must be a number")
	}
	period := int(periodFloat)
	offsetFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("ROC offset parameter must be a number")
	}
	offset := int(offsetFloat)
	return fmt.Sprintf("ROC_%s_%d_%d", strings.ToUpper(field), period, offset), nil
}

func (indicator *Roc) Calculate(
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
	rocValues := talib.Roc(prices, period)
	if len(rocValues) == 0 {
		return 0.0
	}
	roc := rocValues[len(rocValues)-1]
	indicatorKey := fmt.Sprintf("ROC_%s_%d", field, period)
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()
	insertValueToResults(ctx, indicatorKey, currentTimestamp, roc)
	return util.Round(roc)
}

func (indicator *Roc) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 10 // default period for ROC
	}
	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}
	return int32(indicator.period + indicator.offset)
}

func (indicator *Roc) GetLookBack(config []any) int64 {
	period := int(config[1].(float64))
	offset := int(config[2].(float64))
	return int64(5*period + offset)
}

func (indicator *Roc) CalculateAll(
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
	rocValues := talib.Roc(prices, period)
	if len(rocValues) == 0 {
		return []IndicatorValues{}
	}
	var filteredResults []IndicatorValues
	startIndex := period - 1 + offset
	for i, rocValue := range rocValues {
		dataIndex := i + offset
		if dataIndex >= len(stockData) || rocValue != rocValue {
			continue
		}
		if !stockData[dataIndex].Timestamp.Before(startDate) && dataIndex >= startIndex {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     util.Round(rocValue),
			})
		}
	}
	return filteredResults
}
