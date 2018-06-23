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

	"encoding/json"
	"github.com/ChappIO/dashfiles/tools"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your dashfiles workspace",
	Long: `Clear an already existing workspace and set up a
clean workspace that is ready for configuration.

By default this workspace is created in $HOME/.dashfiles/workspace
but you can override this by setting the 'workspace' variable in
your configuration.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Cleaning up workspace")

		cleanDir(workspace, "workspace")

		if len(args) > 0 {
			fmt.Printf("Cloning %s\n", args[0])

			git := tools.Git{}
			if err := git.Clone(args[0], workspace); err != nil {
				fatal(WORKSPACE_CLONE_FAILED, "Failed to clone workspace: "+err.Error())
			}
		}

		// Create the default dashfile
		dashFile := DefaultDashfile()
		dashFileJson := filepath.Join(workspace, DASHFILE_JSON)
		if !fileExists(dashFileJson) {
			if data, err := json.MarshalIndent(dashFile, "", "   "); err == nil {
				ioutil.WriteFile(dashFileJson, data, 0644)
			} else {
				fatal(INIT_FAILED, "Failed to create default "+DASHFILE_JSON)
			}
		} else {
			dashFile.readFromWorkspace(workspace)
		}

		fmt.Println("Your workspace has been initialized at " + workspace)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

func cleanDir(path string, name string) {
	if err := os.RemoveAll(path); err != nil {
		log.Fatalf("Failed to clean up %s: %s", name, err.Error())
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Failed to create %s: %s", name, err.Error())
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
