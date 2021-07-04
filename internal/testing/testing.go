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
