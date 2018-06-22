package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "orbi",
	Short: `Project structure generator`,
}

var orbiPath = fmt.Sprintf("%s/.orbi", os.Getenv("HOME"))
var plansPath = fmt.Sprintf("%s/plans", orbiPath)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Create ~/.orbi directory and sub directories if missing
	if !pathExists(orbiPath) || !pathExists(plansPath) {
		os.MkdirAll(plansPath, 0776)
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("OrbiPath", orbiPath)
	viper.SetDefault("PlansPath", plansPath)
	viper.SetDefault("TemplatesDir", "templates")

	viper.SetConfigName(".config")     // name of config file (without extension)
	viper.AddConfigPath("$HOME/.orbi") // adding home directory as first search path
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
