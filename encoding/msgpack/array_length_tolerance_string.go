// Code generated by "stringer -type=ArrayLengthTolerance -output=array_length_tolerance_string.go"; DO NOT EDIT.

package msgpack

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ArrayLengthToleranceLessThanOrEqual-0]
	_ = x[ArrayLengthToleranceEqualOnly-1]
	_ = x[ArrayLengthToleranceRounding-2]
}

const _ArrayLengthTolerance_name = "ArrayLengthToleranceLessThanOrEqualArrayLengthToleranceEqualOnlyArrayLengthToleranceRounding"

var _ArrayLengthTolerance_index = [...]uint8{0, 35, 64, 92}

func (i ArrayLengthTolerance) String() string {
	if i < 0 || i >= ArrayLengthTolerance(len(_ArrayLengthTolerance_index)-1) {
		return "ArrayLengthTolerance(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ArrayLengthTolerance_name[_ArrayLengthTolerance_index[i]:_ArrayLengthTolerance_index[i+1]]
}