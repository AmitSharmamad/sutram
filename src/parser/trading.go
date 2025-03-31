package parser

import (
	"errors"
	"math"
)

// DefaultTradingFunctions provides implementations of trading functions
type DefaultTradingFunctions struct {
	// Historical data
	Prices     []float64 // Close prices
	OpenPrices []float64 // Open prices
	HighPrices []float64 // High prices
	LowPrices  []float64 // Low prices
	Volumes    []float64 // Volume data

	// Previous state (for cross detection)
	PreviousValues map[string]float64

	// Portfolio state
	PortfolioValue float64
	Cash           float64
	PositionSize   float64
	MarketPosition float64 // 0 = flat, positive = long, negative = short
	AvgEntryPrice  float64
}

// NewDefaultTradingFunctions creates a new instance of trading functions
func NewDefaultTradingFunctions() *DefaultTradingFunctions {
	return &DefaultTradingFunctions{
		PreviousValues: make(map[string]float64),
	}
}

// SMA calculates Simple Moving Average
func (f *DefaultTradingFunctions) SMA(value float64, period int) (float64, error) {
	if len(f.Prices) < period {
		return 0, errors.New("not enough data for SMA calculation")
	}

	sum := 0.0
	for i := len(f.Prices) - period; i < len(f.Prices); i++ {
		sum += f.Prices[i]
	}

	return sum / float64(period), nil
}

// EMA calculates Exponential Moving Average
func (f *DefaultTradingFunctions) EMA(value float64, period int) (float64, error) {
	if len(f.Prices) < period {
		return 0, errors.New("not enough data for EMA calculation")
	}

	// Calculate multiplier
	multiplier := 2.0 / (float64(period) + 1.0)

	// Start with SMA for the first period
	var ema float64
	sum := 0.0
	for i := len(f.Prices) - period; i < len(f.Prices); i++ {
		sum += f.Prices[i]
	}
	ema = sum / float64(period)

	// Apply EMA formula for subsequent periods
	for i := len(f.Prices) - period + 1; i < len(f.Prices); i++ {
		ema = (f.Prices[i]-ema)*multiplier + ema
	}

	return ema, nil
}

