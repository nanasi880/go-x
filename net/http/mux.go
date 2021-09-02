package http

import (
	"net/http"
	"sync"

	"go.nanasi880.dev/x/net/http/middleware"
)

// ServeMux is an HTTP request multiplexer.
type ServeMux struct {
	mutex       sync.Mutex
	middlewares []middleware.Middleware
	mux         *http.ServeMux
	once        sync.Once
	handler     http.Handler
}

// NewServeMux is return new ServeMux instance.
func NewServeMux() *ServeMux {
	return &ServeMux{
		mux: http.NewServeMux(),
	}
}

// UseMiddleware is register middleware. Middleware is processed before all handlers are executed.
func (m *ServeMux) UseMiddleware(middlewares ...middleware.Middleware) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.middlewares = append(m.middlewares, middlewares...)
}

// HandleFunc registers the handler function for the given pattern.
// Middlewares is processed after registered middleware by UseMiddleware and process before handlers are executed.
func (m *ServeMux) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...middleware.Middleware) {
	m.Handle(pattern, handler, middlewares...)
}

// Handle registers the handler for the given pattern.
// Middlewares is processed after registered middleware by UseMiddleware and process before handlers are executed.
// If a handler already exists for pattern, Handle panics.
func (m *ServeMux) Handle(pattern string, handler http.Handler, middlewares ...middleware.Middleware) {
	m.mux.Handle(pattern, middleware.Bind(handler, middlewares...))
}

// ServeHTTP is implementation of http.Handler.
func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.once.Do(m.buildHandlerOnce)
	m.handler.ServeHTTP(w, r)
}

func (m *ServeMux) buildHandlerOnce() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.handler = middleware.Bind(m.mux, m.middlewares...)
	m.middlewares = nil
}
