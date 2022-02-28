package reflectutil

import (
	"reflect"
)

func IsFloat(kind reflect.Kind) bool {
	return kind == reflect.Float32 || kind == reflect.Float64
}
