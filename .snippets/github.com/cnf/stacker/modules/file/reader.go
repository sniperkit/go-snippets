package file

import (
	"errors"
	// "fmt"
	"io"
	"os"
)

// http://golang.org/src/pkg/bufio/bufio.go?s=3623:3673#L143

// Reader reads
type Reader struct {
	buf      []byte
	rd       *os.File
	r, w     int
	err      error
	offset   int
	curfile  int
	filelist []string
}

const minReadBufferSize = 16
const maxConsecutiveEmptyReads = 100

var errNegativeRead = errors.New("reader returned negative count from Read")

func NewReader(filelist []string) io.Reader {
	return &Reader{
		filelist: filelist,
	}
}

// Read reads data into p
func (r *Reader) Read(p []byte) (n int, err error) {
	n = len(p)
	r.buf = make([]byte, n)

	if n == 0 {
		return 0, r.readErr()
	}

	// fmt.Printf("opening [%d]: %s\n", r.curfile, r.filelist[r.curfile])
	r.rd, r.err = os.Open(r.filelist[r.curfile])
	defer r.rd.Close()

	if r.r == r.w {
		// TODO: check for errors
		if r.err != nil {
			return 0, r.readErr()
		}
		r.fill()
		if r.r == r.w {
			err := r.readErr()
			if err == io.EOF && r.curfile < len(r.filelist)-1 {
				r.curfile++
				r.r = 0
				r.w = 0
				// fmt.Printf("next file\n")
				return 0, nil
			}
			return 0, err
		}
	}

	if n > r.w-r.r {
		n = r.w - r.r
	}

	// fmt.Printf("---\n%s", r.buf)
	copy(p[0:n], r.buf[r.r:])
	r.r += n

	return n, nil
}

func (r *Reader) fill() {
	// Read new data: try a limited number of times.
	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		// n, err := r.rd.Read(r.buf[r.w:])
		r.rd.Seek(int64(r.w), 0)
		n, err := r.rd.Read(r.buf)
		// fmt.Printf("%d\n", n)
		if n < 0 {
			panic(errNegativeRead)
		}
		r.w += n
		if err != nil {
			r.err = err
			return
		}
		if n > 0 {
			return
		}
	}
	r.err = io.ErrNoProgress
}

func (r *Reader) readErr() error {
	err := r.err
	r.err = nil
	// fmt.Printf("ERROR: %s\n", err.Error())
	return err
}
