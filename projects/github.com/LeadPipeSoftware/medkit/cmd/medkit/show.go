package medkit

import (
	"github.com/spf13/cobra"
)

// showCmd represents the show command.
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show various resources",
	Long: `
The show command will display various resources based on the sub-command.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	showCmd.AddCommand(showConfigCmd)
	showCmd.AddCommand(showDotfilesCmd)
}
