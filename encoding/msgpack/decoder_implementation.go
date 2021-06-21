package msgpack

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"

	cache "go.nanasi880.dev/x/encoding/msgpack/internal/index"
	xreflect "go.nanasi880.dev/x/reflect"
	"go.nanasi880.dev/x/runtime"
	xunsafe "go.nanasi880.dev/x/unsafe"
)

const (
	readBlockSize = 64
)

func (d *Decoder) decodeValue(rv reflect.Value) error {
	if d.validateUnmarshaler(rv) {
		return d.decodeUnmarshaler(rv)
	}

	format, rawFormat, err := d.readFormat()
	if err != nil {
		return err
	}

	switch format {
	case PositiveFixInt, NegativeFixInt, Uint8, Uint16, Uint32, Uint64, Int8, Int16, Int32, Int64:
		return d.decodeInt(rv, format, rawFormat)
	case Float32, Float64:
		return d.decodeFloat(rv, format)
	case FixStr, Str8, Str16, Str32:
		return d.decodeString(rv, format, rawFormat)
	case Nil:
		return d.decodeNil(rv)
	case Unused:
		return nil
	case False, True:
		return d.decodeBool(rv, format)
	case Bin8, Bin16, Bin32:
		return d.decodeBin(rv, format)
	case Ext8, Ext16, Ext32, FixExt1, FixExt2, FixExt4, FixExt8, FixExt16:
		return d.decodeExt(rv, format)
	case FixArray, Array16, Array32:
		return d.decodeArray(rv, format, rawFormat)
	case FixMap, Map16, Map32:
		return d.decodeMap(rv, format, rawFormat)
	default:
		return fmt.Errorf("unsupported format: %s", format.String())
	}
}

func (d *Decoder) decodeUnmarshaler(rv reflect.Value) error {
	for {
		if xreflect.IsNilable(rv) && rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		if rv.Type().Implements(unmarshalerType) {
			return rv.Interface().(Unmarshaler).UnmarshalMsgPack(d)
		}
		rv = rv.Elem()
	}
}

func (d *Decoder) decodeInt(rv reflect.Value, format FormatName, rawFormat byte) error {
	if !d.validateInt(rv) {
		return fmt.Errorf("%s cannot assign to %s", format.String(), rv.Type().String())
	}

	value, err := d.decodeIntBit(format, rawFormat)
	if err != nil {
		return err
	}

	switch format {
	case NegativeFixInt, Int8, Int16, Int32, Int64:
		return d.setInt(rv, int64(value))
	default:
		return d.setUint(rv, value)
	}
}

func (d *Decoder) decodeIntBit(format FormatName, rawFormat byte) (uint64, error) {
	var value uint64
	switch format {
	case PositiveFixInt:
		value = uint64(rawFormat)
	case NegativeFixInt:
		value = uint64(int8(rawFormat))
	case Int8, Uint8:
		val, err := d.read(1)
		if err != nil {
			return 0, err
		}
		if format == Int8 {
			value = uint64(int8(val[0]))
		} else {
			value = uint64(val[0])
		}
	case Int16, Uint16:
		val, err := d.read(2)
		if err != nil {
			return 0, err
		}
		if format == Int16 {
			value = uint64(int16(binary.BigEndian.Uint16(val)))
		} else {
			value = uint64(binary.BigEndian.Uint16(val))
		}
	case Int32, Uint32:
		val, err := d.read(4)
		if err != nil {
			return 0, err
		}
		if format == Int32 {
			value = uint64(int32(binary.BigEndian.Uint32(val)))
		} else {
			value = uint64(binary.BigEndian.Uint32(val))
		}
	case Int64, Uint64:
		val, err := d.read(8)
		if err != nil {
			return 0, err
		}
		value = binary.BigEndian.Uint64(val)
	default:
		return 0, fmt.Errorf("invalid format: %s", format.String())
	}
	return value, nil
}

