package math

import (
	"math"
)

// Round rounds a float64 to the specified number of decimal places
func Round(val float64) float64 {
	shift := math.Pow(10, float64(2))
	return math.Round(val*shift) / shift
}

// ClampFloat restricts a float value to be within the specified range
func ClampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
