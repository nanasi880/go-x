package testing

import "testing"

// Fail is calls tb.Fail() after logging.
func Fail(tb testing.TB, args ...interface{}) {
	tb.Helper()
	tb.Log(args...)
	tb.Fail()
}

// Failf is calls tb.Fail() after logging.
func Failf(tb testing.TB, format string, args ...interface{}) {
	tb.Helper()
	tb.Logf(format, args...)
	tb.Fail()
}

// Cleanup is register cleanup handler func.
func Cleanup(tb testing.TB, f interface{}) {
	tb.Helper()
	switch f := f.(type) {
	case func():
		tb.Cleanup(f)
	case func() error:
		tb.Cleanup(func() {
			err := f()
			if err != nil {
				tb.Log(err)
			}
		})
	default:
		tb.Fatal("Cleanup() is required `func()` or `func() error` function")
	}
}
