package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	IpfsFsPath = "IpfsFsPath"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dt",
	Short: "data together CLI",
	Long:  `command line client for data together`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		PrintErr(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qri.json)")
	RootCmd.PersistentFlags().BoolVarP(&noColor, "no-color", "c", false, "disable colorized output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home := userHomeDir()
	SetNoColor()

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// if err := os.Mkdir(filepath.Join(userHomeDir(), ".qri"), os.ModePerm); err != nil {
	//  fmt.Errorf("error creating home dir: %s\n", err.Error())
	// }

	// viper.SetConfigName("config") // name of config file (without extension)
	// // viper.AddConfigPath("$QRI_PATH")  // add QRI_PATH env var
	// viper.AddConfigPath("$HOME/.qri") // adding home directory as first search path
	// viper.AddConfigPath(".")          // adding home directory as first search path
	// viper.AutomaticEnv()              // read in environment variables that match

	ipfsFsPath := os.Getenv("IPFS_PATH")
	if ipfsFsPath == "" {
		ipfsFsPath = "$HOME/.ipfs"
	}
	ipfsFsPath = strings.Replace(ipfsFsPath, "~", home, 1)
	viper.SetDefault(IpfsFsPath, ipfsFsPath)

	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	//  // fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}
