package encoding

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	xmath "go.nanasi880.dev/x/math"
)

const (
	contentEncodingKey = "Content-Encoding"
	acceptEncodingKey  = "Accept-Encoding"
)

// ResponseWriteCloser is http.ResponseWriter and io.Closer.
type ResponseWriteCloser interface {
	http.ResponseWriter
	io.Closer
}

// Encoding is an encoding for the compression method used in HTTP.
type Encoding interface {
	// Match returns whether the encoding matches the encoding string.
	Match(encoding string) bool

	// NewReader is creating reader of encoding.
	NewReader(r io.ReadCloser) (io.ReadCloser, error)

	// NewWriter is creating writer of encoding.
	NewWriter(w http.ResponseWriter) (ResponseWriteCloser, error)

	// AddContentEncoding is add Content-Encoding to header.
	AddContentEncoding(header http.Header, requested string)
}

var (
	Compress Encoding = compressEncoding{} // LZW, UNIX compress style.
	Deflate  Encoding = deflateEncoding{}  // deflate of zlib format.
	GZip     Encoding = gzipEncoding{}     // LZ77, UNIX gzip style.
	Identity Encoding = identityEncoding{} // no compression.
)

// AllEncodings returns all supported encoding.
func AllEncodings() []Encoding {
	return []Encoding{
		Compress, Deflate, GZip, Identity,
	}
}

// GetContentEncodings is reading Content-Encoding header.
func GetContentEncodings(header http.Header) []string {
	var result []string

	contentEncodings := header.Values(contentEncodingKey)
	for i := len(contentEncodings); i > 0; i-- {
		encodings := strings.Split(contentEncodings[i-1], ",")
		for j := len(encodings); j > 0; j-- {
			encoding := encodings[j-1]
			encoding = strings.TrimSpace(encoding)
			result = append(result, encoding)
		}
	}

	return result
}

// GetAcceptEncodings is reading Accept-Encoding header.
func GetAcceptEncodings(header http.Header) []string {
	var unordered acceptEncodings

	acceptEncodings := header.Values(acceptEncodingKey)
	for _, encoding := range acceptEncodings {
		encodings := strings.Split(encoding, ",")
		for _, encoding := range encodings {
			if !strings.Contains(encoding, ";") {
				unordered = append(unordered, acceptEncoding{
					encoding:   strings.TrimSpace(encoding),
					weight:     1000,
					incomplete: true,
				})
				continue
			}

			span := strings.Split(encoding, ";")
			if len(span) != 2 {
				// invalid format
				continue
			}

			weight, err := readAcceptEncodingWeight(span[1])
			if err != nil {
				// invalid format
				continue
			}
			weight = xmath.Clamp(weight, 0, 1)

			unordered = append(unordered, acceptEncoding{
				encoding:   strings.TrimSpace(span[0]),
				weight:     int(math.Round(weight * 1000)),
				incomplete: false,
			})
		}
	}

	sort.Stable(unordered)

	result := make([]string, 0, len(unordered))
	for _, a := range unordered {
		result = append(result, a.encoding)
	}

	return result
}

func readAcceptEncodingWeight(s string) (float64, error) {
	span := strings.Split(s, "=")
	if len(span) != 2 {
		return 0, fmt.Errorf("invalid format: %s", s)
	}
	return strconv.ParseFloat(strings.TrimSpace(span[1]), 64)
}
