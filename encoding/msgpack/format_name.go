package msgpack

//go:generate stringer -type=FormatName -output=format_name_string.go
type FormatName byte

const (
	PositiveFixInt FormatName = 0
	FixMap         FormatName = 0x80
	FixArray       FormatName = 0x90
	FixStr         FormatName = 0xa0
	Nil            FormatName = 0xc0
	Unused         FormatName = 0xc1
	False          FormatName = 0xc2
	True           FormatName = 0xc3
	Bin8           FormatName = 0xc4
	Bin16          FormatName = 0xc5
	Bin32          FormatName = 0xc6
	Ext8           FormatName = 0xc7
	Ext16          FormatName = 0xc8
	Ext32          FormatName = 0xc9
	Float32        FormatName = 0xca
	Float64        FormatName = 0xcb
	Uint8          FormatName = 0xcc
	Uint16         FormatName = 0xcd
	Uint32         FormatName = 0xce
	Uint64         FormatName = 0xcf
	Int8           FormatName = 0xd0
	Int16          FormatName = 0xd1
	Int32          FormatName = 0xd2
	Int64          FormatName = 0xd3
	FixExt1        FormatName = 0xd4
	FixExt2        FormatName = 0xd5
	FixExt4        FormatName = 0xd6
	FixExt8        FormatName = 0xd7
	FixExt16       FormatName = 0xd8
	Str8           FormatName = 0xd9
	Str16          FormatName = 0xda
	Str32          FormatName = 0xdb
	Array16        FormatName = 0xdc
	Array32        FormatName = 0xdd
	Map16          FormatName = 0xde
	Map32          FormatName = 0xdf
	NegativeFixInt FormatName = 0xe0
)

func (i FormatName) Byte() byte {
	return byte(i)
}

func decodeFormatName(c byte) FormatName {
	if c < FixMap.Byte() {
		return PositiveFixInt
	}
	if c < FixArray.Byte() {
		return FixMap
	}
	if c < FixStr.Byte() {
		return FixArray
	}
	if c < Nil.Byte() {
		return FixStr
	}
	if c >= NegativeFixInt.Byte() {
		return NegativeFixInt
	}
	return FormatName(c)
}
