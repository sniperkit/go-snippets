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
	"github.com/ChappIO/dashfiles/tools"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your workspace to the latest version",
	Long: `Update your workspace to the latest version. Running
this command is essentially the same as running 'git pull' from
the workspace.''`,
	Run: func(cmd *cobra.Command, args []string) {
		git := tools.Git{
			Dir: workspace,
		}

		fmt.Println("Pulling changes for workspace...")
		if err := git.Pull(); err != nil {
			fatal(UPDATE_FAILED, "Failed to pull changes: "+err.Error())
		}

		fmt.Println("Workspace updated")
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
