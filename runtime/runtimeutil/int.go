package runtimeutil

const (
	// IntBitSize is runtime int type bit size. 32 or 64.
	IntBitSize = 32 << (^uint(0) >> 63)

	// MaxInt is max value of runtime int type.
	MaxInt = 1<<(IntBitSize-1) - 1

	// MinInt is min value of runtime int type.
	MinInt = -1 << (IntBitSize - 1)
)
