package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var cmdInstall = &cobra.Command{
	Use:   "install",
	Short: "install modules",
	Run:   runInstall,
}

func init() {
	cmdRoot.AddCommand(cmdInstall)
}

func readdir(dir string) []os.FileInfo {
	f, err := os.Open(dir)
	ok(err)
	fis, err := f.Readdir(-1)
	ok(err)
	ok(f.Close())
	return fis
}

func subdirs(dir string) (dirs []string) {
	for _, fi := range readdir(dir) {
		if !fi.IsDir() {
			continue
		}

		name := fi.Name()
		if name == ".git" || name == "manage" || name == "old" {
			continue
		}

		dirs = append(dirs, name)
	}

	return dirs
}

func isSymlink(fi os.FileInfo) bool {
	return (fi.Mode() & os.ModeSymlink) != 0
}

func isDir(name string) bool {
	fi, err := os.Stat(name)
	ok(err)
	return fi.Mode().IsDir()
}

// exists checks if the target name points to exists.
func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// linkDirContents links all entries in srcdir to dstdir.
func linkDirContents(srcdir, dstdir string) {
	for _, entry := range readdir(srcdir) {
		link(filepath.Join(srcdir, entry.Name()), dstdir)
	}
}

// resolveConflict resolves a conflict where two sources (src1 exists, src2 not
// yet) are to be linked to the name dst.
func resolveConflict(src1, src2, dst string) {
	v("resolve conflict for %v:\n  %v\n  %v\n", dst, src1, src2)
	if !isDir(src1) {
		warn("unable to resolve conflict for %v: source %v is not a directory\n", dst, src1)
		return
	}
	if !isDir(src2) {
		warn("unable to resolve conflict for %v: source %v is not a directory\n", dst, src2)
	}

	// first, remove the symlink
	ok(os.Remove(dst))

	// create a new directory
	ok(os.Mkdir(dst, 0755))

	// then link the file separately
	linkDirContents(src1, dst)
	linkDirContents(src2, dst)

	v("conflict for %v resolved\n", dst)
}

// link creates a symlink to the file/dir src in targetdir.
func link(src, targetdir string) {
	base := filepath.Base(src)
	if len(base) == 0 {
		warn("invalid item name %v, skipping\n", src)
		return
	}

	dst := filepath.Join(targetdir, base)

	fi, err := os.Lstat(dst)
	if err == nil {
		if !isSymlink(fi) && !fi.IsDir() {
			warn("skipping already existing item %v (neither symlink nor directory)\n", dst)
			return
		}

		if fi.IsDir() {
			v("%v exists, descending into %v\n", dst, src)
			linkDirContents(src, dst)
			return
		}

		linkTarget := readlink(dst)
		switch {
		case linkTarget == src:
			// nothing to do, dst already points to src
			return

		case isSubdir(opts.Base, readlink(dst)):
			resolveConflict(linkTarget, src, dst)
			return

		case exists(dst) && isDir(dst):
			v("descending into %v\n", linkTarget)
			linkDirContents(src, dst)
			return

		case exists(dst):
			v("symlink %v already exists, skipping\n", dst)
			return

		default:
			v("removing dangling symlink\n")
			if !opts.DryRun {
				ok(os.Remove(dst))
			}
		}
	}

	v("link %v -> %v\n", src, dst)
	if !opts.DryRun {
		ok(os.Symlink(src, dst))
	}
}

func install(moduleName, targetdir string) {
	if moduleName == "" {
		return
	}

	modulePath := filepath.Join(opts.Base, moduleName)

	v("install %v\n", moduleName)
	if !exists(modulePath) {
		warn("module %v does not exist, skipping\n", moduleName)
		return
	}

	linkDirContents(modulePath, targetdir)
}

func cleanupBrokenLinks() {
	walkOurSymlinks(opts.Base, opts.Target, func(filename, target string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// if the target exists, do nothing
		if exists(target) {
			return nil
		}

		v("removing broken symlink %v\n", filename)
		return os.Remove(filename)
	})
}

func runInstall(cmd *cobra.Command, args []string) {
	hostname, err := os.Hostname()
	ok(err)

	cleanupBrokenLinks()

	if len(args) == 0 {
		// install only the base module
		install("base", opts.Target)

		if exists(filepath.Join(opts.Base, "base_"+hostname)) {
			install("base_"+hostname, opts.Target)
		}
		return
	}

	if args[0] == "all" {
		// install all modules
		for _, module := range subdirs(opts.Base) {
			install(module, opts.Target)
			if exists(filepath.Join(opts.Base, module+"_"+hostname)) {
				install(module+"_"+hostname, opts.Target)
			}
		}
		return
	}

	// otherwise install the listed modules
	for _, module := range args {
		// install all modules
		install(module, opts.Target)
		if exists(filepath.Join(opts.Base, module+"_"+hostname)) {
			install(module+"_"+hostname, opts.Target)
		}
	}
}

func readlink(file string) string {
	target, err := os.Readlink(file)
	ok(err)
	return target
}

func isSubdir(dir, subdir string) bool {
	for {
		if len(subdir) <= len(dir) {
			return dir == subdir
		}

		subdir = filepath.Dir(subdir)
		if subdir == "." {
			return false
		}
	}
}
