package cache

import (
	"context"
	"os"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/go-timeless-api/util"
	cacheapi "go.polydawn.net/rio/cache"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/fsOp"
	"go.polydawn.net/rio/lib/guid"
	"go.polydawn.net/rio/stitch/placer"
	"go.polydawn.net/rio/transmat/mixins/log"
)

var ShelfFor = cacheapi.ShelfFor

func Lrn2Cache(cacheFs fs.FS, unpackTool rio.UnpackFunc) rio.UnpackFunc {
	return cache{cacheFs, unpackTool}.Unpack
}

type cache struct {
	fs         fs.FS
	unpackTool rio.UnpackFunc
}

/*
	Proxies most args to the cache's unpack tool, except for placementmode and path,
	which it sets to rio.Placement_Direct and a temporary path in the cache filesystem.
	If unpacking completes successfully, the temp path will be moved to a permanent
	location in the cache, which is specified by the public interface `rio/cache.GetShelf`.

	Any behaviors specified by the placementMode -- copying, mounting, etc -- are enacted
	by this func after the unpack finishes and the temp path committed to the cache.
*/
func (c cache) Unpack(
	ctx context.Context,
	wareID api.WareID,
	path string,
	filt api.FilesetFilters,
	placementMode rio.PlacementMode,
	warehouses []api.WarehouseAddr,
	monitor rio.Monitor,
) (_ api.WareID, err error) {
	defer RequireErrorHasCategory(&err, rio.ErrorCategory(""))

	// Zeroth thing: caches are by hash, but remember that filters can give you a
	//  result hash which is different than the requested ware hash.
	//  Right now we deal with this simply/stupidly: if you used filters, no cache for you.
	resultWareID := wareID
	filt2, err := apiutil.ProcessFilters(filt, apiutil.FilterPurposeUnpack)
	if err != nil {
		return api.WareID{}, Errorf(rio.ErrUsage, "invalid filter specification: %s", err)
	}
	if filt2.IsHashAltering() {
		resultWareID = api.WareID{"-", "-"} // This value forces cache miss.
	}

	// First thing: Check if we already have the ware in cache and can jump to placement ASAP.
	//  (This must be first because we're willing to read cache even in "direct" mode, but
	//  yet *not* willing to even initialize empty cache dirs in that mode.)
	shelf := ShelfFor(resultWareID)
	_, err = c.fs.Stat(shelf)
	switch Category(err) {
	case fs.ErrNotExists: // "not exists" is just a cache miss...
		switch placementMode {
		case rio.Placement_Direct: // In direct mode: be direct.  Do nothing to cache.
			return c.unpackTool(ctx, wareID, path, filt, rio.Placement_Direct, warehouses, monitor)
		default: // Everyone else: unpack into cache.
			// pass
		}
		// Unpack into the cache.
		resultWareID, shelf, err = c.populate(ctx, wareID, filt, warehouses, monitor)
		if err != nil {
			return resultWareID, err
		}
		// Now place it from the cache shelf.
		return resultWareID, c.place(ctx, placementMode, shelf, path)
	case nil: // Cache has it!  Reaction varies.
		log.CacheHasIt(monitor, wareID)
		return resultWareID, c.place(ctx, placementMode, shelf, path)
	default:
		// Unknown errors reading cache are mostly considered game over.  Except:
		//  Since direct mode has no responsibility to the cache, it can still go.
		switch placementMode {
		case rio.Placement_Direct:
			return c.unpackTool(ctx, wareID, path, filt, rio.Placement_Direct, warehouses, monitor)
		default:
			return api.WareID{}, Errorf(rio.ErrLocalCacheProblem, "error reading cache: %s", err)
		}
	}
}

func (c cache) place(
	ctx context.Context,
	placementMode rio.PlacementMode,
	shelf fs.RelPath,
	destination string, // still a string at this phase because it's either abs or "-"
) error {
	absShelf := c.fs.BasePath().Join(shelf)
	switch placementMode {
	case rio.Placement_None: // If no placement, cache having it is victory!
		return nil
	case rio.Placement_Direct: // In direct mode, copy.
		_, err := placer.CopyPlacer(absShelf, fs.MustAbsolutePath(destination), true)
		return err
	case rio.Placement_Copy: // In copy mode, ... well obviously copy.
		_, err := placer.CopyPlacer(absShelf, fs.MustAbsolutePath(destination), true)
		return err
	case rio.Placement_Mount: // In mount mode, mount.
		placerFn, err := placer.GetMountPlacer()
		if err != nil {
			return err
		}
		_, err = placerFn(absShelf, fs.MustAbsolutePath(destination), true)
		return err
	default:
		panic("unreachable")
	}
}

func (c cache) populate(
	ctx context.Context,
	wareID api.WareID,
	filt api.FilesetFilters,
	warehouses []api.WarehouseAddr,
	monitor rio.Monitor,
) (_ api.WareID, _ fs.RelPath, err error) {
	defer RequireErrorHasCategory(&err, rio.ErrorCategory(""))

	// Initialize cache.
	//  Ensure the cache commit root dir exists.
	//  Also ensure the cache parent dir exists... no bound on recursion.
	if err := fsOp.MkdirAll(osfs.New(fs.AbsolutePath{}), c.fs.BasePath().CoerceRelative(), 0700); err != nil {
		return api.WareID{}, fs.RelPath{}, Errorf(rio.ErrLocalCacheProblem, "cannot initialize cache dirs: %s", err)
	}
	if err := fsOp.MkdirAll(c.fs, fs.MustRelPath(string(wareID.Type)+"/fileset"), 0700); err != nil {
		return api.WareID{}, fs.RelPath{}, Errorf(rio.ErrLocalCacheProblem, "cannot initialize cache dirs: %s", err)
	}

	// Pick a temp path to unpack into.
	tmpPath := fs.MustRelPath("./.tmp.unpack." + guid.New())
	tmpPathStr := c.fs.BasePath().Join(tmpPath).String()
	// Defer cleanup of the temp path.
	//  (If we're successful, we'll have moved it out of this path before return.)
	defer os.RemoveAll(tmpPathStr)
	// Delegate!
	resultWareID, err := c.unpackTool(ctx, wareID, tmpPathStr, filt, rio.Placement_Direct, warehouses, monitor)
	if err != nil {
		return resultWareID, fs.RelPath{}, err
	}

	// Successful unpack: commit it to its shelf location.
	//  This may also require mkdir'ing the prefix dirs of the shelf.
	//  In case of race: accept our fate, assume the racing party acted in good faith,
	//  return the shelf path anyway, and our defer'd rm will act on our wasted copy.
	shelf := ShelfFor(resultWareID)
	c.fs.Mkdir(shelf.Dir().Dir(), 0755)
	c.fs.Mkdir(shelf.Dir(), 0755)
	if err := os.Rename(tmpPathStr, c.fs.BasePath().Join(shelf).String()); err != nil {
		if _, ok := err.(*os.LinkError); ok && os.IsExist(err) {
			// Oh, fine.  Somebody raced us to it.
			return resultWareID, shelf, nil
		}
		// Any other error: sad.
		return resultWareID, shelf, Errorf(rio.ErrLocalCacheProblem, "error commiting %q into cache: %s", resultWareID, err)
	}
	return resultWareID, shelf, nil
}
