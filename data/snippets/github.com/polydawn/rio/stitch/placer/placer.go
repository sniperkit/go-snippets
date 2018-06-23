package placer

import (
	"go.polydawn.net/rio/fs"
)

type Placer func(srcPath, dstPath fs.AbsolutePath, writable bool) (Janitor, error)

type Janitor interface {
	// Describe, in a shell-like way, what the teardown would do.
	// (e.g. 'rm -rf' or 'umount' plus the absolute path.)
	Description() string

	// Do the teardown.
	Teardown() error

	// Whether or not to always attempt the teardown, even when other teardowns
	// in a group have errored.
	// (This is *false* for e.g. copyPlacer, because recursive removes if
	// an unmount somewhere failed are *extremely* dangerous.)
	AlwaysTry() bool
}
