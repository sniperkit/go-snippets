package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var cmdImport = &cobra.Command{
	Use:   "import FILE|DIR MODULE",
	Short: "import existing files and directories into a given module",
	Run:   runImport,
}

func init() {
	cmdRoot.AddCommand(cmdImport)
}

func pathIsSymlink(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return isSymlink(fi)
}

func runImport(cmd *cobra.Command, args []string) {
	v("runImport\n")
	if len(args) != 2 {
		warn("usage: bundle add FILE|DIR MODULE\n")
		os.Exit(1)
	}

	src, module := args[0], args[1]
	if !filepath.IsAbs(src) {
		s, err := filepath.Abs(src)
		if err != nil {
			warn("unable to find absolute path for %v: %v\n", src, err)
		} else {
			src = s
		}
	}

	if !exists(src) {
		warn("%s: does not exist\n", src)
		os.Exit(2)
	}

	if pathIsSymlink(src) {
		warn("%s: already a symlink, not importing\n")
		os.Exit(2)
	}

	if !isSubdir(opts.Target, filepath.Dir(src)) {
		warn("%s is not under %s\n", src, opts.Target)
		os.Exit(2)
	}

	moduleDir := filepath.Join(opts.Base, module)
	rel, err := filepath.Rel(opts.Target, src)
	ok(err)
	dst := filepath.Join(moduleDir, rel)

	if exists(dst) {
		warn("%s already exists\n", dst)
		os.Exit(2)
	}

	if !exists(moduleDir) {
		v("creating directories for new module %q\n", module)
		ok(os.MkdirAll(moduleDir, 0700))
	}

	ok(os.MkdirAll(filepath.Dir(dst), 0700))
	v("moving %q to %q\n", src, dst)
	ok(os.Rename(src, dst))
	v("creating symlink to %q in %q\n", dst, filepath.Dir(src))
	link(dst, filepath.Dir(src))

	_, err = os.Stat(filepath.Join(dst, ".git"))
	if os.IsNotExist(err) {
		return
	}
	ok(err)
	v("%s seems to be a Git repository, adding it to bundles\n", dst)
	bundle := Bundle{
		Dir:    dst,
		Ref:    strings.TrimSpace(runOutput(dst, "git", "rev-parse", "--abbrev-ref", "HEAD")),
		Commit: strings.TrimSpace(runOutput(dst, "git", "rev-parse", "HEAD")),
	}
	remotes := strings.Split(strings.TrimSpace(runOutput(dst, "git", "remote")), "\n")
	var r string
	switch len(remotes) {
	case 0:
		v("no remote known, giving up")
		return
	case 1:
		r = remotes[0]
	default:
		// find 'origin', if it does not exists, take the first one
		for i := len(remotes) - 1; i >= 0; i-- {
			r = remotes[i]
			if remotes[i] == "origin" {
				break
			}
		}
	}
	bundle.Source = strings.TrimSpace(runOutput(dst, "git", "remote", "get-url", r))
	cfg := loadBundleConfig()
	cfg.Bundles = append(cfg.Bundles, bundle)
	saveBundleConfig(cfg)
	ok(os.RemoveAll(filepath.Join(dst, ".git")))
}
