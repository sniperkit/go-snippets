package profile

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Profile",
		Long:  dataDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newAddCommand(),
		newListCommand(),
		newRemoveCommand(),
	)

	return cmd
}

var dataDescription = `
The **metadata profile** command has subcommands for managing Metadata datas.

To see help for a subcommand, use:

    metadata profile [command] --help

For full details on using metadata profile visit metadata's github repository.

`
