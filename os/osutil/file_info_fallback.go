//go:build !darwin && !windows && !linux
// +build !darwin,!windows,!linux

package osutil

import (
	"os"
	"time"
)

type fileSysInfo struct {
	info os.FileInfo
}

func (f *fileSysInfo) UID() int {
	return 0 // TODO
}

func (f *fileSysInfo) GID() int {
	return 0 // TODO
}

func (f *fileSysInfo) LastModifyTime() time.Time {
	return f.info.ModTime() // TODO
}

func (f *fileSysInfo) LastAccessTime() time.Time {
	return f.info.ModTime() // TODO
}

func (f *fileSysInfo) LastChangeTime() time.Time {
	return f.info.ModTime() // TODO
}

func newFileInfo(info os.FileInfo) FileSysInfo {
	return &fileSysInfo{
		info: info,
	}
}
