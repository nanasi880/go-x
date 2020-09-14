package bits_test

import (
	"testing"

	"go.nanasi880.dev/x/bits"
)

func TestLSB(t *testing.T) {
	if got := bits.LSB(0); got != -1 {
		t.Logf("mask:%016x want:%d got:%d\n ", 0, -1, got)
		t.Fail()
	}
	for i := 1; i <= 64; i++ {
		// make mask
		mask := uint64(1 << (i - 1))
		want := i - 1
		got := bits.LSB(mask)
		if want != got {
			t.Logf("mask:%016x want:%d got:%d\n ", mask, i, got)
			t.Fail()
		}
	}
}

func TestMSB(t *testing.T) {
	if got := bits.MSB(0); got != -1 {
		t.Logf("mask:%016x want:%d got:%d\n ", 0, -1, got)
		t.Fail()
	}
	for i := 1; i <= 64; i++ {
		// make mask
		mask := uint64(0x8000000000000000 >> (64 - i))
		want := i - 1
		got := bits.MSB(mask)
		if want != got {
			t.Logf("mask:%016x want:%d got:%d\n ", mask, i, got)
			t.Fail()
		}
	}
}
