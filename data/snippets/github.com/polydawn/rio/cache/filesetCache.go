package cache

import (
	"fmt"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/rio/fs"
	whutil "go.polydawn.net/rio/warehouse/util"
)

func ShelfFor(wareID api.WareID) fs.RelPath {
	chunk1, chunk2, _ := whutil.ChunkifyHash(wareID)
	return fs.MustRelPath(fmt.Sprintf("%s/fileset/%s/%s/%s",
		wareID.Type,
		chunk1, chunk2, wareID.Hash,
	))
}
