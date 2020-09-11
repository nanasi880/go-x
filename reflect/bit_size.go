package reflect

import (
	"fmt"
	"reflect"
)

// IntegerBitSize returns bit size of integer kind. If kind is Int or Uint, returns 0.
func IntegerBitSize(kind reflect.Kind) int {
	switch kind {
	case reflect.Int8:
		return 8
	case reflect.Int16:
		return 16
	case reflect.Int32:
		return 32
	case reflect.Int64:
		return 64
	case reflect.Int:
		return 0
	case reflect.Uint8:
		return 8
	case reflect.Uint16:
		return 16
	case reflect.Uint32:
		return 32
	case reflect.Uint64:
		return 64
	case reflect.Uint:
		return 0
	default:
		panic(fmt.Sprintf("not an integer kind: %v", kind))
	}
}
