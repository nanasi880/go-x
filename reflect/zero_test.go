package reflect_test

import (
	stdReflect "reflect"
	"testing"

	"go.nanasi880.dev/x/reflect"
)

func TestZeroValue(t *testing.T) {

	typ := stdReflect.TypeOf(0)
	v := reflect.ZeroValue(typ)

	if v.Kind() != stdReflect.Int {
		t.Fatal(v)
	}
	if v.Interface() != 0 {
		t.Fatal(v)
	}
}

func TestZeroValue_Ptr(t *testing.T) {
	v := 42
	p := &v

	rv := stdReflect.ValueOf(&p)
	rv.Elem().Set(reflect.ZeroValue(rv.Elem().Type()))

	if p != nil {
		t.Fatal()
	}
}
