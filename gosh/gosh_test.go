package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MediaMath/cove"
)

func TestGetGoshMapFromArgs(t *testing.T) {
	goshMap, _ := getMap([]string{"foo/bar,git@github.com/MediaMath/foo.git", "salt,git@github.com/MediaMath/salt.git"})

	if len(goshMap) != 2 {
		t.Errorf("Map does not have appropriate number of items |%v|", goshMap)
	}

	if goshMap["foo/bar"].GithubUrl != "git@github.com/MediaMath/foo.git" {
		t.Errorf("Incorrect url for foo/bar", goshMap["foo/bar"])
	}

	if goshMap["salt"].GithubUrl != "git@github.com/MediaMath/salt.git" {
		t.Errorf("Incorrect url for salt", goshMap["salt"])
	}

}

func TestGetGoshMapMixOfImpliedAndExplicit(t *testing.T) {
	goshMap, _ := getMap([]string{"foo/bar/baz", "salt,git@github.com/MediaMath/salt.git"})

	if len(goshMap) != 2 {
		t.Errorf("Map does not have appropriate number of items |%v|", goshMap)
	}

	if goshMap["foo/bar/baz"] == nil {
		t.Errorf("Did not get values |%v|", goshMap)
	}

	if goshMap["foo/bar/baz"].GithubUrl != "git@github.com:foo/bar.git" {
		t.Errorf("Incorrect url for foo/bar/baz |%v|", goshMap["foo/bar/baz"].GithubUrl)
	}

	if goshMap["salt"].GithubUrl != "git@github.com/MediaMath/salt.git" {
		t.Errorf("Incorrect url for salt", goshMap["salt"])
	}
}

func TestImpliedGithubRepo(t *testing.T) {
	implied, _ := impliedGithubRepo("github.com/MediaMath/foo")

	if implied.GithubUrl != "git@github.com:MediaMath/foo.git" {
		t.Errorf("Got:%v", implied)
	}

}

func TestImpliedNoHostProducesGithubUrl(t *testing.T) {
	implied, _ := impliedGithubRepo("bar/foo")

	if implied.GithubUrl != "git@github.com:bar/foo.git" {
		t.Errorf("Got:%v", implied)
	}
}

func TestImpliedSubpackage(t *testing.T) {
	impliedWithHost, _ := impliedGithubRepo("github.com/bar/foo/baz")

	if impliedWithHost.GithubUrl != "git@github.com:bar/foo.git" {
		t.Errorf("Got:%v", impliedWithHost)
	}

	impliedNoHost, _ := impliedGithubRepo("bar/foo/baz")

	if impliedNoHost.GithubUrl != "git@github.com:bar/foo.git" {
		t.Errorf("Got:%v", impliedNoHost)
	}

	impliedDeep, _ := impliedGithubRepo("bar/foo/baz/goose/gander")

	if impliedDeep.GithubUrl != "git@github.com:bar/foo.git" {
		t.Errorf("Got:%v", impliedDeep)
	}
}

func TestNonGithubHostFails(t *testing.T) {
	if _, norepo := impliedGithubRepo("foo.com/bar/baz"); norepo == nil {
		t.Errorf("Didnt get error on norepo")
	}
}

func TestSingleParamImplicationFails(t *testing.T) {
	if _, singleParam := impliedGithubRepo("foo"); singleParam == nil {
		t.Errorf("Didnt get error on only host")
	}
}

func TestToSingleAndMultiGoPath(t *testing.T) {
	oldPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldPath)

	os.Setenv("GOPATH", "/foo/bar:/local/foo/bar")
	if path := to(&Location{"github.com", "/projecta"}); path != "/foo/bar/src/projecta" {
		t.Errorf("Didn't get the correct to path from to for a multi-path: %s", path)
	}

	os.Setenv("GOPATH", "/single/path")
	if path := to(&Location{"github.com", "/projecta"}); path != "/single/path/src/projecta" {
		t.Errorf("Didn't get the correct to path from to for a single path: %s", path)
	}

}

func TestGoshInstallsMainPackagesLikeGoGetExplicit(t *testing.T) {
	geckoPath := filepath.Join(cove.GetFirstGoPath(), "bin", "gecko")
	os.Remove(geckoPath)

	if _, err := os.Stat(geckoPath); !os.IsNotExist(err) {
		t.Errorf("gecko wasn't removed")
	}

	goshMap, parseErr := getMap([]string{"github.com/MediaMath/cove/gecko"})
	if parseErr != nil {
		t.Errorf("%v", parseErr)
	}

	if gosh(false, goshMap) {
		t.Errorf("Gosh failed")
	}

	if _, err := os.Stat(geckoPath); os.IsNotExist(err) {
		t.Errorf("gecko wasn't installed")
	}
}

func TestGoshInstallsMainPackagesLikeGoGetWildCard(t *testing.T) {
	geckoPath := filepath.Join(cove.GetFirstGoPath(), "bin", "gecko")
	os.Remove(geckoPath)

	if _, err := os.Stat(geckoPath); !os.IsNotExist(err) {
		t.Errorf("gecko wasn't removed")
	}

	goshMap, parseErr := getMap([]string{"github.com/MediaMath/cove/..."})
	if parseErr != nil {
		t.Errorf("%v", parseErr)
	}

	if gosh(false, goshMap) {
		t.Errorf("Gosh failed")
	}

	if _, err := os.Stat(geckoPath); os.IsNotExist(err) {
		t.Errorf("gecko wasn't installed")
	}
}
