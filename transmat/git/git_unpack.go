package git

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	. "github.com/warpfork/go-errcat"
	"gopkg.in/src-d/go-git.v4/plumbing/filemode"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/go-timeless-api/util"
	"go.polydawn.net/rio/config"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/fsOp"
	"go.polydawn.net/rio/transmat/mixins/cache"
	"go.polydawn.net/rio/transmat/mixins/filters"
	"go.polydawn.net/rio/transmat/mixins/fshash"
	gitWarehouse "go.polydawn.net/rio/warehouse/impl/git"
)

var (
	_ rio.UnpackFunc = Unpack
)

func Unpack(
	ctx context.Context, // Long-running call.  Cancellable.
	wareID api.WareID, // What wareID to fetch for unpacking.
	path string, // Where to unpack the fileset (absolute path).
	filt api.FilesetFilters, // Optionally: filters we should apply while unpacking.
	placementMode rio.PlacementMode, // Optionally: a placement mode (default is "copy").
	warehouses []api.WarehouseAddr, // Warehouses we can try to fetch from.
	mon rio.Monitor, // Optionally: callbacks for progress monitoring.
) (_ api.WareID, err error) {
	if mon.Chan != nil {
		defer close(mon.Chan)
	}
	defer RequireErrorHasCategory(&err, rio.ErrorCategory(""))

	// Sanitize arguments.
	if wareID.Type != PackType {
		return api.WareID{}, Errorf(rio.ErrUsage, "this transmat implementation only supports packtype %q (not %q)", PackType, wareID.Type)
	}
	if placementMode == "" {
		placementMode = rio.Placement_Copy
	}
	// Wrap the direct unpack func with cache behavior; call that.
	return cache.Lrn2Cache(
		osfs.New(config.GetCacheBasePath()),
		unpack,
	)(ctx, wareID, path, filt, placementMode, warehouses, mon)
}

func unpack(
	ctx context.Context,
	wareID api.WareID,
	path string,
	filt api.FilesetFilters,
	placementMode rio.PlacementMode,
	warehouses []api.WarehouseAddr,
	mon rio.Monitor,
) (_ api.WareID, err error) {
	defer RequireErrorHasCategory(&err, rio.ErrorCategory(""))

	// Sanitize arguments.
	path2 := fs.MustAbsolutePath(path)
	filt2, err := apiutil.ProcessFilters(filt, apiutil.FilterPurposeUnpack)
	if err != nil {
		return api.WareID{}, Errorf(rio.ErrUsage, "invalid filter specification: %s", err)
	}

	// Pick a warehouse and get a reader.
	//  This is a *very* expensive operation for git.  It's less
	//  of "pick a warehouse" and more "download the whole thing and hope we
	//  get what we wanted" (which is very ironic for a system that has
	//  a CAS system on its inside, yes).
	whCtrl, err := pick(ctx,
		wareID,
		warehouses,
		osfs.New(config.GetCacheBasePath().Join(fs.MustRelPath("git/objs"))),
		mon,
	)
	if err != nil {
		return api.WareID{}, err
	}

	// Get submodule config.  Fetch them all.
	submodules, err := whCtrl.Submodules(wareID.Hash)
	if err != nil {
		return api.WareID{}, err
	}
	// We'll organize them by path now; only thing that's useful.
	submoduleCtrls := map[string]*gitWarehouse.Controller{}
	for _, submCfg := range submodules {
		// TODO it would be dreamy to parallelize this.
		whCtrl, err := pick(ctx,
			api.WareID{"git", submCfg.Hash},
			[]api.WarehouseAddr{api.WarehouseAddr(submCfg.URL)},
			osfs.New(config.GetCacheBasePath().Join(fs.MustRelPath("git/objs"))),
			mon,
		)
		if err != nil {
			return api.WareID{}, err
		}
		submoduleCtrls[submCfg.Path] = whCtrl
	}

	// Open a tree to walk in the main repo.
	//  We'll do submodule checkouts somewhere deep in the middle of this.
	tr, err := whCtrl.GetTree(wareID.Hash)
	if err != nil {
		panic(err)
	}

	// Construct filesystem wrapper to use for all our ops.
	afs := osfs.New(path2)

	// Walk.
	if err := unpackOneRepo(ctx, tr, afs, true, filt2, submoduleCtrls, mon); err != nil {
		return api.WareID{}, err
	}

	// That's it.  Checkout should have already checked the hash, so we just return it.
	return wareID, nil
}