// RSI calculates Relative Strength Index
func (f *DefaultTradingFunctions) RSI(period int) (float64, error) {
	if len(f.Prices) < period+1 {
		return 0, errors.New("not enough data for RSI calculation")
	}

	// Calculate price changes
	changes := make([]float64, period)
	for i := 0; i < period; i++ {
		idx := len(f.Prices) - period - 1 + i
		changes[i] = f.Prices[idx+1] - f.Prices[idx]
	}

	// Calculate average gains and losses
	sumGain := 0.0
	sumLoss := 0.0
	for _, change := range changes {
		if change > 0 {
			sumGain += change
		} else {
			sumLoss += math.Abs(change)
		}
	}

	// Calculate RS and RSI
	avgGain := sumGain / float64(period)
	avgLoss := sumLoss / float64(period)

	if avgLoss == 0 {
		return 100, nil
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi, nil
}

// MACD calculates Moving Average Convergence Divergence
func (f *DefaultTradingFunctions) MACD(fastPeriod, slowPeriod, signalPeriod int) (float64, error) {
	if len(f.Prices) < slowPeriod+signalPeriod {
		return 0, errors.New("not enough data for MACD calculation")
	}

	// Calculate fast EMA
	fastEMA, err := f.EMA(0, fastPeriod)
	if err != nil {
		return 0, err
	}

	// Calculate slow EMA
	slowEMA, err := f.EMA(0, slowPeriod)
	if err != nil {
		return 0, err
	}

	// Calculate MACD line
	macdLine := fastEMA - slowEMA

	// For signal line, we would need historical MACD values
	// For simplicity, we return the MACD line value
	return macdLine, nil
}

// BollingerBandUpper calculates the upper Bollinger Band
func (f *DefaultTradingFunctions) BollingerBandUpper(period int, stdDev float64) (float64, error) {
	if len(f.Prices) < period {
		return 0, errors.New("not enough data for Bollinger Bands calculation")
	}

	// Calculate SMA
	sma, err := f.SMA(0, period)
	if err != nil {
		return 0, err
	}

	// Calculate standard deviation
	sumSquaredDiff := 0.0
	for i := len(f.Prices) - period; i < len(f.Prices); i++ {
		diff := f.Prices[i] - sma
		sumSquaredDiff += diff * diff
	}

	standardDeviation := math.Sqrt(sumSquaredDiff / float64(period))

	// Upper band = SMA + (stdDev * standard deviation)
	upperBand := sma + (stdDev * standardDeviation)

	return upperBand, nil
}

// BollingerBandLower calculates the lower Bollinger Band
func (f *DefaultTradingFunctions) BollingerBandLower(period int, stdDev float64) (float64, error) {
	if len(f.Prices) < period {
		return 0, errors.New("not enough data for Bollinger Bands calculation")
	}

	// Calculate SMA
	sma, err := f.SMA(0, period)
	if err != nil {
		return 0, err
	}

	// Calculate standard deviation
	sumSquaredDiff := 0.0
	for i := len(f.Prices) - period; i < len(f.Prices); i++ {
		diff := f.Prices[i] - sma
		sumSquaredDiff += diff * diff
	}

	standardDeviation := math.Sqrt(sumSquaredDiff / float64(period))

	// Lower band = SMA - (stdDev * standard deviation)
	lowerBand := sma - (stdDev * standardDeviation)

	return lowerBand, nil
}

// BollingerBandMiddle calculates the middle Bollinger Band (SMA)
func (f *DefaultTradingFunctions) BollingerBandMiddle(period int, stdDev float64) (float64, error) {
	return f.SMA(0, period)
}

// ATR calculates Average True Range
func (f *DefaultTradingFunctions) ATR(period int) (float64, error) {
	if len(f.Prices) < period+1 || len(f.HighPrices) < period+1 || len(f.LowPrices) < period+1 {
		return 0, errors.New("not enough data for ATR calculation")
	}

	// Calculate true ranges
	trueRanges := make([]float64, period)
	for i := 0; i < period; i++ {
		idx := len(f.Prices) - period - 1 + i

		// Calculate the three differences
		highLow := f.HighPrices[idx+1] - f.LowPrices[idx+1]
		highPrevClose := math.Abs(f.HighPrices[idx+1] - f.Prices[idx])
		lowPrevClose := math.Abs(f.LowPrices[idx+1] - f.Prices[idx])

		// True range is the maximum of the three
		trueRanges[i] = math.Max(highLow, math.Max(highPrevClose, lowPrevClose))
	}

	// Calculate average
	sum := 0.0
	for _, tr := range trueRanges {
		sum += tr
	}

	return sum / float64(period), nil
}

// Stochastic calculates the Stochastic Oscillator
func (f *DefaultTradingFunctions) Stochastic(kPeriod, kSlowing, dPeriod int) (float64, error) {
	if len(f.Prices) < kPeriod || len(f.HighPrices) < kPeriod || len(f.LowPrices) < kPeriod {
		return 0, errors.New("not enough data for Stochastic calculation")
	}

	// Find highest high and lowest low over lookback period
	highestHigh := math.Inf(-1)
	lowestLow := math.Inf(1)

	for i := len(f.Prices) - kPeriod; i < len(f.Prices); i++ {
		if f.HighPrices[i] > highestHigh {
			highestHigh = f.HighPrices[i]
		}

		if f.LowPrices[i] < lowestLow {
			lowestLow = f.LowPrices[i]
		}
	}

	// Calculate %K (Fast Stochastic)
	currentClose := f.Prices[len(f.Prices)-1]

	// Avoid division by zero
	if highestHigh == lowestLow {
		return 50, nil
	}

	percentK := (currentClose - lowestLow) / (highestHigh - lowestLow) * 100

	// For simplicity, we return %K value
	// A full implementation would also calculate %D (the SMA of %K)
	return percentK, nil
}

// CrossAbove checks if value1 crossed above value2 in the last bar
func (f *DefaultTradingFunctions) CrossAbove(value1, value2 float64) (float64, error) {
	// Generate unique key for storing previous values
	key := "crossabove"

	// If we don't have previous values, initialize them
	if _, exists := f.PreviousValues[key+"_val1"]; !exists {
		f.PreviousValues[key+"_val1"] = 0
	}
	if _, exists := f.PreviousValues[key+"_val2"]; !exists {
		f.PreviousValues[key+"_val2"] = 0
	}

	// Get previous values
	prevVal1 := f.PreviousValues[key+"_val1"]
	prevVal2 := f.PreviousValues[key+"_val2"]

	// Update previous values for next time
	f.PreviousValues[key+"_val1"] = value1
	f.PreviousValues[key+"_val2"] = value2

	// Check for cross above
	if prevVal1 <= prevVal2 && value1 > value2 {
		return 1, nil
	}

	return 0, nil
}

// CrossBelow checks if value1 crossed below value2 in the last bar
func (f *DefaultTradingFunctions) CrossBelow(value1, value2 float64) (float64, error) {
	// Generate unique key for storing previous values
	key := "crossbelow"

	// If we don't have previous values, initialize them
	if _, exists := f.PreviousValues[key+"_val1"]; !exists {
		f.PreviousValues[key+"_val1"] = 0
	}
	if _, exists := f.PreviousValues[key+"_val2"]; !exists {
		f.PreviousValues[key+"_val2"] = 0
	}

	// Get previous values
	prevVal1 := f.PreviousValues[key+"_val1"]
	prevVal2 := f.PreviousValues[key+"_val2"]

	// Update previous values for next time
	f.PreviousValues[key+"_val1"] = value1
	f.PreviousValues[key+"_val2"] = value2

	// Check for cross below
	if prevVal1 >= prevVal2 && value1 < value2 {
		return 1, nil
	}

	return 0, nil
}

// PercentChange calculates the percentage change over a given period
func (f *DefaultTradingFunctions) PercentChange(current float64, periods int) (float64, error) {
	if len(f.Prices) < periods+1 {
		return 0, errors.New("not enough data for percent change calculation")
	}

	previousPrice := f.Prices[len(f.Prices)-periods-1]

	if previousPrice == 0 {
		return 0, errors.New("division by zero in percent change calculation")
	}

	percentChange := (current - previousPrice) / previousPrice * 100

	return percentChange, nil
}

// Highest returns the highest value over a given period
func (f *DefaultTradingFunctions) Highest(value float64, period int) (float64, error) {
	if len(f.Prices) < period {
		return 0, errors.New("not enough data for highest calculation")
	}

	highest := math.Inf(-1)

	for i := len(f.Prices) - period; i < len(f.Prices); i++ {
		if f.Prices[i] > highest {
			highest = f.Prices[i]
		}
	}

	return highest, nil
}

// Lowest returns the lowest value over a given period
func (f *DefaultTradingFunctions) Lowest(value float64, period int) (float64, error) {
	if len(f.Prices) < period {
		return 0, errors.New("not enough data for lowest calculation")
	}

	lowest := math.Inf(1)

	for i := len(f.Prices) - period; i < len(f.Prices); i++ {
		if f.Prices[i] < lowest {
			lowest = f.Prices[i]
		}
	}

	return lowest, nil
}
