package msgpack

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"

	xbytes "go.nanasi880.dev/x/bytes"
	xunsafe "go.nanasi880.dev/x/unsafe"
)

// Unmarshaler is the interface implemented by types that
// can unmarshal themselves into valid message pack.
type Unmarshaler interface {
	UnmarshalMsgPack(d *Decoder) error
}

// Unmarshal is decode message pack data to any data.
func Unmarshal(b []byte, v interface{}) error {
	return NewDecoderBytes(b).Decode(v)
}

// UnmarshalStringKey is decode message pack data to any data.
func UnmarshalStringKey(b []byte, v interface{}) error {
	return NewDecoderBytes(b).SetStructKeyType(StructKeyTypeString).Decode(v)
}

// Decoder is message pack decoder.
type Decoder struct {
	data                 []byte
	reader               *xbytes.BinaryReader
	work                 []byte
	structKeyType        StructKeyType
	arrayLengthTolerance ArrayLengthTolerance
	timeZone             *time.Location
}

// NewDecoder is create decoder instance.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		data:                 nil,
		reader:               xbytes.NewBinaryReader(r),
		work:                 make([]byte, readBlockSize),
		structKeyType:        StructKeyTypeInt,
		arrayLengthTolerance: ArrayLengthToleranceLessThanOrEqual,
		timeZone:             nil,
	}
}

// NewDecoderBytes is create decoder instance.
func NewDecoderBytes(in []byte) *Decoder {
	if in == nil {
		in = make([]byte, 0)
	}
	return &Decoder{
		data:                 in,
		reader:               nil,
		work:                 make([]byte, readBlockSize),
		structKeyType:        StructKeyTypeInt,
		arrayLengthTolerance: ArrayLengthToleranceLessThanOrEqual,
		timeZone:             nil,
	}
}

// Reset is reset decoder. However, the work buffer will not be reset.
func (d *Decoder) Reset(r io.Reader) {
	*d = Decoder{
		data:                 nil,
		reader:               d.reader,
		work:                 d.work,
		structKeyType:        StructKeyTypeInt,
		arrayLengthTolerance: ArrayLengthToleranceLessThanOrEqual,
		timeZone:             nil,
	}
	if d.reader == nil {
		d.reader = xbytes.NewBinaryReader(r)
	}
	d.reader.Reset(r)
}

// ResetBytes is reset decoder. However, the work buffer will not be reset.
func (d *Decoder) ResetBytes(in []byte) {
	if in == nil {
		in = make([]byte, 0)
	}
	*d = Decoder{
		data:                 in,
		reader:               d.reader,
		work:                 d.work,
		structKeyType:        StructKeyTypeInt,
		arrayLengthTolerance: ArrayLengthToleranceLessThanOrEqual,
		timeZone:             nil,
	}
	if d.reader != nil {
		d.reader.Reset(nil)
	}
}

// SetStructKeyType is set StructKeyType to Decoder.
func (d *Decoder) SetStructKeyType(t StructKeyType) *Decoder {
	d.structKeyType = t
	return d
}

// SetArrayLengthTolerance is set ArrayLengthTolerance to Decoder.
func (d *Decoder) SetArrayLengthTolerance(tolerance ArrayLengthTolerance) *Decoder {
	d.arrayLengthTolerance = tolerance
	return d
}

// SetTimeZone is set time zone to Decoder.
// The decoder will set this time zone to the time when decoding. If loc is nil, use the UTC time.
func (d *Decoder) SetTimeZone(loc *time.Location) *Decoder {
	d.timeZone = loc
	return d
}

// Decode is decode data from message pack.
func (d *Decoder) Decode(v interface{}) (e error) {
	defer func() {
		r := recover()
		if r != nil {
			e = fmt.Errorf("%v", r)
		}
	}()
	if v == nil {
		return fmt.Errorf("nil")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("rv.kind() != reflect.Ptr")
	}
	return d.decodeValue(rv)
}

// DecodeFormat is decode format from message pack.
func (d *Decoder) DecodeFormat() (Format, error) {
	name, raw, err := d.readFormat()
	return Format{
		Name: name,
		Raw:  raw,
	}, err
}

// DecodeInt64 is decode integer type as int64 from message pack.
func (d *Decoder) DecodeInt64(format Format) (int64, error) {
	i, err := d.decodeIntBit(format.Name, format.Raw)
	return int64(i), err
}

// DecodeUint64 is decode integer type as uint64 from message pack.
func (d *Decoder) DecodeUint64(format Format) (uint64, error) {
	return d.decodeIntBit(format.Name, format.Raw)
}

// DecodeFloat32 is decode float type as float32 from message pack.
func (d *Decoder) DecodeFloat32(format Format) (float32, error) {
	switch format.Name {
	case Float32:
		val, err := d.read(4)
		if err != nil {
			return 0, err
		}
		return math.Float32frombits(binary.BigEndian.Uint32(val)), nil
	case Float64:
		val, err := d.read(8)
		if err != nil {
			return 0, err
		}
		return float32(math.Float64frombits(binary.BigEndian.Uint64(val))), nil
	default:
		return 0, fmt.Errorf("invalid format: %s", format.String())
	}
}

// DecodeFloat64 is decode float type as float64 from message pack.
func (d *Decoder) DecodeFloat64(format Format) (float64, error) {
	switch format.Name {
	case Float32:
		val, err := d.read(4)
		if err != nil {
			return 0, err
		}
		return float64(math.Float32frombits(binary.BigEndian.Uint32(val))), nil
	case Float64:
		val, err := d.read(8)
		if err != nil {
			return 0, err
		}
		return math.Float64frombits(binary.BigEndian.Uint64(val)), nil
	default:
		return 0, fmt.Errorf("invalid format: %s", format.String())
	}
}

// DecodeString is decode string type or bin type as string from message pack.
func (d *Decoder) DecodeString(format Format) (string, error) {
	var (
		length int
		err    error
	)
	switch format.Name {
	case Bin8, Bin16, Bin32:
		length, err = d.decodeBinHeader(format.Name)
	case FixStr, Str8, Str16, Str32:
		length, err = d.decodeStringHeader(format.Name, format.Raw)
	default:
		return "", fmt.Errorf("invalid format: %s", format.String())
	}
	if err != nil {
		return "", err
	}

	tmp, err := d.decodeBytes(length)
	if err != nil {
		return "", err
	}
	return xunsafe.BytesToString(tmp), nil
}

// DecodeTime is decode time from message pack.
func (d *Decoder) DecodeTime(header ExtHeader) (time.Time, error) {
	return d.decodeTime(header.Length)
}

// DecodeExtHeader is decode ext header from message pack.
func (d *Decoder) DecodeExtHeader(format Format) (ExtHeader, error) {
	typeCode, length, err := d.decodeExtHeader(format.Name)
	return ExtHeader{
		Format: format,
		Type:   typeCode,
		Length: length,
	}, err
}

// DecodeArrayHeader is decode array header from message pack.
func (d *Decoder) DecodeArrayHeader(format Format) (ArrayHeader, error) {
	length, err := d.decodeArrayHeader(format.Name, format.Raw)
	return ArrayHeader{
		Format: format,
		Length: length,
	}, err
}

// DecodeMapHeader is decode map header from message pack.
func (d *Decoder) DecodeMapHeader(format Format) (MapHeader, error) {
	length, err := d.decodeMapHeader(format.Name, format.Raw)
	return MapHeader{
		Format: format,
		Length: length,
	}, err
}

// SkipObject is seek current message pack object.
func (d *Decoder) SkipObject(format Format) error {
	return d.skipObject(format.Name, format.Raw)
}
