package httputil

import "net/http"

// MethodHandler is router by HTTP method.
type MethodHandler struct {
	GET     http.Handler // HTTP GET
	HEAD    http.Handler // HTTP HEAD
	POST    http.Handler // HTTP POST
	PUT     http.Handler // HTTP PUT
	PATCH   http.Handler // HTTP PATCH
	DELETE  http.Handler // HTTP DELETE
	CONNECT http.Handler // HTTP CONNECT
	OPTIONS http.Handler // HTTP OPTIONS
	TRACE   http.Handler // HTTP TRACE
	Default http.Handler // Default handler
}

// ServeHTTP is implementation of http.Handler.
func (h *MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	switch r.Method {
	case http.MethodGet:
		handler = h.GET
	case http.MethodHead:
		handler = h.HEAD
	case http.MethodPost:
		handler = h.POST
	case http.MethodPut:
		handler = h.PUT
	case http.MethodPatch:
		handler = h.PATCH
	case http.MethodDelete:
		handler = h.DELETE
	case http.MethodConnect:
		handler = h.CONNECT
	case http.MethodOptions:
		handler = h.OPTIONS
	case http.MethodTrace:
		handler = h.TRACE
	}
	if handler == nil {
		handler = h.Default
	}
	if handler == nil {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
	handler.ServeHTTP(w, r)
}
