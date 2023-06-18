package util

import (
	"math"
)

// round iFloat to iDecimalPlaces decimal points
func RoundToDecimal(iFloat float64, iDecimalPlaces int) float64 {
	var multiplier float64 = 10
	for i := 1; i < iDecimalPlaces; i++ {
		multiplier = multiplier * 10
	}
	return math.Round(iFloat*multiplier) / multiplier
}
