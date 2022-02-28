package middleware

import (
	"net/http"
	"strings"
)

// ContentType is checks the value of the Content-Type header and returns an error if it does not match.
// If the HTTP method is either GET, HEAD, CONNECT, OPTIONS, TRACE, the checking process will be skipped.
func ContentType(types ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodConnect, http.MethodOptions, http.MethodTrace:
				next.ServeHTTP(w, r)
			default:
				contentType := r.Header.Get("Content-Type")

				if r.ContentLength <= 0 && contentType == "" {
					next.ServeHTTP(w, r)
				}

				for _, t := range types {
					if strings.EqualFold(contentType, t) {
						next.ServeHTTP(w, r)
						return
					}
				}
				w.WriteHeader(http.StatusUnsupportedMediaType)
			}
		})
	}
}
