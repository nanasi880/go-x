package msgpack

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"time"
	"unsafe"

	cache "go.nanasi880.dev/x/encoding/msgpack/internal/index"
	xreflect "go.nanasi880.dev/x/reflect"
	xunsafe "go.nanasi880.dev/x/unsafe"
)

func (e *Encoder) encodeValue(rv reflect.Value) error {
	if xreflect.IsNil(rv) {
		return e.EncodeNil()
	}
	if rv.Type().Implements(marshalerType) {
		return e.encodeMarshaler(rv)
	}
	switch rv.Kind() {
	case reflect.Bool:
		return e.encodeBool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.encodeInt(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.encodeUint(rv.Uint())
	case reflect.Float32:
		return e.encodeFloat32(float32(rv.Float()))
	case reflect.Float64:
		return e.encodeFloat64(rv.Float())
	case reflect.Array:
		return e.encodeSlice(rv.Slice(0, rv.Len()))
	case reflect.Slice:
		return e.encodeSlice(rv)
	case reflect.Map:
		return e.encodeMap(rv)
	case reflect.Ptr:
		return e.encodePtr(rv)
	case reflect.String:
		return e.encodeString(rv.String())
	case reflect.Struct:
		return e.encodeStruct(rv)
	default:
		return fmt.Errorf("%v is unsupported type", rv.Kind())
	}
}

func (e *Encoder) encodeNil() error {
	return e.writeByte(Nil.Byte())
}

func (e *Encoder) encodeMarshaler(rv reflect.Value) error {

	// zero allocation pass
	if rv.CanAddr() {
		p := unsafe.Pointer(rv.Convert(marshalerType).UnsafeAddr())
		i := *(*Marshaler)(p)
		return i.MarshalMsgPack(e)
	}

	return rv.Interface().(Marshaler).MarshalMsgPack(e)
}

func (e *Encoder) encodeBool(b bool) error {
	if b {
		return e.writeByte(True.Byte())
	} else {
		return e.writeByte(False.Byte())
	}
}

func (e *Encoder) encodeInt(v int64) error {
	// PositiveFixInt
	if v >= 0 && v <= 0x7F {
		return e.writeByte(byte(int8(v)))
	}
	// NegativeFixInt
	if v >= -32 && v <= -1 {
		return e.writeByte(byte(int8(v)))
	}
	// Int8
	if v >= math.MinInt8 && v <= math.MaxInt8 {
		data := e.work[:2]
		data[0] = Int8.Byte()
		data[1] = byte(int8(v))
		return e.write(data)
	}
	// Int16
	if v >= math.MinInt16 && v <= math.MaxInt16 {
		data := e.work[:3]
		data[0] = Int16.Byte()
		binary.BigEndian.PutUint16(data[1:], uint16(int16(v)))
		return e.write(data)
	}
	// Int32
	if v >= math.MinInt32 && v <= math.MaxInt32 {
		data := e.work[:5]
		data[0] = Int32.Byte()
		binary.BigEndian.PutUint32(data[1:], uint32(int32(v)))
		return e.write(data)
	}
	// Int64
	data := e.work[:9]
	data[0] = Int64.Byte()
	binary.BigEndian.PutUint64(data[1:], uint64(v))
	return e.write(data)
}

func (e *Encoder) encodeUint(v uint64) error {
	// PositiveFixInt
	if v <= 0x7F {
		return e.writeByte(byte(v))
	}
	// Uint8
	if v <= math.MaxUint8 {
		data := e.work[:2]
		data[0] = Uint8.Byte()
		data[1] = byte(v)
		return e.write(data)
	}
	// Uint16
	if v <= math.MaxUint16 {
		data := e.work[:3]
		data[0] = Uint16.Byte()
		binary.BigEndian.PutUint16(data[1:], uint16(v))
		return e.write(data)
	}
	// Uint32
	if v <= math.MaxUint16 {
		data := e.work[:5]
		data[0] = Uint32.Byte()
		binary.BigEndian.PutUint32(data[1:], uint32(v))
		return e.write(data)
	}
	// Uint64
	data := e.work[:9]
	data[0] = Uint64.Byte()
	binary.BigEndian.PutUint64(data[1:], v)
	return e.write(data)
}

func (e *Encoder) encodeFloat32(v float32) error {
	buf := e.work[:5]
	buf[0] = Float32.Byte()
	binary.BigEndian.PutUint32(buf[1:], math.Float32bits(v))
	return e.write(buf)
}

func (e *Encoder) encodeFloat64(v float64) error {
	buf := e.work[:9]
	buf[0] = Float64.Byte()
	binary.BigEndian.PutUint64(buf[1:], math.Float64bits(v))
	return e.write(buf)
}

func (e *Encoder) encodeSlice(rv reflect.Value) error {
	if rv.Type().ConvertibleTo(byteSliceType) {
		return e.encodeBin(rv.Convert(byteSliceType).Bytes())
	}
	return e.encodeSliceToArray(rv)
}

func (e *Encoder) encodeSliceToArray(rv reflect.Value) error {
	length := rv.Len()
	err := e.encodeArrayHeaderInt(length)
	if err != nil {
		return err
	}
	for i := 0; i < length; i++ {
		err := e.encodeValue(rv.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeBin(v []byte) error {
	err := e.encodeBinHeader(len(v))
	if err != nil {
		return err
	}
	return e.write(v)
}

func (e *Encoder) encodeArrayHeaderInt(length int) error {
	if length > math.MaxUint32 {
		return fmt.Errorf("too long")
	}
	return e.encodeArrayHeader(uint32(uint(length)))
}

func (e *Encoder) encodeArrayHeader(length uint32) error {
	if length < 16 {
		return e.writeByte(byte(length) | FixArray.Byte())
	}
	if length <= math.MaxUint16 {
		buf := e.work[:3]
		buf[0] = Array16.Byte()
		binary.BigEndian.PutUint16(buf[1:], uint16(length))
		return e.write(buf)
	}
	buf := e.work[:5]
	buf[0] = Array32.Byte()
	binary.BigEndian.PutUint32(buf[1:], length)
	return e.write(buf)
}

func (e *Encoder) encodeBinHeader(length int) error {
	if length <= math.MaxUint8 {
		bin := e.work[:2]
		bin[0] = Bin8.Byte()
		bin[1] = byte(uint(length))
		return e.write(bin)
	}
	if length <= math.MaxUint16 {
		bin := e.work[:3]
		bin[0] = Bin16.Byte()
		binary.BigEndian.PutUint16(bin[1:], uint16(uint(length)))
		return e.write(bin)
	}
	if length <= math.MaxUint32 {
		bin := e.work[:5]
		bin[0] = Bin32.Byte()
		binary.BigEndian.PutUint32(bin[1:], uint32(uint(length)))
		return e.write(bin)
	}
	return fmt.Errorf("too long")
}

func (e *Encoder) encodeMap(rv reflect.Value) error {
	err := e.encodeMapHeaderInt(rv.Len())
	if err != nil {
		return err
	}
	it := rv.MapRange()
	for it.Next() {
		if err := e.encodeValue(it.Key()); err != nil {
			return err
		}
		if err := e.encodeValue(it.Value()); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeMapHeaderInt(length int) error {
	if length > math.MaxUint32 {
		return fmt.Errorf("too long")
	}
	return e.encodeMapHeader(uint32(uint(length)))
}

func (e *Encoder) encodeMapHeader(length uint32) error {
	if length < 16 {
		buf := byte(uint(length)) | FixMap.Byte()
		return e.writeByte(buf)
	}
	if length <= math.MaxUint16 {
		buf := e.work[:3]
		buf[0] = Map16.Byte()
		binary.BigEndian.PutUint16(buf[1:], uint16(length))
		return e.write(buf)
	}
	buf := e.work[:5]
	buf[0] = Map32.Byte()
	binary.BigEndian.PutUint32(buf[1:], length)
	return e.write(buf)
}

func (e *Encoder) encodePtr(rv reflect.Value) error {
	ptr := rv.Pointer()
	if e.pointers.contains(ptr) {
		return fmt.Errorf("cyclic pointer")
	}

	e.pointers = append(e.pointers, ptr)
	defer func() {
		e.pointers.pop()
	}()

	return e.encodeValue(rv.Elem())
}

func (e *Encoder) encodeString(v string) error {
	err := e.encodeStringHeader(len(v))
	if err != nil {
		return err
	}
	bin := xunsafe.StringToBytes(v)
	return e.write(bin)
}

func (e *Encoder) encodeStringHeader(length int) error {
	if length < 32 {
		buf := FixStr.Byte() | byte(uint(length))
		return e.writeByte(buf)
	}
	if length <= math.MaxUint8 {
		buf := e.work[:2]
		buf[0] = Str8.Byte()
		buf[1] = byte(uint(length))
		return e.write(buf)
	}
	if length <= math.MaxUint16 {
		buf := e.work[:3]
		buf[0] = Str16.Byte()
		binary.BigEndian.PutUint16(buf[1:], uint16(uint(length)))
		return e.write(buf)
	}
	if length <= math.MaxUint32 {
		buf := e.work[:5]
		buf[0] = Str32.Byte()
		binary.BigEndian.PutUint32(buf[1:], uint32(uint(length)))
		return e.write(buf)
	}
	return fmt.Errorf("too long")
}

func (e *Encoder) encodeStruct(rv reflect.Value) error {
	if rv.Type() == timeType {
		return e.encodeTimeFrom(rv)
	}
	if e.structKeyType == StructKeyTypeInt {
		return e.encodeStructToArray(rv)
	}
	return e.encodeStructToMap(rv)
}

func (e *Encoder) encodeStructToArray(rv reflect.Value) error {
	indexes, err := cache.GetInt(rv.Type(), e.structTagName)
	if err != nil {
		return err
	}

	err = e.encodeArrayHeaderInt(len(indexes))
	if err != nil {
		return err
	}

	for _, index := range indexes {
		if index < 0 {
			err := e.encodeNil()
			if err != nil {
				return err
			}
			continue
		}
		err := e.encodeValue(rv.Field(index))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeStructToMap(rv reflect.Value) error {
	indexes, err := cache.GetStringOrdered(rv.Type(), e.structTagName)
	if err != nil {
		return err
	}

	err = e.encodeMapHeaderInt(len(indexes))
	if err != nil {
		return err
	}

	for _, index := range indexes {
		if err := e.encodeString(index.Key); err != nil {
			return err
		}
		if err := e.encodeValue(rv.Field(index.Index)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeTimeFrom(rv reflect.Value) error {

	// zero allocation pass
	if rv.CanAddr() {
		p := unsafe.Pointer(rv.UnsafeAddr())
		t := *(*time.Time)(p)
		return e.encodeTime(t)
	}

	return e.encodeTime(rv.Interface().(time.Time))
}

func (e *Encoder) encodeTime(t time.Time) error {
	var (
		sec  = t.Unix()
		nsec = uint32(t.Nanosecond())
	)
	if sec < 0 {
		return e.encodeTime12(sec, nsec)
	}
	if sec <= math.MaxUint32 && nsec == 0 {
		return e.encodeTime4(uint32(sec))
	}
	if sec <= 0x3FFFFFFFF {
		return e.encodeTime8(sec, nsec)
	}
	return e.encodeTime12(sec, nsec)
}

func (e *Encoder) encodeTime4(sec uint32) error {
	bin := e.work[:6]
	bin[0] = FixExt4.Byte()
	bin[1] = TimestampTypeCode
	binary.BigEndian.PutUint32(bin[2:], sec)
	return e.write(bin)
}

func (e *Encoder) encodeTime8(sec int64, nsec uint32) error {
	if nsec > 999999999 {
		return fmt.Errorf("nsec out of range: %d", nsec)
	}
	enc := (uint64(nsec) << 34) | uint64(sec)

	bin := e.work[:10]
	bin[0] = FixExt8.Byte()
	bin[1] = TimestampTypeCode
	binary.BigEndian.PutUint64(bin[2:], enc)
	return e.write(bin)
}

func (e *Encoder) encodeTime12(sec int64, nsec uint32) error {
	if nsec > 999999999 {
		return fmt.Errorf("nsec out of range: %d", nsec)
	}
	bin := e.work[:15]
	bin[0] = Ext8.Byte()
	bin[1] = 12 // Payload Length
	bin[2] = TimestampTypeCode
	binary.BigEndian.PutUint32(bin[3:], nsec)
	binary.BigEndian.PutUint64(bin[7:], uint64(sec))
	return e.write(bin)
}

func (e *Encoder) writeByte(b byte) error {
	if e.bw != nil {
		return e.bw.WriteByte(b)
	}
	buf := e.work[:1]
	buf[0] = b
	return e.write(buf)
}

func (e *Encoder) write(b []byte) error {
	_, err := e.w.Write(b)
	return err
}
