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
	Constructs a placer which will make files appear in place by use of an AUFS mount.

	If writable=false, the AUFS indirection will be skipped, and a simple bind mount used.
	If writable=true, an AUFS work/layer dir will be created in a tmpdir, and writes
	end up there (meaning the original source remains unmutated).
*/
func NewAufsPlacer(workDir fs.AbsolutePath) (Placer, error) {
	if err := fsOp.MkdirAll(rootFs, workDir.CoerceRelative(), 0700); err != nil {
		return nil, Errorf(rio.ErrLocalCacheProblem, "error creating aufs work area: %s", err)
	}
	return func(srcPath, dstPath fs.AbsolutePath, writable bool) (Janitor, error) {
		// Short-circuit into bind placer if not writable.
		if writable == false {
			return BindPlacer(srcPath, dstPath, writable)
		}

		// Determine desired type.
		//  Jump to copy placer if it's a file!  AUFS doesn't handle plain files.
		//  Jump to bind placer for any special files.
		srcStat, err := rootFs.LStat(srcPath.CoerceRelative())
		if err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error placing with aufs mount: %s", err)
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

		// Make the layer dir.
		//  Note that we're going to fix props on it in just a bit, because they
		//  leak through... but we have to do it *after* mount, because... AUFS.
		//  In doing so, fix props on layerPath; otherwise they instantly leak through.
		layerPath := workDir.Join(fs.MustRelPath("layer-" + guid.New()))
		if err := rootFs.Mkdir(layerPath.CoerceRelative(), 0700); err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error creating aufs layer area: %s", err)
		}

		// Set up AUFS mount.
		//  If you were doing this in a shell, it'd be roughly `mount -t aufs -o br="$layer":"$base" none "$composite"`.
		//  Yes, this may behave oddly in the event of paths containing ":" or "=".
		if err := syscall.Mount("none", dstPath.String(), "aufs", 0, fmt.Sprintf("br:%s=rw:%s=ro", layerPath.String(), srcPath.String())); err != nil {
			return nil, Errorf(rio.ErrAssemblyInvalid, "error placing with aufs mount: %s", err)
		}

		// Repair props on the layer dir.
		//  When we made the mount syscall, AUFS made a bunch of files like '.wh..wh.orph/'
		//  in the layer dir, which bumps its mtime, which leaks through to the final union.
		fmeta, _, err := fsOp.ScanFile(rootFs, srcPath.CoerceRelative())
		if err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error creating aufs layer area: %s", err)
		}
		// (This is usually what you'd use PlaceFile to do, but it errors on existing files.)
		fmeta.Name = layerPath.CoerceRelative()
		if err := rootFs.Lchown(fmeta.Name, fmeta.Uid, fmeta.Gid); err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error creating aufs layer area: %s", err)
		}
		if err := rootFs.Chmod(fmeta.Name, fmeta.Perms); err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error creating aufs layer area: %s", err)
		}
		if err := rootFs.SetTimesNano(fmeta.Name, fmeta.Mtime, fs.DefaultAtime); err != nil {
			return nil, Errorf(rio.ErrLocalCacheProblem, "error creating aufs layer area: %s", err)
		}

		// Return a cleanup func that will gracefully unmount... and also remove layer content.
		return aufsJanitor{
			dstPath,
			layerPath,
		}, nil
	}, nil
}

type aufsJanitor struct {
	mountPath fs.AbsolutePath
	layerPath fs.AbsolutePath
}

func (j aufsJanitor) Description() string {
	return fmt.Sprintf("umount %q; rm -rf %q;", j.mountPath, j.layerPath)
}
func (j aufsJanitor) Teardown() error {
	if err := syscall.Unmount(j.mountPath.String(), 0); err != nil {
		return Errorf(rio.ErrLocalCacheProblem, "error tearing down aufs mount: %s", err)
	}
	if err := os.RemoveAll(j.layerPath.String()); err != nil {
		return Errorf(rio.ErrLocalCacheProblem, "error tearing down aufs placement: %s", err)
	}
	return nil
}
func (j aufsJanitor) AlwaysTry() bool { return true }
