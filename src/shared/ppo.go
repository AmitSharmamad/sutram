package shared

import "algo/src/util"

type PPO struct {
	fastPeriod int
	slowPeriod int
}

func (indicator *PPO) Calculate(
	prices Prices,
	config map[string]any,
	point string,
) float64 {
	fastPeriod := indicator.fastPeriod
	if fastPeriod == 0 {
		if p, ok := config["fastPeriod"].(int); ok {
			fastPeriod = p
		} else {
			fastPeriod = 12 // default fast period
		}
	}

	slowPeriod := indicator.slowPeriod
	if slowPeriod == 0 {
		if p, ok := config["slowPeriod"].(int); ok {
			slowPeriod = p
		} else {
			slowPeriod = 26 // default slow period
		}
	}

	if len(prices) < slowPeriod {
		return 0.0
	}

	// Calculate EMAs
	fastEMA := &EMA{period: fastPeriod}
	slowEMA := &EMA{period: slowPeriod}

	fastValue := fastEMA.Calculate(prices, nil, point)
	slowValue := slowEMA.Calculate(prices, nil, point)

	if slowValue == 0 {
		return 0.0
	}

	// PPO = ((Fast EMA - Slow EMA) / Slow EMA) * 100
	ppo := ((fastValue - slowValue) / slowValue) * 100

	return util.Round(ppo)
}

func (ppo *PPO) GetDataRange() int32 {
	if ppo.slowPeriod > ppo.fastPeriod {
		return int32(ppo.slowPeriod)
	}
	return int32(ppo.fastPeriod)
}
