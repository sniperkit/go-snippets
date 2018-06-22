package placer

import (
	"os"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/fsOp"
)

var rootFs = osfs.New(fs.MustAbsolutePath("/")) // handy, since placers are always absolutized

func mkDest(dstPath fs.AbsolutePath, wantType fs.Type) error {
	// Handle all the cases for existing things at destination.
	dstStat, err := rootFs.LStat(dstPath.CoerceRelative())
	switch Category(err) {
	case nil: // It exists.  But is it the right type?
		if dstStat.Type == wantType {
			return nil // Already matches.  Huzzah.
		}
		return Errorf(rio.ErrAssemblyInvalid, "placer: destination already exists and is different type than source")
	case fs.ErrNotExists:
		// Carry on.  We'll create it.
	default: // Any other error: raise.
		return Errorf(rio.ErrAssemblyInvalid, "placer: destination unusable: %s", err)
	}

	// If we made it this far: dest doesn't exist yet.
	// Capture the parent dir mtime and defer its repair, because we're about to disrupt it.
	defer fsOp.RepairMtime(rootFs, dstPath.Dir().CoerceRelative())()

	// Make the dest node, matching type of the source.
	// The perms don't matter; will be shadowed.
	// We assume the parent dirs are all in place because you're almost
	// certainly using this as part of an assembler.
	switch wantType {
	case fs.Type_Symlink:
		fallthrough
	case fs.Type_NamedPipe:
		fallthrough
	case fs.Type_Socket:
		fallthrough
	case fs.Type_Device:
		fallthrough
	case fs.Type_CharDevice:
		fallthrough
	case fs.Type_File:
		var f *os.File
		f, err = os.OpenFile(dstPath.String(), os.O_CREATE, 0644)
		f.Close()
	case fs.Type_Dir:
		err = os.Mkdir(dstPath.String(), 0644)
	}
	if err != nil {
		return Errorf(rio.ErrAssemblyInvalid, "placer: destination unusable: %s", err)
	}
	return nil
}
