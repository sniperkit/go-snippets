package packages

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/alexandrecarlton/gogurt"
)

type CMake struct{}

func (cmake CMake) Name() string {
	return "cmake"
}

func (cmake CMake) URL(version string) string {
	majorVersion := strings.Join(strings.Split(version, ".")[:2], ".")
	return fmt.Sprintf("https://cmake.org/files/v%s/cmake-%s.tar.gz", majorVersion, version)
}

func (cmake CMake) Build(config gogurt.Config) error {

	bootstrap := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(cmake),
		Args: []string{
			// TODO: Use system libraries instead.
			// We may need to use an existing cmake for this,
			// and fill PathPrefix with these libraries, and set:
			// - CMAKE_EXE_LINKER_FLAGS
			// - CMAKE_USE_SYSTEM_<lib>
			"--parallel=" + strconv.Itoa(int(config.NumCores)),
			"--no-system-bzip2",
			"--no-system-curl",
			"--no-system-expat",
			"--no-system-zlib",
			"--no-system-liblzma",
			"--no-qt-gui",
		},
	}.Cmd()
	bootstrap.Path = "./bootstrap"
	if err := bootstrap.Run(); err != nil {
		return err
	}

	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()
	return make.Run()
}

func (cmake CMake) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{"install"},
	}.Cmd()
	return makeInstall.Run()
}

func (cmake CMake) Dependencies() []gogurt.Package {
	return []gogurt.Package{}
}
