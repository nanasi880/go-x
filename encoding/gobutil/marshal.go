package gobutil

import (
	"bytes"
	"encoding/gob"
)

// Marshal is marshal data as gob.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
