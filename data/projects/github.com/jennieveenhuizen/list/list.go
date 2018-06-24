// © 2018 Jennie Veenhuizen
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the “Software”), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// The list command lists files. If no files are specified, it lists the
// current directory.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

var (
	listDirAsFile bool
	listDotFiles  bool
	recursive     bool
	showSize      bool
	showTime      bool
	showType      bool
)

func init() {
	flag.BoolVar(&listDotFiles,
		"a", false, "Do not ignore files starting with a dot.")
	flag.BoolVar(&listDirAsFile,
		"d", false, "List directories themselves, not their contents.")
	flag.BoolVar(&showType,
		"f", false, "Show file type.")
	flag.BoolVar(&recursive,
		"r", false, "List subdirectories recursively.")
	flag.BoolVar(&showSize,
		"s", false, "Show size.")
	flag.BoolVar(&showTime,
		"t", false, "Show time of last modification.")
}

type printer struct{}

func (printer) Write(b []byte) (int, error) {
	n, err := os.Stdout.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	return n, nil
}

var out = bufio.NewWriter(printer{})

var exitStatus = 0

func logError(err error) {
	// Flush the output first to print the error in the right context.
	out.Flush()
	log.Println(err)
	exitStatus = 1
}

type file struct {
	name string
	os.FileInfo
}

func (f file) writeType() {
	var s string
	switch f.Mode() & os.ModeType {
	case 0:
		s = "f   "
	case os.ModeDir:
		s = "dir "
	case os.ModeSymlink:
		s = "link"
	case os.ModeNamedPipe:
		s = "pipe"
	case os.ModeSocket:
		s = "sock"
	case os.ModeDevice:
		s = "dev "
	default:
		s = "????"
	}
	out.WriteString(s)
}

func (f file) writeTime() {
	out.WriteString(f.ModTime().Format("2006-01-02 15:04"))
}

func (f file) writeSize() {
	fmt.Fprintf(out, "%12d", f.Size())
}

func (f file) writeName() {
	out.WriteString(f.name)
	if showType && f.Mode()&os.ModeSymlink != 0 {
		s, err := os.Readlink(f.name)
		if err != nil {
			return
		}
		out.WriteString(" -> ")
		out.WriteString(s)
	}
}

func (f file) list() {
	if showSize {
		f.writeSize()
		out.WriteRune(' ')
	}
	if showTime {
		f.writeTime()
		out.WriteRune(' ')
	}
	if showType {
		f.writeType()
		out.WriteRune(' ')
	}
	f.writeName()
	out.WriteRune('\n')
}

// toPrefix takes a directory name and returns a prefix for the filenames in
// that directory. The amount of cleanup done here is kept to a minimum in
// order for this to be portable across operating systems.
func toPrefix(s string) string {
	if s == "." {
		return ""
	}
	if strings.HasSuffix(s, string(os.PathSeparator)) {
		return s
	}
	return s + string(os.PathSeparator)
}

func (f file) listDir() {
	v, err := ioutil.ReadDir(f.name)
	if err != nil {
		logError(err)
		return
	}
	prefix := toPrefix(f.name)
	for _, fi := range v {
		if !listDotFiles && strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		g := file{prefix + fi.Name(), fi}
		g.list()
		if recursive && g.IsDir() {
			g.listDir()
		}
	}
}

func list(name string) {
	stat := os.Stat
	if showType {
		stat = os.Lstat
	}
	fi, err := stat(name)
	if err != nil {
		logError(err)
		return
	}
	f := file{name, fi}
	if !listDirAsFile && f.IsDir() {
		f.listDir()
	} else {
		f.list()
	}
}

func main() {
	flag.Parse()
	v := flag.Args()
	if len(v) == 0 {
		list(".")
	} else {
		sort.Strings(v)
		for i := range v {
			list(v[i])
		}
	}
	out.Flush()
	os.Exit(exitStatus)
}
