package main

import "github.com/blang/semver"

var (
	progVersion = semver.Version{
		Major: 0,
		Minor: 1,
		Patch: 0,
		Pre: []semver.PRVersion{
			{VersionStr: "beta"},
		},
	}

	buildVersion string
)

func init() {
	progVersion.Build = []string{
		buildVersion,
	}
}
