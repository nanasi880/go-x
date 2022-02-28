package httputil_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.nanasi880.dev/x/internal/testing/testutil"
	"go.nanasi880.dev/x/net/http/httputil"
)

func TestMethodHandler_ServeHTTP(t *testing.T) {
	handlerFunc := func(resp string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, resp)
		}
	}

	handler := httputil.MethodHandler{
		GET:     handlerFunc("GET"),
		HEAD:    handlerFunc("HEAD"),
		POST:    handlerFunc("POST"),
		PUT:     handlerFunc("PUT"),
		PATCH:   handlerFunc("PATCH"),
		DELETE:  handlerFunc("DELETE"),
		CONNECT: handlerFunc("CONNECT"),
		OPTIONS: handlerFunc("OPTIONS"),
		TRACE:   handlerFunc("TRACE"),
		Default: handlerFunc("DEFAULT"),
	}

	server := httptest.NewServer(&handler)
	defer server.Close()

	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
	for _, method := range methods {
		req, err := http.NewRequest(method, server.URL+"/", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		body := testutil.ReadAllAsString(t, resp.Body)
		if body != method && body != "" {
			t.Fatal(body, method)
		}
		resp.Body.Close()
	}
}
