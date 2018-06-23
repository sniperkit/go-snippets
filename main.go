package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:           "dab",
	Short:         "manage dotfiles and bundles",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var opts = struct {
	DryRun  bool
	Verbose bool
	Target  string
	Base    string
}{}

func findBasedir() string {
	exe := os.Args[0]

	dir, err := filepath.Abs(filepath.Dir(exe))
	ok(err)

	fi, err := os.Lstat(exe)
	ok(err)

	if isSymlink(fi) {
		dir = readlink(exe)
	}

	for {
		if filepath.Dir(dir) == dir {
			warn("unable to find bundles.json\n")
			os.Exit(1)
		}

		filename := filepath.Join(dir, "bundles.json")
		fi, err = os.Stat(filename)
		if err == nil {
			v("found basedir: %v\n", dir)
			return dir
		}

		if err != nil && os.IsNotExist(err) {
			dir = filepath.Dir(dir)
			continue
		}

		ok(err)
	}
}

func init() {
	defaultTarget := os.Getenv("TARGET")
	if defaultTarget == "" {
		defaultTarget = os.Getenv("HOME")
	}

	fs := cmdRoot.PersistentFlags()
	fs.StringVar(&opts.Target, "target", defaultTarget, "set target directory")
	fs.StringVar(&opts.Base, "base", "", "set base directory")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "only print actions, do not execute them")
	fs.BoolVar(&opts.Verbose, "verbose", false, "be verbose")
}

func ok(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "error: %+v\n", err)
	os.Exit(1)
}

func v(s string, args ...interface{}) {
	if !opts.Verbose {
		return
	}
	fmt.Printf(s, args...)
}

func msg(s string, args ...interface{}) {
	fmt.Printf(s, args...)
}

func warn(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
}

func main() {
	if opts.Base == "" {
		opts.Base = findBasedir()
	}

	if opts.Base == "" {
		fmt.Fprintf(os.Stderr, "unable to find basedir\n")
		os.Exit(2)
	}

	err := cmdRoot.Execute()
	var exitCode int
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
