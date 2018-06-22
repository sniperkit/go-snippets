package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/bronzdoc/orbi/plan"
	"github.com/bronzdoc/orbi/vars"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateVars string
var templateVarsPath string

var execCmd = &cobra.Command{
	Use:   "exec PLAN",
	Short: "Executes a plan",
	Run: func(cmd *cobra.Command, args []string) {
		var planName string

		if len(args) > 0 {
			planName = args[0]
		} else {
			log.Fatal("orbi exec expects a plan name, see orbi exec --help")
		}

		if templateVarsPath != "" {
			if varsPathExists(templateVarsPath) {
				fileContent, err := ioutil.ReadFile(templateVarsPath)
				if err != nil {
					log.Fatalf("vars-file: %s", err)
				}

				templateVars = string(fileContent)
			} else {
				log.Fatalf("vars-file: couldn't find %s", templateVarsPath)
			}
		}

		vars, err := vars.Parse(templateVars)
		if err != nil {
			log.Fatal(err)
		}

		// Dinamically build plan templates path
		viper.Set("TemplatesPath", fmt.Sprintf(
			"%s/%s/%s", viper.GetString("PlansPath"), planName, viper.GetString("TemplatesDir"),
		))

		options := map[string]interface{}{
			"vars": vars,
		}

		plan := plan.Factory(planName, options)

		if err := plan.Execute(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&templateVars, "vars", "v", "", "template vars KEY=VALUE")
	execCmd.Flags().StringVarP(&templateVarsPath, "vars-file", "", "", "template vars file path containing KEY=VALUE")
}

func varsPathExists(varsPath string) bool {
	_, err := os.Stat(varsPath)
	return err == nil
}
