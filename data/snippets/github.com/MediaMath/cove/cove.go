// Package cove is a thin wrapper around the go toolchain
package cove

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MediaMath/cove/cmd"
)

//Package is a specific and unique go package
type Package string

//PackagePattern is a pattern, potentially with wildcards, to search for packages in.  It can correlate to 0 to many Packages
type PackagePattern string

//PackagePatternsAsStrings takes a slice of paths and returns a slice of strings
func PackagePatternsAsStrings(paths []PackagePattern) []string {
	var returns []string
	for _, path := range paths {
		returns = append(returns, string(path))
	}

	return returns
}

//PackagesAsStrings takes a slice of packages and returns a slice of strings
func PackagesAsStrings(packs []Package) []string {
	var returns []string
	for _, pack := range packs {
		returns = append(returns, string(pack))
	}

	return returns
}

//PackagesFromStrings creates a slice of packages from a slice of strings
func PackagesFromStrings(str []string) []Package {
	var returns []Package
	for _, s := range str {
		returns = append(returns, Package(s))
	}

	return returns
}

//PackagePatternsFromStrings creates a slice of paths from a slice of strings
func PackagePatternsFromStrings(str []string) []PackagePattern {
	var returns []PackagePattern
	for _, s := range str {
		returns = append(returns, PackagePattern(s))
	}

	return returns
}

//GoCmd takes the sub and args and prepares a command like 'go sub arg1 arg2...'
func GoCmd(sub string, args ...string) *exec.Cmd {
	arguments := append([]string{sub}, args...)
	return exec.Command("go", arguments...)
}

// Get runs 'go get pack'
func Get(pack Package) error {
	return cmd.Run(GoCmd("get", string(pack)))
}

// PackageExists checks to see if a given package name exists
func PackageExists(pack Package) bool {
	err := cmd.Run(GoCmd("list", string(pack)))
	return err == nil
}

// MissingDependencies returns any go packages that are missing for a given package.
func MissingDependencies(pack Package) ([]Package, error) {
	var deps []Package
	var parsed missing
	if err := PackageJSON(pack, &parsed); err != nil {
		return deps, fmt.Errorf("Could not determine missing dependencies for %v: %v", pack, err)
	}

	return missingFromParsed(&parsed)
}

func missingFromParsed(parsed *missing) ([]Package, error) {
	if !parsed.Incomplete {
		return []Package{}, nil
	}

	seen := make(map[string]bool)
	var missing []Package
	for _, errs := range parsed.DepsErrors {
		for _, imports := range errs.ImportStack {
			if _, contains := seen[imports]; !contains {
				seen[imports] = true
				missing = append(missing, Package(imports))
			}
		}
	}
	return missing, nil
}

type missing struct {
	Incomplete bool
	DepsErrors []struct {
		ImportStack []string
	}
}

// Packages gets all packages that match any of the paths.
// The package list will only contain 1 entry per package in sorted order.
// Invalid paths will generate an error, but will not stop the evaluation of the other paths.
func Packages(paths ...PackagePattern) ([]Package, error) {
	packs, err := cmd.Output(GoCmd("list", PackagePatternsAsStrings(paths)...))
	sort.Strings(packs)
	return PackagesFromStrings(packs), err
}

// PackageJSON takes a SINGLE fully qualified package import path and decodes the 'go list -json' response.
// See $GOROOT/src/cmd/go/list.go for documentation on the json output.
func PackageJSON(pack Package, v interface{}) error {
	err := cmd.PipeWith(GoCmd("list", "-json", string(pack)), func(stdout io.Reader) error {
		return json.NewDecoder(stdout).Decode(v)
	})

	if err != nil {
		if !PackageExists(pack) {
			return fmt.Errorf("Cannot get json for '%v' as it is not a go package. Original error: %v", pack, err)
		}

		return err
	}

	return nil
}

// CoverageProfile creates a cover profile file for the provided package
// The file is created in outdir.  The parameter short sets whether to run
// all tests or only the short ones. The mode specifies which style of cover mode is used.
// If a profile is able to be created its file name is returned.
func CoverageProfile(short bool, mode string, outdir string, pack Package) (string, error) {
	if direrr := os.MkdirAll(outdir, 0744); direrr != nil {
		return "", direrr
	}

	profile := getProfileFileName(outdir, pack)

	if err := cmd.Run(GoCmd("test", string(pack), fmt.Sprintf("-covermode=%s", mode), fmt.Sprintf("-coverprofile=%s", profile), getShort(short))); err != nil {
		return "", fmt.Errorf("%s:%v:%v", pack, err, profile)
	}

	if _, err := os.Stat(profile); os.IsNotExist(err) {
		return "", nil
	}

	return profile, nil
}

// CoverageReport turns the profile into a report using 'go tool cover'
func CoverageReport(profile string, outdir string) (string, error) {
	report := getReportFileName(profile, outdir)
	if err := cmd.Run(GoCmd("tool", "cover", fmt.Sprintf("-html=%s", profile), "-o", report)); err != nil {
		return "", err
	}

	return report, nil
}

func getReportFileName(profile string, outdir string) string {
	report := filepath.Base(profile)
	extension := filepath.Ext(report)
	name := report[0 : len(report)-len(extension)]
	fullPath := filepath.Join(outdir, name)
	return fmt.Sprintf("%s.html", fullPath)
}

func getProfileFileName(outdir string, pack Package) string {
	profile := strings.Replace(string(pack), "/", ".", -1)
	fullPath := filepath.Join(outdir, profile)
	return fmt.Sprintf("%s.out", fullPath)
}

func getShort(short bool) string {
	if short {
		return "-short"
	}

	return ""
}
