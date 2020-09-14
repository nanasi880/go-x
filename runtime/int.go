package runtime

// IntBitSize is runtime int type bit size. 32 or 64.
const IntBitSize = 32 << (^uint(0) >> 63)
