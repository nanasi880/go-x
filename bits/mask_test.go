package bits_test

import (
	"testing"

	"go.nanasi880.dev/x/bits"
)

func TestLeastMask(t *testing.T) {
	for i := 0; i <= 64; i++ {
		// make mask
		want := uint64(0)
		for j := 0; j < i; j++ {
			want |= 1 << j
		}
		got := bits.LeastMask(i)
		if got != want {
			t.Logf("want:%d got:%d\n ", want, got)
			t.Fail()
		}
	}
}

func TestMostMask(t *testing.T) {
	for i := 0; i <= 64; i++ {
		// make mask
		want := uint64(0)
		for j := 0; j < i; j++ {
			want |= uint64(0x8000000000000000 >> j)
		}
		got := bits.MostMask(i)
		if got != want {
			t.Logf("want:%d got:%d\n ", want, got)
			t.Fail()
		}
	}
}
