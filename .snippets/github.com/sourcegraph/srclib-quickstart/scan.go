package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"sourcegraph.com/sourcegraph/srclib/unit"
)

func init() {
	_, err := flagParser.AddCommand("scan",
		// TODO: Update this description with your own.
		"scan for source files",
		"Looks for source files and writes source units to STDOUT.",
		&scanCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type ScanCmd struct{}

var scanCmd ScanCmd

func scan(scanDir string) ([]*unit.SourceUnit, error) {
	var units []*unit.SourceUnit
	var files []string

	err := filepath.Walk(scanDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Walking directory %s failed with: %s", scanDir, err)
		}
		if info.Mode().IsRegular() && isSourceFile(path) {
			relpath, err := filepath.Rel(scanDir, path)
			if err != nil {
				return fmt.Errorf("Making path %s relative to %s failed with: %s", path, scanDir, err)
			}
			files = append(files, relpath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Scanning for man pages failed with: %s", err)
	}

	// TODO: Add source units to `units`.

	return units, nil
}

func isSourceFile(path string) bool {
	// TODO: Check if `path` points to a source file of your target language.
	return false
}

// Execute is called by the command line interface when invoked with the `scan` subcommand.
func (c *ScanCmd) Execute(args []string) error {
	scanDir, err := filepath.EvalSymlinks(getCWD())
	if err != nil {
		return fmt.Errorf("Resolving the path to scan failed with: %s", err)
	}

	units, err := scan(scanDir)
	if err != nil {
		return fmt.Errorf("Scanning the path failed with: %s", err)
	}

	bytes, err := json.MarshalIndent(units, "", "  ")
	if err != nil {
		return fmt.Errorf("Marshalling source units failed with: %s, units: %s", err, units)
	}

	if _, err := os.Stdout.Write(bytes); err != nil {
		return fmt.Errorf("Writing output failed with: %s", err)
	}

	return nil
}
