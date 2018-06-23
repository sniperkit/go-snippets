package placer

import (
	"io/ioutil"
	"os/exec"
	"strings"

	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/config"
	"go.polydawn.net/rio/fs"
)

/*
	The copy placer is always defined and always supported and is never swappable.

	The placer used to handle "mount"-mode placements may vary drastically, however.
*/

/*
	Returns the most reasonable mounting placer implementation available on this platform.

	In order, the attempted systems are: overlayfs, aufs.
	(Bind mounts are not considered a valid substitution by default, because their
	behavior on "writable=true" is very different.)

	If the environment var RIO_MOUNT_PLACER is set, we'll either return that or
	an error explaing why it's not available;
	otherwise, autodetection will examine what filesystem drivers and capabilities
	are available, and try to pick the most performant/reliable/sensible thing available.

	For placers that need a working dir, one will be created under RIO_MOUNT_WORKDIR
	if set,	or RIO_BASE/wrk.
*/
func GetMountPlacer() (Placer, error) {
	// TODO the switch for explicit RIO_MOUNT_PLACER
	if isFSAvailable("overlay") {
		return NewOverlayPlacer(config.GetMountWorkPath().Join(fs.MustRelPath("overlay")))
	}
	if isFSAvailable("aufs") {
		return NewAufsPlacer(config.GetMountWorkPath().Join(fs.MustRelPath("aufs")))
	}
	return nil, Errorf(rio.ErrAssemblyInvalid, "placer: no power (cannot find usable mount placer)")
}

/*
	Detect if a filesystem is available according to the kernel.

	Also attempts modprobe in case the modules for the filesystem is available but not loaded.
	This is aggressive, but in practice, we have observed that "apt-get install aufs"
	does not necessarily load it nor mark it for load on future boots, despite that
	pretty likely being what the user wants, so we consider it sensible to do that load ourselves.
*/
func isFSAvailable(fs string) bool {
	// Arguably the greatest thing to do would of course just be to issue the syscall once and see if it flies...
	// but that's a distrubingly stateful and messy operation so we're gonna check a bunch of next-best-things instead.

	// If it's in /proc/filesystems, we should be good to go.
	// (If it's not, the libs might be installed, but not loaded, so we'll try that next.)
	if fss, err := ioutil.ReadFile("/proc/filesystems"); err == nil {
		fssLines := strings.Split(string(fss), "\n")
		for _, line := range fssLines {
			parts := strings.Split(line, "\t")
			if len(parts) < 2 {
				continue
			}
			if parts[1] == fs {
				return true
			}
		}
	}

	// Blindly attempt to modprobe the FS module into the kernel.
	// If the modprobe command exists, we can attempt to use it to load the FS module.
	// If it doesn't, okay, bail; FS is not available and we can't get it.
	// If it works, great.  If it doesn't, okay, we'll move on.
	// Repeatedly installing it if it already exists no-op's correctly.
	if err := exec.Command(
		"modprobe", fs,
	).Run(); err != nil {
		return false
	}

	return true
}
