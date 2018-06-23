package ptype

import "io"

type OutputRewriter interface {
	RewriteOutput(io.Reader, io.Writer) error
}
