package parser

import (
	util "algo/internal/util/math"
	"fmt"
	"math"
	"strings"
	"time"
)

type Vortex struct {
	period int
	offset int
}

func (indicator *Vortex) SliceData(prices Prices) Prices {
	// We need enough data to calculate Vortex indicator properly
	start := max(0, len(prices)-(indicator.offset+2*indicator.period))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// generateIndicatorKey generates a unique key for Vortex indicator based on its configuration
func (indicator *Vortex) GenerateIndicatorKey(
	config []any,
	ctx *EvaluationContext,
) (string, error) {
	// Validate we have the required number of parameters
	if len(config) != 3 {
		return "", fmt.Errorf("Vortex requires 3 parameters: period, signal, offset")
	}

	// Extract period parameter
	periodFloat, ok := config[0].(float64)
	if !ok {
		return "", fmt.Errorf("Vortex period parameter must be a number")
	}
	period := int(periodFloat)

	// Extract signal parameter
	signal, ok := config[1].(string)
	if !ok {
		return "", fmt.Errorf("Vortex signal parameter must be a string")
	}

	// Extract offset parameter
	offsetFloat, ok := config[2].(float64)
	if !ok {
		return "", fmt.Errorf("Vortex offset parameter must be a number")
	}
	offset := int(offsetFloat)

	// Generate the key - include non-default values
	baseKey := fmt.Sprintf("VORTEX_%d_%s", period, strings.ToUpper(signal))

	if offset != 0 {
		baseKey += fmt.Sprintf("_%d", offset)
	}

	return baseKey, nil
}

// Internal helper method to slice all required data
func (indicator *Vortex) sliceAllData(highs, lows, closes Prices) (Prices, Prices, Prices) {
	// We need enough data to calculate Vortex indicator properly
	start := max(0, len(closes)-(indicator.offset+2*indicator.period))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end]
}

func (indicator *Vortex) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	period := int(config[0].(float64)) // default: 14
	signal := config[1].(string)       // default: "POSITIVE_VI"
	offset := int(config[2].(float64)) // default: 0

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

	// Calculate VM+ and VM- for each bar
	var vmPlus = make([]float64, len(closes))
	var vmMinus = make([]float64, len(closes))
	var tr = make([]float64, len(closes))

	for i := 1; i < len(closes); i++ {
		// VM+ = |Current High - Prior Low|
		vmPlus[i] = math.Abs(highs[i] - lows[i-1])

		// VM- = |Current Low - Prior High|
		vmMinus[i] = math.Abs(lows[i] - highs[i-1])

		// Calculate True Range
		tr[i] = max(
			highs[i]-lows[i],
			max(
				math.Abs(highs[i]-closes[i-1]),
				math.Abs(lows[i]-closes[i-1]),
			),
		)
	}

	// Calculate Vortex Indicator
	var viPlus = make([]float64, len(closes))
	var viMinus = make([]float64, len(closes))

	for i := period; i < len(closes); i++ {
		// Calculate sum of VM+ for period
		var sumVMPlus float64
		for j := i - period + 1; j <= i; j++ {
			sumVMPlus += vmPlus[j]
		}

		// Calculate sum of VM- for period
		var sumVMMinus float64
		for j := i - period + 1; j <= i; j++ {
			sumVMMinus += vmMinus[j]
		}

		// Calculate sum of TR for period
		var sumTR float64
		for j := i - period + 1; j <= i; j++ {
			sumTR += tr[j]
		}

		// Calculate VI+ and VI-
		viPlus[i] = sumVMPlus / sumTR
		viMinus[i] = sumVMMinus / sumTR
	}

	// Return the appropriate value based on signal parameter
	lastIndex := len(viPlus) - 1
	if signal == "NEGATIVE_VI" {
		return util.Round(viMinus[lastIndex])
	} else { // Default to "POSITIVE_VI" for any other value
		return util.Round(viPlus[lastIndex])
	}
}

func (indicator *Vortex) GetDataRange(config map[string]any) int32 {
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

	// We need period+1 data points (due to the use of previous values) + offset
	return int32(indicator.period + 1 + indicator.offset)
}

