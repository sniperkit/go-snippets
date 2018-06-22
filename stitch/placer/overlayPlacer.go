package placer

import (
	"fmt"
	"os"
	"syscall"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fsOp"
	"go.polydawn.net/rio/lib/guid"
)

/*
	Constructs a placer which will make files appear in place by use of an overlayfs mount.

	If writable=false, the overlay indirection will be skipped, and a simple bind mount used.
	If writable=true, an overlay work/layer dir will be created in a tmpdir, and writes
	end up there (meaning the original source remains unmutated).
*/
func NewOverlayPlacer(workDir fs.AbsolutePath) (Placer, error) {
	if err := fsOp.MkdirAll(rootFs, workDir.CoerceRelative(), 0700); err != nil {
		return nil, Errorf(rio.ErrLocalCacheProblem, "error creating overlay work area: %s", err)
	}
	return func(srcPath, dstPath fs.AbsolutePath, writable bool) (Janitor, error) {
		// Short-circuit into bind placer if not writable.
		if writable == false {
			return BindPlacer(srcPath, dstPath, writable)
		}

		// Determine desired type.
		//  Jump to copy placer if it's a file!  Overlayfs doesn't handle plain files.
		//  Jump to bind placer for any special files.
		srcStat, err := rootFs.LStat(srcPath.CoerceRelative())
		if err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error placing with overlay mount: %s", err)
		}
		switch srcStat.Type {
		case fs.Type_File:
			// Files don't get a full mount.  A copy placer does the right thing.
			return CopyPlacer(srcPath, dstPath, writable)
		case fs.Type_Dir:
			// pass
		case fs.Type_Symlink, fs.Type_NamedPipe, fs.Type_Socket, fs.Type_Device, fs.Type_CharDevice:
			return BindPlacer(srcPath, dstPath, writable)
		default:
			panic("unreachable file type enum")
		}

		// Make the destination path exist and be the right type to mount over.
		if err := mkDest(dstPath, srcStat.Type); err != nil {
			return nil, err
		}

		// Make the layer and work dirs.
		//  In doing so, fix props on upperPath; otherwise they instantly leak through.
		//  (Notice how this is easier than with AUFS, because Overlay's design of
		//  splitting work versus layer dirs fixes a LOT of systemic stupidity.)
		overlayPath := workDir.Join(fs.MustRelPath("overlay-" + guid.New()))
		workPath := overlayPath.Join(fs.MustRelPath("work"))
		upperPath := overlayPath.Join(fs.MustRelPath("upper"))
		if err := rootFs.Mkdir(overlayPath.CoerceRelative(), 0700); err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error creating overlay work area: %s", err)
		}
		if err := rootFs.Mkdir(workPath.CoerceRelative(), 0700); err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error creating overlay work area: %s", err)
		}
		fmeta, _, err := fsOp.ScanFile(rootFs, srcPath.CoerceRelative())
		if err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error creating overlay work area: %s", err)
		}
		fmeta.Name = upperPath.CoerceRelative()
		if err := fsOp.PlaceFile(rootFs, *fmeta, nil, false); err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error creating overlay work area: %s", err)
		}

		// Set up overlay mount.
		//  If you were doing this in a shell, it'd be roughly `mount -t overlay overlay -o lowerdir=lower,upperdir=upper,workdir=work mntpoint`.
		//  Yes, this may behave oddly in the event of paths containing ":" or "=" or ",".
		if err := syscall.Mount("none", dstPath.String(), "overlay", 0, fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", srcPath, upperPath, workPath)); err != nil {
			return nil, Errorf(rio.ErrAssemblyInvalid, "error placing with overlay mount: %s", err)
		}

		// Return a cleanup func that will gracefully unmount... and also remove layer content.
		return overlayJanitor{
			dstPath,
			upperPath,
			workPath,
		}, nil
	}, nil
}

type overlayJanitor struct {
	mountPath fs.AbsolutePath
	upperPath fs.AbsolutePath
	workPath  fs.AbsolutePath
}

func (j overlayJanitor) Description() string {
	return fmt.Sprintf("umount %q; rm -rf %q; rm -rf %q;", j.mountPath, j.upperPath, j.workPath)
}
func (j overlayJanitor) Teardown() error {
	if err := syscall.Unmount(j.mountPath.String(), 0); err != nil {
		return Errorf(rio.ErrLocalCacheProblem, "error tearing down overlay mount: %s", err)
	}
	if err := os.RemoveAll(j.upperPath.String()); err != nil {
		return Errorf(rio.ErrLocalCacheProblem, "error tearing down overlay placement: %s", err)
	}
	if err := os.RemoveAll(j.workPath.String()); err != nil {
		return Errorf(rio.ErrLocalCacheProblem, "error tearing down overlay placement: %s", err)
	}
	return nil
}
func (j overlayJanitor) AlwaysTry() bool { return true }
