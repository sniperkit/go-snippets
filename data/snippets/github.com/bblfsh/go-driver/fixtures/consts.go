package fixtures

import "unsafe"

const a = 1

const (
	b        = int(1)
	c    int = 1
	d, e     = 1, "e"
	f        = c + d
	g        = unsafe.Sizeof(int(5))
)
