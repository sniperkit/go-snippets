package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

type Lua struct{}

func (lua Lua) Name() string {
	return "lua"
}

func (lua Lua) URL(version string) string {
	return fmt.Sprintf("https://www.lua.org/ftp/lua-%s.tar.gz", version)
}

func (lua Lua) Build(config gogurt.Config) error {
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Args: []string{
			"linux",
			"MYCFLAGS=-I" + config.IncludeDir(ReadLine{}),
			"MYLDFLAGS=-L" + config.LibDir(ReadLine{}) + " -L" + config.LibDir(Ncurses{}),
			"MYLIBS=-ltinfow",
		},
	}.Cmd()
	return make.Run()
}

func (lua Lua) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{
			"install",
			"INSTALL_TOP=" + config.InstallDir(lua),
		},
	}.Cmd()
	return makeInstall.Run()
}

func (lua Lua) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Ncurses{},
		ReadLine{},
	}
}
