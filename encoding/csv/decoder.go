package csv

import (
	"bytes"
	"encoding"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	xreflect "go.nanasi880.dev/x/reflect"
	"go.nanasi880.dev/x/unsafe"
)

var (
	errInvalidDecodeType = fmt.Errorf("pointer of slice of struct (or struct pointer) only")
	intCastMap           = map[reflect.Kind]func(v int64) interface{}{
		reflect.Int8:  func(v int64) interface{} { return int8(v) },
		reflect.Int16: func(v int64) interface{} { return int16(v) },
		reflect.Int32: func(v int64) interface{} { return int32(v) },
		reflect.Int64: func(v int64) interface{} { return v },
		reflect.Int:   func(v int64) interface{} { return int(v) },
	}
	uintCastMap = map[reflect.Kind]func(v uint64) interface{}{
		reflect.Uint8:  func(v uint64) interface{} { return uint8(v) },
		reflect.Uint16: func(v uint64) interface{} { return uint16(v) },
		reflect.Uint32: func(v uint64) interface{} { return uint32(v) },
		reflect.Uint64: func(v uint64) interface{} { return v },
		reflect.Uint:   func(v uint64) interface{} { return uint(v) },
	}
)

// Unmarshaler is the interface implemented by types
// that can unmarshal a CSV description of themselves.
// The input can be assumed to be a valid encoding of
// a CSV value. UnmarshalCSV must copy the CSV data
// if it wishes to retain the data after returning.
type Unmarshaler interface {
	UnmarshalCSV([]byte) error
}

// UnmarshalString is decodes a slice of a structure from CSV string.
func UnmarshalString(csv string, out interface{}) error {
	return Unmarshal(unsafe.StringToBytes(csv), out)
}

// Unmarshal is decodes a slice of a structure from CSV data.
func Unmarshal(csv []byte, out interface{}) error {
	return NewDecoder(bytes.NewReader(csv)).Decode(out)
}

// NewDecoder is create csv decoder.
func NewDecoder(r io.Reader) *Decoder {

	reader := csv.NewReader(r)

	return &Decoder{
		Comma:            reader.Comma,
		Comment:          reader.Comment,
		FieldsPerRecord:  reader.FieldsPerRecord,
		LazyQuotes:       reader.LazyQuotes,
		TrimLeadingSpace: reader.TrimLeadingSpace,
		ReuseRecord:      reader.ReuseRecord,
		UseHeader:        true,
		Nil:              "",
		r:                reader,
	}
}

// An Decoder reads CSV values to an input stream.
type Decoder struct {
	Comma            rune
	Comment          rune
	FieldsPerRecord  int
	LazyQuotes       bool
	TrimLeadingSpace bool
	ReuseRecord      bool
	UseHeader        bool
	Nil              string
	r                *csv.Reader
}

type decodeElement struct {
	elem   reflect.Value
	access reflect.Value
}

// Decode is decodes a slice of a structure from CSV data, CSV data read from the io.Reader specified by NewDecoder.
func (d *Decoder) Decode(out interface{}) error {

	if out == nil {
		return fmt.Errorf("nil")
	}

	d.r.Comma = d.Comma
	d.r.Comment = d.Comment
	d.r.FieldsPerRecord = d.FieldsPerRecord
	d.r.LazyQuotes = d.LazyQuotes
	d.r.TrimLeadingSpace = d.TrimLeadingSpace
	d.r.ReuseRecord = d.ReuseRecord

	sliceElemType, err := d.getValueType(out)
	if err != nil {
		return err
	}

	outSlice := reflect.ValueOf(out).Elem()

	return d.decodeRows(outSlice, sliceElemType)
}

func (d *Decoder) decodeRows(out reflect.Value, elemType reflect.Type) error {

	// csv column index : struct field index
	var fieldIndex map[int]int
	if d.UseHeader {
		header, err := d.r.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fieldIndex = d.getFieldIndexByTag(elemType, header)
	} else {
		fieldIndex = d.getFieldIndexByOrder(elemType)
	}

	result := reflect.MakeSlice(out.Type(), 0, 0)
	for {
		record, err := d.r.Read()
		if err == io.EOF {
			out.Set(result)
			return nil
		}
		if err != nil {
			return err
		}

		elem := d.allocElem(elemType)
		for i, raw := range record {

			fi, ok := fieldIndex[i]
			if !ok {
				continue
			}

			field := elem.access.Field(fi)
			if err := d.decodeValue(raw, field); err != nil {
				return err
			}
		}

		result = reflect.Append(result, elem.elem)
	}
}

