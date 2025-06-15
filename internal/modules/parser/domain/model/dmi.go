package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"math"
	"strings"
	"time"
)

type DMI struct {
	period int
	offset int
}

func (indicator *DMI) SliceData(prices Prices) Prices {
	// We need enough data to calculate DMI properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// generateIndicatorKey generates a unique key for DMI indicator based on its configuration
func (indicator *DMI) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 3 {
		return "", fmt.Errorf("DMI requires 3 parameters: point, period, offset")
	}

	// Extract point parameter
	point, ok := config[0].(string)
	if !ok {
		return "", fmt.Errorf("DMI point parameter must be a string")
	}

	// Extract period parameter
	periodFloat, ok := config[1].(float64)
	if !ok {
		return "", fmt.Errorf("DMI period parameter must be a number")
	}
	period := int(periodFloat)

	// Extract offset parameter
	offsetFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("DMI offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include offset only if non-zero
	baseKey := fmt.Sprintf("DMI_%s_%d", strings.ToUpper(point), period)

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

// Helper method to slice all required data
func (indicator *DMI) sliceAllData(highs, lows, closes Prices) (Prices, Prices, Prices) {
	// We need enough data to calculate DMI properly
	start := max(0, len(closes)-(indicator.offset+2*indicator.period))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end]
}

func (indicator *DMI) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	point := config[0].(string) // Point parameter - DMI always uses High/Low/Close
	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	// Note: DMI calculation doesn't use the point parameter as it always uses High, Low, Close
	_ = point // Acknowledge the parameter but don't use it

	indicator.period = period
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")

	// Check if we have enough data points
	if len(highs) < period+offset+1 || len(lows) < period+offset+1 || len(closes) < period+offset+1 {
		return 0.0
	}

	highs, lows, closes = indicator.sliceAllData(highs, lows, closes)

	// Need at least period+1 data points to calculate
	if len(highs) < period+1 || len(lows) < period+1 || len(closes) < period+1 {
		return 0.0
	}

	// Calculate +DM, -DM, and TR (True Range)
	var posDM = make([]float64, len(closes))
	var negDM = make([]float64, len(closes))
	var tr = make([]float64, len(closes))

	for i := 1; i < len(closes); i++ {
		// Calculate current high - previous high
		upMove := highs[i] - highs[i-1]
		// Calculate previous low - current low
		downMove := lows[i-1] - lows[i]

		// Calculate +DM and -DM
		if upMove > downMove && upMove > 0 {
			posDM[i] = upMove
		} else {
			posDM[i] = 0
		}

		if downMove > upMove && downMove > 0 {
			negDM[i] = downMove
		} else {
			negDM[i] = 0
		}

		// Calculate True Range
		tr[i] = max(
			highs[i]-lows[i],
			max(
				math.Abs(highs[i]-closes[i-1]),
				math.Abs(lows[i]-closes[i-1]),
			),
		)
	}

	// Calculate smoothed +DM, -DM, and TR
	var smoothPosDM = make([]float64, len(closes))
	var smoothNegDM = make([]float64, len(closes))
	var smoothTR = make([]float64, len(closes))

	// Initial values for smoothed series (first period sum)
	var sumPosDM, sumNegDM, sumTR float64
	for i := 1; i <= period; i++ {
		sumPosDM += posDM[i]
		sumNegDM += negDM[i]
		sumTR += tr[i]
	}
	smoothPosDM[period] = sumPosDM
	smoothNegDM[period] = sumNegDM
	smoothTR[period] = sumTR

	// Calculate subsequent smoothed values
	for i := period + 1; i < len(closes); i++ {
		smoothPosDM[i] = smoothPosDM[i-1] - (smoothPosDM[i-1] / float64(period)) + posDM[i]
		smoothNegDM[i] = smoothNegDM[i-1] - (smoothNegDM[i-1] / float64(period)) + negDM[i]
		smoothTR[i] = smoothTR[i-1] - (smoothTR[i-1] / float64(period)) + tr[i]
	}

	// Calculate +DI and -DI
	var posDI = make([]float64, len(closes))
	var negDI = make([]float64, len(closes))

	for i := period; i < len(closes); i++ {
		if smoothTR[i] > 0 {
			posDI[i] = 100 * (smoothPosDM[i] / smoothTR[i])
			negDI[i] = 100 * (smoothNegDM[i] / smoothTR[i])
		}
	}

	// Calculate DX (Directional Index)
	var dx = make([]float64, len(closes))
	for i := period; i < len(closes); i++ {
		if posDI[i]+negDI[i] > 0 {
			dx[i] = 100 * (math.Abs(posDI[i]-negDI[i]) / (posDI[i] + negDI[i]))
		}
	}

	// Calculate ADX (Average Directional Index)
	var adx = make([]float64, len(closes))

	// Initial ADX value as simple average of DX
	var sumDX float64
	for i := period; i < 2*period; i++ {
		sumDX += dx[i]
	}
	adx[2*period-1] = sumDX / float64(period)

	// Calculate subsequent ADX values
	for i := 2 * period; i < len(closes); i++ {
		adx[i] = ((adx[i-1] * float64(period-1)) + dx[i]) / float64(period)
	}

	// Return the last calculated ADX value
	return util.Round(adx[len(adx)-1])
}

