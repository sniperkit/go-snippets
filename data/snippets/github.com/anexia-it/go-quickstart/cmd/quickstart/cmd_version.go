package main

import (
	quickstart "github.com/anexia-it/go-quickstart"
	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Display version information and exit",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("quickstart v%s", quickstart.VersionString())
	},
}

func init() {
	cmdRoot.AddCommand(cmdVersion)
}
