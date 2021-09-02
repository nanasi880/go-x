package encoding_test

import (
	"net/http"
	"testing"

	"go.nanasi880.dev/x/net/http/middleware/encoding"
)

func TestGetContentEncodings(t *testing.T) {
	testSuites := []struct {
		Headers []string
		Want    []string
	}{
		{
			Headers: []string{"gzip, deflate"},
			Want:    []string{"deflate", "gzip"},
		},
	}

	for _, suite := range testSuites {
		header := make(http.Header)
		addHeaderAll(header, "Content-Encoding", suite.Headers)
		encodings := encoding.GetContentEncodings(header)
		if !stringEq(suite.Want, encodings) {
			t.Logf("want: %v got: %v", suite.Want, encodings)
			t.Fail()
		}
	}
}

func TestGetAcceptEncodings(t *testing.T) {
	testSuites := []struct {
		Header string
		Want   []string
	}{
		{
			Header: "gzip",
			Want:   []string{"gzip"},
		},
		{
			Header: "gzip, compress, br",
			Want:   []string{"gzip", "compress", "br"},
		},
		{
			Header: "br;q=1.0, gzip;q=0.8, *;q=0.1",
			Want:   []string{"br", "gzip", "*"},
		},
		{
			Header: "invalid;q=aaa, gzip",
			Want:   []string{"gzip"},
		},
	}

	for _, suite := range testSuites {
		header := make(http.Header)
		header.Set("Accept-Encoding", suite.Header)
		encodings := encoding.GetAcceptEncodings(header)
		if !stringEq(suite.Want, encodings) {
			t.Logf("want: %v got: %v", suite.Want, encodings)
			t.Fail()
		}
	}
}

func addHeaderAll(h http.Header, key string, values []string) {
	for _, v := range values {
		h.Add(key, v)
	}
}

func stringEq(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
