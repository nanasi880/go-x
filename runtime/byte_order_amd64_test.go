package runtime_test

import (
	"testing"

	"go.nanasi880.dev/x/runtime"
)

func TestCurrentEndianAMD64(t *testing.T) {
	e := runtime.CurrentByteOrder()
	if e != runtime.LittleEndian {
		t.Fatal(e)
	}
}
