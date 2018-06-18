// Package rotate implements a Writer that will write to files in
// a directory and rotate them when they reach a specific size.
// It will also only keep a fixed number of files.
// It can be used anywhere an io.Writer is used. for example in
// log.SetOutput

// copy from https://github.com/stathat/rotate

package rotate

import (
	"errors"
	"fmt"
	"os"
	pathPkg "path"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	maxDefault  = 1024 * 1024 * 8
	keepDefault = 10

	currentDefault = "current"
)

// RootPerm defines the permissions that Writer will use if it
// need to create the root directory.
var RootPerm = os.FileMode(0755)

// FilePerm defines the permissions that Writer will use for all
// the files it creates.
var FilePerm = os.FileMode(0666)

var (
	ErrorRootMustDirectory error = errors.New("rotate: root must be a directory")
)

type Writer struct {
	root    string
	prefix  string
	max     int
	keep    int
	current *os.File

	mutex sync.Mutex
	size  int
}

func New(root, prefix string) (*Writer, error) {
	w := &Writer{
		root:   root,
		prefix: prefix,
		max:    maxDefault,
		keep:   keepDefault,
	}
	if err := w.setup(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Writer) SetMax(max int) {
	w.max = max
}

func (w *Writer) SetKeep(keep int) {
	w.keep = keep
}

func (w *Writer) setup() error {
	fi, err := os.Stat(w.root)
	if err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(w.root, RootPerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !fi.IsDir() {
		return ErrorRootMustDirectory
	}

	return w.openCurrent()
}

func (w *Writer) openCurrent() (err error) {
	path := pathPkg.Join(w.root, currentDefault)
	if w.current, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, FilePerm); err != nil {
		return err
	}
	w.size = 0
	return nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if n, err = w.current.Write(p); err != nil {
		return
	}

	w.size += n

	if w.size >= w.max {
		if err = w.rotate(); err != nil {
			return
		}
	}

	return n, nil
}

func (w *Writer) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if err := w.current.Close(); err != nil {
		return err
	}
	return nil
}

func (w *Writer) rotate() error {
	if err := w.current.Close(); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_%d", w.prefix, time.Now().UnixNano())
	if err := os.Rename(pathPkg.Join(w.root, currentDefault), pathPkg.Join(w.root, filename)); err != nil {
		return err
	}

	if err := w.clean(); err != nil {
		return err
	}

	return w.openCurrent()
}

func (w *Writer) clean() error {
	dir, err := os.Open(w.root)
	if err != nil {
		return err
	}

	files, err := dir.Readdirnames(1024)
	if err != nil {
		return err
	}

	var filenames []string
	for _, file := range files {
		if strings.HasPrefix(file, w.prefix+"_") {
			filenames = append(filenames, file)
		}
	}

	if len(filenames) < w.keep {
		return nil
	}

	sort.Strings(filenames)

	for i := 0; i < len(filenames)-w.keep; i++ {
		if err := os.Remove(pathPkg.Join(w.root, filenames[i])); err != nil {
			return err
		}
	}

	return nil
}
