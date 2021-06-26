package xx

import (
	"encoding/binary"
	"math/bits"
)

type xxhUint128 struct {
	high uint64
	low  uint64
}

func mul32To64(lhs uint32, rhs uint32) uint64 {
	return uint64(lhs) * uint64(rhs)
}

func mul64To128(lhs uint64, rhs uint64) xxhUint128 {
	high, low := bits.Mul64(lhs, rhs)
	return xxhUint128{
		high: high,
		low:  low,
	}
}

func mul128Fold64(lhs uint64, rhs uint64) uint64 {
	high, low := bits.Mul64(lhs, rhs)
	return high ^ low
}

// shorthand
func rotl32(x uint32, r int) uint32 {
	return bits.RotateLeft32(x, r)
}

// shorthand
func rotl64(x uint64, r int) uint64 {
	return bits.RotateLeft64(x, r)
}

// shorthand
func readLE32(p []byte) uint32 {
	return binary.LittleEndian.Uint32(p)
}

// shorthand
func readLE64(p []byte) uint64 {
	return binary.LittleEndian.Uint64(p)
}

// shorthand
func writeLE64(p []byte, x uint64) {
	binary.LittleEndian.PutUint64(p, x)
}

// shorthand
func swap32(x uint32) uint32 {
	return bits.ReverseBytes32(x)
}

// shorthand
func swap64(x uint64) uint64 {
	return bits.ReverseBytes64(x)
}
