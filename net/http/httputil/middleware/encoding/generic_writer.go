package encoding

import (
	"io"
	"net/http"
)

type genericWriter struct {
	parent http.ResponseWriter
	writer io.WriteCloser
}

func (w *genericWriter) Header() http.Header {
	return w.parent.Header()
}

func (w *genericWriter) Write(bytes []byte) (int, error) {
	return w.writer.Write(bytes)
}

func (w *genericWriter) WriteHeader(statusCode int) {
	w.parent.WriteHeader(statusCode)
}

func (w *genericWriter) Close() error {
	if w.writer == nil {
		return nil
	}

	err := w.writer.Close()
	w.writer = nil

	return err
}
