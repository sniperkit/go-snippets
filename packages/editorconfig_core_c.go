package packages

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type EditorConfigCoreC struct{}

func (editorconfig EditorConfigCoreC) Name() string {
	return "editorconfig-core-c"
}

func (editorconfig EditorConfigCoreC) URL(version string) string {
	return fmt.Sprintf("https://github.com/editorconfig/editorconfig-core-c/archive/v%s.tar.gz", version)
}

func (editorconfig EditorConfigCoreC) Build(config gogurt.Config) error {
	build := filepath.Join(config.BuildDir(editorconfig), "build")
	os.Mkdir(build, 0755)

	cmake := gogurt.CMakeCmd{
		Prefix: config.InstallDir(editorconfig),
		SourceDir: config.BuildDir(editorconfig),
		BuildDir: build,
		CacheEntries: map[string]string{
			// "CMAKE_BUILD_TYPE": "Release",
			"BUILD_DOCUMENTATION": "OFF",
			"PCRE_LIBRARY": filepath.Join(config.LibDir(Pcre{}), "libpcre.a"),
			"PCRE_INCLUDE_DIR": config.IncludeDir(Pcre{}),
			"PCRE_STATIC": "ON",
			"BUILD_STATICALLY_LINKED_EXE": "ON",
			"CMAKE_EXE_LINKER_FLAGS": "-pthread",
		},
	}.Cmd()
	if err := cmake.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{Jobs: config.NumCores, Args: []string{"-d"}}.Cmd()
	make.Dir = build
	return make.Run()
}

func (editorconfig EditorConfigCoreC) Install(config gogurt.Config) error {
	build := filepath.Join(config.BuildDir(editorconfig), "build")
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	makeInstall.Dir = build
	return makeInstall.Run()
}

func (editorconfig EditorConfigCoreC) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Pcre{},
		// CMake is needed, but CentOS 7 has a new enough version.
	}
}
