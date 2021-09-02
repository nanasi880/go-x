package middleware_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	xhttp "go.nanasi880.dev/x/net/http"
	"go.nanasi880.dev/x/net/http/middleware"
	"go.nanasi880.dev/x/net/http/middleware/encoding"
)

func TestContentEncoding(t *testing.T) {
	mux := xhttp.NewServeMux()
	mux.UseMiddleware(middleware.ContentType("application/json"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, r.Body)
	}, middleware.ContentEncoding(encoding.AllEncodings()...))

	server := httptest.NewServer(mux)
	defer server.Close()

	const reqJSON = `{"hello": "world"}`
	req, err := http.NewRequest(http.MethodPost, server.URL+"/", newGzipStringReader(reqJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := readGzipAsString(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if respBody != reqJSON {
		t.Fatal(respBody)
	}
}

func newGzipReader(bin []byte) io.Reader {
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)
	_, _ = gz.Write(bin)
	_ = gz.Close()
	return buf
}

func newGzipStringReader(s string) io.Reader {
	return newGzipReader([]byte(s))
}

func readGzip(r io.Reader) ([]byte, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	return io.ReadAll(gz)
}

func readGzipAsString(r io.Reader) (string, error) {
	b, err := readGzip(r)
	return string(b), err
}
