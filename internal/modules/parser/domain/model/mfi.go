package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"strings"
	"time"
)

type MFI struct {
	period   int
	maType   string
	maPeriod int
	offset   int
}

func (indicator *MFI) SliceData(prices Prices) Prices {
	// We need enough data to calculate MFI properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// generateIndicatorKey generates a unique key for MFI indicator based on its configuration
func (indicator *MFI) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 4 {
		return "", fmt.Errorf("MFI requires 4 parameters: period, maType, maPeriod, offset")
	}

	// Extract period parameter
	periodFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("MFI period parameter must be a number")
	}
	period := int(periodFloat)

	// Extract maType parameter
	maType, ok := config[1].(string)
	if !ok {
		return "", fmt.Errorf("MFI maType parameter must be a string")
	}

	// Extract maPeriod parameter
	maPeriodFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("MFI maPeriod parameter must be a number")
	}
	maPeriod := int(maPeriodFloat)

	// Extract offset parameter
	offsetFloat, ok := config[3].(float64)
	if !ok {
		return "", fmt.Errorf("MFI offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include non-default values
	baseKey := fmt.Sprintf("MFI_%d_%s_%d", period, strings.ToUpper(maType), maPeriod)

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

// Helper method to slice all required data
func (indicator *MFI) sliceAllData(highs, lows, closes, volumes Prices) (Prices, Prices, Prices, Prices) {
	// We need enough data to calculate MFI properly
	start := max(0, len(closes)-(indicator.offset+2*indicator.period))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end], volumes[start:end]
}

func (indicator *MFI) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	period := int(config[0].(float64))
	maType := config[1].(string)
	maPeriod := int(config[2].(float64))
	offset := int(config[3].(float64))

	indicator.period = period
	indicator.maType = maType
	indicator.maPeriod = maPeriod
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")
	volumes := GetPrices(*ctx.Data, "volume")

	// Check if we have enough data points
	if len(highs) < period+offset+1 || len(lows) < period+offset+1 ||
		len(closes) < period+offset+1 || len(volumes) < period+offset+1 {
		return 0.0
	}

	highs, lows, closes, volumes = indicator.sliceAllData(highs, lows, closes, volumes)

	// Need at least period+1 data points to calculate
	if len(highs) < period+1 || len(lows) < period+1 ||
		len(closes) < period+1 || len(volumes) < period+1 {
		return 0.0
	}

	// Calculate typical price
	var typicalPrices = make([]float64, len(closes))
	for i := 0; i < len(closes); i++ {
		typicalPrices[i] = (highs[i] + lows[i] + closes[i]) / 3.0
	}

	// Calculate raw money flow
	var rawMoneyFlow = make([]float64, len(closes))
	for i := 0; i < len(closes); i++ {
		rawMoneyFlow[i] = typicalPrices[i] * volumes[i]
	}

	// Calculate positive and negative money flow
	var positiveMoneyFlow = make([]float64, len(closes))
	var negativeMoneyFlow = make([]float64, len(closes))

	for i := 1; i < len(closes); i++ {
		if typicalPrices[i] > typicalPrices[i-1] {
			positiveMoneyFlow[i] = rawMoneyFlow[i]
			negativeMoneyFlow[i] = 0.0
		} else if typicalPrices[i] < typicalPrices[i-1] {
			positiveMoneyFlow[i] = 0.0
			negativeMoneyFlow[i] = rawMoneyFlow[i]
		} else {
			positiveMoneyFlow[i] = 0.0
			negativeMoneyFlow[i] = 0.0
		}
	}

	// Calculate money flow ratio and MFI
	var mfi = make([]float64, len(closes))

	for i := period; i < len(closes); i++ {
		var sumPositiveFlow float64
		var sumNegativeFlow float64

		for j := i - period + 1; j <= i; j++ {
			sumPositiveFlow += positiveMoneyFlow[j]
			sumNegativeFlow += negativeMoneyFlow[j]
		}

		// Avoid division by zero
		if sumNegativeFlow == 0 {
			mfi[i] = 100.0
		} else {
			moneyFlowRatio := sumPositiveFlow / sumNegativeFlow
			mfi[i] = 100.0 - (100.0 / (1.0 + moneyFlowRatio))
		}
	}

	// Apply moving average smoothing if maPeriod > 0
	if indicator.maPeriod > 0 {
		// Extract valid MFI values (skip initial zeros)
		validMfi := mfi[period:]

		// Apply the specified moving average type
		smoothedMfi := util.CalculateMovingAverage(validMfi, indicator.maType, indicator.maPeriod)

		// Return the last value from the smoothed MFI
		if len(smoothedMfi) > 0 {
			return util.Round(smoothedMfi[len(smoothedMfi)-1])
		}
	}

	// Return the last calculated MFI value if no smoothing applied
	return util.Round(mfi[len(mfi)-1])
}

