package cmd

import (
	"fmt"
	"log"

	"github.com/bronzdoc/orbi/plan"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit PLAN_NAME",
	Short: "Edit a plan definition, $EDITOR will be used to open the file.",
	Run: func(cmd *cobra.Command, args []string) {
		var planName string

		if len(args) > 0 {
			planName = args[0]
		} else {
			err := fmt.Errorf("orbi plan edit expects a plan name, see orbi plan edit --help")
			log.Fatal(err)
		}

		if err := plan.Edit(planName); err != nil {
			log.Fatalf("edit command failed: %v\n", err)
		}
	},
}

func init() {
	planCmd.AddCommand(editCmd)
}
