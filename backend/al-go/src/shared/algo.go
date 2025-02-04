package shared

import "math"

func Round(f float64) float64 {
	pow := math.Pow(10, float64(2))
	return math.Round(f*pow) / pow
}
