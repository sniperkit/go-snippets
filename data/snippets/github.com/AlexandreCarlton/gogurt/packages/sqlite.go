package packages

import (
	"fmt"
	"strings"
	"strconv"
	"time"
	"github.com/alexandrecarlton/gogurt"
)

type SQLite struct{}

func (sqlite SQLite) Name() string {
	return "sqlite"
}

type sqliteVersion struct {
	major int;
	minor int;
	patch int;
}

func newSqliteVersion(version string) *sqliteVersion {
	v := strings.Split(version, ".")
	major, _ := strconv.Atoi(v[0])
	minor, _ := strconv.Atoi(v[1])
	patch, _ := strconv.Atoi(v[2])
	return &sqliteVersion{major, minor, patch}
}

func (v sqliteVersion) zeroed() string {
	return fmt.Sprintf(
		"%d%d%s%d%s",
		v.major,
		v.minor,
		strings.Repeat("0", 3 - len(strconv.Itoa(v.minor))),
		v.patch,
		strings.Repeat("0", 3 - len(strconv.Itoa(v.patch))))
}

// TODO: Implement comparison function
// so if v < ... => 2017
// and default, return current year.
func (v sqliteVersion) year() int {
	switch v {
	case sqliteVersion{3, 20, 1}:
		return 2017
	default:
		return time.Now().Year()
	}
}

func (sqlite SQLite) URL(version string) string {
	sqliteVersion := newSqliteVersion(version)
	return fmt.Sprintf("https://sqlite.org/%d/sqlite-autoconf-%s.tar.gz", sqliteVersion.year(), sqliteVersion.zeroed())
}

func (sqlite SQLite) Build(config gogurt.Config) error {
	// Termcap has been succeeded by terminfo.
	// TODO: Create symlink from termcap to terminfo ?
	gogurt.ReplaceInFile("configure", "termcap", "tinfow")

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(sqlite),
		Args: []string{
			"--disable-shared",
			"--enable-static",
			"--enable-readline",
		},
		CFlags: []string{
			"-I" + config.IncludeDir(ReadLine{}),
			"-I" + config.IncludeDir(Ncurses{}),
		},
		CppFlags: []string{
			"-I" + config.IncludeDir(ReadLine{}),
			"-I" + config.IncludeDir(Ncurses{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(ReadLine{}),
			"-L" + config.LibDir(Ncurses{}),
		},
		Libs: []string{
			"-lreadline",
			"-lncursesw",
			"-ltinfow",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores}.Cmd()
	return make.Run()
}

func (sqlite SQLite) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (sqlite SQLite) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Ncurses{},
		ReadLine{},
	}
}
