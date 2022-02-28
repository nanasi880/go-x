package reflectutil

import "reflect"

// AllocateTo is allocate zero value and assign to value.
func AllocateTo(value reflect.Value) {
	value.Set(reflect.New(value.Type().Elem()))
}
