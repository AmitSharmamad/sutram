package indicator

import (
	"algo/src/shared"
	"algo/src/util"
)

type KST struct {
	roc1Period int
	roc2Period int
	roc3Period int
	roc4Period int
	sma1Period int
	sma2Period int
	sma3Period int
	sma4Period int
}

func (indicator *KST) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	// Set default periods if not configured
	if indicator.roc1Period == 0 {
		indicator.roc1Period = 10 // default ROC1 period
		indicator.roc2Period = 15 // default ROC2 period
		indicator.roc3Period = 20 // default ROC3 period
		indicator.roc4Period = 30 // default ROC4 period
		indicator.sma1Period = 10 // default SMA1 period
		indicator.sma2Period = 10 // default SMA2 period
		indicator.sma3Period = 10 // default SMA3 period
		indicator.sma4Period = 15 // default SMA4 period
	}

	maxPeriod := max(
		indicator.roc1Period+indicator.sma1Period,
		indicator.roc2Period+indicator.sma2Period,
		indicator.roc3Period+indicator.sma3Period,
		indicator.roc4Period+indicator.sma4Period,
	)

	if len(prices) < maxPeriod {
		return 0.0
	}

	// Calculate ROCs
	roc1 := &ROC{period: indicator.roc1Period}
	roc2 := &ROC{period: indicator.roc2Period}
	roc3 := &ROC{period: indicator.roc3Period}
	roc4 := &ROC{period: indicator.roc4Period}

	rocValue1 := roc1.Calculate(prices, nil, point)
	rocValue2 := roc2.Calculate(prices, nil, point)
	rocValue3 := roc3.Calculate(prices, nil, point)
	rocValue4 := roc4.Calculate(prices, nil, point)

	// Calculate SMAs of ROCs
	sma1 := &SMA{period: indicator.sma1Period}
	sma2 := &SMA{period: indicator.sma2Period}
	sma3 := &SMA{period: indicator.sma3Period}
	sma4 := &SMA{period: indicator.sma4Period}

	smaValue1 := sma1.Calculate(Prices{rocValue1}, nil, point)
	smaValue2 := sma2.Calculate(Prices{rocValue2}, nil, point)
	smaValue3 := sma3.Calculate(Prices{rocValue3}, nil, point)
	smaValue4 := sma4.Calculate(Prices{rocValue4}, nil, point)

	// Calculate KST
	kst := (smaValue1 * 1) + (smaValue2 * 2) + (smaValue3 * 3) + (smaValue4 * 4)

	return util.Round(kst)
}

func (kst *KST) GetDataRange() int32 {
	return int32(max(
		kst.roc1Period+kst.sma1Period,
		kst.roc2Period+kst.sma2Period,
		kst.roc3Period+kst.sma3Period,
		kst.roc4Period+kst.sma4Period,
	))
}

func max(numbers ...int) int {
	if len(numbers) == 0 {
		return 0
	}
	maxNum := numbers[0]
	for _, num := range numbers[1:] {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}

func (indicator *KST) CalculateAll(data shared.StockData, config map[string]any, point string) []float64 {
	var results []float64
	return results
}
