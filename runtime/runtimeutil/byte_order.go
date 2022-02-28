package runtimeutil

//go:generate stringer -type=ByteOrder -output=byte_order_string.go
type ByteOrder int

const (
	LittleEndian ByteOrder = iota
	BigEndian
)

// CurrentByteOrder is get runtime machine endian.
func CurrentByteOrder() ByteOrder {
	return currentByteOrder()
}
