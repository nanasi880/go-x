package mathutil

func Clamp(value, min, max float64) float64 {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}
