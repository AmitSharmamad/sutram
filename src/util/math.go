package util

import "math"

func Max(v1 int32, v2 int32) int32 {
	if v1 > v2 {
		return v1
	}
	return v2
}

func Round(f float64) float64 {
	pow := math.Pow(10, float64(2))
	return math.Round(f*pow) / pow
}
