package reflectutil

import "reflect"

// IsInt returns whether the value like int.
func IsInt(kind reflect.Kind) bool {
	return IsSignedInt(kind) || IsUnsignedInt(kind)
}

func IsSignedInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func IsUnsignedInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}
