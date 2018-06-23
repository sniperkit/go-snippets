package warehouse

import (
	"context"
	"io"

	"go.polydawn.net/go-timeless-api"
)

/*
	A blobstore-style warehouse supports opening reads and writes
	which return simple binary io.Reader and io.Writer streams.

	Blobstore backing implementations are typically simple key-value stores.
	Examples are 'kvfs' (using a local filesystem),
	'kvhttp' (readonly, aiming at http(s) URLs),
	'kvgs' (using Google Cloud Storage as a k/v bucket),
	'kvs3' (using AWS S3 as a k/v bucket), etc.

	Transmats using a blobstore warehouse have some packing format which
	reduces filesets down to a single binary stream; for example, the tar
	packing format.
*/
type BlobstoreController interface {
	OpenReader(wareID api.WareID) (io.ReadCloser, error)
	OpenWriter() (BlobstoreWriteController, error)
}

/*
	Blobstore-style warehouses return a "write controller", which is both
	a simple `io.Writer`, and also carries a `Commit` function which must
	be called when the write is complete and the hash known.

	Using Blobstore.OpenWriter causes temp space to be allocated in the
	warehouse to accept the incoming binary data.
	Calling `Commit` moves the data into final position and makes it available
	for reading, and closes the writer.
	Calling `Close` on the write controller before commit aborts the write,
	freeing the temp space used.

	The WareID given to the `Commit` call is assumed to be correct -- warehouses
	are a transport layer, and understand nothing of the packing format.
*/
type BlobstoreWriteController interface {
	io.WriteCloser
	Commit(wareID api.WareID) error
}

/*
	A no-op implementation of BlobstoreWriteController.
	You can use this to invoke a PackFunc as "scan only" -- it'll produce
	a wareID without actually saving the packed data anywhere.
*/
type NullBlobstoreWriteController struct{}

func (NullBlobstoreWriteController) Write(bs []byte) (int, error)   { return len(bs), nil }
func (NullBlobstoreWriteController) Close() error                   { return nil }
func (NullBlobstoreWriteController) Commit(wareID api.WareID) error { return nil }

/*
	A repository-style warehouse generally supports multiple versions of files
	stored in a custom format. We generally won't _write_ to these repositories
	because they tend to not support idempotent commits.

	Repository backing implementations typically have a cache of the current
	contents and a method of fetching updates. They will have a method to
	retrieve contents via hash.
	Examples include 'git'
*/
type RepositoryController interface {
	Clone(context.Context) error
}
