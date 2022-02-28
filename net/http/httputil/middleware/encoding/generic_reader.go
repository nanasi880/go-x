package encoding

import "io"

type genericReader struct {
	parent io.ReadCloser
	reader io.ReadCloser
}

func (r *genericReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *genericReader) Close() error {
	if r.reader == nil {
		return nil
	}

	err1 := r.parent.Close()
	err2 := r.reader.Close()
	r.parent = nil
	r.reader = nil

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	return nil
}
