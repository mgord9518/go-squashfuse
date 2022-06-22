package squashfuse

import (
	"context"

	"syscall"
	"github.com/hanwen/go-fuse/v2/fuse"
    "github.com/anuvu/squashfs"
	"github.com/hanwen/go-fuse/v2/fs"
	"strings"
	"path/filepath"

)

var _ = (fs.NodeOnAdder)((*SquashFS)(nil))
var _ = (fs.NodeGetattrer)((*File)(nil))

func (s *SquashFS) OnAdd(ctx context.Context) {

	if s.sqfs == nil {
		panic("SquashFS must be created with squashfuse.Open()!")
	}

    s.sqfs.Walk("/", func(path string, info squashfs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		dir, base := filepath.Split(filepath.Clean(path))

		p := &s.Inode
		for _, str := range strings.Split(dir, "/") {
			if len(str) == 0 {
				continue
			}

			ch := p.GetChild(str)
			if ch == nil {
				ch = p.NewPersistentInode(
				ctx, &fs.Inode{}, fs.StableAttr{
					Mode: fuse.S_IFDIR,
				})
				p.AddChild(str, ch, true)
			}

			p = ch
		}

		f, err := squashfs.Open(path, s.sqfs)
		defer f.Close()
		if err != nil { return err }

		if len(base) < 1 {
			return nil
		}

		if len(info.SymlinkTarget) == 0 {
			ch := p.NewPersistentInode(
			ctx, &File{
					Path: filepath.Join(dir, base),
					Attr: fuse.Attr{
						Mode:  uint32(info.Mode()) & 07777,
						Mtime: uint64(info.ModTime().Unix()),
						Size:  uint64(f.Size()),
					},
					SquashFS: s,
				}, fs.StableAttr{})
			p.AddChild(base, ch, true)
		} else {
			ch := p.NewPersistentInode(
			ctx, &fs.MemSymlink{
					Data: []byte(info.SymlinkTarget),
				}, fs.StableAttr{Mode: syscall.S_IFLNK})
			p.AddChild(base, ch, true)
		}

		return nil

		})
}

// TODO: Properly report attr
func (r *SquashFS) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0755
	return 0
}
