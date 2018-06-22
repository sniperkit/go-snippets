package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/alexandrecarlton/gogurt"
	"github.com/alexandrecarlton/gogurt/packages"
)

// have table of latest versions?
// and checksums...


func main() {

	// TODO: Create a default config object, merge that in.
	// TODO: look at mitchellh's homedir instead of using os.GetEnv(HOME)
	// TODO: Viper this.
	config := gogurt.Config{
		Prefix: filepath.Join(os.Getenv("HOME"), ".local"),
		CacheFolder: filepath.Join(os.Getenv("HOME"), ".cache/gogurt"),
		BuildFolder: filepath.Join(os.Getenv("HOME"),  ".gogurt/build"),
		InstallFolder: filepath.Join(os.Getenv("HOME"), ".gogurt/install"),
		NumCores: 3,
		// TODO We'll have a set of default versions, and our config will override the defaults.
		PackageVersions: map[string]string {
			"asciidoc": "8.6.10",
			"aspell": "0.60.6.1",
			"autoconf": "2.69",
			"automake": "1.15",
			"bzip2": "1.0.6",
			"ccache": "3.3.4",
			"cmake": "3.9.1",
			"curl": "7.54.1",
			"editorconfig-core-c": "0.12.1",
			"expat": "2.2.0",
			"fish": "2.6.0",
			"fuse": "3.1.1",
			"gcc": "7.1.0",
			"gdb": "8.0",
			"git": "2.13.2",
			"glib": "2.53.6",
			"global": "6.5.7",
			"gmp": "6.1.2",
			"go": "1.9",
			"help2man": "1.47.4",
			"libevent": "2.1.8-stable",
			"libtool": "2.4.6",
			"libxml2": "2.9.4",
			"libxslt": "1.1.32",
			"lua": "5.3.4",
			"llvm": "4.0.1",
			"mpc": "1.0.3",
			"mpfr": "3.1.5",
			"ncurses": "6.0",
			"neovim": "0.2.0",
			"ninja": "1.7.2",
			"nodejs": "8.1.3",
			"openssl": "1.0.2k",
			"pbzip2": "1.1.13",
			"pkg-config": "0.29.2",
			"pcre": "8.40",
			"pigz": "2.3.4",
			"readline": "7.0",
			"rr": "4.5.0",
			"sqlite": "3.20.1",
			"sshfs": "3.2.0",
			"stow": "2.2.2",
			"texinfo": "6.3",
			"the_silver_searcher": "2.0.0",
			"tmux": "2.5",
			"tree": "1.7.0",
			"universal-ctags": "b5c9b76",
			"valgrind": "3.13.0",
			"libffi": "3.2.1",
			"xz": "5.2.3",
			"vim": "8.0.0045",
			"zsh": "5.3.1",
			"zlib": "1.2.11",
		},
	}

	name := os.Args[1]
	mappings := map[string]gogurt.Package{
		"asciidoc": packages.AsciiDoc{},
		"aspell": packages.Aspell{},
		"autoconf": packages.AutoConf{},
		"automake": packages.AutoMake{},
		"bzip2": packages.Bzip2{},
		"ccache": packages.CCache{},
		"cmake": packages.CMake{},
		"curl": packages.Curl{},
		"editorconfig-core-c": packages.EditorConfigCoreC{},
		"expat": packages.Expat{},
		"fish": packages.Fish{},
		"fuse": packages.FUSE{},
		"gcc": packages.GCC{},
		"gdb": packages.GDB{},
		"gettext" : packages.GetText{},
		"git": packages.Git{},
		"glib": packages.GLib{},
		"global": packages.Global{},
		"gmp": packages.GMP{},
		"go": packages.Go{},
		"help2man": packages.Help2Man{},
		"libevent": packages.Libevent{},
		"libffi": packages.LibFFI{},
		"libtool": packages.LibTool{},
		"libxml2": packages.LibXML2{},
		"libxslt": packages.LibXSLT{},
		"llvm": packages.LLVM{},
		"lua": packages.Lua{},
		"mpc": packages.MPC{},
		"mpfr": packages.MPFR{},
		"neovim": packages.NeoVim{},
		"ninja": packages.Ninja{},
		"openssl": packages.OpenSSL{},
		"pbzip2": packages.PBZip2{},
		"pcre": packages.Pcre{},
		"pigz": packages.Pigz{},
		"pkg-config": packages.PkgConfig{},
		"python2": packages.Python2{},
		"ncurses": packages.Ncurses{},
		"nodejs": packages.NodeJS{},
		"readline": packages.ReadLine{},
		"rr": packages.RR{},
		"sqlite": packages.SQLite{},
		"sshfs": packages.SSHFS{},
		"stow": packages.Stow{},
		"tree": packages.Tree{},
		"texinfo": packages.TexInfo{},
		"the_silver_searcher": packages.TheSilverSearcher{},
		"tmux": packages.Tmux{},
		"universal-ctags": packages.UniversalCTags{},
		"valgrind": packages.Valgrind{},
		"vim": packages.Vim{},
		"xz": packages.XZ{},
		"zlib": packages.Zlib{},
		"zsh": packages.Zsh{},
	}
	pac := mappings[name]

	installPackage(pac, config)
}

func installPackage(pac gogurt.Package, config gogurt.Config) {

	if _, err := os.Stat(config.InstallDir(pac)); err == nil {
		log.Printf("Package '%s' already installed, skipping.", pac.Name())
		return
	}

	for _, dependency := range pac.Dependencies() {
		installPackage(dependency, config)
	}
	version := config.PackageVersions[pac.Name()]

	// Download tarball
	url := pac.URL(version)
	cacheFilename := filepath.Join(config.CacheDir(pac), filepath.Base(url))

	if _, err := os.Stat(cacheFilename); err == nil {
		log.Printf("File '%s' already exists, not downloading a new copy.", cacheFilename)
	} else if err := gogurt.Download(url, cacheFilename); err != nil {
		log.Fatalf("Could not download url '%s' to file '%s': %s\n", url, cacheFilename, err.Error())
	}

	buildDirname := config.BuildDir(pac)
	if err := gogurt.DecompressSourceArchive(cacheFilename, buildDirname); err != nil {
		log.Fatalf("Error extracting %s to directory %s: %s", cacheFilename, buildDirname, err.Error())
	}
	if err := os.Chdir(buildDirname); err != nil {
		log.Fatalf("Error changing to directory '%s': %s", buildDirname, err.Error())
	}
	if err := pac.Build(config); err != nil {
		log.Fatalf("Error building package '%s': %s", pac.Name(), err.Error())
	}
	if err := pac.Install(config); err != nil {
		log.Fatalf("Error installing package '%s': %s", pac.Name(), err.Error())
	}
}