// GetLookBack returns the number of historical data points needed for Vortex calculation
// config[0] = period (float64), config[1] = signal (string), config[2] = offset (float64)
func (indicator *Vortex) GetLookBack(config []any) int64 {
	// Parse config parameters same as Calculate function
	period := int(config[0].(float64)) // default: 14
	// config[1] is signal (string) - not needed for lookback calculation
	offset := int(config[2].(float64)) // default: 0

	// For Vortex calculation, we need:
	// - Period data points for VM+/VM- and TR calculation
	// - 1 additional for previous values calculation
	// - Additional buffer for stability (using 2*period like SliceData)
	// - Plus any offset
	return int64(2*period + 1 + offset)
}

// CalculateAll calculates Vortex values for all data points
func (indicator *Vortex) CalculateAll(
	config []any,
	ctx *EvaluationContext,
	startDate time.Time,
) []IndicatorValues {
	// Parse config parameters
	period := int(config[0].(float64))
	signal := config[1].(string)
	offset := int(config[2].(float64))

	indicator.period = period
	indicator.offset = offset

	// Get stock data
	stockData := *ctx.Data
	if len(stockData) == 0 {
		return []IndicatorValues{}
	}

	// Check if we have enough data (need period + 1 for previous values)
	if len(stockData) < period+1 {
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

	// Step 1: Calculate VM+ (Vortex Movement Plus), VM- (Vortex Movement Minus), and TR (True Range)
	vmPlus := make([]float64, len(stockData))
	vmMinus := make([]float64, len(stockData))
	tr := make([]float64, len(stockData))

	for i := 1; i < len(stockData); i++ {
		// VM+ = |Current High - Prior Low|
		vmPlus[i] = math.Abs(highs[i] - lows[i-1])

		// VM- = |Current Low - Prior High|
		vmMinus[i] = math.Abs(lows[i] - highs[i-1])

		// Calculate True Range: max(H-L, |H-C_prev|, |L-C_prev|)
		hl := highs[i] - lows[i]
		hc := math.Abs(highs[i] - closes[i-1])
		lc := math.Abs(lows[i] - closes[i-1])

		tr[i] = math.Max(hl, math.Max(hc, lc))
	}

	// Step 2: Calculate VI+ and VI- (Vortex Indicators)
	viPlus := make([]float64, len(stockData))
	viMinus := make([]float64, len(stockData))

	for i := period; i < len(stockData); i++ {
		// Calculate sum of VM+ for the period
		var sumVMPlus float64
		for j := i - period + 1; j <= i; j++ {
			sumVMPlus += vmPlus[j]
		}

		// Calculate sum of VM- for the period
		var sumVMMinus float64
		for j := i - period + 1; j <= i; j++ {
			sumVMMinus += vmMinus[j]
		}

		// Calculate sum of TR for the period
		var sumTR float64
		for j := i - period + 1; j <= i; j++ {
			sumTR += tr[j]
		}

		// Calculate Vortex Indicators
		// VI+ = Sum(VM+) / Sum(TR)
		// VI- = Sum(VM-) / Sum(TR)
		if sumTR > 0 {
			viPlus[i] = sumVMPlus / sumTR
			viMinus[i] = sumVMMinus / sumTR
		}
	}

	// Step 3: Build results starting from valid calculation points
	var result []IndicatorValues

	for i := period; i < len(stockData); i++ {
		// Apply offset
		if i < period+offset {
			continue
		}

		// Check start date filter
		if stockData[i].Timestamp.Before(startDate) {
			continue
		}

		// Select the appropriate value based on signal parameter
		var value float64
		if strings.ToUpper(signal) == "NEGATIVE_VI" {
			value = viMinus[i]
		} else { // Default to "POSITIVE_VI" for any other value
			value = viPlus[i]
		}

		result = append(result, IndicatorValues{
			Timestamp: stockData[i].Timestamp,
			Value:     util.Round(value),
		})
	}

	return result
}
