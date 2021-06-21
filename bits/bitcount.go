package bits

// MSB is returns the position of the most significant bit. If the value is 0, returns -1.
func MSB(v uint64) int {
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	return popCount(v) - 1
}

// LSB is returns the position of the least significant bit. If the value is 0, returns -1.
func LSB(v uint64) int {
	v |= v << 1
	v |= v << 2
	v |= v << 4
	v |= v << 8
	v |= v << 16
	v |= v << 32
	pos := 64 - popCount(v)
	mask1 := pos & 0b01000000
	mask2 := (pos & 0b01000000) >> 6
	return pos - mask1 - mask2
}
