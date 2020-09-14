package msgpack

// ArrayLengthTolerance is tolerance of array length when decode time.
//go:generate stringer -type=ArrayLengthTolerance -output=array_length_tolerance_string.go
type ArrayLengthTolerance int

const (
	// ArrayLengthToleranceLessThanOrEqual is allowed if the array length of message pack is the same or less than or equal to the length of go.
	ArrayLengthToleranceLessThanOrEqual ArrayLengthTolerance = iota

	// ArrayLengthToleranceEqualOnly is allowed if the array length of message pack is the only equal to the length of go.
	ArrayLengthToleranceEqualOnly

	// ArrayLengthToleranceRounding is allowed always.
	// If the array length of  is the same or greater than to the length of go, overloaded data is discarded.
	ArrayLengthToleranceRounding
)
