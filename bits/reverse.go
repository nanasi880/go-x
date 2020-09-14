package bits

// Reverse64 is arrange the bits in reverse order.
func Reverse64(v uint64) uint64 {
	return reverse64(v)
}

func reverse64(v uint64) uint64 {
	v = (v&0x5555555555555555)<<1 | (v >> 1 & 0x5555555555555555)
	v = (v&0x3333333333333333)<<2 | (v >> 2 & 0x3333333333333333)
	v = (v&0x0f0f0f0f0f0f0f0f)<<4 | (v >> 4 & 0x0f0f0f0f0f0f0f0f)
	v = (v&0x00ff00ff00ff00ff)<<8 | (v >> 8 & 0x00ff00ff00ff00ff)
	v = (v&0x0000ffff0000ffff)<<16 | (v >> 16 & 0x0000ffff0000ffff)
	v = (v&0x00000000ffffffff)<<32 | (v >> 32 & 0x00000000ffffffff)
	return v
}