func (d *Decoder) decodeFloat(rv reflect.Value, format FormatName) error {
	if !d.validateFloat(rv) {
		return fmt.Errorf("%s cannot assign to %s", format.String(), rv.Type().String())
	}
	switch format {
	case Float32:
		val, err := d.read(4)
		if err != nil {
			return err
		}
		return d.setFloat32(rv, math.Float32frombits(binary.BigEndian.Uint32(val)))
	case Float64:
		val, err := d.read(8)
		if err != nil {
			return err
		}
		return d.setFloat64(rv, math.Float64frombits(binary.BigEndian.Uint64(val)))
	default:
		return fmt.Errorf("%s cannot assign to %s", format.String(), rv.Type().String())
	}
}

func (d *Decoder) decodeString(rv reflect.Value, format FormatName, rawFormat byte) error {
	if !d.validateString(rv) {
		return fmt.Errorf("%s cannot assign to %s", format.String(), rv.Type().String())
	}

	length, err := d.decodeStringHeader(format, rawFormat)
	if err != nil {
		return err
	}

	tmp, err := d.decodeBytes(length)
	if err != nil {
		return err
	}
	return d.setString(rv, tmp)
}

func (d *Decoder) decodeStringHeader(format FormatName, b byte) (int, error) {
	var (
		length uint32
	)
	switch format {
	case FixStr:
		length = uint32(b & (^FixStr.Byte()))
	case Str8:
		val, err := d.read(1)
		if err != nil {
			return 0, err
		}
		length = uint32(val[0])
	case Str16:
		val, err := d.read(2)
		if err != nil {
			return 0, err
		}
		length = uint32(binary.BigEndian.Uint16(val))
	case Str32:
		val, err := d.read(4)
		if err != nil {
			return 0, err
		}
		length = binary.BigEndian.Uint32(val)
	default:
		return 0, fmt.Errorf("%s is not a string type", format.String())
	}

	if runtime.MaxInt < uint64(length) {
		return 0, fmt.Errorf("the string is too large (string length exceeds the maximum value of int)")
	}

	return int(length), nil
}

func (d *Decoder) decodeNil(rv reflect.Value) error {
	// ポインタ型にもNilを入れるために変則的な見方をする
	// var v *int
	// Decode(&v)
	// のような呼び出しを行った場合、最初のポインタ値はCanSet() == falseだが、2階層目はCanSet() == trueとなるので
	// それを利用してセット対象がポインタ型なのかどうかを見分ける
	for rv.Kind() == reflect.Ptr && !rv.CanSet() {
		rv = rv.Elem()
	}

	// already nil
	if xreflect.IsNilable(rv) && rv.IsNil() {
		return nil
	}

	// interface{} special
	if rv.Type() == interfaceType {
		rv.Set(zeroInterfaceValue)
		return nil
	}

	// zero value fallback
	rv.Set(reflect.Zero(rv.Type()))
	return nil
}

func (d *Decoder) decodeBool(rv reflect.Value, format FormatName) error {
	if !d.validateBool(rv) {
		return fmt.Errorf("%s cannot assign to %s", format.String(), rv.Type().String())
	}
	return d.setBool(rv, format == True)
}

func (d *Decoder) decodeBin(rv reflect.Value, format FormatName) error {
	if !d.validateBin(rv) {
		return fmt.Errorf("%s cannot assign to %s", format.String(), rv.Type().String())
	}

	length, err := d.decodeBinHeader(format)
	if err != nil {
		return err
	}

	tmp, err := d.decodeBytes(length)
	if err != nil {
		return err
	}
	return d.setBin(rv, tmp)
}

