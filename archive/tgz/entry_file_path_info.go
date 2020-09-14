package tgz

import (
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"

	"go.nanasi880.dev/x/io"
)

// FilePathInfo is file path of archive entry.
type FilePathInfo struct {
	filePath string
	opt      *FilePathOption
}

// FilePathOption is file path entry option.
type FilePathOption struct {
	DontRecurse bool
}

func NewFilePath(p string, opt *FilePathOption) *FilePathInfo {

	if opt == nil {
		opt = new(FilePathOption)
	}

	return &FilePathInfo{
		filePath: p,
		opt:      opt,
	}
}

func (p *FilePathInfo) process(w *tar.Writer) error {

	if p.opt.DontRecurse {
		info, err := os.Stat(p.filePath)
		if err != nil {
			return err
		}
		return p.processSingleFile(w, p.filePath, info)
	}

	return filepath.Walk(p.filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return p.processSingleFile(w, path, info)
	})
}

func (p *FilePathInfo) processSingleFile(w *tar.Writer, path string, info os.FileInfo) error {

	if err := p.validateFile(info); err != nil {
		return err
	}

	fInfo := newFileInfo(info)
	header := fInfo.tarHeader(path)

	if err := w.WriteHeader(header); err != nil {
		return err
	}

	if header.Typeflag == tar.TypeDir {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return io.Copy(w, f)
}

func (p *FilePathInfo) validateFile(info os.FileInfo) error {

	const unsupported = os.ModeType & (^os.ModeDir) & (^os.ModeSymlink)

	mode := info.Mode()
	if mode&unsupported > 0 {
		return fmt.Errorf("unsupported file mode: %v", mode)
	}
	return nil
}
