package encoding

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipEncoding struct{}

// Match is implementation of Encoding.
func (e gzipEncoding) Match(encoding string) bool {
	return strings.EqualFold(encoding, "gzip") || strings.EqualFold(encoding, "x-gzip")
}

// NewReader is implementation of Encoding.
func (e gzipEncoding) NewReader(r io.ReadCloser) (io.ReadCloser, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &genericReader{
		parent: r,
		reader: gz,
	}, nil
}

// NewWriter is implementation of Encoding.
func (e gzipEncoding) NewWriter(w http.ResponseWriter) (ResponseWriteCloser, error) {
	return &genericWriter{
		parent: w,
		writer: gzip.NewWriter(w),
	}, nil
}

// AddContentEncoding is implementation of Encoding.
func (e gzipEncoding) AddContentEncoding(header http.Header, requested string) {
	header.Add(contentEncodingKey, requested)
}
