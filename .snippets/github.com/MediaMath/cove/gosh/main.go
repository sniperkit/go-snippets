package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MediaMath/cove"
	"github.com/MediaMath/cove/cmd"
)

//gosh - Get Over SsH - gets go packages via specific uri's.
//the arguments are of the form <package_name>,<uri> - github.com/MediaMath/cove,git@github.com:MediaMath/cove
func main() {
	overwrite := flag.Bool("f", false, "Overwrite existing packages")
	flag.Parse()

	goshMap, parseErr := getMap(flag.Args())
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	if gosh(*overwrite, goshMap) {
		os.Exit(1)
	}
}

type Location struct {
	GithubUrl string
	Path      string
}

type packToLocation map[cove.Package]*Location

func impliedGithubRepo(pack cove.Package) (*Location, error) {
	components := strings.Split(string(pack), "/")

	var owner, repo string
	if len(components) < 2 {
		return nil, fmt.Errorf("Cannot get implied github repo from %v", pack)
	} else if len(components) == 2 {

		owner = components[0]
		repo = components[1]
	} else {
		if strings.Contains(components[0], ".") && components[0] != "github.com" {
			return nil, fmt.Errorf("Only able to imply github.com repository urls: %v", pack)
		} else if strings.Contains(components[0], ".") && components[0] == "github.com" {

			owner = components[1]
			repo = components[2]
		} else {
			owner = components[0]
			repo = components[1]
		}
	}

	return &Location{fmt.Sprintf("git@github.com:%s/%s.git", owner, repo), fmt.Sprintf("github.com/%s/%s", owner, repo)}, nil
}

func parsePair(arg string) (cove.Package, *Location, error) {
	pair := strings.Split(arg, ",")
	if len(pair) == 0 || len(pair) > 2 {
		return cove.Package(""), nil, fmt.Errorf("Arguments are unparseable: %v", arg)
	}

	pack := cove.Package(pair[0])
	if len(pair) == 1 {
		repo, err := impliedGithubRepo(pack)
		if err != nil {
			return pack, nil, err
		}

		fmt.Printf("Using Github url %s for %v\n", repo.GithubUrl, pack)
		return pack, repo, nil
	} else {
		return pack, &Location{pair[1], string(pack)}, nil
	}

}

func getMap(args []string) (packToLocation, error) {
	goshMap := make(packToLocation)
	for _, arg := range args {
		pack, repo, err := parsePair(arg)
		if err != nil {
			return nil, err
		}

		goshMap[pack] = repo
	}

	return goshMap, nil
}

func gosh(overwrite bool, goshMap packToLocation) bool {
	hadError := false
	for pack, location := range goshMap {

		if overwrite || !cove.PackageExists(pack) {
			if logError(clone(location.GithubUrl, to(location))) {
				hadError = true
			}
		}
	}

	for pack, _ := range goshMap {
		if cove.PackageExists(pack) {
			fmt.Printf("go get %v\n", pack)
			if logError(cove.Get(pack)) {
				hadError = true
			}
		}
	}

	return hadError
}

func logError(err error) bool {
	if err != nil {
		log.Printf("%s\n", err)
		return true
	}

	return false
}

func clone(src string, destination string) error {
	vcsCopy := func(copyTemp string) error { return vcsClone(src, copyTemp) }
	return overwrite(destination, vcsCopy)
}

func vcsClone(src string, destination string) error {
	fmt.Printf("git clone %v %v\n", src, destination)
	return cmd.Run(exec.Command("git", "clone", "--depth=1", src, destination))
}

func to(location *Location) string {
	return filepath.Join(cove.GetFirstGoPath(), "src", location.Path)
}

func overwrite(destination string, copyTo func(string) error) error {
	copyTemp, copyTempErr := ioutil.TempDir("", "overwrite-copy")
	if copyTempErr != nil {
		return copyTempErr
	}
	defer os.RemoveAll(copyTemp)

	if copyErr := copyTo(copyTemp); copyErr != nil {
		return copyErr
	}

	oldDestinationTemp, oldTempErr := ioutil.TempDir("", "overwrite-old")
	if oldTempErr != nil {
		return oldTempErr
	}

	//if destination does not already exist create it
	os.MkdirAll(destination, 0744)

	if renameErr := os.Rename(destination, oldDestinationTemp); renameErr != nil {
		return renameErr
	}

	if replaceErr := os.Rename(copyTemp, destination); replaceErr != nil {
		if undoErr := os.Rename(oldDestinationTemp, destination); undoErr != nil {
			return fmt.Errorf("%v\n%v", replaceErr, undoErr)
		}

		return replaceErr
	}

	os.RemoveAll(oldDestinationTemp)
	return nil
}
