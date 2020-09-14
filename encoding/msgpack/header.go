package msgpack

import "fmt"

// Format is message pack format data.
type Format struct {
	Name FormatName
	Raw  byte
}

func (f Format) String() string {
	return fmt.Sprintf("Name: %s Raw: %02x", f.Name.String(), f.Raw)
}

// ExtHeader is header of ext data.
type ExtHeader struct {
	Format Format
	Type   byte
	Length uint32
}

// ArrayHeader is header of array data.
type ArrayHeader struct {
	Format Format
	Length uint32
}

// MapHeader is header of map data.
type MapHeader struct {
	Format Format
	Length uint32
}
