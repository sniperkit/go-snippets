package main

import (
	"encoding/json"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"path/filepath"
	"sourcegraph.com/sourcegraph/srclib/unit"
	"strings"
)

// YAML does not require special configuration, so all srcfileconfig logic
// is not included in this file
type ScanCmd struct{}

var (
	parser  = flags.NewNamedParser("srclib-yaml", flags.Default)
	scanCmd = ScanCmd{}
)

func init() {
	_, err := parser.AddCommand("scan",
		"scan for YAML files",
		"Scan the directory tree rooted at the current directory for YAML files.",
		&scanCmd)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}

func (c *ScanCmd) Execute(args []string) error {
	CWD, err := os.Getwd()
	if err != nil {
		return err
	}
	units, err := scan(CWD)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(units, "", "  ")
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(b); err != nil {
		return err
	}
	return nil
}

func scan(scanDir string) ([]*unit.SourceUnit, error) {
	u := unit.SourceUnit{}
	u.Key.Name = filepath.Base(scanDir)
	u.Key.Type = "yaml"
	u.Key.Repo = strings.Join([]string{filepath.Base(filepath.Dir(scanDir)), filepath.Base(scanDir)}, "/")
	units := []*unit.SourceUnit{&u}

	if err := filepath.Walk(scanDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if isYAMLFile(path) {
			rp, err := filepath.Rel(scanDir, path)
			if err != nil {
				return err
			}
			u.Files = append(u.Files, filepath.ToSlash(rp))
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return units, nil
}

func isYAMLFile(filename string) bool {
	return filepath.Ext(filename) == ".yml"
}
