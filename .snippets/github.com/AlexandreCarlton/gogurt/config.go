package gogurt

import (
	"path/filepath"
)

// Shit. We need to put in the latest versions here as well.
// Have a like a subsection, e.g. 'bzip2': '1.0.6'

type Config struct {
	// where we assume applications will finally be installed.
	Prefix string

	// Where we store our tarballs.
	CacheFolder string

	// Where we extract and build our applications
	BuildFolder string

	// Where we install applications to keep them separate.
	InstallFolder string

	// How many cores are allocated to building applications.
	NumCores uint

  // We load from defaults, overwriting with the entries found in the config file.
	PackageVersions map[string]string

}

func (config Config) CacheDir(archive SourceArchive) string {
	return filepath.Join(config.CacheFolder, archive.Name(), config.PackageVersions[archive.Name()])
}

func (config Config) BuildDir(p Package) string {
  return filepath.Join(config.BuildFolder, p.Name(), config.PackageVersions[p.Name()])
}

func (config Config) InstallDir(p Package) string {
  return filepath.Join(config.InstallFolder, p.Name(), config.PackageVersions[p.Name()])
}

func (config Config) BinDir(p Package) string {
	return filepath.Join(config.InstallDir(p), "bin")
}

func (config Config) IncludeDir(p Package) string {
	return filepath.Join(config.InstallDir(p), "include")
}

func (config Config) LibDir(p Package) string {
	return filepath.Join(config.InstallDir(p), "lib")
}

func (config Config) PkgConfigLibDir(p Package) string {
	return filepath.Join(config.LibDir(p), "pkgconfig")
}

func (config Config) PkgConfigShareDir(p Package) string {
	return filepath.Join(config.InstallDir(p), "share", "pkgconfig")
}
