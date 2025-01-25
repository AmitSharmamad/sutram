package indicator_test

import (
	"algo/src/algo/indicator"
	"algo/src/stock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEMA_calculate(t *testing.T) {
	// Test case 1: Empty data
	data := stock.StockData{}
	config := map[string]any{}
	point := "close"
	ema := &indicator.EMA{}
	result := ema.Calculate(data, config, point)
	assert.Equal(t, 0.0, result)
}
