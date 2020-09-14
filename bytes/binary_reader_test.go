package bytes_test

import (
	"io"
	"strings"
	"testing"

	"go.nanasi880.dev/x/bytes"
)

type nullReader struct{}

func (r nullReader) Read(p []byte) (int, error) {
	return len(p), nil
}

func TestBinaryReader_Seek(t *testing.T) {
	r := bytes.NewBinaryReaderBuffer(strings.NewReader("0123456789"), make([]byte, 3))
	n, err := r.Seek(5, io.SeekCurrent)
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Fatal(n)
	}
	b, err := r.Read1()
	if err != nil {
		t.Fatal(err)
	}
	if b != '5' {
		t.Fatalf("%c", rune(b))
	}
}

func BenchmarkReader_Read(b *testing.B) {

	b.Run("Read1", func(b *testing.B) {
		reader := bytes.NewBinaryReader(nullReader{})
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := reader.Read1()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Read2", func(b *testing.B) {
		reader := bytes.NewBinaryReader(nullReader{})
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := reader.Read2()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Read4", func(b *testing.B) {
		reader := bytes.NewBinaryReader(nullReader{})
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := reader.Read4()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Read8", func(b *testing.B) {
		reader := bytes.NewBinaryReader(nullReader{})
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := reader.Read8()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
