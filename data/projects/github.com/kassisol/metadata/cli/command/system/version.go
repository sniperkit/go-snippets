package system

import (
	"github.com/kassisol/metadata/version"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the Metadata version information",
		Long:  versionDescription,
		Run: func(cmd *cobra.Command, args []string) {
			info := version.New()
			info.ShowVersion()
		},
	}

	return cmd
}

var versionDescription = `
All software has versions. This is Metadata's

`
