package tartrans

import (
	"archive/tar"
	"context"
	"crypto/sha512"
	"fmt"
	"io"
	"strings"

	"github.com/polydawn/refmt/misc"
	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/go-timeless-api/util"
	"go.polydawn.net/rio/config"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/fsOp"
	"go.polydawn.net/rio/lib/treewalk"
	"go.polydawn.net/rio/transmat/mixins/cache"
	"go.polydawn.net/rio/transmat/mixins/filters"
	"go.polydawn.net/rio/transmat/mixins/fshash"
	"go.polydawn.net/rio/transmat/mixins/log"
	"go.polydawn.net/rio/transmat/util"
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
	reader, err := PickReader(wareID, warehouses, false, mon)
	if err != nil {
		return api.WareID{}, err
	}
	defer reader.Close()

	// Construct filesystem wrapper to use for all our ops.
	afs := osfs.New(path2)

	// Extract.
	prefilterWareID, unpackWareID, err := unpackTar(ctx, afs, filt2, reader, mon)
	if err != nil {
		return unpackWareID, err
	}

	// Check for hash mismatch before returning, because that IS an error,
	//  but also return the hash we got either way.
	if prefilterWareID != wareID {
		return unpackWareID, ErrorDetailed(
			rio.ErrWareHashMismatch,
			fmt.Sprintf("hash mismatch: expected %q, got %q (filtered %q)", wareID, prefilterWareID, unpackWareID),
			map[string]string{
				"expected": wareID.String(),
				"actual":   prefilterWareID.String(),
				"filtered": unpackWareID.String(),
			},
		)
	}
	return unpackWareID, nil
}

