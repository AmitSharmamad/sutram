package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"time"

	"github.com/markcheno/go-talib"
)

// AvgPrice represents the Average Price indicator.
type AvgPrice struct {
	offset int
}

func (indicator *AvgPrice) SliceData(prices Prices) Prices {
	start := max(0, len(prices)-(indicator.offset+1))
	end := len(prices) - indicator.offset
	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}
	return prices[start:end]
}

func (indicator *AvgPrice) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	if len(config) != 1 {
		return "", fmt.Errorf("AvgPrice requires 1 parameter: offset")
	}
	offsetFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("AvgPrice offset parameter must be a number")
	}
	offset := int(offsetFloat)
	if offset != 0 {
		return fmt.Sprintf("AVGPRICE_%d", offset), nil
	}
	return "AVGPRICE", nil
}

func (indicator *AvgPrice) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	if len(config) != 1 {
		panic("AvgPrice requires 1 parameter: offset")
	}
	offset := int(config[0].(float64))
	indicator.offset = offset
	data := *ctx.Data
	if len(data) <= offset {
		return 0.0
	}
	high := GetPrices(data, "HIGH")
	low := GetPrices(data, "LOW")
	open := GetPrices(data, "OPEN")
	close := GetPrices(data, "CLOSE")

	high = indicator.SliceData(high)
	low = indicator.SliceData(low)
	open = indicator.SliceData(open)
	close = indicator.SliceData(close)

	if len(high) == 0 || len(low) == 0 || len(open) == 0 || len(close) == 0 {
		return 0.0
	}

	avgPriceValues := talib.AvgPrice(open, high, low, close)
	if len(avgPriceValues) == 0 {
		return 0.0
	}
	avgPrice := avgPriceValues[len(avgPriceValues)-1]
	indicatorKey := fmt.Sprintf("AVGPRICE_%d", offset)
	dataLen := len(*ctx.Data)
	if dataLen == 0 {
		return 0.0
	}
	currentTimestamp := (*ctx.Data)[dataLen-1].Timestamp.Unix()
	insertValueToResults(ctx, indicatorKey, currentTimestamp, avgPrice)
	return util.Round(avgPrice)
}

func (indicator *AvgPrice) GetDataRange(config map[string]any) int32 {
	if p, ok := config["offset"].(float64); ok {
		indicator.offset = int(p)
	} else {
		indicator.offset = 0
	}
	return int32(1 + indicator.offset)
}

func (indicator *AvgPrice) GetLookBack(config []any) int64 {
	if len(config) != 1 {
		return 1
	}
	offset := int(config[0].(float64))
	return int64(1 + offset)
}

func (indicator *AvgPrice) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	if len(config) != 1 {
		return []IndicatorValues{}
	}
	offset := int(config[0].(float64))
	indicator.offset = offset
	stockData := *ctx.Data
	if len(stockData) <= offset {
		return []IndicatorValues{}
	}
	high := GetPrices(stockData, "HIGH")
	low := GetPrices(stockData, "LOW")
	open := GetPrices(stockData, "OPEN")
	close := GetPrices(stockData, "CLOSE")
	avgPriceValues := talib.AvgPrice(open, high, low, close)
	if len(avgPriceValues) == 0 {
		return []IndicatorValues{}
	}
	var filteredResults []IndicatorValues
	for i, avgPrice := range avgPriceValues {
		dataIndex := i + offset
		if dataIndex >= len(stockData) || avgPrice != avgPrice {
			continue
		}
		if !stockData[dataIndex].Timestamp.Before(startDate) {
			filteredResults = append(filteredResults, IndicatorValues{
				Timestamp: stockData[dataIndex].Timestamp,
				Value:     util.Round(avgPrice),
			})
		}
	}
	return filteredResults
}
