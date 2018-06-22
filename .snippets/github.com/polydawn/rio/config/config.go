/*
	Helpers for loading contextual config.

	Config for Rio means "things that are the host machine operator's concerns".
	So, things like cache paths and preferred mounting systems are considered "config",
	as opposed to parameters for function calls.
	(This distinction is meaningful because config is generally not passed in calls,
	because it wouldn't be correct to do so when using commands via remote RPC; in
	such a situation, the *remote* Rio will read its *local* config in order to
	comply with the operator's rules there on that machine and environment.)
*/
package config

import (
	"os"
	"path/filepath"

	"go.polydawn.net/rio/fs"
)

/*
	Return the path that is the root for rio's fileset caches.

	The default value is `"$RIO_BASE/cache"`;
	this can be overriden by the `RIO_CACHE` environment variable.
*/
func GetCacheBasePath() fs.AbsolutePath {
	pth := os.Getenv("RIO_CACHE")
	if pth == "" {
		return GetRioBasePath().Join(fs.MustRelPath("cache"))
	}
	pth, err := filepath.Abs(pth)
	if err != nil {
		panic(err)
	}
	return fs.MustAbsolutePath(pth)
}

/*
	Return the path prefix that will be used as a workspace for mount subsystems.

	The default value is `"$RIO_BASE/mount"`;
	this can be overriden by the `RIO_MOUNT_WORKDIR` environment variable.
*/
func GetMountWorkPath() fs.AbsolutePath {
	pth := os.Getenv("RIO_MOUNT_WORKDIR")
	if pth == "" {
		return GetRioBasePath().Join(fs.MustRelPath("mount"))
	}
	pth, err := filepath.Abs(pth)
	if err != nil {
		panic(err)
	}
	return fs.MustAbsolutePath(pth)
}

/*
	Return the home-base path prefix that is the default root for all other Rio paths.

	The default value is `"/var/lib/timeless/rio"`;
	this can be overriden by the `RIO_BASE` environment variable.
*/
func GetRioBasePath() fs.AbsolutePath {
	pth := os.Getenv("RIO_BASE")
	if pth == "" {
		pth = "/var/lib/timeless/rio"
	}
	pth, err := filepath.Abs(pth)
	if err != nil {
		panic(err)
	}
	return fs.MustAbsolutePath(pth)
}
