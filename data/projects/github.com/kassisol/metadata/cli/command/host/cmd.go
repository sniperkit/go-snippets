package host

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "Host",
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
The **metadata host** command has subcommands for managing Hosts.

To see help for a subcommand, use:

    metadata host [command] --help

For full details on using metadata host visit metadata's github repository.

`
