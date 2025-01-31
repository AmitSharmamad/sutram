package indicator_test

import (
	"testing"

	"algo/src/shared"

	"github.com/stretchr/testify/assert"
)

func TestEMA_calculate(t *testing.T) {
	// Test case 1: Empty data
	data := shared.StockData{}
	config := map[string]any{}
	point := "close"
	ema := &shared.EMA{}
	prices := shared.GetPrices(data, point)
	result := ema.Calculate(prices, config, point)
	assert.Equal(t, 0.0, result)
}
