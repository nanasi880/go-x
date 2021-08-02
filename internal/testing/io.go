package testing

import (
	"io"
	"testing"
)

func Close(tb testing.TB, c io.Closer) {
	tb.Helper()
	err := c.Close()
	if err != nil {
		tb.Log(err)
	}
}
