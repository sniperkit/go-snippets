// +build linux

// We needed linx-specific syscalls not exported by the standard lib in order to get
// chtimes on symlinks with nano precision to work correctly.
// (Stdlib only provides 'chtimes', no 'lchtimes'.)

package osfs

import (
	"syscall"
	"time"
	"unsafe"

	"go.polydawn.net/rio/fs"
)

func (afs *osFS) SetTimesLNano(path fs.RelPath, mtime time.Time, atime time.Time) error {
	rpath, err := afs.realpath(path, false)
	if err != nil {
		return err
	}

	var utimes [2]syscall.Timespec
	utimes[0] = syscall.NsecToTimespec(atime.UnixNano())
	utimes[1] = syscall.NsecToTimespec(mtime.UnixNano())

	// These are not currently available in syscall
	AT_FDCWD := -100
	AT_SYMLINK_NOFOLLOW := 0x100

	var _path *byte
	_path, err = syscall.BytePtrFromString(rpath)
	if err != nil { // EINVAL if the path string contains NUL bytes.
		return fs.NormalizeIOError(err)
	}

	// Note this does depend on kernel 2.6.22 or newer.  Fallbacks are available but we haven't implemented them and they lose nano precision.
	if _, _, err := syscall.Syscall6(syscall.SYS_UTIMENSAT, uintptr(AT_FDCWD), uintptr(unsafe.Pointer(_path)), uintptr(unsafe.Pointer(&utimes[0])), uintptr(AT_SYMLINK_NOFOLLOW), 0, 0); err != 0 {
		return fs.NormalizeIOError(err)
	}

	return nil
}

func (afs *osFS) SetTimesNano(path fs.RelPath, mtime time.Time, atime time.Time) error {
	rpath, err := afs.realpath(path, true)
	if err != nil {
		return err
	}

	// Note that this is disambiguated from plain `os.Chtimes` only in that it refuses to fall back to lower precision on old kernels.
	// Like LUtimesNano, it depends on kernel 2.6.22 or newer.
	var utimes [2]syscall.Timespec
	utimes[0] = syscall.NsecToTimespec(atime.UnixNano())
	utimes[1] = syscall.NsecToTimespec(mtime.UnixNano())
	if err := syscall.UtimesNano(rpath, utimes[0:]); err != nil {
		return fs.NormalizeIOError(err)
	}
	return nil
}
