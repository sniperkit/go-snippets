package main

import "github.com/spf13/cobra"

var mainCmd = &cobra.Command{
	Use: "pedsnet-vocab",

	Short: "Commands for validating the integrity of PEDSnet vocabularies.",

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	mainCmd.AddCommand(versionCmd)
	mainCmd.AddCommand(summaryCmd)
	mainCmd.AddCommand(compareCmd)

	mainCmd.Execute()
}
