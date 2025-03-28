package util

import "math"

func Max[T int32 | float32 | int64 | float64](v1 T, v2 T) T {
    if v1 > v2 {
        return v1
    }
    return v2
}

func Round(f float64) float64 {
	pow := math.Pow(10, float64(2))
	return math.Round(f*pow) / pow
}

// Pow is a wrapper for math.Pow
func Pow(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

// Sqrt is a wrapper for math.Sqrt
func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}
