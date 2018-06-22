package kvfs

import (
	"io"
	"net/url"
	"os"
	"path/filepath"

	. "github.com/warpfork/go-errcat"
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/lib/guid"
	"go.polydawn.net/rio/warehouse"
	"go.polydawn.net/rio/warehouse/util"
)

var (
	_ warehouse.BlobstoreController      = Controller{}
	_ warehouse.BlobstoreWriteController = &WriteController{}
)

type Controller struct {
	addr     api.WarehouseAddr // user's string retained for messages
	basePath fs.AbsolutePath
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
	case "file":
	case "ca+file":
		whCtrl.ctntAddr = true
	default:
		return whCtrl, Errorf(rio.ErrUsage, "unsupported scheme in warehouse addr: %q (valid options are 'file' or 'ca+file')", u.Scheme)
	}
	absPth, err := filepath.Abs(filepath.Join(u.Host, u.Path))
	if err != nil {
		panic(err)
	}
	whCtrl.basePath = fs.MustAbsolutePath(absPth)

	// Check that the warehouse exists.
	//  If it does, we're good: return happily.
	checkPath := whCtrl.basePath
	if !whCtrl.ctntAddr {
		// In non-CA mode, the check for warehouse existence is a little strange;
		//  for reading, we could declare 404 if the path doesn't exist... but we don't
		//  know whether this is going to be used for reading or writing yet!
		//  So we have to look at the path segment above it to see if a write might be valid.
		checkPath = checkPath.Dir()
	}
	stat, err := os.Stat(checkPath.String())
	switch {
	case os.IsNotExist(err):
		return whCtrl, Errorf(rio.ErrWarehouseUnavailable, "warehouse does not exist (%s)", err)
	case err != nil: // normally we'd style this as the default cause, but, we must check it before the IsDir option
		return whCtrl, Errorf(rio.ErrWarehouseUnavailable, "warehouse unavailable (%s)", err)
	case !stat.IsDir():
		return whCtrl, Errorf(rio.ErrWarehouseUnavailable, "warehouse does not exist (%s is not a dir)", checkPath)
	default: // only thing left is err == nil
		return whCtrl, nil
	}
}

func (whCtrl Controller) OpenReader(wareID api.WareID) (io.ReadCloser, error) {
	finalPath := whCtrl.basePath
	if whCtrl.ctntAddr {
		chunkA, chunkB, _ := util.ChunkifyHash(wareID)
		finalPath = finalPath.
			Join(fs.MustRelPath(chunkA)).
			Join(fs.MustRelPath(chunkB)).
			Join(fs.MustRelPath(wareID.Hash))
	}
	file, err := os.OpenFile(finalPath.String(), os.O_RDONLY, 0)
	switch {
	case err == nil:
		return file, nil
	case os.IsNotExist(err):
		return nil, Errorf(rio.ErrWareNotFound, "ware %s not found in warehouse %s", wareID, whCtrl.addr)
	default:
		return nil, Errorf(rio.ErrWarehouseUnavailable, "ware %s could not be retrieved from warehouse %s: %s", wareID, whCtrl.addr, err)
	}
}

func (whCtrl Controller) OpenWriter() (warehouse.BlobstoreWriteController, error) {
	wc := &WriteController{whCtrl: whCtrl}
	// Pick a random upload path.
	if whCtrl.ctntAddr {
		tmpName := fs.MustRelPath(".tmp.upload." + guid.New())
		wc.stagePath = whCtrl.basePath.Join(tmpName)
	} else {
		// In non-CA mode, "base" path isn't really "base"; it's the final destination.
		tmpName := fs.MustRelPath(".tmp.upload." + whCtrl.basePath.Last() + "." + guid.New())
		wc.stagePath = whCtrl.basePath.Dir().Join(tmpName)
	}
	// Open file the file for write.
	file, err := os.OpenFile(wc.stagePath.String(), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return wc, Errorf(rio.ErrWarehouseUnwritable, "failed to reserve temp space in warehouse: %s", err)
	}
	wc.stream = file
	// Return the controller -- which has methods to either commit+close, or cancel+close.
	return wc, nil
}

type WriteController struct {
	stream    io.WriteCloser  // Write to this.
	whCtrl    Controller      // Needed for the final move-into-place.
	stagePath fs.AbsolutePath // Needed for the final move-into-place.
}

func (wc *WriteController) Write(bs []byte) (int, error) {
	return wc.stream.Write(bs)
}

/*
	Cancel the current write.  Close the stream, and remove any temporary files.
*/
func (wc *WriteController) Close() error {
	wc.stream.Close()
	return os.Remove(wc.stagePath.String())
}

/*
	Commit the current data as the given hash.
	Caller must be an adult and specify the hash truthfully.
	Closes the writer and invalidates any future use.
*/
func (wc *WriteController) Commit(wareID api.WareID) error {
	// Close the file.
	if err := wc.stream.Close(); err != nil {
		return Errorf(rio.ErrWarehouseUnwritable, "failed to commit to file: %s", err)
	}
	// Compute final path.
	// Make parent dirs if necessary in content-addr mode.
	finalPath := wc.whCtrl.basePath
	if wc.whCtrl.ctntAddr {
		chunkA, chunkB, _ := util.ChunkifyHash(wareID)
		finalPath = finalPath.Join(fs.MustRelPath(chunkA))
		if err := os.Mkdir(finalPath.String(), 0755); err != nil && !os.IsExist(err) {
			return Errorf(rio.ErrWarehouseUnwritable, "failed to commit to file: %s", err)
		}
		finalPath = finalPath.Join(fs.MustRelPath(chunkB))
		if err := os.Mkdir(finalPath.String(), 0755); err != nil && !os.IsExist(err) {
			return Errorf(rio.ErrWarehouseUnwritable, "failed to commit to file: %s", err)
		}
		finalPath = finalPath.Join(fs.MustRelPath(wareID.Hash))
	}
	// Move into place.
	if err := os.Rename(wc.stagePath.String(), finalPath.String()); err != nil {
		return Errorf(rio.ErrWarehouseUnwritable, "failed to commit to file: %s", err)
	}
	return nil
}