func (indicator *DMI) GetDataRange(config map[string]any) int32 {
	if p, ok := config["period"].(float64); ok {
		indicator.period = int(p)
	} else {
		indicator.period = 14 // default period
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// We need 2*period data points for ADX calculation + 1 for previous values + offset
	return int32(2*indicator.period + 1 + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for DMI calculation
// config[0] = point (string), config[1] = period (float64), config[2] = offset (float64)
func (indicator *DMI) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	// config[0] is point (string) - not needed for lookback calculation
	period := int(config[1].(float64)) // default: 14
	offset := int(config[2].(float64)) // default: 0

	// For DMI/ADX calculation, we need:
	// - Previous values for directional movement calculation (1 extra)
	// - Period for initial smoothed values
	// - Another period for ADX calculation from DX values
	// - Plus any offset
	// Using the same calculation as GetDataRange: 2*period + 1 + offset
	return int64(5*period + 1 + offset)
}

// CalculateAll calculates DMI/ADX values for all data points
func (indicator *DMI) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	// Parse config parameters
	point := config[0].(string) // Point parameter - DMI always uses High/Low/Close
	period := int(config[1].(float64))
	offset := int(config[2].(float64))

	// Note: DMI calculation doesn't use the point parameter as it always uses High, Low, Close
	_ = point // Acknowledge the parameter but don't use it

	indicator.period = period
	indicator.offset = offset

	// Get stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	// Check if we have enough data (need 2*period + 1 for ADX calculation)
	if len(stockData) < 2*period+1 {
		return []IndicatorValues{}
	}

	// Extract price arrays
	highs := make([]float64, len(stockData))
	lows := make([]float64, len(stockData))
	closes := make([]float64, len(stockData))

	for i, data := range stockData {
		highs[i] = data.High
		lows[i] = data.Low
		closes[i] = data.Close
	}

	// Step 1: Calculate +DM, -DM, and TR (True Range)
	posDM := make([]float64, len(stockData))
	negDM := make([]float64, len(stockData))
	tr := make([]float64, len(stockData))

	for i := 1; i < len(stockData); i++ {
		// Calculate directional movements
		upMove := highs[i] - highs[i-1]
		downMove := lows[i-1] - lows[i]

		// Calculate +DM and -DM
		if upMove > downMove && upMove > 0 {
			posDM[i] = upMove
		} else {
			posDM[i] = 0
		}

		if downMove > upMove && downMove > 0 {
			negDM[i] = downMove
		} else {
			negDM[i] = 0
		}

		// Calculate True Range
		hl := highs[i] - lows[i]
		hc := math.Abs(highs[i] - closes[i-1])
		lc := math.Abs(lows[i] - closes[i-1])

		tr[i] = math.Max(hl, math.Max(hc, lc))
	}

	// Step 2: Calculate smoothed +DM, -DM, and TR using Wilder's smoothing
	smoothPosDM := make([]float64, len(stockData))
	smoothNegDM := make([]float64, len(stockData))
	smoothTR := make([]float64, len(stockData))

	// Initial values for smoothed series (first period sum)
	var sumPosDM, sumNegDM, sumTR float64
	for i := 1; i <= period; i++ {
		sumPosDM += posDM[i]
		sumNegDM += negDM[i]
		sumTR += tr[i]
	}
	smoothPosDM[period] = sumPosDM
	smoothNegDM[period] = sumNegDM
	smoothTR[period] = sumTR

	// Calculate subsequent smoothed values using Wilder's smoothing
	for i := period + 1; i < len(stockData); i++ {
		smoothPosDM[i] = smoothPosDM[i-1] - (smoothPosDM[i-1] / float64(period)) + posDM[i]
		smoothNegDM[i] = smoothNegDM[i-1] - (smoothNegDM[i-1] / float64(period)) + negDM[i]
		smoothTR[i] = smoothTR[i-1] - (smoothTR[i-1] / float64(period)) + tr[i]
	}

	// Step 3: Calculate +DI and -DI
	posDI := make([]float64, len(stockData))
	negDI := make([]float64, len(stockData))

	for i := period; i < len(stockData); i++ {
		if smoothTR[i] > 0 {
			posDI[i] = 100 * (smoothPosDM[i] / smoothTR[i])
			negDI[i] = 100 * (smoothNegDM[i] / smoothTR[i])
		}
	}

	// Step 4: Calculate DX (Directional Index)
	dx := make([]float64, len(stockData))
	for i := period; i < len(stockData); i++ {
		diSum := posDI[i] + negDI[i]
		if diSum > 0 {
			dx[i] = 100 * (math.Abs(posDI[i]-negDI[i]) / diSum)
		}
	}

	// Step 5: Calculate ADX (Average Directional Index)
	adx := make([]float64, len(stockData))

	// Initial ADX value as simple average of DX for the first period
	var sumDX float64
	for i := period; i < 2*period; i++ {
		sumDX += dx[i]
	}
	if 2*period-1 < len(stockData) {
		adx[2*period-1] = sumDX / float64(period)
	}

	// Calculate subsequent ADX values using Wilder's smoothing
	for i := 2 * period; i < len(stockData); i++ {
		adx[i] = ((adx[i-1] * float64(period-1)) + dx[i]) / float64(period)
	}

	// Step 6: Build results starting from valid calculation points
	var result []IndicatorValues

	for i := 2*period - 1; i < len(stockData); i++ {
		// Apply offset
		if i < 2*period-1+offset {
			continue
		}

		// Check start date filter
		if stockData[i].Timestamp.Before(startDate) {
			continue
		}

		result = append(result, IndicatorValues{
			Timestamp: stockData[i].Timestamp,
			Value:     util.Round(adx[i]),
		})
	}

	return result
}
