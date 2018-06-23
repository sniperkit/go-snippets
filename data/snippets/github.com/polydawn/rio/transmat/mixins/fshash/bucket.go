package fshash

import (
	"fmt"

	"go.polydawn.net/go-timeless-api/util"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/lib/treewalk"
)

/*
	Bucket keeps hashes of file content and the set of metadata per file and dir.
	This is to make it possible to range over the filesystem out of order and
	construct a total hash of the system in order later.
*/
type Bucket interface {
	AddRecord(metadata fs.Metadata, contentHash []byte) // record a file into the bucket
	Iterator() (rootRecord RecordIterator)              // return a treewalk root that does a traversal ordered by path
	Root() Record                                       // return the 0'th record; is the root if `Iterator` has already been invoked.
	Length() int
}

type Record struct {
	// Name field is *almost* identical to Metadata.Name.String(), but includes trailing slashes for dirs.
	Name        string
	Metadata    fs.Metadata
	ContentHash []byte
}

/*
	Error raised while the bucket is being iterated and we discover it has either
	a duplicate entry or a missing parent dir for some entry.

	If you're getting this after walking a filesystem, you had a bug in your walk;
	if you're getting it after walking a tar or some format, it may be missing entries as well.

	This is not meant to be user facing; map it to something meaningful in your caller.
*/
type ErrInvalidFilesystem struct {
	Msg string
}

func (e ErrInvalidFilesystem) Error() string {
	return fmt.Sprintf("invariant broken while traversing fshash bucket: %s", e.Msg)
}

// for sorting
type recordsByFilename []Record

func (a recordsByFilename) Len() int           { return len(a) }
func (a recordsByFilename) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a recordsByFilename) Less(i, j int) bool { return a[i].Name < a[j].Name }

/*
	RecordIterator is used for walking Bucket contents in hash-ready order.

	May panic with ErrInvalidFilesystem during traversal.
	TODO : reconsider if maybe you wouldn't like a "sancheck" method that does that scan first?
*/
type RecordIterator interface {
	treewalk.Node
	Record() Record
}

/*
	Returns a new, consistent, "blank" metadata for a directory.
	You must assign the `Name` metadata.
*/
func DefaultDirMetadata() fs.Metadata {
	return fs.Metadata{
		Type:  fs.Type_Dir,
		Perms: 0755,
		Mtime: apiutil.DefaultMtime,
	}
}
