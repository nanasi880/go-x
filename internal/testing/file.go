package testing

import (
	"io/ioutil"
	"testing"
)

// MustReadFile is read file and returns content. In case of error, call Fatal.
func MustReadFile(tb testing.TB, filename string) []byte {
	tb.Helper()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		tb.Fatal(err)
	}
	return b
}
