package tartrans

import (
	"context"
	"fmt"
	"io"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/go-timeless-api/util"
	"go.polydawn.net/rio/fs/nilfs"
	"go.polydawn.net/rio/transmat/mixins/log"
)

var (
	_ rio.MirrorFunc = Mirror
)

func Mirror(
	ctx context.Context, // Long-running call.  Cancellable.
	wareID api.WareID, // What wareID to mirror.
	target api.WarehouseAddr, // Warehouse to ensure the ware is mirrored into.
	sources []api.WarehouseAddr, // Warehouses we can try to fetch from.
	mon rio.Monitor, // Optionally: callbacks for progress monitoring.
) (_ api.WareID, err error) {
	defer RequireErrorHasCategory(&err, rio.ErrorCategory(""))
	if mon.Chan != nil {
		defer close(mon.Chan)
	}

	// Try to read the ware from the target first; if successfull, no-op out.
	//  We don't fully re-verify the content, because that requires a time
	//  committment, and we want this command to be fast when run repeatedly.
	reader, err := PickReader(wareID, []api.WarehouseAddr{target}, false, mon)
	if err == nil {
		log.MirrorNoop(mon, target, wareID)
		reader.Close()
		return wareID, nil
	}

	// Connect to target warehouse, and get write controller opened.
	//  During mirroring, unlike unpacking, we actually *do* know the hash
	//  of what we'll be uploading... but there's nothing dramatically better
	//  we can do with that knowledge.
	wc, err := OpenWriteController(target, wareID.Type, mon)
	if err != nil {
		return api.WareID{}, err
	}
	defer wc.Close()

	// Pick a source warehouse and get a reader.
	reader, err = PickReader(wareID, sources, false, mon)
	if err != nil {
		return api.WareID{}, err
	}
	defer reader.Close()

	// Prepare to scan this as we process.
	//  It would be unfortunate to accidentally foist corrupted or
	//  wrongly identified content onto a mirror.
	reader = flippingReader{reader, wc}
	afs := nilFS.New()

	// "unpack", scanningly.  This drives the copy.
	filt, _ := apiutil.ProcessFilters(api.Filter_NoMutation, apiutil.FilterPurposeUnpack)
	// We can ignore the pre/post filter wareIDs, since we know its a no-mutation filter.
	gotWare, _, err := unpackTar(ctx, afs, filt, reader, mon)
	if err != nil {
		// If errors at this stage: still return a blank wareID, because
		//  we haven't finished *uploading* it.
		return api.WareID{}, err
	}

	// Check for hash mismatch; abort if detected.
	//  Again, return a blank wareID, because we haven't *uploaded* it.
	//  You'll have no choice but to inspect the error details if you need the value.
	if gotWare != wareID {
		return api.WareID{}, ErrorDetailed(
			rio.ErrWareHashMismatch,
			fmt.Sprintf("hash mismatch: expected %q, got %q", wareID, gotWare),
			map[string]string{
				"expected": wareID.String(),
				"actual":   gotWare.String(),
			},
		)
	}

	// All's quiet: flush and commit.
	return gotWare, wc.Commit(wareID)
}

// Proxy read calls, also copying each buffer into another write.
type flippingReader struct {
	read io.ReadCloser
	dup  io.Writer
}

func (fr flippingReader) Read(b []byte) (int, error) {
	n, err := fr.read.Read(b)
	if err == nil || err == io.EOF {
		n2, err2 := fr.dup.Write(b[:n])
		if n2 < n {
			return n, io.ErrShortWrite
		}
		if err2 != nil {
			return n, err2
		}
	}
	return n, err
}

func (fr flippingReader) Close() error {
	return fr.read.Close()
}
