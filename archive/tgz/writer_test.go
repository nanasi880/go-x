package tgz_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"path/filepath"
	"testing"

	"go.nanasi880.dev/x/archive/tgz"
	"go.nanasi880.dev/x/internal/testing/testutil"
)

func TestEncoder(t *testing.T) {
	testSuites := []struct {
		Entries []tgz.Entry
		Names   []string
	}{
		{
			Entries: []tgz.Entry{
				tgz.NewFilePathInfo("testdata", nil),
			},
			Names: []string{
				"testdata",
				"testdata/sample.txt",
			},
		},
		{
			Entries: []tgz.Entry{
				tgz.FilePath("testdata"),
			},
			Names: []string{
				"testdata",
				"testdata/sample.txt",
			},
		},
		{
			Entries: []tgz.Entry{
				tgz.NewFilePathInfo("testdata", &tgz.FilePathOption{
					DontRecurse: true,
				}),
			},
			Names: []string{
				"testdata",
			},
		},
	}
	for suiteNo, suite := range testSuites {
		func() {
			buf := new(bytes.Buffer)

			w := tgz.NewWriter(buf)
			w.Add(suite.Entries...)
			err := w.Write()
			if err != nil {
				testutil.Failf(t, "suiteNo:%d err:%v", suiteNo, err)
				return
			}

			gz, err := gzip.NewReader(buf)
			if err != nil {
				testutil.Failf(t, "suiteNo:%d err:%v", suiteNo, err)
				return
			}
			defer testutil.Close(t, gz)
			r := tar.NewReader(gz)

			for {
				header, err := r.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					testutil.Failf(t, "suiteNo:%d err:%v", suiteNo, err)
					return
				}
				found := false
				for _, name := range suite.Names {
					if header.Name == filepath.FromSlash(name) {
						found = true
						break
					}
				}
				if !found {
					testutil.Failf(t, "suiteNo:%d err:%v", suiteNo, err)
				}
			}
		}()
	}
}
