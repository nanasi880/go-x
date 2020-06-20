package csv

import (
	"bytes"
	"encoding"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"

	xreflect "go.nanasi880.dev/x/reflect"
	"go.nanasi880.dev/x/unsafe"
)

var (
	errInvalidType         = fmt.Errorf("slice of struct or array of struct or struct or pointer of struct only")
	wellKnownEncodingKinds = []reflect.Kind{
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String,
	}
)

func Marshal(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// NewEncoder is create csv encoder.
func NewEncoder(w io.Writer) *Encoder {

	writer := csv.NewWriter(w)

	return &Encoder{
		Comma:          writer.Comma,
		UseCRLF:        writer.UseCRLF,
		UseHeader:      true,
		Nil:            "",
		w:              writer,
		alreadyWritten: false,
		typeCache:      nil,
	}
}

type Marshaler interface {
	MarshalCSV() ([]byte, error)
}

type Encoder struct {
	Comma          rune
	UseCRLF        bool
	UseHeader      bool
	Nil            string
	w              *csv.Writer
	alreadyWritten bool
	typeCache      reflect.Type
}

func (e *Encoder) Encode(v interface{}) (err error) {

	if v == nil {
		return fmt.Errorf("`v` is nil")
	}

	defer func() {
		if err == nil {
			e.w.Flush()
			err = e.w.Error()
		}
	}()

	e.w.Comma = e.Comma
	e.w.UseCRLF = e.UseCRLF

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return e.encodeSlice(rv)
	case reflect.Ptr:
		return e.encodePtr(rv)
	case reflect.Struct:
		return e.encodeStruct(rv)
	default:
		return errInvalidType
	}
}

func (e *Encoder) encodeSlice(v reflect.Value) error {

	l := v.Len()
	for i := 0; i < l; i++ {

		elem := v.Index(i)

		switch elem.Kind() {
		case reflect.Ptr:
			if err := e.encodePtr(elem); err != nil {
				return err
			}
		case reflect.Struct:
			if err := e.encodeStruct(elem); err != nil {
				return err
			}
		default:
			return errInvalidType
		}
	}

	return nil
}

func (e *Encoder) encodePtr(v reflect.Value) error {

	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return errInvalidType
	}

	return e.encodeStruct(elem)
}

func (e *Encoder) encodeStruct(v reflect.Value) error {

	t := v.Type()
	if e.typeCache == nil {
		e.typeCache = t
	}

	if e.typeCache != t {
		return fmt.Errorf("the type cannot be changed during writing")
	}

	if !e.alreadyWritten && e.UseHeader {
		if err := e.writeHeader(t); err != nil {
			return err
		}
	}
	e.alreadyWritten = true

	return e.writeValue(v, t)
}

func (e *Encoder) writeHeader(t reflect.Type) error {

	var header []string

	numFiled := t.NumField()
	for i := 0; i < numFiled; i++ {
		field := t.Field(i)

		headerName := field.Name
		if tag, ok := field.Tag.Lookup("csv"); ok {
			if tag == "-" {
				continue
			}
			headerName = tag
		}

		header = append(header, headerName)
	}

	return e.w.Write(header)
}

func (e *Encoder) writeValue(v reflect.Value, t reflect.Type) error {

	var values []string

	numField := t.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		if tag, ok := field.Tag.Lookup("csv"); ok {
			if tag == "-" {
				continue
			}
		}

		encoded, err := e.encodeValue(v.Field(i))
		if err != nil {
			return err
		}

		values = append(values, encoded)
	}

	return e.w.Write(values)
}

func (e *Encoder) encodeValue(rv reflect.Value) (string, error) {

	v := rv.Interface()

	{
		var (
			marshalFunc func() ([]byte, error)
		)

		switch v := v.(type) {
		case Marshaler:
			marshalFunc = v.MarshalCSV
		case encoding.TextMarshaler:
			marshalFunc = v.MarshalText
		}

		if marshalFunc != nil {
			if xreflect.IsNilable(rv) && rv.IsNil() {
				return e.Nil, nil
			} else {
				encoded, err := marshalFunc()
				if err != nil {
					return "", err
				}
				return unsafe.BytesToString(encoded), nil
			}
		}
	}

	for _, kind := range wellKnownEncodingKinds {
		if rv.Kind() != kind {
			continue
		}

		return fmt.Sprint(v), nil
	}

	return "", errInvalidType
}
