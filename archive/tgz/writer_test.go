package tgz_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"go.nanasi880.dev/x/archive/tgz"
)

func TestEncoder(t *testing.T) {

	buf := new(bytes.Buffer)

	err := tgz.NewWriter(buf).
		Add(tgz.FilePath("./testdata")).
		Write()

	if err != nil {
		t.Fatal(err)
	}

	gzReader, err := gzip.NewReader(buf)
	if err != nil {
		t.Fatal(err)
	}
	reader := tar.NewReader(gzReader)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		t.Log(header.Name)
	}
}