func (indicator *MFI) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 20 // default period
	}

	if t, ok := config["type"].(string); ok {
		indicator.maType = t
	} else {
		indicator.maType = "SIMPLE" // default type
	}

	if mp, ok := config["maPeriod"].(float64); ok {
		indicator.maPeriod = int(mp)
	} else {
		indicator.maPeriod = 20 // default MA period
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// We need period+1 data points (due to the use of previous values) + offset
	return int32(indicator.period + 1 + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for MFI calculation
// config[0] = period (float64), config[1] = maType (string), config[2] = maPeriod (float64), config[3] = offset (float64)
func (indicator *MFI) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	period := int(config[0].(float64))
	// config[1] is maType (string) - not needed for lookback calculation
	// config[2] is maPeriod (float64) - not needed for lookback calculation
	offset := int(config[3].(float64))

	// For MFI calculation, we need:
	// - Period data points for money flow calculation
	// - 1 additional for typical price comparison
	// - Additional buffer for stability (using 2*period like SliceData)
	// - Plus any offset
	return int64(2*period + 1 + offset)
}

// CalculateAll calculates MFI values for all data points
func (indicator *MFI) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	// Parse config parameters
	period := int(config[0].(float64))
	maType := config[1].(string)
	maPeriod := int(config[2].(float64))
	offset := int(config[3].(float64))

	indicator.period = period
	indicator.maType = maType
	indicator.maPeriod = maPeriod
	indicator.offset = offset

	// Get stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	// Check if we have enough data
	if len(stockData) < period+1 {
		return []IndicatorValues{}
	}

	// Extract price arrays
	highs := make([]float64, len(stockData))
	lows := make([]float64, len(stockData))
	closes := make([]float64, len(stockData))
	volumes := make([]float64, len(stockData))

	for i, data := range stockData {
		highs[i] = data.High
		lows[i] = data.Low
		closes[i] = data.Close
		volumes[i] = data.Volume
	}

	// Step 1: Calculate Typical Price = (High + Low + Close) / 3
	typicalPrices := make([]float64, len(stockData))
	for i := 0; i < len(stockData); i++ {
		typicalPrices[i] = (highs[i] + lows[i] + closes[i]) / 3.0
	}

	// Step 2: Calculate Raw Money Flow = Typical Price × Volume
	rawMoneyFlow := make([]float64, len(stockData))
	for i := 0; i < len(stockData); i++ {
		rawMoneyFlow[i] = typicalPrices[i] * volumes[i]
	}

	// Step 3: Determine Positive and Negative Money Flow
	positiveMoneyFlow := make([]float64, len(stockData))
	negativeMoneyFlow := make([]float64, len(stockData))

	for i := 1; i < len(stockData); i++ {
		if typicalPrices[i] > typicalPrices[i-1] {
			positiveMoneyFlow[i] = rawMoneyFlow[i]
			negativeMoneyFlow[i] = 0.0
		} else if typicalPrices[i] < typicalPrices[i-1] {
			positiveMoneyFlow[i] = 0.0
			negativeMoneyFlow[i] = rawMoneyFlow[i]
		} else {
			positiveMoneyFlow[i] = 0.0
			negativeMoneyFlow[i] = 0.0
		}
	}

	// Step 4: Calculate MFI for each valid period
	mfi := make([]float64, len(stockData))

	for i := period; i < len(stockData); i++ {
		var sumPositiveFlow float64
		var sumNegativeFlow float64

		// Sum positive and negative money flows for the period
		for j := i - period + 1; j <= i; j++ {
			sumPositiveFlow += positiveMoneyFlow[j]
			sumNegativeFlow += negativeMoneyFlow[j]
		}

		// Calculate Money Flow Ratio and MFI
		if sumNegativeFlow == 0 {
			mfi[i] = 100.0
		} else {
			moneyFlowRatio := sumPositiveFlow / sumNegativeFlow
			mfi[i] = 100.0 - (100.0 / (1.0 + moneyFlowRatio))
		}
	}

	// Step 5: Apply moving average smoothing if maPeriod > 0
	var finalMFI []float64
	if maPeriod > 0 {
		// Extract valid MFI values (skip initial zeros)
		validMfi := make([]float64, 0)
		for i := period; i < len(stockData); i++ {
			validMfi = append(validMfi, mfi[i])
		}

		// Apply the specified moving average type using utility function
		smoothedMfi := util.CalculateMovingAverage(validMfi, maType, maPeriod)

		// Create full-length array with proper positioning
		finalMFI = make([]float64, len(stockData))
		startIdx := period + maPeriod - 1
		for i, val := range smoothedMfi {
			if startIdx+i < len(finalMFI) {
				finalMFI[startIdx+i] = val
			}
		}
	} else {
		finalMFI = mfi
	}

	// Step 6: Build results starting from valid calculation points
	var result []IndicatorValues
	startIdx := period
	if maPeriod > 0 {
		startIdx = period + maPeriod - 1
	}

	for i := startIdx; i < len(stockData); i++ {
		// Apply offset
		if i < startIdx+offset {
			continue
		}

		// Check start date filter
		if stockData[i].Timestamp.Before(startDate) {
			continue
		}

		result = append(result, IndicatorValues{
			Timestamp: stockData[i].Timestamp,
			Value:     util.Round(finalMFI[i]),
		})
	}

	return result
}
