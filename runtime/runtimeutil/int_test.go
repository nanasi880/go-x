package runtimeutil_test

import (
	"testing"

	"go.nanasi880.dev/x/runtime/runtimeutil"
)

func TestIntSize(t *testing.T) {
	if runtimeutil.IntBitSize != 32 && runtimeutil.IntBitSize != 64 {
		t.Fatal(runtimeutil.IntBitSize)
	}
}
