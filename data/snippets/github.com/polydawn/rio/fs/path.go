package fs

import (
	"fmt"
	"path"
	"strings"
)

/*
	Relative paths and absolute paths are not interchangeable.
	Using these path structures aids in keeping these ideas separate and
	hopefully prevent certain kinds of errors caused by accidentally using the incorrect path type.
	It is recommended to use absolute paths when possible.
*/
type RelPath struct {
	path      string
	lastSplit int
}

/*
	Converts a string to an relative path struct.
	Will panic if the given path is absolute.
*/
func MustRelPath(p string) RelPath {
	p = path.Clean(p)
	if p[0] == '/' {
		panic(fmt.Errorf("fs: not a relative path (%q)", p))
	}
	if p == "." { // We can't stop people from using the zero value, so, use it.
		return RelPath{}
	}
	return RelPath{p, strings.LastIndexByte(p, '/')}
}

/*
	An relative path string with a zero value of "." for current directory
*/
func (p RelPath) String() string {
	if p.path == "" {
		return "."
	} else if len(p.path) == 2 && p.path[0:2] == ".." { // a '..' prefix
		return p.path
	} else if len(p.path) > 2 && p.path[0:3] == "../" { // a '..' prefix
		return p.path
	} else {
		return "./" + p.path
	}
}

/*
	Returns the relative path of the directory containing this path
*/
func (p RelPath) Dir() RelPath {
	if p.path == "" {
		return p
	} else if p.lastSplit == -1 {
		return RelPath{}
	} else {
		p2 := p.path[0:p.lastSplit]
		return RelPath{p2, strings.LastIndexByte(p2, '/')}
	}
}

/*
	Returns the filename of this path
*/
func (p RelPath) Last() string {
	if p.path == "" {
		return "."
	} else if p.lastSplit == -1 {
		return p.path
	} else {
		return p.path[p.lastSplit+1:]
	}
}

/*
	Creates a new relative path by following the relative path from this path.
*/
func (p RelPath) Join(p2 RelPath) RelPath {
	switch {
	case p2.path == "":
		return p
	case p.path == "":
		return p2
	case p2.path[0] == '.': // '..' prefix requires cleaning again.
		pj := path.Clean(p.path + "/" + p2.path)
		if pj == "." { // We can't stop people from using the zero value, so, use it.
			return RelPath{}
		}
		return RelPath{pj, strings.LastIndexByte(pj, '/')}
	default:
		return RelPath{p.path + "/" + p2.path, len(p.path) + p2.lastSplit + 1}
	}
}

/*
	Returns a slice of relative paths for each segment in the path.

	E.g. for "./a/b/c" it will return [".", "./a", "./a/b", "./a/b/c"].
*/
func (p RelPath) Split() []RelPath {
	if p.path == "" {
		return []RelPath{RelPath{}}
	}
	if p.lastSplit == -1 {
		return []RelPath{RelPath{}, p}
	}
	n := strings.Count(p.path, "/") + 1
	slice := make([]RelPath, n+1)
	slice[n] = p
	for i := n - 1; i >= 0; i-- {
		slice[i] = slice[i+1].Dir()
	}
	return slice
}

/*
	Returns a slice of relative paths for each segment in the path, minus the
	path itself.

	E.g. for "./a/b/c" it will return [".", "./a", "./a/b"].
	Importantly, for "." it will return [] -- an empty list (note that this
	is different than the behavior of `path.Dir().Split()`, which would still
	yield ["."].)

	For some path where you wish to create all parent paths, it may be useful
	to `range path.SplitParent()`.
*/
func (p RelPath) SplitParent() []RelPath {
	if p.path == "" {
		return []RelPath{}
	}
	if p.lastSplit == -1 {
		return []RelPath{RelPath{}}
	}
	n := strings.Count(p.path, "/")
	slice := make([]RelPath, n+1)
	slice[n] = p.Dir()
	for i := n - 1; i >= 0; i-- {
		slice[i] = slice[i+1].Dir()
	}
	return slice
}

/*
	Predicate for if this path goes "up" -- in other words, if it starts with
	"..".
*/
func (p RelPath) GoesUp() bool {
	return len(p.path) >= 2 && p.path[0:2] == ".."
}

/*
	Relative paths and absolute paths are not interchangeable.
	Using these path structures aids in keeping these ideas separate and
	hopefully prevent certain kinds of errors caused by accidentally using the incorrect path type.
	It is recommended to use absolute paths when possible.
*/
type AbsolutePath struct {
	path      string
	lastSplit int
}

/*
	Converts a string to an absolute path struct.
	Will panic if the given path is not absolute.
*/
func MustAbsolutePath(p string) AbsolutePath {
	path, err := ParseAbsolutePath(p)
	if err != nil {
		panic(err)
	}
	return path
}

/*
	Converts a string to an absolute path struct,
	returning an error if the given path string is not absolute.
*/
func ParseAbsolutePath(p string) (AbsolutePath, error) {
	p = path.Clean(p)
	if p[0] != '/' {
		return AbsolutePath{}, fmt.Errorf("fs: not an absolute path (%q)", p)
	}
	if p == "/" { // We can't stop people from using the zero value, so, use it.
		return AbsolutePath{}, nil
	}
	return AbsolutePath{p, strings.LastIndexByte(p, '/')}, nil

}

/*
	An absolute path string with a zero value of "/" for the root directory
*/
func (p AbsolutePath) String() string {
	if p.path == "" {
		return "/"
	}
	return p.path
}

/*
	Returns the absolute path of the directory containing this path
*/
func (p AbsolutePath) Dir() AbsolutePath {
	if p.path == "" {
		return p
	} else if p.lastSplit == 0 {
		return AbsolutePath{}
	} else {
		p2 := p.path[0:p.lastSplit]
		return AbsolutePath{p2, strings.LastIndexByte(p2, '/')}
	}
}

/*
	Returns the filename of this path
*/
func (p AbsolutePath) Last() string {
	if p.path == "" {
		return "/"
	} else {
		return p.path[p.lastSplit+1:]
	}
}

/*
	Creates a new absolute path by following the relative path from this path.
*/
func (p AbsolutePath) Join(p2 RelPath) AbsolutePath {
	switch {
	case p2.path == "":
		return p
	//case p.path == "": // Comes out the same as the math in the default case.
	//	return AbsolutePath{"/" + p2.path, p2.lastSplit + 1}
	case p2.path[0] == '.': // '..' prefix requires cleaning again.
		pj := path.Clean(p.path + "/" + p2.path)
		if pj == "/" { // We can't stop people from using the zero value, so, use it.
			return AbsolutePath{}
		}
		return AbsolutePath{pj, strings.LastIndexByte(pj, '/')}
	default:
		return AbsolutePath{p.path + "/" + p2.path, len(p.path) + p2.lastSplit + 1}
	}
}

func (p AbsolutePath) CoerceRelative() RelPath {
	return MustRelPath("." + p.path)
}
