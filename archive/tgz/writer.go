package tgz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

// Writer is tar.gz encoder.
type Writer struct {
	w       io.Writer
	entries []Entry
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

// Add is add Entry to archive.
func (e *Writer) Add(entries ...Entry) *Writer {
	e.entries = append(e.entries, entries...)
	return e
}

// Write is archive and write to tar.gz.
func (e *Writer) Write() (err error) {

	gz := gzip.NewWriter(e.w)
	defer func() {
		e := gz.Close()
		if e != nil {
			err = fmt.Errorf("%w: %v", err, e)
		}
	}()

	tw := tar.NewWriter(gz)
	defer func() {
		e := tw.Close()
		if e != nil {
			err = fmt.Errorf("%w: %v", err, e)
		}
	}()

	for _, entry := range e.entries {
		err := entry.process(tw)
		if err != nil {
			return err
		}
	}

	return nil
}
