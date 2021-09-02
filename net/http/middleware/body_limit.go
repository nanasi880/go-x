package middleware

import (
	"errors"
	"io"
	"net/http"

	"go.nanasi880.dev/x/bytes"
)

var (
	ErrStatusRequestEntityTooLarge = errors.New("request entity too large")
)

type bodyLimitReader struct {
	body    io.ReadCloser
	n       int64
	limit   int64
	onError func()
}

// Read is implementation of io.Reader.
func (r *bodyLimitReader) Read(p []byte) (int, error) {
	n, err := r.body.Read(p)
	r.n += int64(n)
	if r.n > r.limit {
		r.onError()
		return n, ErrStatusRequestEntityTooLarge
	}
	return n, err
}

// Close is implementation of io.ReadCloser.
func (r *bodyLimitReader) Close() error {
	return r.body.Close()
}

// BodyLimitConfig is configuration of BodyLimit middleware.
type BodyLimitConfig struct {
	Limit   bytes.Size       // limit byte size of request body
	OnError http.HandlerFunc // custom error handler
}

// BodyLimit is limiting body size.
func BodyLimit(limit bytes.Size) Middleware {
	return BodyLimitWithConfig(BodyLimitConfig{
		Limit: limit,
	})
}

// BodyLimitWithConfig is limiting body size.
func BodyLimitWithConfig(config BodyLimitConfig) Middleware {
	if config.OnError == nil {
		config.OnError = func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusRequestEntityTooLarge)
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > config.Limit.Int64() {
				config.OnError(w, r)
				return
			}
			r.Body = &bodyLimitReader{
				body:  r.Body,
				limit: config.Limit.Int64(),
				onError: func() {
					config.OnError(w, r)
				},
			}
			next.ServeHTTP(w, r)
		})
	}
}
