package testutil

import (
	"io"
	"testing"
)

// Close is closing io.Closer and logging if error happen.
func Close(tb testing.TB, c io.Closer) {
	tb.Helper()
	err := c.Close()
	if err != nil {
		tb.Log(err)
	}
}

// ReadAllAsString is read all bytes as string.
func ReadAllAsString(tb testing.TB, r io.Reader) string {
	tb.Helper()
	return string(ReadAll(tb, r))
}

// ReadAll is read all bytes.
func ReadAll(tb testing.TB, r io.Reader) []byte {
	tb.Helper()
	bin, err := io.ReadAll(r)
	if err != nil {
		tb.Error(err)
	}
	return bin
}
