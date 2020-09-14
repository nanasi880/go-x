package runtime_test

import (
	"testing"

	"go.nanasi880.dev/x/runtime"
)

func TestIntSize(t *testing.T) {
	if runtime.IntBitSize != 32 && runtime.IntBitSize != 64 {
		t.Fatal(runtime.IntBitSize)
	}
}
