package fs

import (
	"io"
	"time"
)

/*
	Interface for all primitive functions we expect to be able to perform
	on a filesystem.

	All paths accepted are RelPath types; typically the FS instance
	is constructed with an AbsolutePath, and all further operations are
	joined with that base path.
*/
type FS interface {
	// The basepath this filesystem was constructed with.
	// This may be useful for debugging messages, but should not otherwise
	// be used; there's no reason, since all methods are consistently
	// using it internally.
	// May be '/' (the zero AbsolutePath) if not applicable.
	BasePath() AbsolutePath

	OpenFile(path RelPath, flag int, perms Perms) (File, error)

	Mkdir(path RelPath, perms Perms) error

	Mklink(path RelPath, target string) error

	Mkfifo(path RelPath, perms Perms) error

	MkdevBlock(path RelPath, major int64, minor int64, perms Perms) error

	MkdevChar(path RelPath, major int64, minor int64, perms Perms) error

	Lchown(path RelPath, uid uint32, gid uint32) error

	Chmod(path RelPath, perms Perms) error

	SetTimesLNano(path RelPath, mtime time.Time, atime time.Time) error

	SetTimesNano(path RelPath, mtime time.Time, atime time.Time) error

	Stat(path RelPath) (*Metadata, error)

	LStat(path RelPath) (*Metadata, error)

	ReadDirNames(path RelPath) ([]string, error)

	Readlink(path RelPath) (target string, isSymlink bool, err error)

	/*
		Resolve a symlink (within the confines of the basepath!), returning
		the path to the final result.

		The resolution acts as if the FS basepath is the real root of the
		filesystem: rooted symlinks or excessive '..' segments still
		yield paths within the FS basepath.
		For example, `ResolveLink("../../..", "./lnk")` yields "./".

		This may opt to be bug-for-bug compatible with linux path resolution:
		specifically, for paths like "/untraversable/../target", where the first
		dir is 0000, we will return a permissions error... despite being able
		to clearly see that the following ".." will make the question irrelevant.
	*/
	ResolveLink(symlink string, startingAt RelPath) (RelPath, error)
}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt
}