func unpackTar(
	ctx context.Context,
	afs fs.FS,
	filt apiutil.FilesetFilters,
	reader io.Reader,
	mon rio.Monitor,
) (
	prefilterWareID api.WareID,
	actualWareID api.WareID,
	err error,
) {
	defer RequireErrorHasCategory(&err, rio.ErrorCategory(""))

	// Wrap input stream with decompression as necessary.
	//  Which kind of decompression to use can be autodetected by magic bytes.
	reader2, err := Decompress(reader)
	if err != nil {
		return api.WareID{}, api.WareID{}, Errorf(rio.ErrWareCorrupt, "corrupt tar compression: %s", err)
	}

	// Convert the raw byte reader to a tar stream.
	tr := tar.NewReader(reader2)

	// Allocate bucket for keeping each metadata entry and content hash;
	// the full tree hash will be computed from this at the end.
	// We keep one for the raw ware data as we consume it, so we can verify no fuckery;
	// we keep a second, separate one for the filtered data, which will compute a different hash.
	prefilterBucket := &fshash.MemoryBucket{}
	filteredBucket := &fshash.MemoryBucket{}

	// Also allocate a map for keeping records of which dirs we've created.
	// This is necessary for correct bookkeepping in the face of the tar format's
	// allowance for implicit parent dirs.
	dirs := map[fs.RelPath]struct{}{}

	// Iterate over each tar entry, mutating filesystem as we go.
	for {
		fmeta := fs.Metadata{}
		thdr, err := tr.Next()

		// Check for done.
		if err == io.EOF {
			break // sucess!  end of archive.
		}
		if err != nil {
			return api.WareID{}, api.WareID{}, Errorf(rio.ErrWareCorrupt, "corrupt tar: %s", err)
		}
		if ctx.Err() != nil {
			return api.WareID{}, api.WareID{}, Errorf(rio.ErrCancelled, "cancelled")
		}

		// Reshuffle metainfo to our default format.
		if err := TarHdrToMetadata(thdr, &fmeta); err != nil {
			return api.WareID{}, api.WareID{}, err
		}
		if strings.HasPrefix(fmeta.Name.String(), "..") {
			return api.WareID{}, api.WareID{}, Errorf(rio.ErrWareCorrupt, "corrupt tar: paths that use '../' to leave the base dir are invalid")
		}

		// Infer parents, if necessary.  The tar format allows implicit parent dirs.
		//
		// Note that if any of the implicitly conjured dirs is specified later, unpacking won't notice,
		// but bucket hashing iteration will (correctly) blow up for repeat entries.
		// It may well be possible to construct a tar like that, but it's already well established that
		// tars with repeated filenames are just asking for trouble and shall be rejected without
		// ceremony because they're just a ridiculous idea.
		for _, parent := range fmeta.Name.SplitParent() {
			// If we already initialized this parent, superb; move along.
			if _, exists := dirs[parent]; exists {
				continue
			}
			// If we're missing a dir, conjure a node with defaulted values.
			log.DirectoryInferred(mon, parent, fmeta.Name)
			conjuredFmeta := fshash.DefaultDirMetadata()
			conjuredFmeta.Name = parent
			prefilterBucket.AddRecord(conjuredFmeta, nil)
			filters.Apply(filt, &conjuredFmeta)
			filteredBucket.AddRecord(conjuredFmeta, nil)
			dirs[conjuredFmeta.Name] = struct{}{}
			if err := fsOp.PlaceFile(afs, conjuredFmeta, nil, filt.SkipChown); err != nil {
				return api.WareID{}, api.WareID{}, Errorf(rio.ErrInoperablePath, "error while unpacking: %s", err)
			}
		}

		// Apply filters.
		//  ... uck, to one copy of the meta.  We can't add either to their buckets
		//  until after the file is placed because we need the content hash.
		filteredFmeta := fmeta
		filters.Apply(filt, &filteredFmeta)

		// Place the file.
		switch fmeta.Type {
		case fs.Type_File:
			reader := &util.HashingReader{tr, sha512.New384()}
			if err := fsOp.PlaceFile(afs, filteredFmeta, reader, filt.SkipChown); err != nil {
				return api.WareID{}, api.WareID{}, Errorf(rio.ErrInoperablePath, "error while unpacking: %s", err)
			}
			prefilterBucket.AddRecord(fmeta, reader.Hasher.Sum(nil))
			filteredBucket.AddRecord(filteredFmeta, reader.Hasher.Sum(nil))
		case fs.Type_Dir:
			dirs[fmeta.Name] = struct{}{}
			fallthrough
		default:
			if err := fsOp.PlaceFile(afs, filteredFmeta, nil, filt.SkipChown); err != nil {
				return api.WareID{}, api.WareID{}, Errorf(rio.ErrInoperablePath, "error while unpacking: %s", err)
			}
			prefilterBucket.AddRecord(fmeta, nil)
			filteredBucket.AddRecord(filteredFmeta, nil)
		}
	}

	// Cleanup dir times with a post-order traversal over the bucket.
	//  Files and dirs placed inside dirs cause the parent's mtime to update, so we have to re-pave them.
	if err := treewalk.Walk(filteredBucket.Iterator(), nil, func(node treewalk.Node) error {
		record := node.(fshash.RecordIterator).Record()
		if record.Metadata.Type != fs.Type_Dir {
			return nil
		}
		return afs.SetTimesNano(record.Metadata.Name, record.Metadata.Mtime, fs.DefaultAtime)
	}); err != nil {
		return api.WareID{}, api.WareID{}, Errorf(rio.ErrInoperablePath, "error while unpacking: %s", err)
	}

	// Hash the thing!
	prefilterHash := misc.Base58Encode(fshash.HashBucket(prefilterBucket, sha512.New384))
	filteredHash := misc.Base58Encode(fshash.HashBucket(filteredBucket, sha512.New384))
	if !filt.IsHashAltering() {
		// Paranoia check for new feature.
		//  When paranoia reduced, replace with skipping the double computation.
		if prefilterHash != filteredHash {
			panic(fmt.Errorf("prefilterHash %q != filteredHash %q", prefilterHash, filteredHash))
		}
	}

	return api.WareID{"tar", prefilterHash}, api.WareID{"tar", filteredHash}, nil
}
