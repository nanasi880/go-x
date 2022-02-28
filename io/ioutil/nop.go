package ioutil

import "io"

type nopWriteCloser struct {
	w io.Writer
}

func (w nopWriteCloser) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w nopWriteCloser) Close() error {
	return nil
}

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w: w}
}
