package stitch

import (
	"context"
	"sort"
	"sync"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/fs"
)

/*
	Struct to gather the args for a single rio.Pack func call.
	(The context object and monitors are handled in a different band.)

	It may be interesting to note the similarity to Formula.Outputs from
	the go-timeless-api packages, but they are distinct:
	this one is internal, not serializable, contains the list of warehouses,
	as well as the path inline rather than in a map key, so we can sort slices.
*/
type PackSpec struct {
	Path      fs.AbsolutePath
	PackType  api.PackType
	Filters   api.FilesetFilters
	Warehouse api.WarehouseAddr
	Monitor   rio.Monitor
}

// Cast slices to this type to sort by target path (which is effectively mountability order).
type PackSpecByPath []PackSpec

func (a PackSpecByPath) Len() int           { return len(a) }
func (a PackSpecByPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PackSpecByPath) Less(i, j int) bool { return a[i].Path.String() < a[j].Path.String() }

type packResult struct {
	WareID api.WareID
	Error  error
}

func PackMulti(ctx context.Context, packTool rio.PackFunc, targetFs fs.FS, parts []PackSpec) (map[api.AbsPath]api.WareID, error) {
	// Since packfuncs do not mutate their target path, the order we launch them
	//  is not actually important.  But we sort it anyway, just for consistency.
	sort.Sort(PackSpecByPath(parts))

	// Fan out packing in parallel.
	packResults := make([]packResult, len(parts))
	var wg sync.WaitGroup
	wg.Add(len(parts))
	for i, part := range parts {
		go func(i int, part PackSpec) {
			defer wg.Done()
			res := &packResults[i]
			// Unpack with placement=none to populate cache.
			rerootedPath := targetFs.BasePath().Join(part.Path.CoerceRelative())
			res.WareID, res.Error = packTool(
				ctx, // TODO fork em out
				part.PackType,
				rerootedPath.String(),
				part.Filters,
				part.Warehouse,
				part.Monitor,
			)
			// TODO if any error, fan out cancellations
		}(i, part)
	}
	wg.Wait()

	// Gather results; any errors from individual packs error all.
	//  In the event of multiple errors, only the first is reported.
	var err error
	var results = make(map[api.AbsPath]api.WareID, len(parts))
	for i, result := range packResults {
		results[api.AbsPath(parts[i].Path.String())] = result.WareID
		if result.Error != nil && err == nil {
			err = result.Error
		}
	}
	return results, err
}
