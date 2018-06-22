package packages

// NOTE:
// There are several things that make Vim
// Vim hardcodes the prefix used to look for its runtime files when configuring (in src/auto/pathdef.c).
// These files are in PREFIX/share/vim/vim<ver>. which is problematic if your prefix is in /home.
// This is normally not a problem for personal use, but if your prefix is in /home then it is problematic to distribute this to your friends.
// They'll need to set the environment variable VIMRUNTIME=<path/to>/share/vim/vim<ver>. Then everything should just work.

import (
	"fmt"
	"path/filepath"
	// "strings"
	"github.com/alexandrecarlton/gogurt"
)

type Vim struct{}

func (vim Vim) Name() string {
	return "vim"
}

func (vim Vim) URL(version string) string {
	return fmt.Sprintf("https://github.com/vim/vim/archive/v%s.tar.gz", version)
}

func (vim Vim) Build(config gogurt.Config) error {

	//pythonMajorVersion := strings.Join(strings.Split(config.PackageVersions["python2"], ".")[:2], ".")

	configure := gogurt.ConfigureCmd{
		// Note: Vim cares where it thinks it's prefix is, as it uses this to find its
		// accompanying *.vim files
		Prefix: config.Prefix,
		Args: []string{
			"--with-features=huge",
			"--enable-multibyte",
			"--enable-cscope",
			"--enable-rubyinterp", // TODO: Need to install ruby.
			"--enable-pythoninterp",
			"--with-python-config-dir=" + filepath.Join(config.LibDir(Python2{}), "python2.7", "config"), // TODO: Get major the version out of the config.
			"--enable-luainterp",
			"--disable-darwin",
			"--disable-gui",
			"--disable-netbeans",
			"--without-x",
			"--with-tlib=tinfow",
		},
		CFlags: []string{
			"-Os", // Vim tries to compile with _FORTIFY_SOURCE, which requires -O
			"-I" + config.IncludeDir(OpenSSL{}),
			"-I" + config.IncludeDir(Ncurses{}),
			"-I" + config.IncludeDir(Python2{}),
		},
		CppFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
		},
		LdFlags: []string{
			"-static",
			"-L" + config.LibDir(Ncurses{}),
			"-L" + config.LibDir(Python2{}),

			// Included due to Vim's linking of Python.
			"-L" + config.LibDir(OpenSSL{}),
			"-L" + config.LibDir(Expat{}),
			"-L" + config.LibDir(LibFFI{}),
			"-L" + config.LibDir(ReadLine{}),
			"-L" + config.LibDir(Zlib{}),
		},
		Libs: []string{
			"-lncursesw",
			"-ltinfow",
			"-lpython2.7",
			"-lffi",
			"-lreadline",
			"-lssl",
			"-lcrypto",
			"-ldl",
			"-lz",
		},
		Paths: []string{
			config.BinDir(Python2{}), // We want our static Python2 to be picked up.
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()
	return make.Run()
}

func (vim Vim) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{
			"install",
			"prefix=" + config.InstallDir(vim),
		},
	}.Cmd()
	return makeInstall.Run()
}

func (vim Vim) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Ncurses{},
		Python2{},
	}
}