func (d *Decoder) decodeBinHeader(format FormatName) (int, error) {
	var (
		length uint32
	)
	switch format {
	case Bin8:
		val, err := d.read(1)
		if err != nil {
			return 0, err
		}
		length = uint32(val[0])
	case Bin16:
		val, err := d.read(2)
		if err != nil {
			return 0, err
		}
		length = uint32(binary.BigEndian.Uint16(val))
	case Bin32:
		val, err := d.read(4)
		if err != nil {
			return 0, err
		}
		length = binary.BigEndian.Uint32(val)
	default:
		return 0, fmt.Errorf("%s is not a bin type", format.String())
	}

	if runtime.MaxInt < uint64(length) {
		return 0, fmt.Errorf("the bin is too large (bin length exceeds the maximum value of int)")
	}

	return int(length), nil
}

func (d *Decoder) decodeExt(rv reflect.Value, format FormatName) error {
	typeCode, length, err := d.decodeExtHeader(format)
	if err != nil {
		return err
	}
	switch typeCode {
	case TimestampTypeCode:
		return d.decodeTimeTo(rv, format, length)
	default:
		return fmt.Errorf("unsupported ext type code: %d", typeCode)
	}
}

func (d *Decoder) decodeExtHeader(format FormatName) (byte, uint32, error) {
	switch format {
	case FixExt1, FixExt2, FixExt4, FixExt8, FixExt16:
		val, err := d.read(1)
		if err != nil {
			return 0, 0, err
		}
		typeCode := val[0]
		switch format {
		case FixExt1:
			return typeCode, 1, nil
		case FixExt2:
			return typeCode, 2, nil
		case FixExt4:
			return typeCode, 4, nil
		case FixExt8:
			return typeCode, 8, nil
		default:
			return typeCode, 16, nil
		}
	case Ext8:
		val, err := d.read(2)
		if err != nil {
			return 0, 0, err
		}
		return val[1], uint32(val[0]), nil
	case Ext16:
		length, err := d.read(2)
		if err != nil {
			return 0, 0, err
		}
		typeCode, err := d.read(1)
		if err != nil {
			return 0, 0, err
		}
		return typeCode[0], uint32(binary.BigEndian.Uint16(length)), nil
	case Ext32:
		length, err := d.read(4)
		if err != nil {
			return 0, 0, err
		}
		typeCode, err := d.read(1)
		if err != nil {
			return 0, 0, err
		}
		return typeCode[0], binary.BigEndian.Uint32(length), nil
	default:
		return 0, 0, fmt.Errorf("%s is not a ext type", format.String())
	}
}

func (d *Decoder) decodeExtHeaderAsInt(format FormatName) (byte, int, error) {
	typeCode, length, err := d.decodeExtHeader(format)
	if err != nil {
		return 0, 0, err
	}
	if runtime.MaxInt < uint64(length) {
		return 0, 0, fmt.Errorf("the ext data is too long (ext data length exceeds the maximum value of int)")
	}
	return typeCode, int(length), nil
}

func (d *Decoder) decodeTimeTo(rv reflect.Value, format FormatName, length uint32) error {
	if !d.validateTime(rv) {
		return fmt.Errorf("%s cannot assign to %s", format.String(), rv.Type().String())
	}

	t, err := d.decodeTime(length)
	if err != nil {
		return err
	}

	return d.setTime(rv, t)
}

func (d *Decoder) decodeTime(length uint32) (time.Time, error) {
	var t time.Time
	switch length {
	case 4:
		sec, err := d.read(4)
		if err != nil {
			return t, err
		}
		t = time.Unix(int64(binary.BigEndian.Uint32(sec)), 0)
	case 8:
		payload, err := d.read(8)
		if err != nil {
			return t, err
		}
		var (
			mix  = binary.BigEndian.Uint64(payload)
			nsec = (mix >> 34) & 0x3FFFFFFF
			sec  = mix & 0x3FFFFFFFF
		)
		t = time.Unix(int64(sec), int64(nsec))
	case 12:
		tmp, err := d.read(12)
		if err != nil {
			return t, err
		}
		var (
			nsec = tmp[:4]
			sec  = tmp[4:]
		)
		t = time.Unix(
			int64(binary.BigEndian.Uint64(sec)),
			int64(binary.BigEndian.Uint32(nsec)),
		)
	default:
		return t, fmt.Errorf("timestamp data is not a well known encode")
	}

	if d.timeZone == nil {
		return t.UTC(), nil
	} else {
		return t.In(d.timeZone), nil
	}
}

