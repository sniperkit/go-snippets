package packages


// Ncurses terminology:
// - terminfo:
// - termcap:
// - tic:

// Find the fallbacks in misc/terminfo.src
// Each entry is a line that begins with a non #, non space character, finishing at the pipe.


import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/alexandrecarlton/gogurt"
)

type Ncurses struct{}

func (ncurses Ncurses) Name() string {
	return "ncurses"
}

func (ncurses Ncurses) URL(version string) string {
	return fmt.Sprintf("http://ftp.gnu.org/gnu/ncurses/ncurses-%s.tar.gz", version)
}

func (ncurses Ncurses) Build(config gogurt.Config) error {
	terminals, err := getTerminals(config)
	if err != nil {
		return err
	}

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(ncurses),
		Args: []string{
			"--with-static",
			"--with-termlib", // generate separate terminfo library
			"--with-xterm-new", // specify if xterm terminfo should be new version
			"--with-fallbacks=" + strings.Join(terminals, ","),
			"--with-ticlib", // generate separate ticlib library (incompatible with termcap)
			"--without-ada",
			"--without-debug",
			"--without-getcap-cache",
			"--without-libtool",
			"--without-progs",
			"--without-shared",
			"--without-tests",
			"--without-termpath",
			"--enable-const", // compile with extra/non-standard const
			// "--enable-getcap", // fast termcap load, no xrefs to terminfo
			"--enable-ext-colors", // compile for 256-color support
			"--enable-overwrite",
			"--enable-pc-files", // generate and install .pc files for pkg-config
			"--enable-sigwinch", // compile with SIGWINCH handler
			// "--enable-termcap", // compile in termcap fallback support
			"--enable-wgetch-events", // compile with wgetch-events code
			"--enable-widec", // compile with wide-char/UTF-8 code
			"--disable-database", // Do not use terminfo, only fallbacks/termcap
			"--disable-db-install", // suppress install of terminal database
			"--with-pkg-config-libdir=" + config.InstallDir(ncurses) + "/share/pkgconfig",
		},
		CFlags: []string{
			"-fPIC",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (ncurses Ncurses) Install(config gogurt.Config) error {
	make := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	if err := make.Run(); err != nil {
		return err
	}
	if err := createNonWideSymlinks(config.LibDir(ncurses), "w.a", ".a"); err != nil {
		return err
	}
	if err := createNonWideSymlinks(config.PkgConfigShareDir(ncurses), "w.pc", ".pc"); err != nil {
		return err
	}
	return createCursesSymlinks(config.LibDir(ncurses))
}

// Many applications still expect the regular, non-wide libraries.
// We therefore create symlinks of them pointing to their wide equivalents.
func createNonWideSymlinks(directory, oldSuffix, newSuffix string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, oldSuffix) {
			 nonWideLink := strings.TrimSuffix(path, oldSuffix) + newSuffix
			if err := os.Symlink(path, nonWideLink); err != nil {
				return err
			}
		}
		return nil
	})
}

// Many applications still expect regular curses.
// We therefore create symlinks of them pointing to their curses equivalents.
// TODO: handle pkgconfig.
func createCursesSymlinks(directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		lib := filepath.Base(path)
		if strings.HasPrefix(lib, "libn") {
			cursesLink := filepath.Join(filepath.Dir(path), "lib" + strings.TrimPrefix(lib, "libn"))
			if err := os.Symlink(path, cursesLink); err != nil {
				return err
			}
		}
		return nil
	})
}

func (ncurses Ncurses) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}

// from misc/terminfo.src:
// # Entries with embedded plus signs are designed to be included through use/tc
// # capabilities, not used as standalone entries.
// Now, when we create fallback.c The '+' and '-' are replaced with '_', leading to multiple definitions.
// For now, we omit entries containing '+'
func getTerminals(config gogurt.Config) ([]string, error) {
	ncurses := Ncurses{}
	termInfoFile, err := os.Open(filepath.Join(config.BuildDir(ncurses), "misc", "terminfo.src"))
	if err != nil {
		return []string{}, err
	}
	defer termInfoFile.Close()

	scanner := bufio.NewScanner(termInfoFile)
	terminals := make([]string, 0, 2000) // There were 1667 by my count in 6.0, but 2000 is a nice number.
	for scanner.Scan() {
		line := scanner.Text()
		first, _ := utf8.DecodeRuneInString(line) // TODO: Can assume misc/terminfo.src is ASCII, it even says so.
		// We do not check for numerals; this will give us variables in fallback.c beginning with them.
		if unicode.IsLetter(first) {
			terminal := strings.TrimRight(strings.SplitN(line, "|", 2)[0], ",")

			// HACK - just check that the line doesn't have quotes ".
			if !strings.Contains(terminal, "+") && !strings.Contains(terminal, "tvi") && !strings.Contains(terminal, "wyse") {
				terminals = append(terminals, terminal)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}
	return terminals, nil
}
