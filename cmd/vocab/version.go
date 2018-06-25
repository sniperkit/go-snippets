package main

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
)

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

var versionCmd = &cobra.Command{
	Use: "version",

	Short: "Prints the version of the program.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", progVersion)
	},
}

func init() {
	progVersion.Build = []string{
		buildVersion,
	}
}
