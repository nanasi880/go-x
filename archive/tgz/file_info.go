package tgz

import (
	"archive/tar"
	"os"
	"os/user"
	"strconv"

	"go.nanasi880.dev/x/os/osutil"
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

	sysInfo := osutil.FileInfo(info)

	header.Uid = sysInfo.UID()
	header.Gid = sysInfo.GID()

	if u, err := user.LookupId(strconv.Itoa(header.Uid)); err == nil {
		header.Uname = u.Name
	}
	if g, err := user.LookupGroupId(strconv.Itoa(header.Gid)); err == nil {
		header.Gname = g.Name
	}

	header.ModTime = sysInfo.LastModifyTime()
	header.AccessTime = sysInfo.LastAccessTime()
	header.ChangeTime = sysInfo.LastChangeTime()

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
