// Cgo squashfuse build with SquashFS-tools-ng using anuvu's
// bindings

package squashfuse

import (
	"github.com/anuvu/squashfs"
	"github.com/hanwen/go-fuse/v2/fs"
	"context"
	"github.com/hanwen/go-fuse/v2/fuse"
	"syscall"
//	"io/ioutil"
	"sync"
	//"os"
)

type SquashFS struct {
	fs.Inode
	Path      string
	sqfs     *squashfs.SquashFs
}

type File struct {
	Path       string
	SquashFS   *SquashFS
	Data     []byte
	Attr       fuse.Attr
	mu         sync.Mutex
	offset     int64
	fs.Inode
	file       *squashfs.File
}

type WalkFunc squashfs.WalkFunc

func (s *SquashFS) Walk(root string, fn WalkFunc) error {
	return s.sqfs.Walk(root, squashfs.WalkFunc(fn))
}

func Open(fname string) (*SquashFS, error) {
	s := &SquashFS{}

	sqfs, err := squashfs.OpenSquashfs(fname)

	s.Path = sqfs.Filename
	s.sqfs = &sqfs

	return s, err
}

func (f *File) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var err error

	target := f.Path

	info, _ := f.SquashFS.sqfs.Lstat(target)

	if len(info.SymlinkTarget) > 0 {
		target = info.SymlinkTarget[1:]
	}

	// Need to close this file somewhere
	f.file, _ = squashfs.Open(target, f.SquashFS.sqfs)

	return err, fuse.FOPEN_KEEP_CACHE, 0
}

func (f *File) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Size  = f.Attr.Size
	out.Mode  = f.Attr.Mode
	out.Mtime = f.Attr.Mtime
	out.Blksize = 512
	return 0
}

// TODO: fix crash when reading multiple files at once
// TODO: make files seekable, it can currently only return streams
func (f *File) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	var errno syscall.Errno

	n, err := f.file.Read(dest)
	if err != nil {
		errno = 1
	}

	return fuse.ReadResultData(dest[:n]), errno
}
