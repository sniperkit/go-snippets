// Copyright © 2017 Thomas Biesaart <thomas.biesaart@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/ChappIO/dashfiles/tools"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

var dashHome string
var cfgFile string
var workspace string

var RootCmd = &cobra.Command{
	Use:   "dashfiles",
	Short: "Manage your dotfiles with ease",
	Long: `The dashfiles client lets you store your configurations in git
and install them with ease on many systems. This client provides a
collection of commands that allow your to easily manage those configurations.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize Config
	cobra.OnInitialize(func() {

		// Initialize viper
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {

			viper.AddConfigPath(dashHome)
			viper.SetConfigName("config")
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}

		// Initialize dash home
		if !filepath.IsAbs(dashHome) {
			dashHome = filepath.Join(getUserHome(), dashHome)
		}

		// Initialize Workspace
		if !filepath.IsAbs(workspace) {
			workspace = filepath.Join(dashHome, workspace)
		}
	})

	// Verify Git
	cobra.OnInitialize(func() {
		if err := tools.CheckGitInstall(); err != nil {
			fatal(GIT_NOT_INSTALLED, err.Error())
		}
	})

	// Dash Home
	defaultDashHome := filepath.Join(getUserHome(), ".dashfiles")
	RootCmd.PersistentFlags().StringVar(&dashHome, "home", defaultDashHome, "dashfiles home folder")
	viper.BindPFlag("home", RootCmd.PersistentFlags().Lookup("home"))

	// Config File
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default \""+filepath.Join(defaultDashHome, "config.yaml")+"\")")

	// Workspace
	RootCmd.PersistentFlags().StringVar(&workspace, "workspace", "workspace", "the path to your workspace relative to your dashfiles home folder")
	viper.BindPFlag("workspace", RootCmd.PersistentFlags().Lookup("workspace"))

}

func getUserHome() string {
	home, err := homedir.Dir()
	if err != nil {
		fatal(UNKNOWN_ERROR, err.Error())
	}
	return home
}
