package httputil_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.nanasi880.dev/x/internal/testing/testutil"
	xhttp "go.nanasi880.dev/x/net/http/httputil"
	"go.nanasi880.dev/x/net/http/httputil/middleware"
)

func TestServeMux_HandleFunc(t *testing.T) {
	mux := xhttp.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "OK")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer testutil.Close(t, resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Fatal(string(body))
	}
}

func TestServeMux_UseMiddleware(t *testing.T) {
	mux := xhttp.NewServeMux()
	mux.UseMiddleware(testMiddleware("X-Test-1", "Hello"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-1") != "Hello" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "X-Test-1: %s", r.Header.Get("X-Test-1"))
			return
		}
		if r.Header.Get("X-Test-2") != "World" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "X-Test-2: %s", r.Header.Get("X-Test-2"))
			return
		}
		w.WriteHeader(http.StatusOK)
	}, testMiddleware("X-Test-2", "World"))

	server := httptest.NewServer(mux)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer testutil.Close(t, resp.Body)

	if resp.StatusCode != http.StatusOK {
		testutil.Fail(t, testutil.ReadAllAsString(t, resp.Body))
	}

	if resp.Header.Get("X-Test-1") != "Hello" {
		testutil.Failf(t, "X-Test-1: %s", resp.Header.Get("X-Test-1"))
	}
	if resp.Header.Get("X-Test-2") != "World" {
		testutil.Failf(t, "X-Test-2: %s", resp.Header.Get("X-Test-2"))
	}
}

func testMiddleware(key string, value string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Add(key, value)
			w.Header().Add(key, value)
			next.ServeHTTP(w, r)
		})
	}
}