func (d *Decoder) decodeValue(raw string, rv reflect.Value) error {

	if rv.Kind() == reflect.Ptr && raw == d.Nil {
		if !rv.IsNil() {
			rv.Set(xreflect.ZeroValue(rv.Type()))
		}
		return nil
	}

	rawBytes := unsafe.StringToBytes(raw)

	access := rv
	if rv.Kind() == reflect.Ptr {
		rv.Set(reflect.New(rv.Type().Elem()))
		access = rv.Elem()
	}

	if i, ok := rv.Interface().(Unmarshaler); ok {
		return i.UnmarshalCSV(rawBytes)
	}
	if i, ok := rv.Interface().(encoding.TextUnmarshaler); ok {
		return i.UnmarshalText(rawBytes)
	}

	switch kind := access.Kind(); kind {

	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		access.Set(reflect.ValueOf(v))
		return nil

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v, err := strconv.ParseInt(raw, 10, xreflect.IntegerBitSize(kind))
		if err != nil {
			return err
		}
		access.Set(reflect.ValueOf(intCastMap[kind](v)))
		return nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v, err := strconv.ParseUint(raw, 10, xreflect.IntegerBitSize(kind))
		if err != nil {
			return err
		}
		access.Set(reflect.ValueOf(uintCastMap[kind](v)))
		return nil

	case reflect.Float32:
		v, err := strconv.ParseFloat(raw, 32)
		if err != nil {
			return err
		}
		access.Set(reflect.ValueOf(float32(v)))
		return nil

	case reflect.Float64:
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return err
		}
		access.Set(reflect.ValueOf(v))
		return nil

	case reflect.Complex64:
		var v complex64
		if err := d.decodeComplex(raw, &v); err != nil {
			return err
		}
		access.Set(reflect.ValueOf(v))
		return nil

	case reflect.Complex128:
		var v complex128
		if err := d.decodeComplex(raw, &v); err != nil {
			return err
		}
		access.Set(reflect.ValueOf(v))
		return nil

	case reflect.String:
		access.Set(reflect.ValueOf(raw))
		return nil

	default:
		return fmt.Errorf("cannot decode")
	}
}

func (d *Decoder) getFieldIndexByOrder(t reflect.Type) map[int]int {

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var (
		numField = t.NumField()
		result   = make(map[int]int)
		j        = 0
	)
	for i := 0; i < numField; i++ {

		field := t.Field(i)
		if tag, ok := field.Tag.Lookup("csv"); ok {
			if tag == "-" {
				continue
			}
		}

		result[j] = i
		j++
	}

	return result
}

func (d *Decoder) getFieldIndexByTag(t reflect.Type, header []string) map[int]int {

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var (
		numField = t.NumField()
		result   = make(map[int]int)
		last     = make(map[string]int)
	)
	for i := 0; i < numField; i++ {

		field := t.Field(i)
		name := field.Name
		if tag, ok := field.Tag.Lookup("csv"); ok {
			if tag == "-" {
				continue
			}
			name = tag
		}

		j, ok := last[name]
		if ok {
			j++
		}
		for ; j < len(header); j++ {
			if header[j] == name {
				result[j] = i
				last[name] = j
				break
			}
		}
	}

	return result
}

func (d *Decoder) getValueType(v interface{}) (reflect.Type, error) {

	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr {
		return nil, errInvalidDecodeType
	}

	slice := rv.Elem()
	if slice.Kind() != reflect.Slice {
		return nil, errInvalidDecodeType
	}

	sliceElemType := slice.Type().Elem()
	if err := d.validateValueType(sliceElemType); err != nil {
		return nil, err
	}

	return sliceElemType, nil
}

func (d *Decoder) validateValueType(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Struct:
		return nil
	case reflect.Ptr:
		if t.Elem().Kind() == reflect.Struct {
			return nil
		}
		return errInvalidDecodeType
	default:
		return errInvalidDecodeType
	}
}

func (d *Decoder) allocElem(t reflect.Type) decodeElement {

	if t.Kind() == reflect.Ptr {
		rv := reflect.New(t.Elem())
		return decodeElement{
			elem:   rv,
			access: rv.Elem(),
		}
	}

	rv := xreflect.ZeroValue(t)
	return decodeElement{
		elem:   rv,
		access: rv,
	}
}

func (d *Decoder) parseComplex(raw string) (string, string, error) {

	// (1+2i)
	s := strings.TrimLeft(raw, "(")
	s = strings.TrimRight(s, "i)")

	part := strings.Split(s, "+")
	if len(part) != 2 {
		return "", "", fmt.Errorf("invalid complex format: %s", raw)
	}

	return part[0], part[1], nil
}

func (d *Decoder) decodeComplex(raw string, out interface{}) error {

	rv := reflect.ValueOf(out).Elem()

	bitSize := 64
	if rv.Kind() == reflect.Complex64 {
		bitSize = 32
	}

	realS, imagS, err := d.parseComplex(raw)
	if err != nil {
		return err
	}

	realF, err := strconv.ParseFloat(realS, bitSize)
	if err != nil {
		return err
	}
	imagF, err := strconv.ParseFloat(imagS, bitSize)
	if err != nil {
		return err
	}

	if rv.Kind() == reflect.Complex64 {
		var v = complex(float32(realF), float32(imagF))
		rv.Set(reflect.ValueOf(v))
	} else {
		var v = complex(realF, imagF)
		rv.Set(reflect.ValueOf(v))
	}

	return nil
}
