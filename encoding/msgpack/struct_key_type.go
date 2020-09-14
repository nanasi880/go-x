package msgpack

// StructKeyType is determines how to serialize/deserialize the structure when it is encode/decode.
//go:generate stringer -type=StructKeyType -output=struct_key_type_string.go
type StructKeyType byte

const (
	StructKeyTypeInt    StructKeyType = iota // StructKeyTypeInt is struct serialize/deserialize as int key array format.
	StructKeyTypeString                      // StructKeyTypeString is struct serialize/deserialize as string key map format.
)
