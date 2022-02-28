package unsafeutil

import (
	"unsafe"
)

// BytesToString returns string from bytes.
// Note that the returned string shares the backing store with the bytes.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes returns bytes from string.
// Note that the returned bytes shares the backing store with the string.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		S   string
		Cap int
	}{
		S:   s,
		Cap: len(s),
	}))
}
