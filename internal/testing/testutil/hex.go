package testutil

import (
	"encoding/hex"
	"testing"
)

// MustDecodeHexString is decode hex string. In case of error, call Fatal.
func MustDecodeHexString(tb testing.TB, h string) []byte {
	tb.Helper()

	b, err := hex.DecodeString(h)
	if err != nil {
		tb.Fatal(err)
	}
	return b
}
