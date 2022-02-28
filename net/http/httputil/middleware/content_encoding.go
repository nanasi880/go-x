package middleware

import (
	"net/http"

	"go.nanasi880.dev/x/net/http/httputil/middleware/encoding"
)

// ContentEncoding decodes the request data according to the Content-Encoding header, and encodes the response data according to the Accept-Encoding header.
func ContentEncoding(encodings ...encoding.Encoding) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			body := r.Body
			for _, enc := range encoding.GetContentEncodings(r.Header) {
				for _, encoder := range encodings {
					if !encoder.Match(enc) {
						continue
					}

					r, err := encoder.NewReader(body)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					body = r
					break
				}
			}
			r.Body = body

			var rw encoding.ResponseWriteCloser
			for _, enc := range encoding.GetAcceptEncodings(r.Header) {
				for _, encoder := range encodings {
					if !encoder.Match(enc) {
						continue
					}

					writer, err := encoder.NewWriter(w)
					if err != nil {
						break
					}
					rw = writer
					w = rw

					encoder.AddContentEncoding(w.Header(), enc)
					break
				}
				if rw != nil {
					break
				}
			}
			if rw != nil {
				defer func() {
					_ = rw.Close()
				}()
			}

			next.ServeHTTP(w, r)
		})
	}
}
