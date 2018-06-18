//+build linux

package fileutil

import (
	"os"
	"syscall"
)

// dirEntBuffer is the number of bytes allocated when calling the low-level
// ReadDirent function.  32k is the value C uses, but this *must* be set to a
// value that's higher than the CIFS responses, which we are currently seeing
// as roughly 6k.
const dirEntBuffer = 32768

// Readdir calls low-level syscall functions in order to avoid what appears to
// be a kernel bug in CentOS 7 when reading from CIFS filesystems.  See
// https://github.com/golang/go/issues/24015 and/or
// https://stackoverflow.com/questions/46719753/golang-bug-ioutil-readdir-listing-files-on-cifs-share-or-doing-something-wro
// for some context.
func Readdir(path string) ([]os.FileInfo, error) {
	var d *os.File
	var err error
	var bufp, nbuf int

	d, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	var buf = make([]byte, dirEntBuffer)
	var names = make([]string, 0, 100)
	for {
		if bufp >= nbuf {
			bufp = 0
			nbuf, err = syscall.ReadDirent(int(d.Fd()), buf)
			if err != nil {
				return nil, err
			}
		}

		// Check for EOF
		if nbuf <= 0 {
			break
		}

		var nb int
		nb, _, names = syscall.ParseDirent(buf[bufp:nbuf], -1, names)
		bufp += nb
	}

	var items []os.FileInfo
	for _, fname := range names {
		var item, err = os.Stat(path + "/" + fname)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, err
}
