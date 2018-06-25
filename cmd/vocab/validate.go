package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validateCmd = &cobra.Command{
	Use: "validate <path>",

	Short: "Evaluates various aspects of a concept.csv for integrity.",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}

		_, err := os.Open(args[0])

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	flags := validateCmd.Flags()

	flags.String("vocab", "", "The accompanying vocabularies file to check the integrity of the `vocabulary_id` field.")

	viper.BindPFlag("validate.vocab", flags.Lookup("vocab"))
}
