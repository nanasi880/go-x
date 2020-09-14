package tgz

import "archive/tar"

// FilePath is file path of archive entry
type FilePath string

func (p FilePath) process(w *tar.Writer) error {
	info := &FilePathInfo{
		filePath: string(p),
		opt:      new(FilePathOption),
	}
	return info.process(w)
}
