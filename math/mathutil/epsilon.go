package mathutil

import "math"

var (
	epsilon32 = math.Nextafter32(1, 2) - 1.0
	epsilon   = math.Nextafter(1, 2) - 1.0
)

// Epsilon32 return the epsilon as float32.
func Epsilon32() float32 {
	return epsilon32
}

// Epsilon return the epsilon.
func Epsilon() float64 {
	return epsilon
}
