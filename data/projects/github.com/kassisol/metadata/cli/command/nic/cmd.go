package nic

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nic",
		Short: "NIC",
		Long:  dataDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newAddCommand(),
		newListCommand(),
		newRemoveCommand(),
		newUpdateCommand(),
	)

	return cmd
}

var dataDescription = `
The **metadata nic** command has subcommands for managing Metadata datas.

To see help for a subcommand, use:

    metadata nic [command] --help

For full details on using metadata nic visit metadata's github repository.

`
