package base64util

import (
	"encoding/base64"
	"io"
	"os"
)

// Decode is decode base64 data from byte slice.
func Decode(enc *base64.Encoding, data []byte) ([]byte, error) {
	length := enc.DecodedLen(len(data))
	result := make([]byte, length)
	_, err := enc.Decode(result, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeFromReader is decode base64 data from io.Reader.
func DecodeFromReader(enc *base64.Encoding, r io.Reader) ([]byte, error) {
	decoder := base64.NewDecoder(enc, r)
	return io.ReadAll(decoder)
}

// DecodeFromFile is decode base64 data from file.
func DecodeFromFile(enc *base64.Encoding, filename string) (result []byte, e error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if er := f.Close(); er != nil {
			result = nil
			e = er
		}
	}()
	return DecodeFromReader(enc, f)
}
