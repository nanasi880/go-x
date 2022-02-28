//go:build !amd64 && !386
// +build !amd64,!386

package runtimeutil

import "unsafe"

func currentByteOrder() ByteOrder {
	x := uint16(0xFF00)
	p := (*byte)(unsafe.Pointer(&x))
	if *p == 0 {
		return LittleEndian
	} else {
		return BigEndian
	}
}
