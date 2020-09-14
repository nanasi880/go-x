package io_test

import (
	"bytes"
	"testing"

	"go.nanasi880.dev/x/io"
)

func TestAsyncWriter_Await(t *testing.T) {

	out := new(bytes.Buffer)

	w := io.NewAsyncWriter(io.NopWriteCloser(out))

	_, _ = w.Write([]byte("Hello"))
	_, _ = w.Write([]byte("World"))

	err := w.Await()
	if err != nil {
		t.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	if out.String() != "HelloWorld" {
		t.Fatal(out.String())
	}
}
