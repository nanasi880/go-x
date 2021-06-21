package bytes

import (
	"fmt"
	"io"
)

// BinaryReader is id.BinaryReader
type BinaryReader struct {
	reader io.Reader
	buf    []byte
	bufLen int
	pos    int
	err    error
}

// NewBinaryReader is create BinaryReader instance
func NewBinaryReader(r io.Reader) *BinaryReader {
	return NewBinaryReaderBuffer(r, nil)
}

// NewBinaryReaderBuffer is create BinaryReader instance with work buffer
func NewBinaryReaderBuffer(r io.Reader, buf []byte) *BinaryReader {
	if len(buf) == 0 {
		buf = make([]byte, 1024)
	}
	return &BinaryReader{
		reader: r,
		buf:    buf,
		bufLen: 0,
		pos:    0,
		err:    nil,
	}
}

// Reset is reset binary reader.
func (r *BinaryReader) Reset(reader io.Reader) {
	*r = BinaryReader{
		reader: reader,
		buf:    r.buf,
		bufLen: 0,
		pos:    0,
		err:    nil,
	}
}

// Read is implements io.Reader
func (r *BinaryReader) Read(p []byte) (int, error) {
	var (
		total int
	)
	for {
		n := copy(p[total:], r.buf[r.pos:r.bufLen])
		total += n
		r.pos += n

		if total == len(p) {
			return len(p), nil
		}

		err := r.read()
		if err != nil {
			return total, err
		}
		if r.bufLen == 0 {
			return total, io.ErrUnexpectedEOF
		}
	}
}

// Read1 is read 1 byte
func (r *BinaryReader) Read1() (byte, error) {
	if r.pos+1 <= r.bufLen {
		b := r.buf[r.pos]
		r.pos++
		return b, nil
	}

	if err := r.read(); err != nil {
		return 0, err
	}

	if r.bufLen < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	r.pos = 1
	return r.buf[0], nil
}

// Read2 is read 2 byte
func (r *BinaryReader) Read2() ([2]byte, error) {
	var (
		buf   [2]byte
		total int
	)
	for {
		n := copy(buf[total:], r.buf[r.pos:r.bufLen])
		total += n
		r.pos += n

		if total >= len(buf) {
			return buf, nil
		}

		err := r.read()
		if err != nil {
			return buf, err
		}
		if r.bufLen == 0 {
			return buf, io.ErrUnexpectedEOF
		}
	}
}

// Read4 is read 4 byte
func (r *BinaryReader) Read4() ([4]byte, error) {
	var (
		buf   [4]byte
		total int
	)
	for {
		n := copy(buf[total:], r.buf[r.pos:r.bufLen])
		total += n
		r.pos += n

		if total >= len(buf) {
			return buf, nil
		}

		err := r.read()
		if err != nil {
			return buf, err
		}
		if r.bufLen == 0 {
			return buf, io.ErrUnexpectedEOF
		}
	}
}

// Read8 is read 8 byte
func (r *BinaryReader) Read8() ([8]byte, error) {
	var (
		buf   [8]byte
		total int
	)
	for {
		n := copy(buf[total:], r.buf[r.pos:r.bufLen])
		total += n
		r.pos += n

		if total >= len(buf) {
			return buf, nil
		}

		err := r.read()
		if err != nil {
			return buf, err
		}
		if r.bufLen == 0 {
			return buf, io.ErrUnexpectedEOF
		}
	}
}

// Read16 is read 16 byte
func (r *BinaryReader) Read16() ([16]byte, error) {
	var (
		buf   [16]byte
		total int
	)
	for {
		n := copy(buf[total:], r.buf[r.pos:r.bufLen])
		total += n
		r.pos += n

		if total >= len(buf) {
			return buf, nil
		}

		err := r.read()
		if err != nil {
			return buf, err
		}
		if r.bufLen == 0 {
			return buf, io.ErrUnexpectedEOF
		}
	}
}

// Read32 is read 32 byte
func (r *BinaryReader) Read32() ([32]byte, error) {
	var (
		buf   [32]byte
		total int
	)
	for {
		n := copy(buf[total:], r.buf[r.pos:r.bufLen])
		total += n
		r.pos += n

		if total >= len(buf) {
			return buf, nil
		}

		err := r.read()
		if err != nil {
			return buf, err
		}
		if r.bufLen == 0 {
			return buf, io.ErrUnexpectedEOF
		}
	}
}

// Read64 is read 64 byte
func (r *BinaryReader) Read64() ([64]byte, error) {
	var (
		buf   [64]byte
		total int
	)
	for {
		n := copy(buf[total:], r.buf[r.pos:r.bufLen])
		total += n
		r.pos += n

		if total >= len(buf) {
			return buf, nil
		}

		err := r.read()
		if err != nil {
			return buf, err
		}
		if r.bufLen == 0 {
			return buf, io.ErrUnexpectedEOF
		}
	}
}

func (r *BinaryReader) Seek(n int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		return 0, fmt.Errorf("support `io.SeekCurrent` only")
	}

	total := n
	for {
		remaining := r.bufLen - r.pos
		if total <= int64(remaining) {
			r.pos += int(total)
			return n, nil
		}
		total -= int64(remaining)

		err := r.read()
		if err != nil {
			return n - total, err
		}
		if r.bufLen == 0 {
			return n - total, io.ErrUnexpectedEOF
		}
	}
}

func (r *BinaryReader) read() error {
	if r.err != nil {
		return r.err
	}

	n, err := r.reader.Read(r.buf)
	r.pos = 0
	r.bufLen = n
	r.err = err

	if err != io.EOF {
		return err
	}
	return nil
}
