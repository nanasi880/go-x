package binary

import (
	"encoding/binary"

	"go.nanasi880.dev/x/runtime"
)

// CurrentByteOrder is get current byte order.
func CurrentByteOrder() binary.ByteOrder {
	if runtime.CurrentByteOrder() == runtime.LittleEndian {
		return binary.LittleEndian
	}
	return binary.BigEndian
}
