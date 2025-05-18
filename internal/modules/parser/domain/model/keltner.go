package parser

import (
	util "algo/internal/util/math"
	"math"
)

type KeltnerChannels struct {
	emaPeriod  int
	atrPeriod  int
	multiplier float64
	offset     int
}

func (indicator *KeltnerChannels) SliceData(prices Prices) Prices {
	// We need enough data to calculate Keltner Channels properly
	maxPeriod := max(indicator.emaPeriod, indicator.atrPeriod)
	start := max(0, len(prices)-(indicator.offset+2*maxPeriod))
	end := len(prices) - indicator.offset

	if start < 0 || end <= start || end > len(prices) {
		return Prices{}
	}

	return prices[start:end]
}

// Helper method to slice all required data
func (indicator *KeltnerChannels) sliceAllData(highs, lows, closes Prices) (Prices, Prices, Prices) {
	// We need enough data to calculate Keltner Channels properly
	maxPeriod := max(indicator.emaPeriod, indicator.atrPeriod)
	start := max(0, len(closes)-(indicator.offset+2*maxPeriod))
	end := len(closes) - indicator.offset

	if start < 0 || end <= start || end > len(closes) {
		return Prices{}, Prices{}, Prices{}
	}

	return highs[start:end], lows[start:end], closes[start:end]
}

func (indicator *KeltnerChannels) Calculate(
	config []any,
	ctx *EvaluationContext,
) float64 {
	emaPeriod := int(config[0].(float64))
	atrPeriod := int(config[1].(float64))
	multiplier := float64(config[2].(float64))
	offset := int(config[3].(float64))

	indicator.emaPeriod = emaPeriod
	indicator.atrPeriod = atrPeriod
	indicator.multiplier = multiplier
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")

	// Check if we have enough data points
	maxPeriod := max(emaPeriod, atrPeriod)
	if len(highs) < maxPeriod+offset+1 || len(lows) < maxPeriod+offset+1 || len(closes) < maxPeriod+offset+1 {
		return 0.0
	}

	highs, lows, closes = indicator.sliceAllData(highs, lows, closes)

	// Need at least maxPeriod+1 data points to calculate
	if len(highs) < maxPeriod+1 || len(lows) < maxPeriod+1 || len(closes) < maxPeriod+1 {
		return 0.0
	}

	// Calculate EMA of typical price (middle line)
	var typicalPrices = make([]float64, len(closes))
	for i := 0; i < len(closes); i++ {
		typicalPrices[i] = (highs[i] + lows[i] + closes[i]) / 3.0
	}

	// Calculate EMA
	multiplierEMA := 2.0 / float64(indicator.emaPeriod+1)
	var emaSum float64
	for i := 0; i < indicator.emaPeriod; i++ {
		emaSum += typicalPrices[i]
	}
	ema := emaSum / float64(indicator.emaPeriod)

	var emas = make([]float64, len(closes))
	emas[indicator.emaPeriod-1] = ema

	for i := indicator.emaPeriod; i < len(closes); i++ {
		ema = (typicalPrices[i]-ema)*multiplierEMA + ema
		emas[i] = ema
	}

	// Calculate ATR (Average True Range)
	var tr = make([]float64, len(closes))

	// Calculate True Range
	for i := 1; i < len(closes); i++ {
		tr[i] = max(
			highs[i]-lows[i],
			max(
				math.Abs(highs[i]-closes[i-1]),
				math.Abs(lows[i]-closes[i-1]),
			),
		)
	}

	// Calculate initial ATR as simple average of TR
	var atrSum float64
	for i := 1; i <= indicator.atrPeriod; i++ {
		atrSum += tr[i]
	}
	atr := atrSum / float64(indicator.atrPeriod)

	var atrs = make([]float64, len(closes))
	atrs[indicator.atrPeriod] = atr

	// Calculate subsequent ATR values
	for i := indicator.atrPeriod + 1; i < len(closes); i++ {
		atr = ((atr * float64(indicator.atrPeriod-1)) + tr[i]) / float64(indicator.atrPeriod)
		atrs[i] = atr
	}

	// Calculate upper and lower bands
	var upperBand = make([]float64, len(closes))
	var lowerBand = make([]float64, len(closes))

	maxIndex := max(indicator.emaPeriod, indicator.atrPeriod+1)
	for i := maxIndex; i < len(closes); i++ {
		upperBand[i] = emas[i] + indicator.multiplier*atrs[i]
		lowerBand[i] = emas[i] - indicator.multiplier*atrs[i]
	}

	// Return the middle line (EMA)
	return util.Round(emas[len(emas)-1])
}