func (d *Decoder) decodeArray(rv reflect.Value, format FormatName, b byte) error {
	if d.structKeyType == StructKeyTypeInt && d.validateStruct(rv) {
		return d.decodeArrayToStruct(rv, format, b)
	}
	return d.decodeArrayToArray(rv, format, b)
}

func (d *Decoder) decodeArrayToStruct(rv reflect.Value, format FormatName, b byte) error {
	length, err := d.decodeArrayHeaderAsInt(format, b)
	if err != nil {
		return fmt.Errorf("%s cannot decode to %s: %w", format, rv.Type().String(), err)
	}

	// allocate
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}

	structIndexes := cache.GetInt(rv.Type())

	// Arrayの場合はNilでパディングしてあるはず かつ
	// ずれている場合に補正が効かないのでエラーとする
	if len(structIndexes) != length {
		return fmt.Errorf("structure indexes not match")
	}

	for i := 0; i < length; i++ {
		index := structIndexes[i]
		if index < 0 {
			nilByte, err := d.read(1)
			if err != nil {
				return err
			}
			if decodeFormatName(nilByte[0]) != Nil {
				return fmt.Errorf("struct padding must be Nil")
			}
		}
		field := rv.Field(index)
		if err := d.decodeValue(field); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) decodeArrayToArray(rv reflect.Value, format FormatName, b byte) error {
	if !d.validateArray(rv) {
		return fmt.Errorf("%s cannot assign to %s", format.String(), rv.Type().String())
	}

	length, err := d.decodeArrayHeaderAsInt(format, b)
	if err != nil {
		return err
	}

	// allocate
	for rv.Kind() == reflect.Ptr {
		if xreflect.IsNilable(rv) && rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}

	// interface
	if rv.Type() == interfaceType {
		rv.Set(reflect.MakeSlice(interfaceSliceType, length, length))
	}

	// grow capacity
	if rv.Kind() == reflect.Slice {
		if rv.Len() < length {
			if rv.Cap() >= length {
				// fast pass
				rv.SetLen(length)
			} else {
				// slow pass
				rv.Set(reflect.MakeSlice(rv.Type(), length, length))
			}
		}
	}
	if xreflect.IsNilable(rv) && rv.IsNil() {
		rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
	}

	switch d.arrayLengthTolerance {
	case ArrayLengthToleranceLessThanOrEqual:
		if rv.Len() > length {
			return fmt.Errorf("the message pack array is too long (exceeds the maximum value of go array)")
		}
	case ArrayLengthToleranceEqualOnly:
		if rv.Len() != length {
			return fmt.Errorf("the array lengths do not match")
		}
	case ArrayLengthToleranceRounding:
		break
	}

	for i := 0; i < length; i++ {
		if i < rv.Len() {
			err = d.decodeValue(rv.Index(i))
			if err != nil {
				return err
			}
		} else {
			if err := d.skipCurrentObject(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Decoder) decodeArrayHeader(format FormatName, b byte) (uint32, error) {
	var (
		length uint32
	)
	switch format {
	case FixArray:
		length = uint32(b & (^FixArray.Byte()))
	case Array16:
		val, err := d.read(2)
		if err != nil {
			return 0, err
		}
		length = uint32(binary.BigEndian.Uint16(val))
	case Array32:
		val, err := d.read(4)
		if err != nil {
			return 0, err
		}
		length = binary.BigEndian.Uint32(val)
	default:
		return 0, fmt.Errorf("%s is not a array type", format.String())
	}

	return length, nil
}

func (d *Decoder) decodeArrayHeaderAsInt(format FormatName, b byte) (int, error) {
	length, err := d.decodeArrayHeader(format, b)
	if err != nil {
		return 0, err
	}
	if runtime.MaxInt < uint64(length) {
		return 0, fmt.Errorf("the array is too long (array length exceeds the maximum value of int)")
	}
	return int(length), nil
}

func (d *Decoder) decodeMap(rv reflect.Value, format FormatName, b byte) error {
	if d.structKeyType == StructKeyTypeString && d.validateStruct(rv) {
		return d.decodeMapToStruct(rv, format, b)
	}
	return d.decodeMapToMap(rv, format, b)
}

func (d *Decoder) decodeMapToStruct(rv reflect.Value, format FormatName, b byte) error {
	length, err := d.decodeMapHeaderAsInt(format, b)
	if err != nil {
		return fmt.Errorf("%s cannot decode to %s: %w", format, rv.Type().String(), err)
	}

	// allocate
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}

	structIndexes := cache.GetString(rv.Type())

	for i := 0; i < length; i++ {
		index, ok, err := d.lookupStringIndex(structIndexes)
		if err != nil {
			return fmt.Errorf("struct key decode failed: %w", err)
		}

		if !ok {
			if err := d.skipCurrentObject(); err != nil {
				return err
			}
			continue
		}

		field := rv.Field(index)
		if err := d.decodeValue(field); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) decodeMapToMap(rv reflect.Value, format FormatName, b byte) error {
	if !d.validateMap(rv) {
		return fmt.Errorf("%s cannot assign to %s", format.String(), rv.Type().String())
	}

	length, err := d.decodeMapHeaderAsInt(format, b)
	if err != nil {
		return fmt.Errorf("%s cannot decode to %s: %w", format, rv.Type().String(), err)
	}

	// allocate
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}
	if rv.Type() == interfaceType {
		rv.Set(reflect.MakeMap(interfaceMapType))
		rv = rv.Elem()
	}

	var (
		keyType   = rv.Type().Key()
		valueType = rv.Type().Elem()
	)
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(reflect.MapOf(keyType, valueType)))
	}
	for i := 0; i < length; i++ {
		var (
			key   = reflect.New(keyType)
			value = reflect.New(valueType)
		)
		if err := d.decodeValue(key); err != nil {
			return err
		}
		if err := d.decodeValue(value); err != nil {
			return err
		}

		rv.SetMapIndex(key.Elem(), value.Elem())
	}

	return nil
}

