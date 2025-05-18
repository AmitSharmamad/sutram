package parser

import (
	util "algo/internal/util/math"
	"math"
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
	period := int(config[0].(float64))
	offset := int(config[1].(float64))

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
