// Cgo squashfuse build with SquashFS-tools-ng using anuvu's
// bindings

package squashfuse

import (
	"github.com/anuvu/squashfs"
	"github.com/hanwen/go-fuse/v2/fs"
	"context"
	"github.com/hanwen/go-fuse/v2/fuse"
	"syscall"
	"io/ioutil"
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
	fs.Inode
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

// TODO: Figure out how to read data without copying into RAM
func (f *File) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var err error

	target := f.Path

	info, _ := f.SquashFS.sqfs.Lstat(target)

	if len(info.SymlinkTarget) > 0 {
		target = info.SymlinkTarget[1:]
	}

	f2, _ := squashfs.Open(target, f.SquashFS.sqfs)
	defer f2.Close()

	if f.Data == nil {
		f.Data, err = ioutil.ReadAll(f2)
	}

	return err, fuse.FOPEN_KEEP_CACHE, 0
}

func (f *File) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Size  = f.Attr.Size
	out.Mode  = f.Attr.Mode
	out.Mtime = f.Attr.Mtime
	out.Blksize = 512
	return 0
}

func (f *File) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	end := int(off) + len(dest)

	if end > len(f.Data) {
		end = len(f.Data)
	}

	return fuse.ReadResultData(f.Data[off:end]), 0
}