func (d *Decoder) decodeMapHeader(format FormatName, b byte) (uint32, error) {
	var (
		length uint32
	)
	switch format {
	case FixMap:
		length = uint32(b & (^FixMap.Byte()))
	case Map16:
		val, err := d.read(2)
		if err != nil {
			return 0, err
		}
		length = uint32(binary.BigEndian.Uint16(val))
	case Map32:
		val, err := d.read(4)
		if err != nil {
			return 0, err
		}
		length = binary.BigEndian.Uint32(val)
	default:
		return 0, fmt.Errorf("%s is not a map type", format.String())
	}

	return length, nil
}

func (d *Decoder) decodeMapHeaderAsInt(format FormatName, b byte) (int, error) {
	length, err := d.decodeMapHeader(format, b)
	if err != nil {
		return 0, err
	}

	if runtime.MaxInt < uint64(length) {
		return 0, fmt.Errorf("the map is too large (map length exceeds the maximum value of int)")
	}

	return int(length), nil
}

func (d *Decoder) decodeBytes(length int) ([]byte, error) {
	buf := make([]byte, length)
	if err := d.readTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (d *Decoder) validateUnmarshaler(rv reflect.Value) bool {
	typ := rv.Type()
	for {
		if typ.Implements(unmarshalerType) {
			return true
		}
		if typ.Kind() != reflect.Ptr {
			return false
		}
		typ = typ.Elem()
	}
}

func (d *Decoder) validateInt(rv reflect.Value) bool {
	typ := rv.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if xreflect.IsInt(typ.Kind()) {
		return true
	}
	if xreflect.IsFloat(typ.Kind()) {
		return true
	}
	if typ.Kind() == reflect.Interface {
		return typ == interfaceType
	}
	return false
}

func (d *Decoder) validateFloat(rv reflect.Value) bool {
	typ := rv.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if xreflect.IsFloat(typ.Kind()) {
		return true
	}
	if xreflect.IsInt(typ.Kind()) {
		return true
	}
	if typ.Kind() == reflect.Interface {
		return typ == interfaceType
	}
	return false
}

func (d *Decoder) validateString(rv reflect.Value) bool {
	return d.validateBin(rv)
}

func (d *Decoder) validateBin(rv reflect.Value) bool {
	typ := rv.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.String {
		return true
	}
	if typ == byteSliceType {
		return true
	}
	if typ.Kind() == reflect.Array {
		return typ.Elem() == byteType
	}
	return typ == interfaceType
}

func (d *Decoder) validateBool(rv reflect.Value) bool {
	typ := rv.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Bool {
		return true
	}
	if typ == interfaceType {
		return true
	}
	return false
}

func (d *Decoder) validateTime(rv reflect.Value) bool {
	typ := rv.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ == timeType {
		return true
	}
	if typ == interfaceType {
		return true
	}
	return false
}

func (d *Decoder) validateStruct(rv reflect.Value) bool {
	typ := rv.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Struct {
		return true
	}
	return false
}

func (d *Decoder) validateArray(rv reflect.Value) bool {
	typ := rv.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Slice {
		return true
	}
	if typ.Kind() == reflect.Array {
		return true
	}
	return false
}

func (d *Decoder) validateMap(rv reflect.Value) bool {
	typ := rv.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Map {
		return true
	}
	return typ == interfaceType
}

func (d *Decoder) setInt(rv reflect.Value, v int64) error {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}
	if xreflect.IsSignedInt(rv.Kind()) {
		rv.SetInt(v)
		return nil
	}
	if xreflect.IsUnsignedInt(rv.Kind()) {
		rv.SetUint(uint64(v))
		return nil
	}
	if xreflect.IsFloat(rv.Kind()) {
		rv.SetFloat(float64(v))
		return nil
	}
	if rv.Type() == interfaceType {
		rv.Set(reflect.ValueOf(v))
		return nil
	}
	return fmt.Errorf("internal: unsupported kind: %s", rv.Kind())
}

