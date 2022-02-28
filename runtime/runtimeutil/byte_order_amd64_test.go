package runtimeutil_test

import (
	"testing"

	"go.nanasi880.dev/x/runtime/runtimeutil"
)

func TestCurrentEndianAMD64(t *testing.T) {
	e := runtimeutil.CurrentByteOrder()
	if e != runtimeutil.LittleEndian {
		t.Fatal(e)
	}
}
