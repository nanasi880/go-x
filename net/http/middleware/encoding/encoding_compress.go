package encoding

import (
	"compress/lzw"
	"io"
	"net/http"
	"strings"
)

type compressEncoding struct{}

// Match is implementation of Encoding.
func (e compressEncoding) Match(encoding string) bool {
	return strings.EqualFold(encoding, "compress")
}

// NewReader is implementation of Encoding.
func (e compressEncoding) NewReader(r io.ReadCloser) (io.ReadCloser, error) {
	return &genericReader{
		parent: r,
		reader: lzw.NewReader(r, lzw.LSB, 8),
	}, nil
}

// NewWriter is implementation of Encoding.
func (e compressEncoding) NewWriter(w http.ResponseWriter) (ResponseWriteCloser, error) {
	return &genericWriter{
		parent: w,
		writer: lzw.NewWriter(w, lzw.LSB, 8),
	}, nil
}

// AddContentEncoding is implementation of Encoding.
func (e compressEncoding) AddContentEncoding(header http.Header, _ string) {
	header.Add(contentEncodingKey, "compress")
}