func (d *Decoder) setUint(rv reflect.Value, v uint64) error {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}
	if xreflect.IsUnsignedInt(rv.Kind()) {
		rv.SetUint(v)
		return nil
	}
	if xreflect.IsSignedInt(rv.Kind()) {
		rv.SetInt(int64(v))
		return nil
	}
	if xreflect.IsFloat(rv.Kind()) {
		rv.SetFloat(float64(v))
		return nil
	}
	if rv.Type() == interfaceType {
		rv.Set(reflect.ValueOf(v))
		return nil
	}
	return fmt.Errorf("internal: unsupported kind: %s", rv.Kind())
}

func (d *Decoder) setFloat32(rv reflect.Value, v float32) error {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}
	if xreflect.IsFloat(rv.Kind()) {
		rv.SetFloat(float64(v))
		return nil
	}
	if xreflect.IsInt(rv.Kind()) {
		if xreflect.IsSignedInt(rv.Kind()) {
			rv.SetInt(int64(v))
		} else {
			rv.SetUint(uint64(v))
		}
		return nil
	}
	if rv.Type() == interfaceType {
		rv.Set(reflect.ValueOf(v))
		return nil
	}
	return fmt.Errorf("internal: unsupported kind: %s", rv.Kind())
}

func (d *Decoder) setFloat64(rv reflect.Value, v float64) error {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}
	if xreflect.IsFloat(rv.Kind()) {
		rv.SetFloat(v)
		return nil
	}
	if xreflect.IsInt(rv.Kind()) {
		if xreflect.IsSignedInt(rv.Kind()) {
			rv.SetInt(int64(v))
		} else {
			rv.SetUint(uint64(v))
		}
		return nil
	}
	if rv.Type() == interfaceType {
		rv.Set(reflect.ValueOf(v))
		return nil
	}
	return fmt.Errorf("internal: unsupported kind: %s", rv.Kind())
}

