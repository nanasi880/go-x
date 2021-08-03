package testing

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
