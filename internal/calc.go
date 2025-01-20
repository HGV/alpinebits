package internal

import "math"

func CalculateMinFull(mco, std, max int) int {
	if mco == 0 {
		return std
	}
	return int(math.Min(float64(max-mco), float64(std)))
}
