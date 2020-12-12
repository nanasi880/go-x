package env

import "strconv"

// Value is the interface to the dynamic value stored in a value.
type Value interface {
	Set(v string) error
}

// stringValue is Value type of string.
type stringValue struct {
	ptr *string
}

func newStringValue(v *string) *stringValue {
	return &stringValue{
		ptr: v,
	}
}

func (val *stringValue) Set(v string) error {
	*val.ptr = v
	return nil
}

// intValue is Value type of int.
type intValue struct {
	ptr *int
}

func newIntValue(v *int) *intValue {
	return &intValue{
		ptr: v,
	}
}

func (val *intValue) Set(v string) error {

	iv, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return err
	}

	*val.ptr = int(iv)
	return nil
}

// boolValue is Value type of bool.
type boolValue struct {
	ptr *bool
}

func newBoolValue(v *bool) *boolValue {
	return &boolValue{
		ptr: v,
	}
}

func (val *boolValue) Set(v string) error {
	b, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}
	*val.ptr = b
	return nil
}
