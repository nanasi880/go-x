package middleware

import "net/http"

// Middleware is http server middleware function.
type Middleware func(next http.Handler) http.Handler

// Bind is bind middleware to handler.
func Bind(handler http.Handler, m ...Middleware) http.Handler {
	for i := len(m); i > 0; i-- {
		handler = m[i-1](handler)
	}
	return handler
}
