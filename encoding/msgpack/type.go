package msgpack

import (
	"reflect"
	"time"
)

var (
	marshalerType      = reflect.TypeOf((*Marshaler)(nil)).Elem()
	unmarshalerType    = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
	byteSliceType      = reflect.TypeOf(([]byte)(nil))
	timeType           = reflect.TypeOf(time.Time{})
	byteType           = reflect.TypeOf(byte(0))
	interfaceType      = reflect.TypeOf((*interface{})(nil)).Elem()
	interfaceSliceType = reflect.TypeOf(([]interface{})(nil))
	interfaceMapType   = reflect.TypeOf((map[interface{}]interface{})(nil))
)
