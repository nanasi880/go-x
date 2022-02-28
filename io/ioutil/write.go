package ioutil

import "io"

// Write is io.Write ignore written byte.
func Write(w io.Writer, p []byte) error {
	_, err := w.Write(p)
	return err
}
