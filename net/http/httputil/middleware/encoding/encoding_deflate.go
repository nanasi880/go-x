package encoding

import (
	"compress/zlib"
	"io"
	"net/http"
	"strings"
)

type deflateEncoding struct{}

// Match is implementation of Encoding.
func (e deflateEncoding) Match(encoding string) bool {
	return strings.EqualFold(encoding, "deflate")
}

// NewReader is implementation of Encoding.
func (e deflateEncoding) NewReader(r io.ReadCloser) (io.ReadCloser, error) {
	z, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &genericReader{
		parent: r,
		reader: z,
	}, nil
}

// NewWriter is implementation of Encoding.
func (e deflateEncoding) NewWriter(w http.ResponseWriter) (ResponseWriteCloser, error) {
	return &genericWriter{
		parent: w,
		writer: zlib.NewWriter(w),
	}, nil
}

// AddContentEncoding is implementation of Encoding.
func (e deflateEncoding) AddContentEncoding(header http.Header, _ string) {
	header.Add(contentEncodingKey, "deflate")
}
