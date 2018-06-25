package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var cmdStatus = &cobra.Command{
	Use:   "status",
	Short: "displays the status of all modules",
	Run: func(*cobra.Command, []string) {
		state := currentState()
		keys := make([]string, 0, len(state))
		for name := range state {
			keys = append(keys, name)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, name := range keys {
			installed := state[name]
			if !installed {
				msg("[ ---- ] %v\n", name)
			} else {
				msg("[ INST ] %v\n", name)
			}
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdStatus)
}

// State contains the state whether a module is installed or not.
type State map[string]bool

func firstdir(dir string) string {
	dirs := strings.Split(dir, string(filepath.Separator))
	return dirs[0]
}

func installedModules() State {
	state := make(State)
	walkOurSymlinks(opts.Base, opts.Target, func(filename, target string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(opts.Base, target)
		if err != nil {
			return err
		}

		module := firstdir(rel)

		state[module] = true
		return nil
	})

	return state
}

func allModules() (modules []string) {
	for _, dir := range subdirs(opts.Base) {
		module := filepath.Base(dir)
		modules = append(modules, module)
	}
	return modules
}

func currentState() State {
	state := installedModules()
	for _, module := range allModules() {
		if _, ok := state[module]; !ok {
			state[module] = false
		}
	}
	return state
}
