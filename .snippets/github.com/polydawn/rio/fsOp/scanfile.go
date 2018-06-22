package fsOp

import (
	"io"
	"os"

	"go.polydawn.net/rio/fs"
)

/*
	Scan file attributes into an `fs.Metadata` struct, and return an
	`io.ReadCloser` for the file content.

	The reader is nil if the path is any type other than a file.  If a
	reader is returned, the caller is expected to close it.
*/
func ScanFile(afs fs.FS, path fs.RelPath) (fmeta *fs.Metadata, body io.ReadCloser, err error) {
	// most of the heavy work is already done by fs.Lstat; this method just adds the file content.
	fmeta, err = afs.LStat(path)
	if err != nil {
		return fmeta, nil, err
	}
	switch fmeta.Type {
	case fs.Type_File:
		var err error
		body, err = afs.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			return fmeta, body, err
		}
	}
	return
}
