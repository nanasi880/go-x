package msgpack

import (
	"bytes"
	"io"
	"reflect"
	"time"
)

// Marshaler is the interface implemented by types that
// can marshal themselves into valid message pack.
type Marshaler interface {
	MarshalMsgPack(e *Encoder) error
}

// Marshal is encode any data to message pack.
func Marshal(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalStringKey is encode any data to message pack.
func MarshalStringKey(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := NewEncoder(buf).SetStructKeyType(StructKeyTypeString).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encoder is message pack encoder.
type Encoder struct {
	w             io.Writer
	bw            io.ByteWriter
	work          []byte
	pointers      pointerSlice
	structKeyType StructKeyType
	structTagName string
}

// NewEncoder is create encoder instance.
func NewEncoder(w io.Writer) *Encoder {
	bw, _ := w.(io.ByteWriter)
	return &Encoder{
		w:             w,
		bw:            bw,
		work:          make([]byte, 16),
		pointers:      make([]uintptr, 0, 10),
		structKeyType: StructKeyTypeInt,
		structTagName: "msgpack",
	}
}

// SetStructKeyType is set StructKeyType to Encoder.
func (e *Encoder) SetStructKeyType(t StructKeyType) *Encoder {
	e.structKeyType = t
	return e
}

// SetStructTagName is set struct tag name to Encoder.
// If tagName is empty, use `msgpack` tag.
func (e *Encoder) SetStructTagName(tagName string) *Encoder {
	if tagName == "" {
		tagName = "msgpack"
	}
	e.structTagName = tagName
	return e
}

// Reset is reset encoder. However, the work buffer will not be reset.
func (e *Encoder) Reset(w io.Writer) {
	bw, _ := w.(io.ByteWriter)
	*e = Encoder{
		w:             w,
		bw:            bw,
		work:          e.work,
		pointers:      e.pointers,
		structKeyType: StructKeyTypeInt,
		structTagName: "msgpack",
	}
}

// Encode is encode data as message pack.
func (e *Encoder) Encode(v interface{}) error {
	if v == nil {
		return e.encodeNil()
	}
	rv := reflect.ValueOf(v)
	return e.encodeValue(rv)
}

// EncodeNil is encode Nil as message pack.
func (e *Encoder) EncodeNil() error {
	return e.encodeNil()
}

// EncodeBool is encode bool as message pack.
func (e *Encoder) EncodeBool(v bool) error {
	return e.encodeBool(v)
}

// EncodeInt64 is encode int64 as message pack.
func (e *Encoder) EncodeInt64(v int64) error {
	return e.encodeInt(v)
}

// EncodeUint64 is encode uint64 as message pack.
func (e *Encoder) EncodeUint64(v uint64) error {
	return e.encodeUint(v)
}

// EncodeFloat32 is encode float32 as message pack.
func (e *Encoder) EncodeFloat32(v float32) error {
	return e.encodeFloat32(v)
}

// EncodeFloat64 is encode float64 as message pack.
func (e *Encoder) EncodeFloat64(v float64) error {
	return e.encodeFloat64(v)
}

// EncodeBin is encode binary as message pack.
func (e *Encoder) EncodeBin(v []byte) error {
	return e.encodeBin(v)
}

// EncodeString is encode string as message pack.
func (e *Encoder) EncodeString(v string) error {
	return e.encodeString(v)
}

// EncodeTime is encode time.Time as message pack.
func (e *Encoder) EncodeTime(v time.Time) error {
	return e.encodeTime(v)
}

// EncodeArrayHeader is encode array header as message pack.
func (e *Encoder) EncodeArrayHeader(length uint32) error {
	return e.encodeArrayHeader(length)
}

// EncodeMapHeader is encode map header as message pack.
func (e *Encoder) EncodeMapHeader(length uint32) error {
	return e.encodeMapHeader(length)
}
