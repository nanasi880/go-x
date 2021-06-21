package bits

// LeastMask is creates a mask value with the specified number of bits set, starting with the Least significant bit.
func LeastMask(n int) uint64 {
	return (1 << n) - 1
}

// MostMask is creates a mask value with the specified number of bits set, starting with the Most significant bit.
func MostMask(n int) uint64 {
	return reverse64((1 << n) - 1)
}
