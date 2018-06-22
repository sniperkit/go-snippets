package cmd

import (
	"github.com/spf13/cobra"
	"pet/config"
	petSync "pet/sync"
)

// syncCmd represents the sync command
var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove current snippet",
	Long:  `Remove current snippet local and remote`,
	RunE:  remove,
}

func remove(cmd *cobra.Command, args []string) (err error) {
	return petSync.RemoveSync(config.Conf.General.SnippetFile)
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