func (d *Decoder) setString(rv reflect.Value, v []byte) error {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.String {
		rv.SetString(xunsafe.BytesToString(v))
		return nil
	}
	if rv.Type() == byteSliceType {
		rv.SetBytes(v)
		return nil
	}
	if rv.Kind() == reflect.Array {
		return d.copyBinToArray(rv, v)
	}
	if rv.Type() == interfaceType {
		rv.Set(reflect.ValueOf(xunsafe.BytesToString(v)))
		return nil
	}
	return fmt.Errorf("internal: unsupported kind: %s", rv.Kind())
}

func (d *Decoder) setBool(rv reflect.Value, v bool) error {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Bool {
		rv.SetBool(v)
		return nil
	}
	if rv.Type() == interfaceType {
		rv.Set(reflect.ValueOf(v))
		return nil
	}
	return fmt.Errorf("internal: unsupported kind: %s", rv.Kind())
}

func (d *Decoder) setBin(rv reflect.Value, v []byte) error {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}
	if rv.Type() == byteSliceType {
		rv.SetBytes(v)
		return nil
	}
	if rv.Kind() == reflect.String {
		rv.SetString(xunsafe.BytesToString(v))
		return nil
	}
	if rv.Kind() == reflect.Array {
		return d.copyBinToArray(rv, v)
	}
	if rv.Type() == interfaceType {
		rv.Set(reflect.ValueOf(v))
		return nil
	}
	return fmt.Errorf("internal: unsupported kind: %s", rv.Kind())
}

func (d *Decoder) setTime(rv reflect.Value, v time.Time) error {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			xreflect.AllocateTo(rv)
		}
		rv = rv.Elem()
	}
	if rv.Type() == timeType {
		rv.Set(reflect.ValueOf(v))
		return nil
	}
	if rv.Type() == interfaceType {
		rv.Set(reflect.ValueOf(v))
		return nil
	}
	return fmt.Errorf("internal: unsupported kind: %s", rv.Kind())
}

func (d *Decoder) lookupStringIndex(indexes map[string]int) (int, bool, error) {
	format, b, err := d.readFormat()
	if err != nil {
		return 0, false, err
	}

	var length int
	switch format {
	case Bin8, Bin16, Bin32:
		length, err = d.decodeBinHeader(format)
	default:
		length, err = d.decodeStringHeader(format, b)
	}
	if err != nil {
		return 0, false, err
	}

	if length < readBlockSize {
		return d.lookupShortStringIndex(indexes, length)
	}
	return d.lookupLongStringIndex(indexes, length)
}

func (d *Decoder) lookupShortStringIndex(indexes map[string]int, length int) (int, bool, error) {
	buf, err := d.read(length)
	if err != nil {
		return 0, false, err
	}
	index, ok := indexes[string(buf)]
	return index, ok, nil
}

func (d *Decoder) lookupLongStringIndex(indexes map[string]int, length int) (int, bool, error) {
	key, err := d.decodeBytes(length)
	if err != nil {
		return 0, false, err
	}
	index, ok := indexes[string(key)]
	return index, ok, nil
}

func (d *Decoder) skipCurrentObject() error {
	b, err := d.read(1)
	if err != nil {
		return err
	}
	return d.skipObject(decodeFormatName(b[0]), b[0])
}

