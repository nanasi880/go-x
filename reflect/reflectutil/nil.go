package reflectutil

import "reflect"

// IsNilable returns whether the value can be nil
func IsNilable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}

// IsNil returns IsNilable(v) && v.IsNil
func IsNil(v reflect.Value) bool {
	if IsNilable(v) && v.IsNil() {
		return true
	}
	return false
}
