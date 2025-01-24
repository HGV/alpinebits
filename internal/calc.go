package internal

import (
	"math"
)

func CalculateMinFull(mco, std, max int) int {
	if mco == 0 {
		return std
	}
	return int(math.Min(float64(max-mco), float64(std)))
}

func CalculateDiscountPattern(nightsRequired, nightsDiscounted int) string {
	pattern := make([]byte, nightsRequired)

	// Fill the slice with '0' for non-discounted nights
	for i := 0; i < nightsRequired-nightsDiscounted; i++ {
		pattern[i] = '0'
	}

	// Fill the remaining with '1' for discounted nights
	for i := nightsRequired - nightsDiscounted; i < nightsRequired; i++ {
		pattern[i] = '1'
	}

	return string(pattern)
}
