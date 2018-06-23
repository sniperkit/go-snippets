package cmd

import (
	"fmt"
	"log"

	"github.com/bronzdoc/orbi/plan"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get REPO_URL",
	Short: "Get a plan from a git repository url",
	Run: func(cmd *cobra.Command, args []string) {
		var repoURL string

		if len(args) > 0 {
			repoURL = args[0]
		} else {
			err := fmt.Errorf("orbi plan get expects a plan repository url, see orbi plan get --help")
			log.Fatal(err)
		}

		err := plan.Get(repoURL)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	planCmd.AddCommand(getCmd)
}
