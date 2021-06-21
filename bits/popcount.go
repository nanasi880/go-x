package bits

// PopCount is counts the number of bits.
func PopCount(v uint64) int {
	return popCount(v)
}

func popCount(v uint64) int {
	v = (v & 0x5555555555555555) + ((v & 0xaaaaaaaaaaaaaaaa) >> 1)
	v = (v & 0x3333333333333333) + ((v & 0xcccccccccccccccc) >> 2)
	v = (v & 0x0f0f0f0f0f0f0f0f) + ((v & 0xf0f0f0f0f0f0f0f0) >> 4)
	v = (v & 0x00ff00ff00ff00ff) + ((v & 0xff00ff00ff00ff00) >> 8)
	v = (v & 0x0000ffff0000ffff) + ((v & 0xffff0000ffff0000) >> 16)
	v = (v & 0x00000000ffffffff) + ((v & 0xffffffff00000000) >> 32)
	return int(v)
}