func (d *Decoder) skipObject(format FormatName, b byte) error {
	switch format {
	case PositiveFixInt, NegativeFixInt, Uint8, Uint16, Uint32, Uint64, Int8, Int16, Int32, Int64:
		return d.skipInt(format)
	case Float32, Float64:
		return d.skipFloat(format)
	case FixStr, Str8, Str16, Str32:
		return d.skipString(format, b)
	case Nil:
		return nil
	case Unused:
		return nil
	case False, True:
		return nil
	case Bin8, Bin16, Bin32:
		return d.skipString(format, b)
	case Ext8, Ext16, Ext32, FixExt1, FixExt2, FixExt4, FixExt8, FixExt16:
		return d.skipExt(format)
	case FixArray, Array16, Array32:
		return d.skipArray(format, b)
	case FixMap, Map16, Map32:
		return d.skipMap(format, b)
	default:
		return fmt.Errorf("unsupported format: %s", format.String())
	}
}

func (d *Decoder) skipInt(format FormatName) error {
	var seekSize int64
	switch format {
	case PositiveFixInt, NegativeFixInt:
		seekSize = 0
	case Int8, Uint8:
		seekSize = 1
	case Int16, Uint16:
		seekSize = 2
	case Int32, Uint32:
		seekSize = 4
	case Int64, Uint64:
		seekSize = 8
	default:
		panic("internal")
	}
	if seekSize == 0 {
		return nil
	}
	return d.seek(seekSize)
}

func (d *Decoder) skipFloat(format FormatName) error {
	var size int64 = 4
	if format == Float64 {
		size = 8
	}
	return d.seek(size)
}

func (d *Decoder) skipString(format FormatName, b byte) error {
	length, err := d.decodeStringHeader(format, b)
	if err != nil {
		return err
	}
	return d.seek(int64(length))
}

func (d *Decoder) skipExt(format FormatName) error {
	_, length, err := d.decodeExtHeader(format)
	if err != nil {
		return err
	}
	return d.seek(int64(length))
}

func (d *Decoder) skipArray(format FormatName, b byte) error {
	length, err := d.decodeArrayHeader(format, b)
	if err != nil {
		return err
	}
	for i := uint32(0); i < length; i++ {
		err = d.skipCurrentObject()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) skipMap(format FormatName, b byte) error {
	length, err := d.decodeMapHeader(format, b)
	if err != nil {
		return err
	}
	for i := uint32(0); i < length; i++ {
		// key
		if err := d.skipCurrentObject(); err != nil {
			return err
		}
		// value
		if err := d.skipCurrentObject(); err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) copyBinToArray(rv reflect.Value, v []byte) error {
	switch d.arrayLengthTolerance {
	case ArrayLengthToleranceLessThanOrEqual:
		if !(len(v) <= rv.Len()) {
			return fmt.Errorf("the message pack array is too long (exceeds the maximum value of go array)")
		}
	case ArrayLengthToleranceEqualOnly:
		if len(v) != rv.Len() {
			return fmt.Errorf("the array lengths do not match")
		}
	case ArrayLengthToleranceRounding:
		break
	}

	reflect.Copy(rv, reflect.ValueOf(v))
	return nil
}

func (d *Decoder) readFormat() (FormatName, byte, error) {
	b, err := d.read(1)
	if err != nil {
		return 0, 0, err
	}
	return decodeFormatName(b[0]), b[0], nil
}

func (d *Decoder) read(length int) ([]byte, error) {
	if length > readBlockSize {
		panic("internal")
	}

	if d.data != nil {
		if len(d.data) < length {
			return nil, io.ErrUnexpectedEOF
		}
		v := d.data[:length]
		d.data = d.data[length:]
		return v, nil
	}

	_, err := d.reader.Read(d.work[:length])
	return d.work[:length], err
}

func (d *Decoder) readTo(p []byte) error {
	if d.data != nil {
		if len(d.data) < len(p) {
			return io.ErrUnexpectedEOF
		}
		copy(p, d.data)
		d.data = d.data[len(p):]
		return nil
	}

	_, err := d.reader.Read(p)
	return err
}

func (d *Decoder) seek(l int64) error {
	if d.data != nil {
		if int64(len(d.data)) < l {
			return io.ErrUnexpectedEOF
		}
		d.data = d.data[l:]
		return nil
	}

	_, err := d.reader.Seek(l, io.SeekCurrent)
	return err
}
