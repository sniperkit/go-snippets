package commands

import (
	"github.com/kassisol/metadata/cli/command/data"
	"github.com/kassisol/metadata/cli/command/host"
	"github.com/kassisol/metadata/cli/command/ip"
	"github.com/kassisol/metadata/cli/command/nic"
	"github.com/kassisol/metadata/cli/command/profile"
	"github.com/kassisol/metadata/cli/command/system"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		data.NewCommand(),
		host.NewCommand(),
		ip.NewCommand(),
		nic.NewCommand(),
		profile.NewCommand(),
		system.NewInfoCommand(),
		system.NewVersionCommand(),
	)
}
