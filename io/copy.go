package io

import "io"

// Copy is io.Copy ignore written byte.
func Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}
