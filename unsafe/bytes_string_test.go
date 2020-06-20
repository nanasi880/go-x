package unsafe_test

import (
	"bytes"
	"testing"

	"go.nanasi880.dev/x/unsafe"
)

func TestBytesToString(t *testing.T) {

	want := "abc"
	got := unsafe.BytesToString([]byte(want))

	if want != got {
		t.Fatalf("want: %s got: %s", want, got)
	}
}

func TestStringToBytes(t *testing.T) {

	want := []byte{'a', 'b', 'c'}
	got := unsafe.StringToBytes("abc")

	if bytes.Compare(want, got) != 0 {
		t.Fatal(got)
	}

	if len(got) != 3 {
		t.Fatal(len(got))
	}
	if cap(got) != 3 {
		t.Fatal(cap(got))
	}
}
