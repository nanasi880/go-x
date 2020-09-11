package reflect

import "reflect"

// ZeroValue returns zero value of type t.
func ZeroValue(t reflect.Type) reflect.Value {
	return reflect.New(t).Elem()
}
