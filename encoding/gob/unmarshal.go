package gob

import (
	"bytes"
	"encoding/gob"
)

// Unmarshal is unmarshal data from gob.
func Unmarshal(b []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(b)).Decode(v)
}
