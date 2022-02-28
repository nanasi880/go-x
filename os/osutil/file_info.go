package osutil

import (
	"os"
	"time"
)

// FileSysInfo is accessor of FileInfo.Sys().
type FileSysInfo interface {

	// UID is get UNIX UID.
	// return 0, if not supported filesystem.
	UID() int

	// GID is get UNIX GID.
	// return 0, if not supported filesystem.
	GID() int

	// LastModifyTime is get last modify timestamp.
	// LastModifyTime is update when call write(2) syscall to file.
	LastModifyTime() time.Time

	// LastAccessTime is get last access timestamp.
	// LastAccessTime is update when call read(2) syscall to file.
	LastAccessTime() time.Time

	// LastChangeTime is get last access timestamp.
	// LastChangeTime is update when call write(2) syscall to file or update file metadata.
	LastChangeTime() time.Time
}

// FileInfo is create FileSysInfo instance. panic if `info` is nil.
func FileInfo(info os.FileInfo) FileSysInfo {
	if info == nil {
		panic("nil")
	}
	return newFileInfo(info)
}
