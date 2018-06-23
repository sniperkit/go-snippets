package cmd

import (
	"fmt"

	"github.com/bronzdoc/orbi/plan"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List plans",
	Run: func(cmd *cobra.Command, args []string) {
		for _, plan := range plan.List() {
			fmt.Printf("⚫ %s\n", plan)
		}
	},
}

func init() {
	planCmd.AddCommand(listCmd)
}
