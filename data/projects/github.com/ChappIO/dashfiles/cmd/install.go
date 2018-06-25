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
	"github.com/ChappIO/dashfiles/managers"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var updatePackages bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the current workspace to your system",
	Long: `This will (irreversibly) install the current state of your
workspace to your system.`,
	Run: func(cmd *cobra.Command, args []string) {

		dashFile := DefaultDashfile()
		dashFile.readFromWorkspace(workspace)

		// Install dotfiles
		dotfileDir := filepath.Join(workspace, dashFile.Dotfiles)

		if !fileExists(dotfileDir) {
			fmt.Println("The dotfile directory (" + dotfileDir + ") does not exists. Skipping dotfiles.")
		} else {
			outputDir := getUserHome()

			filepath.Walk(dotfileDir, func(path string, info os.FileInfo, err error) error {
				relativePath, _ := filepath.Rel(dotfileDir, path)
				targetPath := filepath.Join(outputDir, relativePath)

				if info.IsDir() {
					err = os.MkdirAll(targetPath, info.Mode())
				} else {
					var data []byte
					data, err = ioutil.ReadFile(path)

					if err == nil {
						fmt.Println("Installing ~/" + relativePath)
						err = ioutil.WriteFile(targetPath, data, info.Mode())
					}
				}

				return err
			})
		}

		// Run custom scripts
		runCustomScripts("dashfiles-pre-install")

		// Install packages
		for _, packageManager := range managers.GetManagers() {
			managerFile := filepath.Join(workspace, packageManager.GetName()+"."+packageManager.GetFileFormat())

			if !fileExists(managerFile) {
				continue
			}

			if !packageManager.IsInstalled() {
				fmt.Println("Found " + managerFile + " but " + packageManager.GetName() + " is not installed")
			}

			fmt.Println("Installing packages using " + packageManager.GetName())
			if err := packageManager.InstallFromFile(managerFile); err != nil {
				fatal(INSTALL_PACKAGE_FAILED, "Failed to install packages: "+err.Error())
			}

			if updatePackages {
				fmt.Println("Updating all " + packageManager.GetName() + " packages")
				packageManager.Update()
			}
		}

		// Run custom scripts
		runCustomScripts("dashfiles-post-install")
	},
}

func runCustomScripts(name string) {
	files, _ := filepath.Glob(filepath.Join(workspace, name+".*"))
	for _, file := range files {
		if info, _ := os.Stat(file); info.Mode()&0111 != 0 {
			command := exec.Command(file, getUserHome())
			command.Stderr = os.Stderr
			command.Stdout = os.Stdout
			command.Stdin = os.Stdin
			fmt.Println("Executing '" + file + "'")
			if err := command.Run(); err != nil {
				fatal(EXTERNAL_SCRIPT_FAILED, "Could not run '"+file+"': "+err.Error()+"\nMake sure you have added the correct shebang")
			}
		} else {

			fmt.Println("Skipping execution of '" + file + "' because it is not executable")
		}

	}
}

func init() {
	RootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVar(&updatePackages, "update", false, "update all installed packages")
}
