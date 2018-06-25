package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// ModuleWalkFunc is called by walkOurSymlinks for each directory found to be a module.
type ModuleWalkFunc func(filename, target string, fi os.FileInfo, err error) error

var cmdRemove = &cobra.Command{
	Use:   "remove",
	Short: "unistall modules",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = append(args, "")
		}

		for i := range args {
			args[i] = filepath.Join(opts.Base, args[i])
		}

		v("removing base dirs %v\n", args)

		walkOurSymlinks(opts.Base, opts.Target, func(filename, target string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			for _, dir := range args {
				if isSubdir(dir, target) {
					if !opts.DryRun {
						ok(os.Remove(filename))
					}
					v("removed %v\n", filename)
					break
				}
			}

			return nil
		})
	},
}

func init() {
	cmdRoot.AddCommand(cmdRemove)
}

func walkOurSymlinks(base, dir string, fn ModuleWalkFunc) {
	ok(filepath.Walk(opts.Target, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			v("error checking %v: %v\n", filename, err)
			return nil
		}

		if !isSymlink(fi) {
			return err
		}

		target := readlink(filename)
		if !isSubdir(base, target) {
			return err
		}

		return fn(filename, target, fi, err)
	}))
}
