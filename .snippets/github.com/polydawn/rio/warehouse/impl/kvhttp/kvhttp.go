package kvhttp

import (
	"io"
	"net/http"
	"net/url"
	"path"

	. "github.com/warpfork/go-errcat"
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/warehouse"
	"go.polydawn.net/rio/warehouse/util"
)

var (
	_ warehouse.BlobstoreController = Controller{}
)

type Controller struct {
	addr     api.WarehouseAddr // user's string retained for messages
	baseUrl  *url.URL
	ctntAddr bool
}

/*
	Initialize a new warehouse controller that operates on a local filesystem.

	May return errors of category:

	  - `rio.ErrUsage` -- for unsupported addressses
	  - `rio.ErrWarehouseUnavailable` -- if the warehouse doesn't exist
*/
func NewController(addr api.WarehouseAddr) (warehouse.BlobstoreController, error) {
	// Stamp out a warehouse handle.
	//  More values will be accumulated in shortly.
	whCtrl := Controller{
		addr: addr,
	}

	// Verify that the addr is sensible up front, and extract features.
	//  - We parse things mostly like URLs.
	//  - We extract whether or not it's content-addressible mode here;
	//  - and extract the filesystem path, and normalize it to its absolute form.
	u, err := url.Parse(string(addr))
	if err != nil {
		return whCtrl, Errorf(rio.ErrUsage, "failed to parse URI: %s", err)
	}
	switch u.Scheme {
	case "http":
	case "ca+http":
		u.Scheme = "http"
		whCtrl.ctntAddr = true
	case "https":
	case "ca+https":
		u.Scheme = "https"
		whCtrl.ctntAddr = true
	default:
		return whCtrl, Errorf(rio.ErrUsage, "unsupported scheme in warehouse addr: %q (valid options are 'http', 'ca+http', 'https', or 'ca+https')", u.Scheme)
	}
	whCtrl.baseUrl = u

	// We skip checking that the warehouse exists.
	//  It's as costly as just starting the actual download.

	return whCtrl, nil
}

func (whCtrl Controller) OpenReader(wareID api.WareID) (io.ReadCloser, error) {
	u := whCtrl.baseUrl
	if whCtrl.ctntAddr {
		chunkA, chunkB, _ := util.ChunkifyHash(wareID)
		u.Path = path.Join(u.Path, chunkA, chunkB, wareID.Hash)
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, Errorf(rio.ErrWarehouseUnavailable, "error connecting to warehouse %s: %s", whCtrl.addr, err)
	}
	switch resp.StatusCode {
	case 200:
		return resp.Body, nil
	case 404:
		resp.Body.Close()
		return nil, Errorf(rio.ErrWareNotFound, "ware %s not found in warehouse %s", wareID, whCtrl.addr)
	default:
		resp.Body.Close()
		return nil, Errorf(rio.ErrWarehouseUnavailable, "unexpected HTTP code from warehouse %s: %s", whCtrl.addr, resp.Status)
	}
}

func (whCtrl Controller) OpenWriter() (warehouse.BlobstoreWriteController, error) {
	return nil, Errorf(rio.ErrUsage, "http warehouses are readonly!")
}
