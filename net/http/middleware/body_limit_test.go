package middleware_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	xbytes "go.nanasi880.dev/x/bytes"
	xhttp "go.nanasi880.dev/x/net/http"
	"go.nanasi880.dev/x/net/http/middleware"
)

func TestBodyLimit(t *testing.T) {

	mux := xhttp.NewServeMux()
	mux.UseMiddleware(middleware.BodyLimit(1 * xbytes.KB))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	testSuites := [...]struct {
		size xbytes.Size
		want bool
	}{
		{
			size: xbytes.KB - 1,
			want: true,
		},
		{
			size: xbytes.KB,
			want: true,
		},
		{
			size: xbytes.KB + 1,
			want: false,
		},
	}

	for _, suite := range testSuites {
		func() {
			body := make([]byte, suite.size)
			req, err := http.NewRequest(http.MethodPost, server.URL+"/", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			_, _ = io.Copy(io.Discard, resp.Body)

			if !suite.want && resp.StatusCode == 200 {
				dump, err := httputil.DumpResponse(resp, false)
				if err != nil {
					t.Fatal(err)
				}
				t.Fatal(string(dump))
			}
			if suite.want && resp.StatusCode != 200 {
				dump, err := httputil.DumpResponse(resp, false)
				if err != nil {
					t.Fatal(err)
				}
				t.Fatal(string(dump))
			}
		}()
	}
}
