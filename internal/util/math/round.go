package math

import (
	"math"
)

// Round rounds a float64 to the specified number of decimal places
func Round(val float64) float64 {
	shift := math.Pow(10, float64(2))
	return math.Round(val*shift) / shift
}
