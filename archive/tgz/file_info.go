package tgz

import (
	"archive/tar"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

type fileInfo struct {
	info os.FileInfo
}

func newFileInfo(info os.FileInfo) *fileInfo {
	return &fileInfo{
		info: info,
	}
}

func (f *fileInfo) tarHeader(path string) *tar.Header {

	info := f.info

	header := new(tar.Header)
	header.Typeflag = f.typeFlag()
	if header.Typeflag == tar.TypeSymlink {
		header.Linkname = path
	} else {
		header.Name = path
	}
	header.Size = info.Size()
	header.Mode = int64(info.Mode())

	header.ModTime = info.ModTime()

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		header.Uid = int(stat.Uid)
		header.Gid = int(stat.Gid)

		u, err := user.LookupId(strconv.Itoa(header.Uid))
		if err == nil {
			header.Uname = u.Name
		}

		g, err := user.LookupGroupId(strconv.Itoa(header.Gid))
		if err == nil {
			header.Gname = g.Name
		}

		header.ModTime = time.Unix(stat.Mtimespec.Unix()).UTC()
		header.AccessTime = time.Unix(stat.Atimespec.Unix()).UTC()
		header.ChangeTime = time.Unix(stat.Ctimespec.Unix()).UTC()
	}

	return header
}

func (f *fileInfo) typeFlag() byte {
	mode := f.info.Mode()
	switch {
	case mode&os.ModeSymlink > 0:
		return tar.TypeSymlink
	case mode&os.ModeDir > 0:
		return tar.TypeDir
	default:
		return tar.TypeReg
	}
}
