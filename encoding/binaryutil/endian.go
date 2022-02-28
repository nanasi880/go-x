package binaryutil

import (
	"encoding/binary"

	"go.nanasi880.dev/x/runtime/runtimeutil"
)

// CurrentByteOrder is get current byte order.
func CurrentByteOrder() binary.ByteOrder {
	if runtimeutil.CurrentByteOrder() == runtimeutil.LittleEndian {
		return binary.LittleEndian
	}
	return binary.BigEndian
}
