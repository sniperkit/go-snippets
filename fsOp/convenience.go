package fsOp

import (
	"os"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/rio/fs"
)

/*
	Makes dirs recursively so the requested path exists, applying the assigned metadata
	to each one that needed to be produced.

	Existing dirs are not mutated.

	Symlinks will be traversed without comment (i.e. this will never emit ErrBreakout).
	(Note that this means this function is *not* used in either transmats nor stitch.)
*/
func MkdirAll(afs fs.FS, path fs.RelPath, perms fs.Perms) error {
	// Check if the path already exists.
	stat, err := afs.Stat(path)
	// Switch on status of the (derefenced) file.
	//  Recurse and mkdir if necessary.
	switch Category(err) {
	case nil:
		if stat.Type == fs.Type_Dir {
			return nil
		}
		return Errorf(fs.ErrNotDir, "%s already exists and is a %s not %s", afs.BasePath().Join(path), stat.Type, fs.Type_Dir)
	case fs.ErrNotExists:
		if path == (fs.RelPath{}) {
			return Errorf(fs.ErrNotExists, "base path %s does not exist!", afs.BasePath())
		}
		if err := MkdirAll(afs, path.Dir(), perms); err != nil {
			return err
		}
		if err := afs.Mkdir(path, perms); err != nil {
			// Re-check if it exists and is a dir; we may have been raced, and we're okay with that.
			stat, err1 := afs.LStat(path)
			if err1 == nil && stat.Type == fs.Type_Dir {
				return nil
			}
			switch Category(err) {
			case fs.ErrAlreadyExists:
				// this seemingly-contradictory message means the path does exist... it's just that stat said it didn't, beacuse it's a dangling symlink.
				return Errorf(fs.ErrNotDir, "%s already exists and is a %s not %s", afs.BasePath().Join(path), fs.Type_Symlink, fs.Type_Dir)
			default:
				return err
			}
		}
		return nil
	case fs.ErrNotDir:
		// Reformat the error a tad to not say "lstat", which is distracting.
		return Errorf(fs.ErrNotDir, "%s has parents which are not a directory", afs.BasePath().Join(path))
	default:
		return err
	}
}

/*
	A form of MkdirAll, recursively creating parent directories as necessary,
	but with more behaviors: the final directory will have all its metadata
	and permissions forced to the given uid+gid and its permissions bitwise
	`|0700` (owner rw), and all parents bitwise `|0001` (everyone-traversable),
	and all affected dirs and their parents will have their mtimes repaired.

	The preferredProps metadata is only partially followed, as the rules
	above take precidence; but the preferredProps will be used for any
	dirs that need be created.

	If intermediate path segements are not dirs, errors will be returned.

	Long story short, it makes sure the given path is read-write *usable* to
	the given uid+gid.
	(Repeatr leans on this in the "cradle" functionality.)
*/
func MkdirUsable(afs fs.FS, path fs.RelPath, preferredProps fs.Metadata) error {
	preferredProps.Type = fs.Type_Dir
	for _, segment := range path.Split() {
		defer RepairMtime(afs, segment.Dir())()
		mergeBits := fs.Perms(01)
		if segment == path {
			mergeBits = 0700
		}
		stat, err := afs.LStat(segment)
		switch Category(err) {
		case nil: // already exists
			if stat.Type == fs.Type_Dir {
				newMode := stat.Perms | mergeBits
				if newMode != stat.Perms {
					if err := afs.Chmod(segment, newMode); err != nil {
						return err
					}
				}
			} else {
				return Errorf(fs.ErrNotDir, "%s already exists and is a %s not %s", afs.BasePath().Join(path), stat.Type, fs.Type_Dir)
			}
		case fs.ErrNotExists:
			preferredProps.Name = segment
			if err := PlaceFile(afs, preferredProps, nil, false); err != nil {
				return err
			}
		default:
			return err
		}
	}
	defer RepairMtime(afs, path)()
	if err := afs.Lchown(path, preferredProps.Uid, preferredProps.Gid); err != nil {
		return err
	}
	return nil
}

/*
	Records the mtime property currently set on a path and returns a function which
	will force it to that value again.

	The typical use for this is `defer RepairMtime(fs, "./somedir")()` right
	before invoking some funcs which will mutate the contents of that dir.
	When the rest of the changes are done and the function returns, the mtime
	of the dir will be forced back to its previous value, seemingly unchanged.
*/
func RepairMtime(afs fs.FS, path fs.RelPath) func() {
	fmeta, err := afs.LStat(path)
	if err != nil {
		return func() {}
	}
	return func() {
		afs.SetTimesLNano(path, fmeta.Mtime, fs.DefaultAtime)
	}
}

/*
	Remove all files and dirs in a dir (if it exists; if not, no-op), recursing
	as necessary.

	This can be useful to do before e.g. some 'unpack' operation, to make sure
	no content collides, while also avoiding unlinking the top level dir, in case
	the current process would lack permission to create it again.

	As an additional edge case, if the given path is a file and not a dir, the file
	will also be removed (so, the consistent logic here is you will be ready to have
	a dir in this path (aka, target an unpack here) when this function returns).
*/
func RemoveDirContent(afs fs.FS, path fs.RelPath) error {
	// Lazy implementation... assumes osfs, and slinks back out to stdlib, because
	// it so happens all of our real usage is fine with that.
	children, err := afs.ReadDirNames(path)
	switch Category(err) {
	case nil:
		// pass
	case fs.ErrNotExists:
		return nil // great
	default:
		return err
	}
	for _, child := range children {
		err := os.RemoveAll(afs.BasePath().Join(path).String() + "/" + child)
		if err != nil {
			return fs.NormalizeIOError(err)
		}
	}
	return nil
}
