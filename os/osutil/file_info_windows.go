package osutil

import (
	"os"
	"syscall"
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
	return time.Unix(0, f.sys().LastWriteTime.Nanoseconds()).UTC()
}

func (f *fileSysInfo) LastAccessTime() time.Time {
	return time.Unix(0, f.sys().LastAccessTime.Nanoseconds()).UTC()
}

func (f *fileSysInfo) LastChangeTime() time.Time {
	return time.Unix(0, f.sys().LastWriteTime.Nanoseconds()).UTC()
}

func (f *fileSysInfo) sys() *syscall.Win32FileAttributeData {
	return f.info.Sys().(*syscall.Win32FileAttributeData)
}

func newFileInfo(info os.FileInfo) FileSysInfo {
	return &fileSysInfo{
		info: info,
	}
}
