// print gathers all tools for formatting output
package cmd

import (
	"time"

	sp "github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var noColor bool
var printPrompt = color.New(color.FgWhite).PrintfFunc()
var spinner = sp.New(sp.CharSets[24], 100*time.Millisecond)

func SetNoColor() {
	color.NoColor = noColor
}

func PrintSuccess(msg string, params ...interface{}) {
	color.Green(msg, params...)
}

func PrintInfo(msg string, params ...interface{}) {
	color.White(msg, params...)
}

func PrintWarning(msg string, params ...interface{}) {
	color.Yellow(msg, params...)
}

func PrintErr(err error, params ...interface{}) {
	color.Red(err.Error(), params...)
}

func PrintNotYetFinished(cmd *cobra.Command) {
	color.Yellow("%s command is not yet implemented", cmd.Name())
}

// func PrintTree(ds *dataset.Dataset, indent int) {
//  fmt.Println(strings.Repeat(" ", indent), ds.Address.String())
//  for i, d := range ds.Datasets {
//    if i < len(ds.Datasets)-1 {
//      fmt.Println(strings.Repeat(" ", indent), "├──", d.Address.String())
//    } else {
//      fmt.Println(strings.Repeat(" ", indent), "└──", d.Address.String())

//    }
//  }
// }
