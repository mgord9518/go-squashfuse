module main

go 1.18

replace github.com/mgord9518/go-squashfuse => ../../

require (
	github.com/hanwen/go-fuse/v2 v2.1.0
	github.com/mgord9518/go-squashfuse v0.0.0-00010101000000-000000000000
)

require (
	github.com/anuvu/squashfs v0.0.0-20220404153901-d496132b2781 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
)
