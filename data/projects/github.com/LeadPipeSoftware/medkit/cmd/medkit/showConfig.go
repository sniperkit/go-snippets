package medkit

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showConfigCmd represents the showConfig command
var showConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "display the MEDKIT configuration",
	Long: `
Display the current MEDKIT configuration. This command takes into account
the config file, environment variables, and command line flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		showConfig()
	},
}

func init() {
}

// showConfig displays the program's configuration settings.
func showConfig() {
	configKeys := []string{HomeDirectory, Bundles, DotFilesDirectory, BackupExtension}

	fmt.Println("MEDKIT configuration settings:")
	fmt.Println()

	for _, configKey := range configKeys {
		var configValue = viper.GetString(configKey)

		if configValue == "" {
			fmt.Printf(" %s: (no value set)\n", configKey)
		} else {
			fmt.Printf(" %s: %s\n", configKey, configValue)
		}
	}
}