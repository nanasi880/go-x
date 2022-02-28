package encoding

import (
	"io"
	"net/http"
	"strings"
)

type identityWriter struct {
	w http.ResponseWriter
}

// Header is implementation of http.ResponseWriter.
func (w *identityWriter) Header() http.Header {
	return w.w.Header()
}

// Write is implementation of http.ResponseWriter.
func (w *identityWriter) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

// WriteHeader is implementation of http.ResponseWriter.
func (w *identityWriter) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
}

// Close is implementation of io.Closer.
func (w *identityWriter) Close() error {
	return nil
}

type identityEncoding struct{}

// Match is implementation of Encoding.
func (e identityEncoding) Match(encoding string) bool {
	return encoding == "" || strings.EqualFold(encoding, "identity")
}

// NewReader is implementation of Encoding.
func (e identityEncoding) NewReader(r io.ReadCloser) (io.ReadCloser, error) {
	return r, nil
}

// NewWriter is implementation of Encoding.
func (e identityEncoding) NewWriter(w http.ResponseWriter) (ResponseWriteCloser, error) {
	return &identityWriter{w: w}, nil
}

// AddContentEncoding is implementation of Encoding.
func (e identityEncoding) AddContentEncoding(_ http.Header, _ string) {
	// nop
}
