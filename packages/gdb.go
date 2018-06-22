package packages

import (
	"fmt"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type GDB struct{}

func (gdb GDB) Name() string {
	return "gdb"
}

func (gdb GDB) URL(version string) string {
	return fmt.Sprintf("https://ftp.gnu.org/gnu/gdb/gdb-%s.tar.gz", version)
}

func (gdb GDB) Build(config gogurt.Config) error {
	// We only want to build what's in gdb/
	gdbDir := filepath.Join(config.BuildDir(gdb), "gdb")

	// Python was statically linked with OpenSSL and Zlib.
	// This means that the python-config.py should link this in.
	// However, 'python' should come BEFORE the libraries it links to.
	// TODO: Raise PR.
	gogurt.ReplaceInFile(
		filepath.Join(gdbDir, "python", "python-config.py"),
		"libs = \\[\\]",
		"libs = ['-lpython' + pyver]")
	gogurt.ReplaceInFile(
		filepath.Join(gdbDir, "python", "python-config.py"),
		".*libs\\.append\\('-lpython'.*",
		"")
	// Ncurses also generates a tinfo library as well, which we must link in.
	gogurt.ReplaceInFile(
		filepath.Join(gdbDir, "configure"),
		"'' ncurses",
		"'ncursesw -ltinfow' ncurses")

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(gdb),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--enable-tui",
			"--with-curses",
			"--with-python=" + config.InstallDir(Python2{}),
			"--with-lzma",
			"--with-liblzma-prefix=" + config.InstallDir(XZ{}),
			"--with-libexpat",
			"--with-libexpat-prefix=" + config.InstallDir(Expat{}),
			//"--with-system-zlib",
			"--with-system-readline",
		},
		CFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
			"-I" + config.IncludeDir(ReadLine{}),
			"-I" + config.IncludeDir(Zlib{}),
		},
		CxxFlags: []string{
			"-I" + config.IncludeDir(Ncurses{}),
			"-I" + config.IncludeDir(ReadLine{}),
			"-I" + config.IncludeDir(Zlib{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(Ncurses{}),
			"-L" + config.LibDir(Zlib{}),
			// For linking of Python.
			"-L" + config.LibDir(Python2{}),
			"-L" + config.LibDir(ReadLine{}),
			"-L" + config.LibDir(OpenSSL{}),
			"-L" + config.LibDir(Expat{}),
			"-L" + config.LibDir(LibFFI{}),
			"-L" + config.LibDir(Zlib{}),
		},
		Libs: []string{
			"-lreadline",
			"-ltinfow",
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

func (gdb GDB) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{
			"install",
		},
	}.Cmd()
	return makeInstall.Run()
}

func (gdb GDB) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Expat{},
		Ncurses{},
		XZ{},
		Python2{},
	}
}
