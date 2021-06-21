package bits_test

import (
	"encoding/binary"
	"testing"

	"go.nanasi880.dev/x/bits"
)

func TestReverse64(t *testing.T) {
	testSuite := []struct {
		base [8]byte
		want [8]byte
	}{
		{
			base: [8]byte{0, 0, 0, 0, 0, 0, 0, 0b00000001},
			want: [8]byte{0b10000000, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			base: [8]byte{0, 0, 0, 0, 0, 0, 0b00000010, 0},
			want: [8]byte{0, 0b01000000, 0, 0, 0, 0, 0, 0},
		},
		{
			base: [8]byte{0, 0, 0, 0, 0, 0b00000100, 0, 0},
			want: [8]byte{0, 0, 0b00100000, 0, 0, 0, 0, 0},
		},
		{
			base: [8]byte{0, 0, 0, 0, 0b00001000, 0, 0, 0},
			want: [8]byte{0, 0, 0, 0b00010000, 0, 0, 0, 0},
		},
	}

	for no, suite := range testSuite {
		base := binary.BigEndian.Uint64(suite.base[:])
		want := binary.BigEndian.Uint64(suite.want[:])

		got := bits.Reverse64(base)

		if got != want {
			t.Logf("no:%d want: %b got: %b", no, want, got)
			t.Fail()
		}
	}
}
