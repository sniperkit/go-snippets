package packages

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type NeoVim struct{}

func (neovim NeoVim) Name() string {
	return "neovim"
}

func (neovim NeoVim) URL(version string) string {
	return fmt.Sprintf("https://github.com/neovim/neovim/archive/v%s.tar.gz", version)
}

func (neovim NeoVim) Build(config gogurt.Config) error {
	buildDepsDir := filepath.Join(config.BuildDir(neovim), "build-deps")
	os.Mkdir(buildDepsDir, 0755)

	cmakeDeps := gogurt.CMakeCmd{
		SourceDir: filepath.Join(config.BuildDir(neovim), "third-party"),
		BuildDir: buildDepsDir,
	}.Cmd()
	if err := cmakeDeps.Run(); err != nil {
		return err
	}
	makeDeps := gogurt.MakeCmd{Dir: buildDepsDir}.Cmd()
	if err := makeDeps.Run(); err != nil {
		return err
	}

	buildDir := filepath.Join(config.BuildDir(neovim), "build")
	os.Mkdir(buildDir, 0755)

	cmake := gogurt.CMakeCmd{
		Prefix: config.InstallDir(neovim),
		SourceDir: config.BuildDir(neovim),
		BuildDir: buildDir,
		CacheEntries: map[string]string{
			"CMAKE_BUILD_TYPE": "Release",
			"BUILD_SHARED_LIBS": "OFF",
			"DEPS_PREFIX": filepath.Join(buildDepsDir, "usr"),
		},
	}.Cmd()
	if err := cmake.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Dir: buildDir,
	}.Cmd()
	return make.Run()
}

func (neovim NeoVim) Install(config gogurt.Config) error {
	buildDir := filepath.Join(config.BuildDir(neovim), "build")
	makeInstall := gogurt.MakeCmd{
		Args: []string{"install"},
		Dir: buildDir,
	}.Cmd()
	return makeInstall.Run()
}

func (neovim NeoVim) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}

