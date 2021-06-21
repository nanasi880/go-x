package io

import "io"

// ReadFull is io.ReadFull ignore read byte.
func ReadFull(r io.Reader, buf []byte) error {
	_, err := io.ReadFull(r, buf)
	return err
}
