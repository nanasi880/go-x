package tgz

import "archive/tar"

// Entry is archive entry.
type Entry interface {
	process(w *tar.Writer) error
}
