package packages

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/alexandrecarlton/gogurt"
)

type LLVM struct{}
type Clang struct {}
type ClangToolsExtra struct {}
type CompilerRT struct {}
type LLDB struct{}

func (llvm LLVM) Name() string {
	return "llvm"
}

func (llvm LLVM) URL(version string) string {
	return fmt.Sprintf("https://releases.llvm.org/%s/llvm-%s.src.tar.xz", version, version)
}

func (clang Clang) Name() string {
	return "clang"
}

func (clang Clang) URL(version string) string {
	return fmt.Sprintf("https://releases.llvm.org/%s/cfe-%s.src.tar.xz", version, version)
}

func (clangToolsExtra ClangToolsExtra) Name() string {
	return "clang-tools-extra"
}

func (clangToolsExtra ClangToolsExtra) URL(version string) string {
	return fmt.Sprintf("https://releases.llvm.org/%s/clang-tools-extra-%s.src.tar.xz", version, version)
}

func (compilerRT CompilerRT) Name() string {
	return "compiler-rt"
}

func (compilerRT CompilerRT) URL(version string) string {
	return fmt.Sprintf("https://releases.llvm.org/%s/compiler-rt-%s.src.tar.xz", version, version)
}

func(lldb LLDB) Name() string {
	return "lldb"
}

func (lldb LLDB) URL(version string) string {
	return fmt.Sprintf("http://releases.llvm.org/%s/lldb-%s.src.tar.xz", version, version)
}

func (llvm LLVM) Build(config gogurt.Config) error {

	clang := Clang{}
	clangDst := filepath.Join(config.BuildDir(llvm), "tools", "clang")
	if err := downloadArchiveTo(clang, clangDst, config); err != nil {
		return err
	}

	compilerRT := CompilerRT{}
	compilerRTDst := filepath.Join(config.BuildDir(llvm), "projects", "compiler-rt")
	if err := downloadArchiveTo(compilerRT, compilerRTDst, config); err != nil {
		return err
	}

	clangToolsExtra := ClangToolsExtra{}
	clangToolsExtraDst := filepath.Join(config.BuildDir(llvm), "tools", "clang", "tools", "extra")
	if err := downloadArchiveTo(clangToolsExtra, clangToolsExtraDst, config); err != nil {
		return err
	}

	// TODO: Implement LLDB
	// lldb := LLDB{}
	// lldbDst := filepath.Join(config.BuildDir(llvm), "tools", "lldb")
	// if err := downloadArchiveTo(lldb, lldbDst, config); err != nil {
	// 	return err
	// }

	buildDir := filepath.Join(config.BuildDir(llvm), "build")
	os.Mkdir(buildDir, 0755)

	cmake := gogurt.CMakeCmd{
		Path: filepath.Join(config.BinDir(CMake{}), "cmake"),
		Prefix: config.InstallDir(llvm),
		SourceDir: config.BuildDir(llvm),
		BuildDir: buildDir,
		CacheEntries: map[string]string{
			"CMAKE_BUILD_TYPE": "Release",
			"BUILD_SHARED_LIBS": "OFF",
			"LIBCLANG_BUILD_STATIC": "ON",
			// "LLDB_DISABLE_PYTHON": "ON", // TODO: Enable Python (which I imagine involves including openssl, readline, etc... :( )
		},
		PathPrefix: []string{
			config.InstallDir(Ncurses{}),
		},
		Paths: []string{
			config.BinDir(CMake{}),
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

func downloadArchiveTo(archive gogurt.SourceArchive, dst string, config gogurt.Config) error {
	version := config.PackageVersions[LLVM{}.Name()]
	archiveFilename := filepath.Join(config.CacheDir(LLVM{}), filepath.Base(archive.URL(version)))

	if _, err := os.Stat(archiveFilename); err == nil {
		log.Printf("File '%s' already exists, not downloading a new copy.", archiveFilename)
	} else if err := gogurt.Download(archive.URL(version), archiveFilename); err != nil {
		log.Fatalf("Could not download url '%s' to file '%s': %s\n", archive.URL(version), archiveFilename, err.Error())
	}
	if err := gogurt.DecompressSourceArchive(archiveFilename, dst); err != nil {
		log.Fatalf("Could not extract '%s' to '%s'", archiveFilename, dst)
	}
	return nil
}

func (llvm LLVM) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{
			"install",
		},
		Dir: filepath.Join(config.BuildDir(llvm), "build"),
	}.Cmd()
	return makeInstall.Run()
}

func (llvm LLVM) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		CMake{},
		// Ncurses{},
	}
}
