package filters

import (
	"os"

	"go.polydawn.net/go-timeless-api/util"
	"go.polydawn.net/rio/fs"
)

var (
	myUid = uint32(os.Getuid())
	myGid = uint32(os.Getgid())
)

/*
	Mutate the given fmeta handle to apply filters.

	Since this is the apiutil package's version of FilesetFilters,
	we can trust the values have been validated to reasonable ranges already,
	and defaults (for either pack or unpack mode) have already been mapped in.
*/
func Apply(filters apiutil.FilesetFilters, fmeta *fs.Metadata) {
	// Apply UID.
	switch filters.Uid {
	case apiutil.FilterKeep:
		// pass
	case apiutil.FilterMine:
		fmeta.Uid = myUid
	default:
		fmeta.Uid = uint32(filters.Uid)
	}

	// Apply GID.
	switch filters.Gid {
	case apiutil.FilterKeep:
		// pass
	case apiutil.FilterMine:
		fmeta.Gid = myGid
	default:
		fmeta.Gid = uint32(filters.Gid)
	}

	// Apply Mtime.
	if filters.Mtime != nil {
		fmeta.Mtime = *filters.Mtime
	}

	// Apply Sticky.
	if !filters.Sticky {
		fmeta.Perms &= 0777
	}
}