func (indicator *KeltnerChannels) CalculateAll(
	config []any,
	ctx *EvaluationContext,
) (float64, float64, float64) {
	emaPeriod := int(config[0].(float64))
	atrPeriod := int(config[1].(float64))
	multiplier := float64(config[2].(float64))
	offset := int(config[3].(float64))

	indicator.emaPeriod = emaPeriod
	indicator.atrPeriod = atrPeriod
	indicator.multiplier = multiplier
	indicator.offset = offset

	highs := GetPrices(*ctx.Data, "high")
	lows := GetPrices(*ctx.Data, "low")
	closes := GetPrices(*ctx.Data, "close")

	// Check if we have enough data points
	maxPeriod := max(emaPeriod, atrPeriod)
	if len(highs) < maxPeriod+offset+1 || len(lows) < maxPeriod+offset+1 || len(closes) < maxPeriod+offset+1 {
		return 0.0, 0.0, 0.0
	}

	highs, lows, closes = indicator.sliceAllData(highs, lows, closes)

	// Need at least maxPeriod+1 data points to calculate
	if len(highs) < maxPeriod+1 || len(lows) < maxPeriod+1 || len(closes) < maxPeriod+1 {
		return 0.0, 0.0, 0.0
	}

	// Calculate EMA of typical price (middle line)
	var typicalPrices = make([]float64, len(closes))
	for i := 0; i < len(closes); i++ {
		typicalPrices[i] = (highs[i] + lows[i] + closes[i]) / 3.0
	}

	// Calculate EMA
	multiplierEMA := 2.0 / float64(indicator.emaPeriod+1)
	var emaSum float64
	for i := 0; i < indicator.emaPeriod; i++ {
		emaSum += typicalPrices[i]
	}
	ema := emaSum / float64(indicator.emaPeriod)

	var emas = make([]float64, len(closes))
	emas[indicator.emaPeriod-1] = ema

	for i := indicator.emaPeriod; i < len(closes); i++ {
		ema = (typicalPrices[i]-ema)*multiplierEMA + ema
		emas[i] = ema
	}

	// Calculate ATR (Average True Range)
	var tr = make([]float64, len(closes))

	// Calculate True Range
	for i := 1; i < len(closes); i++ {
		tr[i] = max(
			highs[i]-lows[i],
			max(
				math.Abs(highs[i]-closes[i-1]),
				math.Abs(lows[i]-closes[i-1]),
			),
		)
	}

	// Calculate initial ATR as simple average of TR
	var atrSum float64
	for i := 1; i <= indicator.atrPeriod; i++ {
		atrSum += tr[i]
	}
	atr := atrSum / float64(indicator.atrPeriod)

	var atrs = make([]float64, len(closes))
	atrs[indicator.atrPeriod] = atr

	// Calculate subsequent ATR values
	for i := indicator.atrPeriod + 1; i < len(closes); i++ {
		atr = ((atr * float64(indicator.atrPeriod-1)) + tr[i]) / float64(indicator.atrPeriod)
		atrs[i] = atr
	}

	// Calculate upper and lower bands
	var upperBand, lowerBand float64

	lastIndex := len(closes) - 1

	upperBand = emas[lastIndex] + indicator.multiplier*atrs[lastIndex]
	lowerBand = emas[lastIndex] - indicator.multiplier*atrs[lastIndex]

	return util.Round(upperBand), util.Round(emas[lastIndex]), util.Round(lowerBand)
}

func (indicator *KeltnerChannels) GetDataRange(config map[string]any) int32 {
	if p, ok := config["emaPeriod"].(float64); ok {
		indicator.emaPeriod = int(p)
	} else {
		indicator.emaPeriod = 20 // default EMA period
	}

	if p, ok := config["atrPeriod"].(float64); ok {
		indicator.atrPeriod = int(p)
	} else {
		indicator.atrPeriod = 10 // default ATR period
	}

	if m, ok := config["multiplier"].(float64); ok {
		indicator.multiplier = m
	} else {
		indicator.multiplier = 2.0 // default multiplier
	}

	if o, ok := config["offset"].(float64); ok {
		indicator.offset = int(o)
	} else {
		indicator.offset = 0
	}

	// Return the maximum range required
	maxPeriod := max(indicator.emaPeriod, indicator.atrPeriod)
	// We need maxPeriod+1 data points (due to the use of previous values) + offset
	return int32(maxPeriod + 1 + indicator.offset)
}
