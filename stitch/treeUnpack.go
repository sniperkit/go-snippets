package stitch

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/cache"
	"go.polydawn.net/rio/config"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/fsOp"
	"go.polydawn.net/rio/stitch/placer"
)

/*
	Struct to gather the args for a single rio.Unpack func call.
	(The context object and monitors are handled in a different band.)

	It may be interesting to note the similarity to Formula.Inputs from
	the go-timeless-api packages, but they are distinct:
	this one is internal, not serializable, contains the list of warehouses,
	as well as the path inline rather than in a map key, so we can sort slices.
*/
type UnpackSpec struct {
	Path       fs.AbsolutePath
	WareID     api.WareID
	Filters    api.FilesetFilters
	Warehouses []api.WarehouseAddr
	Monitor    rio.Monitor
}

// Cast slices to this type to sort by target path (which is effectively mountability order).
type UnpackSpecByPath []UnpackSpec

func (a UnpackSpecByPath) Len() int           { return len(a) }
func (a UnpackSpecByPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UnpackSpecByPath) Less(i, j int) bool { return a[i].Path.String() < a[j].Path.String() }

type unpackResult struct {
	Path     fs.AbsolutePath // cache path or mount source path
	Writable bool
	Error    error
}

type Assembler struct {
	cache      fs.FS
	unpackTool rio.UnpackFunc
	placerTool placer.Placer
}

func NewAssembler(unpackTool rio.UnpackFunc) (*Assembler, error) {
	placerTool, err := placer.GetMountPlacer()
	if err != nil {
		return nil, err
	}
	return &Assembler{
		cache:      osfs.New(config.GetCacheBasePath()),
		unpackTool: unpackTool,
		placerTool: placerTool,
	}, nil
}

