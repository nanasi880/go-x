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
	return int(f.sys().Uid)
}

func (f *fileSysInfo) GID() int {
	return int(f.sys().Gid)
}

func (f *fileSysInfo) LastModifyTime() time.Time {
	return time.Unix(f.sys().Mtim.Unix()).UTC()
}

func (f *fileSysInfo) LastAccessTime() time.Time {
	return time.Unix(f.sys().Atim.Unix()).UTC()
}

func (f *fileSysInfo) LastChangeTime() time.Time {
	return time.Unix(f.sys().Ctim.Unix()).UTC()
}

func (f *fileSysInfo) sys() *syscall.Stat_t {
	return f.info.Sys().(*syscall.Stat_t)
}

func newFileInfo(info os.FileInfo) FileSysInfo {
	return &fileSysInfo{
		info: info,
	}
}
