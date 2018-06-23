package medkit

import (
	"github.com/spf13/cobra"
)

// installCmd represents the install command.
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install various resources",
	Long: `
The install command will install various resources based on the sub-command.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	installCmd.AddCommand(installDotfilesCmd)
}
