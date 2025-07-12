package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
)

// Mama represents the MESA Adaptive Moving Average indicator.
type Mama struct {
	fastLimit float64
	slowLimit float64
	offset    int
}

func (indicator *Mama) SliceData(prices Prices) Prices {
	// MAMA needs at least 2*len(prices) for stability, but talib doesn't document a strict minimum
	start := max(0, len(prices)-(indicator.offset+20))
	end := len(prices) - indicator.offset
	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}
	return prices[start:end]
}

func (indicator *Mama) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	if len(config) != 4 {
		return "", fmt.Errorf("MAMA requires 4 parameters: field, fastLimit, slowLimit, offset")
	}
	field, ok := config[0].(string)
	if !ok {
		return "", fmt.Errorf("MAMA field parameter must be a string")
	}
	fastLimit, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("MAMA fastLimit parameter must be a number")
	}
	slowLimit, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("MAMA slowLimit parameter must be a number")
	}
	offsetFloat, ok := config[3].(float64)
	if !ok {
		return "", fmt.Errorf("MAMA offset parameter must be a number")
	}
	offset := int(offsetFloat)
	if offset != 0 {
		return fmt.Sprintf("MAMA_%s_%.2f_%.2f_%d", strings.ToUpper(field), fastLimit, slowLimit, offset), nil
	}
	return fmt.Sprintf("MAMA_%s_%.2f_%.2f", strings.ToUpper(field), fastLimit, slowLimit), nil
}

func (indicator *Mama) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	field := config[0].(string)
	fastLimit := config[1].(float64)
	slowLimit := config[2].(float64)
	offset := int(config[3].(float64))
	indicator.fastLimit = fastLimit
	indicator.slowLimit = slowLimit
	indicator.offset = offset
	prices := GetPrices(*ctx.Data, field)
	if len(prices) <= offset {
		return 0.0
	}
	prices = indicator.SliceData(prices)
	if len(prices) == 0 {
		return 0.0
	}
	mama, _ := talib.Mama(prices, fastLimit, slowLimit)
	if len(mama) == 0 {
		return 0.0
	}
	value := mama[len(mama)-1]
	indicatorKey := fmt.Sprintf("MAMA_%s_%.2f_%.2f", field, fastLimit, slowLimit)
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()
	insertValueToResults(ctx, indicatorKey, currentTimestamp, value)
	return util.Round(value)
}

func (indicator *Mama) GetDataRange(config map[string]any) int32 {
	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}
	return int32(20 + indicator.offset)
}

func (indicator *Mama) GetLookBack(config []any) int64 {
	offset := int(config[3].(float64))
	return int64(20 + offset)
}

func (indicator *Mama) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	field := config[0].(string)
	fastLimit := config[1].(float64)
	slowLimit := config[2].(float64)
	offset := int(config[3].(float64))
	indicator.fastLimit = fastLimit
	indicator.slowLimit = slowLimit
	indicator.offset = offset
	stockData := *ctx.Data
	if len(stockData) <= offset {
		return []IndicatorValues{}
	}
	prices := GetPrices(stockData, field)
	mama, _ := talib.Mama(prices, fastLimit, slowLimit)
	if len(mama) == 0 {
		return []IndicatorValues{}
	}
	var filteredResults []IndicatorValues
	for i, value := range mama {
		dataIndex := i + offset
		if dataIndex >= len(stockData) || value != value {
			continue
		}
		if !stockData[dataIndex].Timestamp.Before(startDate) {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     util.Round(value),
			})
		}
	}
	return filteredResults
}
