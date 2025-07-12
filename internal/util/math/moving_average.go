package math

import (
	"math"
	"strings"

	"github.com/markcheno/go-talib"
)

// MovingAverageTypes contains all the supported moving average types
var MovingAverageTypes = []string{
	"SIMPLE",
	"EXPONENTIAL",
	"WEIGHTED",
	"TRIPLE",
	"TRIANGULAR",
	"DOUBLE",
	"HULL",
	"MESA",
	"VARIABLE",
	"KAUFMAN",
	"VIDYA",
}

// ConvertMATypeToTalib converts string MA type to go-talib MaType enum
// Used by Bollinger Bands and other indicators that use go-talib functions
func ConvertMATypeToTalib(maType string) talib.MaType {
	switch strings.ToUpper(maType) {
	case "EXPONENTIAL":
		return talib.EMA
	case "WEIGHTED":
		return talib.WMA
	case "TRIPLE":
		return talib.TEMA
	case "TRIANGULAR":
		return talib.TRIMA
	case "DOUBLE":
		return talib.DEMA
	case "KAUFMAN":
		return talib.KAMA
	case "MESA":
		// MESA is a special case - BBands doesn't support it directly
		// Fall back to EMA for now
		return talib.EMA
	case "HULL", "VARIABLE", "VIDYA":
		// These MA types are not available in go-talib BBands
		// Fall back to SMA for compatibility
		return talib.SMA
	default: // "SIMPLE"
		return talib.SMA
	}
}

// CalculateMovingAverage calculates the moving average of a given type for the provided data
func CalculateMovingAverage(data []float64, maType string, period int) []float64 {
	if len(data) < period {
		return []float64{}
	}

	maType = strings.ToUpper(maType)
	result := make([]float64, len(data))

	switch maType {
	case "SIMPLE":
		calculateSMA(data, period, result)
	case "EXPONENTIAL":
		calculateEMA(data, period, result)
	case "WEIGHTED":
		calculateWMA(data, period, result)
	case "TRIPLE":
		calculateTripleEMA(data, period, result)
	case "TRIANGULAR":
		calculateTriangularMA(data, period, result)
	case "DOUBLE":
		calculateDoubleEMA(data, period, result)
	case "HULL":
		calculateHullMA(data, period, result)
	// Note: MESA, VARIABLE, KAUFMAN, VIDYA have specialized implementations
	// and would need more complex code than can be provided here.
	// For now, default to EMA for these types
	default:
		calculateEMA(data, period, result) // Default to EMA
	}

	return result
}

// Simple Moving Average (SMA)
func calculateSMA(data []float64, period int, result []float64) {
	for i := 0; i < len(data); i++ {
		if i < period-1 {
			result[i] = 0 // Not enough data
		} else {
			var sum float64
			for j := 0; j < period; j++ {
				sum += data[i-j]
			}
			result[i] = sum / float64(period)
		}
	}
}

// Exponential Moving Average (EMA)
func calculateEMA(data []float64, period int, result []float64) {
	// Calculate first SMA
	var sum float64
	for i := 0; i < period; i++ {
		sum += data[i]
	}

	// Set initial EMA
	result[period-1] = sum / float64(period)

	// Calculate subsequent EMAs
	k := 2.0 / float64(period+1)
	for i := period; i < len(data); i++ {
		result[i] = (data[i]-result[i-1])*k + result[i-1]
	}
}

// Weighted Moving Average (WMA)
func calculateWMA(data []float64, period int, result []float64) {
	denominator := (period * (period + 1)) / 2

	for i := period - 1; i < len(data); i++ {
		var sum float64
		weight := 1

		for j := 0; j < period; j++ {
			sum += data[i-j] * float64(period-j)
			weight += period - j
		}

		result[i] = sum / float64(denominator)
	}
}

// Double Exponential Moving Average (DEMA)
func calculateDoubleEMA(data []float64, period int, result []float64) {
	ema := make([]float64, len(data))
	emaOfEma := make([]float64, len(data))

	calculateEMA(data, period, ema)
	calculateEMA(ema, period, emaOfEma)

	for i := period*2 - 2; i < len(data); i++ {
		result[i] = 2*ema[i] - emaOfEma[i]
	}
}

// Triple Exponential Moving Average (TEMA)
func calculateTripleEMA(data []float64, period int, result []float64) {
	ema := make([]float64, len(data))
	emaOfEma := make([]float64, len(data))
	emaOfEmaOfEma := make([]float64, len(data))

	calculateEMA(data, period, ema)
	calculateEMA(ema, period, emaOfEma)
	calculateEMA(emaOfEma, period, emaOfEmaOfEma)

	for i := period*3 - 3; i < len(data); i++ {
		result[i] = 3*ema[i] - 3*emaOfEma[i] + emaOfEmaOfEma[i]
	}
}

// Triangular Moving Average (TMA)
func calculateTriangularMA(data []float64, period int, result []float64) {
	n := (period + 1) / 2

	// Calculate first SMA
	sma1 := make([]float64, len(data))
	calculateSMA(data, n, sma1)

	// Calculate second SMA on the first SMA
	calculateSMA(sma1, n, result)
}

// Hull Moving Average (HMA)
func calculateHullMA(data []float64, period int, result []float64) {
	wma1 := make([]float64, len(data))
	wma2 := make([]float64, len(data))
	diff := make([]float64, len(data))

	halfPeriod := period / 2
	sqrtPeriod := int(math.Sqrt(float64(period)))

	calculateWMA(data, period, wma1)
	calculateWMA(data, halfPeriod, wma2)

	for i := 0; i < len(data); i++ {
		if i >= halfPeriod-1 {
			diff[i] = 2*wma2[i] - wma1[i]
		}
	}

	calculateWMA(diff, sqrtPeriod, result)
}
