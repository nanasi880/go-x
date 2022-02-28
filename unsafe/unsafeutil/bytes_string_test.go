package unsafeutil_test

import (
	"bytes"
	"testing"

	"go.nanasi880.dev/x/unsafe/unsafeutil"
)

func TestBytesToString(t *testing.T) {

	want := "abc"
	got := unsafeutil.BytesToString([]byte(want))

	if want != got {
		t.Fatalf("want: %s got: %s", want, got)
	}
}

func TestStringToBytes(t *testing.T) {

	want := []byte{'a', 'b', 'c'}
	got := unsafeutil.StringToBytes("abc")

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