func (a *Assembler) Run(ctx context.Context, targetFs fs.FS, parts []UnpackSpec, fillerDirProps fs.Metadata) (func() error, error) {
	sort.Sort(UnpackSpecByPath(parts))

	// Unpacking either wares or more mounts into paths under mounts is seriously illegal.
	//  It's a massive footgun, entirely strange, and just No.
	//  Doing it into paths under other wares is fine because it's not *leaving* our zone.
	mounts := map[fs.AbsolutePath]struct{}{}
	for _, part := range parts {
		for mount := range mounts {
			if strings.HasPrefix(part.Path.String(), mount.String()) {
				return nil, Errorf(rio.ErrAssemblyInvalid, "invalid inputs config: "+
					"cannot stitch additional inputs under a mount (%q is under mount at %q)",
					part.Path, mount)
			}
		}
		// If this one is a mount, mark it for the rest.
		//  (Paths under it must come after it, due to the sort.)
		if part.WareID.Type == "mount" {
			mounts[part.Path] = struct{}{}
		}
	}

	// Fan out materialization into cache paths.
	unpackResults := make([]unpackResult, len(parts))
	var wg sync.WaitGroup
	wg.Add(len(parts))
	for i, part := range parts {
		go func(i int, part UnpackSpec) {
			defer wg.Done()
			res := &unpackResults[i]
			// If it's a mount, do some parsing, and that's it for prep work.
			//  Also close the monitor channel, because every other unpack tool would.
			if part.WareID.Type == "mount" {
				if part.Monitor.Chan != nil {
					close(part.Monitor.Chan)
				}
				ss := strings.SplitN(part.WareID.Hash, ":", 2)
				if len(ss) != 2 {
					res.Error = Errorf(rio.ErrAssemblyInvalid, "invalid inputs config: mounts must specify mode (e.g. \"ro:/path\" or \"rw:/path\"")
					return
				}
				switch ss[0] {
				case "rw":
					res.Writable = true
				case "ro":
					res.Writable = false
				default:
					res.Error = Errorf(rio.ErrAssemblyInvalid, "invalid inputs config: mounts must specify mode (e.g. \"ro:/path\" or \"rw:/path\"")
					return
				}
				res.Path, res.Error = fs.ParseAbsolutePath(ss[1])
				return
			}
			// Unpack with placement=none to populate cache.
			resultWareID, err := a.unpackTool(
				ctx, // TODO fork em out
				part.WareID,
				"-",
				part.Filters,
				rio.Placement_None,
				part.Warehouses,
				part.Monitor,
			)
			if err != nil {
				res.Error = err
				return
			}
			// Yield the cache path.
			res.Path = config.GetCacheBasePath().Join(cache.ShelfFor(resultWareID))
			res.Writable = true
			// TODO if any error, fan out cancellations
		}(i, part)
	}
	wg.Wait()
	// Yield up any errors from individual unpacks.
	for _, result := range unpackResults {
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// Zip up all placements, in order.
	//  Parent dirs are made as necessary along the way.
	hk := &housekeeping{}
	for i, part := range parts {
		path := part.Path.CoerceRelative()

		// Ensure parent dirs.
		//  We do it in a closure to give convenint scope for defers.
		if err := func() error {
			for _, parentPath := range path.Dir().Split() {
				target, isSymlink, err := targetFs.Readlink(parentPath)
				if isSymlink {
					// Future hackers: if you ever try to make this check cleverer,
					//  also make sure you include a check for host mount crossings.
					return Recategorize(rio.ErrAssemblyInvalid,
						fs.NewBreakoutError(
							targetFs.BasePath(),
							path,
							parentPath,
							target,
						))
				} else if err == nil {
					continue
				} else if Category(err) == fs.ErrNotExists {
					// Capture the parent mtime for restoration before issuing a syscall that bonks it.
					defer fsOp.RepairMtime(targetFs, parentPath.Dir())()
					// Make the parent dir if it does not exist.
					fillerDirProps.Name = parentPath
					// Could be cleaner: this PlaceFile call rechecks the symlink thing, but it's the shortest call for "make all props right plz".
					if err := fsOp.PlaceFile(targetFs, fillerDirProps, nil, false); err != nil {
						return Errorf(rio.ErrAssemblyInvalid, "error creating parent dirs in tree unpack: %s", err)
					}
				} else {
					// Halt assembly attempt for any unhandlable errors that come up during parent path establishment.
					return err
				}
			}
			return nil
		}(); err != nil {
			return nil, err
		}

		// Invoke placer.
		//  Accumulate the individual cleanup funcs into a mega func we'll return.
		//  If errors occur during any placement, fire the cleanups so far before returning.
		targetPath := targetFs.BasePath().Join(part.Path.CoerceRelative())
		var janitor placer.Janitor
		var err error
		switch part.WareID.Type {
		case "mount":
			janitor, err = placer.BindPlacer(unpackResults[i].Path, targetPath, unpackResults[i].Writable)
		default:
			janitor, err = a.placerTool(unpackResults[i].Path, targetPath, unpackResults[i].Writable)
		}
		if err != nil {
			hk.Teardown()
			return nil, err
		}
		hk.append(janitor)
	}
	return hk.Teardown, nil
}

type housekeeping struct {
	CleanupStack []placer.Janitor
}

func (hk *housekeeping) append(janitor placer.Janitor) {
	hk.CleanupStack = append(hk.CleanupStack, janitor)
}

func (hk housekeeping) Teardown() error {
	progress := make([]string, len(hk.CleanupStack))
	var firstError error
	for i := len(hk.CleanupStack) - 1; i >= 0; i-- {
		janitor := hk.CleanupStack[i]
		if firstError != nil && !janitor.AlwaysTry() {
			progress[i] = "\tskipped: " + janitor.Description()
			continue
		}
		err := hk.CleanupStack[i].Teardown()
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			progress[i] = "\tfailed:  " + janitor.Description()
			continue
		}
		progress[i] = "\tsuccess: " + janitor.Description()
	}
	if firstError != nil {
		// Keep the category of the first one, but also fold in
		//  the string of everything that did or did not get cleaned up.
		cleanupReport := strings.Join(progress, "\n")
		firstError = ErrorDetailed(
			Category(firstError),
			fmt.Sprintf("%s.  The following cleanups were attempted:\n%s", firstError, cleanupReport),
			map[string]string{"cleanupReport": cleanupReport},
		)
	}
	return firstError
}
