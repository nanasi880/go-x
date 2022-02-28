package ioutil

import "io"

// Copy is io.Copy ignore written byte.
func Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

// CopyN is io.CopyN ignore written byte.
func CopyN(dst io.Writer, src io.Reader, n int64) error {
	_, err := io.CopyN(dst, src, n)
	return err
}
