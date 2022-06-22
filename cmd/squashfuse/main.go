package main

import (
    "github.com/mgord9518/go-squashfuse"
	"github.com/hanwen/go-fuse/v2/fs"
)

func main() {
	opts := &fs.Options{}

	// Just for testing currently, eventually will have arguments and error
	// handling
	s, _ := squashfuse.Open("sfs")

	server, _ := fs.Mount("test", s, opts)
	server.Wait()
}
