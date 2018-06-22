package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"sourcegraph.com/sourcegraph/srclib"
	"sourcegraph.com/sourcegraph/srclib-docker/df"
	"sourcegraph.com/sourcegraph/srclib/unit"
)

const DockerfileUnitType = "Dockerfile"

func init() {
	_, err := parser.AddCommand("scan",
		"scan for Dockerfiles",
		"Scan the directory tree rooted at the current directory for Dockerfiles.",
		&scanCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type ScanCmd struct {
	Repo   string `long:"repo" description:"repository URI" value-name:"URI"`
	Subdir string `long:"subdir" description:"subdirectory in repository" value-name:"DIR"`
}

var scanCmd ScanCmd

func ignoreWalkError(err error) bool {
	return os.IsPermission(err) || os.IsNotExist(err)
}

func (c *ScanCmd) Execute(args []string) error {
	// We don't need any config yet, but read it anyway so that we
	// don't silently allow people to send us bad JSON.
	var config interface{}
	if err := json.NewDecoder(os.Stdin).Decode(&config); err != nil && err != io.EOF {
		return err
	}
	if err := os.Stdin.Close(); err != nil {
		return err
	}

	var dockerfiles []string
	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil && !ignoreWalkError(err) {
			return err
		}

		if info.Mode().IsDir() {
			// Don't traverse into VCS dirs.
			if isVCSDir[info.Name()] {
				return filepath.SkipDir
			}

			// Don't traverse into git submodules.
			if path != root {
				fs, err := ioutil.ReadDir(path)
				if err != nil && !ignoreWalkError(err) {
					return err
				}
				for _, f := range fs {
					if isVCSDir[f.Name()] {
						return filepath.SkipDir
					}
				}
			}
		}

		// Collect Dockerfiles.
		if info.Mode().IsRegular() && info.Name() == "Dockerfile" {
			dockerfiles = append(dockerfiles, path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if GlobalOpt.Verbose && len(dockerfiles) > 0 {
		log.Printf("Found %d Dockerfiles: %v", len(dockerfiles), dockerfiles)
	}

	units := make([]*unit.SourceUnit, len(dockerfiles))
	for i, dfpath := range dockerfiles {
		data, err := ioutil.ReadFile(dfpath)
		if err != nil {
			return err
		}

		// Parse Dockerfile.
		var baseImages []interface{}
		d, err := df.Decode(bytes.NewReader(data))
		if err == nil {
			if d.From != "" {
				baseImages = append(baseImages, d.From)
			}
		} else {
			log.Printf("Error parsing Dockerfile %q: %s. Adding without dependencies.", dfpath, err)
		}

		units[i] = &unit.SourceUnit{
			Name:         filepath.Dir(dfpath),
			Type:         DockerfileUnitType,
			Dir:          filepath.Dir(dfpath),
			Files:        []string{dfpath},
			Data:         string(data),
			Dependencies: baseImages,
			Ops:          map[string]*srclib.ToolRef{"depresolve": nil, "graph": nil},
		}
	}

	b, err := json.MarshalIndent(units, "", "  ")
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(b); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

var isVCSDir = map[string]bool{".git": true, ".hg": true}
