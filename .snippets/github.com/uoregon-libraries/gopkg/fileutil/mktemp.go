package fileutil

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextSuffix() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// TempFile is a rewrite of the built-in Go function with a very important
// addition: we can specify the file's extension.  This is crucial for
// third-party binaries which rely on the extension of a file.
func TempFile(dir, prefix, ext string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextSuffix()) + ext
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}

// TempNamedFile creates *and then closes* a file, simplifying the common
// pattern of creating a temp file, closing it, and passing the filename to
// external processes / scripts
func TempNamedFile(dir, prefix, ext string) (string, error) {
	var f, err = TempFile(dir, prefix, ext)
	if err != nil {
		return "", err
	}

	var n = f.Name()
	f.Close()
	return n, nil
}
