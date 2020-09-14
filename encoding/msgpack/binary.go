package msgpack

type bigEndian [8]byte

func (e bigEndian) ReadUInt8() byte {
	return e[0]
}

func (e bigEndian) ReadInt8() int8 {
	return int8(e.ReadUInt8())
}

func (e bigEndian) ReadUInt16() uint16 {
	return uint16(e[1]) | uint16(e[0])<<8
}

func (e bigEndian) ReadInt16() int16 {
	return int16(e.ReadUInt16())
}

func (e bigEndian) ReadUInt32() uint32 {
	return uint32(e[3]) | uint32(e[2])<<8 | uint32(e[1])<<16 | uint32(e[0])<<24
}

func (e bigEndian) ReadInt32() int32 {
	return int32(e.ReadUInt32())
}

func (e bigEndian) ReadUInt64() uint64 {
	return uint64(e[7]) | uint64(e[6])<<8 | uint64(e[5])<<16 | uint64(e[4])<<24 |
		uint64(e[3])<<32 | uint64(e[2])<<40 | uint64(e[1])<<48 | uint64(e[0])<<56
}

func (e bigEndian) ReadInt64() int64 {
	return int64(e.ReadUInt64())
}

func newBigEndian2(b [2]byte) bigEndian {
	var e bigEndian
	copy(e[:], b[:])
	return e
}

func newBigEndian4(b [4]byte) bigEndian {
	var e bigEndian
	copy(e[:], b[:])
	return e
}
