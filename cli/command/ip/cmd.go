package ip

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ip",
		Short: "IP Address",
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
The **metadata ip** command has subcommands for managing Metadata IP addresses.

To see help for a subcommand, use:

    metadata ip [command] --help

For full details on using metadata ip visit metadata's github repository.

`
