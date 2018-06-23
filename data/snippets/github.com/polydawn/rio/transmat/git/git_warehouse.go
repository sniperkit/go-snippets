package git

import (
	"context"
	"fmt"
	"net/url"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/transmat/mixins/log"
	gitWarehouse "go.polydawn.net/rio/warehouse/impl/git"
)

// Pick a warehouse.
//
// This function is not cheap; it actually does a ton of work,
// including fetching full content and leaving a lot of it on disk.
// We're not thrilled with this, but it's the way git works.
func pick(
	ctx context.Context,
	wareID api.WareID,
	warehouses []api.WarehouseAddr,
	objcacheWorkdir fs.FS,
	mon rio.Monitor,
) (whCtrl *gitWarehouse.Controller, err error) {
	defer RequireErrorHasCategory(&err, rio.ErrorCategory(""))

	var anyWarehouses bool // for clarity in final error messages
	for _, addr := range warehouses {
		u, err := url.Parse(string(addr))
		if err != nil {
			return nil, Errorf(rio.ErrUsage, "failed to parse URI: %s", err)
		}
		switch u.Scheme {
		case "git":
			fallthrough
		case "ssh":
			fallthrough
		case "http", "https":
			fallthrough
		case "file":
			whCtrl, err = gitWarehouse.NewController(objcacheWorkdir, addr)
		default:
			return nil, Errorf(rio.ErrUsage, "this fetch operation doesn't support %q scheme (valid options are 'git', 'ssh', 'http', 'https', or 'file')", u.Scheme)
		}
		switch Category(err) {
		case nil:
			anyWarehouses = true
			// pass
		case rio.ErrWarehouseUnavailable:
			log.WarehouseUnavailable(mon, err, addr, wareID, "read")
			continue // okay!  skip to the next one.
		default:
			return nil, err
		}
		// Check if the local object store has the hash already; return early if so
		if whCtrl.Contains(wareID.Hash) {
			log.WareObjCacheHit(mon, wareID)
			return whCtrl, nil // happy path return!
		}
		// Fetch from the remote.
		err = whCtrl.Clone(ctx)
		if err == nil {
			err = whCtrl.Update(ctx)
		}
		switch Category(err) {
		case nil:
			anyWarehouses = true
			// pass
		case rio.ErrWarehouseUnavailable:
			log.WarehouseUnavailable(mon, err, addr, wareID, "read")
			continue // okay!  skip to the next one.
		default:
			return nil, err
		}
		// Check again if we have the object now after fetching.
		if whCtrl.Contains(wareID.Hash) {
			log.WareReaderOpened(mon, addr, wareID)
			return whCtrl, nil // happy path return!
		} else {
			log.WareNotFound(mon, fmt.Errorf("not in this repo"), addr, wareID)
			continue // okay!  skip to the next one.
		}
	}
	if !anyWarehouses {
		return nil, Errorf(rio.ErrWarehouseUnavailable, "no warehouses were available!")
	}
	return nil, Errorf(rio.ErrWareNotFound, "none of the available warehouses have ware %q!", wareID)
}
