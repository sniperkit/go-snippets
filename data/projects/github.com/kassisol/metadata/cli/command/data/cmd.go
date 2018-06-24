package data

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "Data",
		Long:  dataDescription,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	cmd.AddCommand(
		newAddCommand(),
		newUpdateCommand(),
		newListCommand(),
		newMemberCommand(),
		newRemoveCommand(),
	)

	return cmd
}

var dataDescription = `
The **metadata data** command has subcommands for managing Metadata datas.

To see help for a subcommand, use:

    metadata data [command] --help

For full details on using metadata data visit metadata's github repository.

`
