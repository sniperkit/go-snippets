package nilFS

import (
	"time"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/rio/fs"
)

func New() fs.FS {
	return &nilFS{fs.MustAbsolutePath("/-")}
}

type nilFS struct {
	basePath fs.AbsolutePath
}

func (afs *nilFS) BasePath() fs.AbsolutePath {
	return afs.basePath
}

func (afs *nilFS) OpenFile(path fs.RelPath, flag int, perms fs.Perms) (fs.File, error) {
	_, err := afs.realpath(path, false)
	if err != nil {
		return nil, err
	}
	return nilFile{}, nil
}

func (afs *nilFS) Mkdir(path fs.RelPath, perms fs.Perms) error {
	_, err := afs.realpath(path, false)
	if err != nil {
		return err
	}
	return nil
}

func (afs *nilFS) Mklink(path fs.RelPath, target string) error {
	_, err := afs.realpath(path, false)
	if err != nil {
		return err
	}
	return nil
}

func (afs *nilFS) Mkfifo(path fs.RelPath, perms fs.Perms) error {
	_, err := afs.realpath(path, false)
	if err != nil {
		return err
	}
	return nil
}

func (afs *nilFS) MkdevBlock(path fs.RelPath, major int64, minor int64, perms fs.Perms) error {
	_, err := afs.realpath(path, false)
	if err != nil {
		return err
	}
	return nil
}

func (afs *nilFS) MkdevChar(path fs.RelPath, major int64, minor int64, perms fs.Perms) error {
	_, err := afs.realpath(path, false)
	if err != nil {
		return err
	}
	return nil
}

func (afs *nilFS) Lchown(path fs.RelPath, uid uint32, gid uint32) error {
	_, err := afs.realpath(path, false)
	if err != nil {
		return err
	}
	return nil
}

func (afs *nilFS) Chmod(path fs.RelPath, perms fs.Perms) error {
	_, err := afs.realpath(path, false)
	if err != nil {
		return err
	}
	return nil
}

func (afs *nilFS) SetTimesLNano(path fs.RelPath, mtime time.Time, atime time.Time) error {
	_, err := afs.realpath(path, false)
	if err != nil {
		return err
	}
	return nil
}

func (afs *nilFS) SetTimesNano(path fs.RelPath, mtime time.Time, atime time.Time) error {
	_, err := afs.realpath(path, false)
	if err != nil {
		return err
	}
	return nil
}

func (afs *nilFS) Stat(path fs.RelPath) (*fs.Metadata, error) {
	_, err := afs.realpath(path, false)
	if err != nil {
		return nil, err
	}
	return &fs.Metadata{}, nil
}

func (afs *nilFS) LStat(path fs.RelPath) (*fs.Metadata, error) {
	_, err := afs.realpath(path, false)
	if err != nil {
		return nil, err
	}
	return &fs.Metadata{}, nil
}

func (afs *nilFS) ReadDirNames(path fs.RelPath) ([]string, error) {
	_, err := afs.realpath(path, false)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (afs *nilFS) Readlink(path fs.RelPath) (string, bool, error) {
	_, err := afs.realpath(path, false)
	if err != nil {
		return "", false, err
	}
	return "", false, nil
}

// resolves a path.
// resolving a path can have errors traversing things and still return nil error,
//  because failure to resolve the path doesn't necessarily mean you shouldn't try.
// (it does however return real errors in case of ErrRecurse and ErrBreakout.)
func (afs *nilFS) realpath(path fs.RelPath, resolveLast bool) (string, error) {
	if path.GoesUp() {
		return "", Errorf(fs.ErrBreakout, "fs: invalid path %q: must not depart basepath", path)
	}
	return afs.BasePath().Join(path).String(), nil
}

func (afs *nilFS) ResolveLink(symlink string, startingAt fs.RelPath) (fs.RelPath, error) {
	if startingAt.GoesUp() {
		return startingAt, Errorf(fs.ErrBreakout, "fs: invalid path %q: must not depart basepath", startingAt)
	}
	return startingAt, nil
}

var _ fs.File = nilFile{}

type nilFile struct{}

func (nilFile) Close() error                            { return nil }
func (nilFile) Read(bs []byte) (int, error)             { return len(bs), nil }
func (nilFile) ReadAt(bs []byte, _ int64) (int, error)  { return len(bs), nil }
func (nilFile) Seek(int64, int) (int64, error)          { return 0, nil }
func (nilFile) Write(bs []byte) (int, error)            { return len(bs), nil }
func (nilFile) WriteAt(bs []byte, _ int64) (int, error) { return len(bs), nil }
