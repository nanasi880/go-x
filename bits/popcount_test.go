package bits_test

import (
	"math"
	"testing"

	"go.nanasi880.dev/x/bits"
)

func TestPopCount(t *testing.T) {
	for i := 0; i <= 64; i++ {
		// make mask
		mask := uint64(0)
		for j := 0; j < i; j++ {
			mask = mask | (1 << j)
		}
		got := bits.PopCount(mask)
		if i != got {
			t.Logf("mask:%016x want:%d got:%d\n ", mask, i, got)
			t.Fail()
		}
	}
}

func BenchmarkPopCount(b *testing.B) {
	b.Run("PopCount", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bits.PopCount(math.MaxUint64)
		}
	})
	b.Run("PopCountByLoop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			popCountByLoop(math.MaxUint64)
		}
	})
}

func popCountByLoop(v uint64) int {
	count := 0
	for i := 0; i < 64; i++ {
		if (v & (1 << i)) > 0 {
			count++
		}
	}
	return count
}