func unpackOneRepo(
	ctx context.Context,
	tr *object.Tree,
	afs fs.FS,
	isRoot bool, // if true, will recurse for submodules (with this set to false).
	filt apiutil.FilesetFilters,
	submoduleCtrls map[string]*gitWarehouse.Controller,
	mon rio.Monitor,
) (err error) {
	tw := object.NewTreeWalker(tr, true, nil)

	// Make the root dir.  Git doesn't have metadata for the tree root.
	conjuredFmeta := fshash.DefaultDirMetadata()
	filters.Apply(filt, &conjuredFmeta)
	if err := fsOp.PlaceFile(afs, conjuredFmeta, nil, filt.SkipChown); err != nil {
		return Errorf(rio.ErrInoperablePath, "error while unpacking: %s", err)
	}

	// Extract.
	// Iterate over each entry, mutating filesystem as we go.
	dirs := make([]fs.RelPath, 1, 200) // Keep for dir time repair at end.
	dirs[0] = fs.RelPath{}
	for {
		fmeta := fs.Metadata{}
		name, te, err := tw.Next()

		// Check for done.
		if err == io.EOF {
			break // sucess!  end of archive.
		}
		if err != nil {
			return Errorf(rio.ErrWareCorrupt, "corrupt git tree: %s", err)
		}
		if ctx.Err() != nil {
			return Errorf(rio.ErrCancelled, "cancelled")
		}
		//fmt.Fprintf(os.Stderr, "walking git tree %s -- %#v\n", name, te)

		// Reshuffle metainfo to our default format.
		fmeta.Name = fs.MustRelPath(name)
		fmeta.Uid = 1000
		fmeta.Gid = 1000
		switch te.Mode {
		case filemode.Dir:
			fmeta.Type = fs.Type_Dir
			fmeta.Perms = 0755
			dirs = append(dirs, fmeta.Name)
		case filemode.Regular:
			fmeta.Type = fs.Type_File
			fmeta.Perms = 0644
		case filemode.Executable:
			fmeta.Type = fs.Type_File
			fmeta.Perms = 0755
		case filemode.Symlink:
			fmeta.Type = fs.Type_Symlink
			fmeta.Perms = 0644
			// Hang on, extracting a symlink is actually rough.
			tf, err := tr.TreeEntryFile(&te)
			if err != nil {
				return Errorf(rio.ErrWareCorrupt, "corrupt git tree: %s", err)
			}
			reader, err := tf.Blob.Reader()
			if err != nil {
				return Errorf(rio.ErrWareCorrupt, "corrupt git tree: %s", err)
			}
			blob, err := ioutil.ReadAll(reader)
			if err != nil {
				return Errorf(rio.ErrWareCorrupt, "corrupt git tree: %s", err)
			}
			fmeta.Linkname = string(blob)
		case filemode.Submodule:
			// Ooowee!  Recurse time!
			if !isRoot {
				// Except of course if we're already a submodule, in which case no.
				// Like git, we will make the empty dir, though.
				fmeta.Type = fs.Type_Dir
				fmeta.Perms = 0755
				dirs = append(dirs, fmeta.Name)
				break
			}
			submCtrl, ok := submoduleCtrls[name]
			if !ok {
				return Errorf(rio.ErrWareCorrupt, "gitlink found at path %q but no matching config in .gitmodules", name)
			}
			submTr, err := submCtrl.GetTree(te.Hash.String())
			if err != nil {
				panic(err)
			}
			submFs := osfs.New(afs.BasePath().Join(fmeta.Name))
			if err := unpackOneRepo(ctx, submTr, submFs, false, filt, nil, mon); err != nil {
				return err
			}
			continue
		case filemode.Empty:
			fallthrough
		case filemode.Deprecated:
			fallthrough
		default:
			panic(fmt.Errorf("unknown git filemode %#v", te.Mode))
		}
		fmeta.Mtime = apiutil.DefaultMtime // git doesn't have time info

		// Apply filters.
		filters.Apply(filt, &fmeta)

		// Place the file.
		switch fmeta.Type {
		case fs.Type_File:
			tf, err := tr.TreeEntryFile(&te)
			if err != nil {
				return Errorf(rio.ErrWareCorrupt, "corrupt git tree: %s", err)
			}
			reader, err := tf.Blob.Reader()
			if err != nil {
				return Errorf(rio.ErrWareCorrupt, "corrupt git tree: %s", err)
			}
			if err := fsOp.PlaceFile(afs, fmeta, reader, filt.SkipChown); err != nil {
				return Errorf(rio.ErrInoperablePath, "error while unpacking: %s", err)
			}
			reader.Close()
		default:
			if err := fsOp.PlaceFile(afs, fmeta, nil, filt.SkipChown); err != nil {
				return Errorf(rio.ErrInoperablePath, "error while unpacking: %s", err)
			}
		}
	}

	// Cleanup dir times with a post-order traversal over the bucket.
	//  Files and dirs placed inside dirs cause the parent's mtime to update, so we have to re-pave them.
	for i := len(dirs) - 1; i >= 0; i-- {
		if err := afs.SetTimesNano(dirs[i], fs.DefaultAtime, fs.DefaultAtime); err != nil {
			return Errorf(rio.ErrInoperablePath, "error while unpacking: %s", err)
		}
	}

	return nil
}
