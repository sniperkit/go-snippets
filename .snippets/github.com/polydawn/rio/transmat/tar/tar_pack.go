package tartrans

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha512"
	"io"
	"time"

	"github.com/polydawn/refmt/misc"

	. "github.com/warpfork/go-errcat"
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/go-timeless-api/util"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/fsOp"
	"go.polydawn.net/rio/transmat/mixins/filters"
	"go.polydawn.net/rio/transmat/mixins/fshash"
)

var (
	_ rio.PackFunc = Pack
)

func Pack(
	ctx context.Context, // Long-running call.  Cancellable.
	packType api.PackType, // The name of pack format.
	pathStr string, // The fileset to scan and pack (absolute path).
	filt api.FilesetFilters, // Optionally: filters we should apply while unpacking.
	warehouseAddr api.WarehouseAddr, // Warehouse to save into (or blank to just scan).
	mon rio.Monitor, // Optionally: callbacks for progress monitoring.
) (_ api.WareID, err error) {
	if mon.Chan != nil {
		defer close(mon.Chan)
	}
	defer RequireErrorHasCategory(&err, rio.ErrorCategory(""))

	// Sanitize arguments.
	if packType != PackType {
		return api.WareID{}, Errorf(rio.ErrUsage, "this transmat implementation only supports packtype %q (not %q)", PackType, packType)
	}
	path, err := fs.ParseAbsolutePath(pathStr)
	if err != nil {
		return api.WareID{}, Errorf(rio.ErrUsage, "pack must be called with absolute path: %s", err)
	}
	filt2, err := apiutil.ProcessFilters(filt, apiutil.FilterPurposePack)
	if err != nil {
		return api.WareID{}, Errorf(rio.ErrUsage, "invalid filter specification: %s", err)
	}

	// Short-circuit exit if the path does not exist.
	//  We could let the errors later bubble, but, why bother opening a writeController,
	//  etc, if we're just going to have to rm the resource a millisecond later?
	afs := osfs.New(path)
	_, err = afs.Stat(fs.RelPath{})
	switch Category(err) {
	case nil:
		// pass
	case fs.ErrNotExists:
		return api.WareID{PackType, ""}, nil
	default:
		return api.WareID{}, Errorf(rio.ErrPackInvalid, "cannot read path for packing: %s", err)
	}

	// Connect to warehouse, and get write controller opened.
	wc, err := OpenWriteController(warehouseAddr, packType, mon)
	if err != nil {
		return api.WareID{}, err
	}
	defer wc.Close()

	// Wrap writer stream to do compress on the way out.
	//  Note on compression levels: The default is 6; and per http://tukaani.org/lzma/benchmarks.html
	//  this appears quite reasonable: higher levels appear to have minimal size payoffs, but significantly rising compress time costs;
	//  decompression time does not vary with compression level.
	// Save a gzip reference just to close it; tar.Writer doesn't passthru its own close.
	gzWriter := gzip.NewWriter(wc)

	// Construct tar writer.
	tarWriter := tar.NewWriter(gzWriter)

	// Scan and tarify!
	wareID, err := packTar(ctx, afs, filt2, tarWriter)
	if err != nil {
		return wareID, err
	}
	// Close all the intermediate writer layers to ensure they've flushed.
	tarWriter.Close()
	gzWriter.Close()

	// If we made it all the way with no errors, commit.
	//  (Otherwise, the write controller will be closed by default by our defers.)
	return wareID, wc.Commit(wareID)
}

func packTar(
	ctx context.Context,
	afs fs.FS,
	filt apiutil.FilesetFilters,
	tw *tar.Writer,
) (api.WareID, error) {
	// Allocate bucket for keeping each metadata entry and content hash;
	// the full tree hash will be computed from this at the end.
	bucket := &fshash.MemoryBucket{}

	// Walk the filesystem, emitting tar entries and filling the bucket as we go.
	tarHeader := &tar.Header{}
	preVisit := func(filenode *fs.FilewalkNode) error {
		if filenode.Err != nil {
			return filenode.Err
		}

		// Consider cancellation.
		if ctx.Err() != nil {
			return Errorf(rio.ErrCancelled, "cancelled")
		}

		// Open file.
		fmeta, file, err := fsOp.ScanFile(afs, filenode.Info.Name) // FIXME : we already have the full metadata loaded; give ScanFile option to accept it!
		if err != nil {
			return err
		}

		// Apply filters.
		filters.Apply(filt, fmeta)

		// Flatten time to seconds.  The tar writer impl doesn't do subsecond precision.
		//  The writer will always flatten it internally, but we need to do it here as well
		//  so that the hash and the serial form are describing the same thing.
		fmeta.Mtime = fmeta.Mtime.Truncate(time.Second)

		// Flip our metadata to tar header format, and flush it.
		MetadataToTarHdr(fmeta, tarHeader)
		if err := tw.WriteHeader(tarHeader); err != nil {
			return Errorf(rio.ErrWarehouseUnwritable, "error while writing pack: %s", err)
		}

		// If it's a file, stream the body into the tar while hashing; for all,
		//  record the metadata in the bucket for the total hash.
		if file == nil {
			bucket.AddRecord(*fmeta, nil)
		} else {
			defer file.Close()
			hasher := sha512.New384()
			tee := io.MultiWriter(tw, hasher)
			_, err := io.Copy(tee, file)
			if err != nil {
				return err
			}
			bucket.AddRecord(*fmeta, hasher.Sum(nil))
		}
		return nil
	}
	if err := fs.Walk(afs, preVisit, nil); err != nil {
		return api.WareID{}, err
	}

	// Hash the thing!
	hash := fshash.HashBucket(bucket, sha512.New384)
	return api.WareID{"tar", misc.Base58Encode(hash)}, nil
}
