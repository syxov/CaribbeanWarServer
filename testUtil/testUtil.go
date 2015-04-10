package testUtil

import (
	"math"
)

func FloatEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0001
}
