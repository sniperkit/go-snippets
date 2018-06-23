package main

import (
	"flag"
	"fmt"
	"github.com/mgutz/ansi"

//	"github.com/gophergala/bron/filters"
)

var (
	blessedPtr   string
	dashboardPtr string
	forcePtr     bool
	quietPtr     bool
	repoPtr      string
	repoPathPtr  string
	statusPtr    bool
	verbosePtr   bool
	vizPtr       bool
)

func colorize(msg string) string {

	lime  := ansi.ColorCode("green+h:black")
	reset := ansi.ColorCode("reset")

	return(lime + msg + reset)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	flag.StringVar(&blessedPtr, "blessedPath", "/go/src/github.com/yaronn/blessed-contrib", "Path where blessed-contrib is installed")
	flag.StringVar(&dashboardPtr, "dashboard", "example", "Name of dashboard to use for visualization")
	flag.BoolVar(&quietPtr, "q", false, "Supresses all output, except status")
	flag.StringVar(&repoPtr, "repo", "github.com/gophergala/bron", "Git repository to scan")
	flag.StringVar(&repoPathPtr, "path", "", "Git repository file path to scan (not currently implemented)")
	flag.BoolVar(&statusPtr, "s", true, "Outputs status of completion of analysis")
	flag.BoolVar(&verbosePtr, "v", false, "verbosity level")
	flag.BoolVar(&vizPtr, "viz", false, "Visualize the results, requires blessed")
	flag.BoolVar(&forcePtr, "f", false, "Force update the data")

	flag.Parse()

	if repoPtr == "" && repoPathPtr == "" {
		fmt.Println("please specify either a repo or a path to a git repo to scan")
	} else if repoPtr != "" && repoPathPtr != "" {
		fmt.Println("please specify only either a repo or a path to a git repo to scan, not both")
	} else {
		if !quietPtr {
			if verbosePtr {
				fmt.Println("going to scan repository", repoPtr, "...")
			}
		}
	}

	if repoPtr != "" {
		if !forcePtr {
			if checkData(repoPtr, dashboardPtr, blessedPtr) {
				if !quietPtr {
					if verbosePtr {
						fmt.Println("found existing data for", repoPtr, "using that ...")
					}
					if vizPtr {
						showDashboard()
					}
				}
			} else {
				clonePath := "https://"+repoPtr+".git"
				uuidRepo := cloneRepo(clonePath)
				updateDashboardData(uuidRepo, repoPtr, dashboardPtr, verbosePtr, vizPtr, quietPtr, statusPtr)
			}
		} else {
			clonePath := "https://"+repoPtr+".git"
			uuidRepo := cloneRepo(clonePath)
			updateDashboardData(uuidRepo, repoPtr, dashboardPtr, verbosePtr, vizPtr, quietPtr, statusPtr)
		}
	}

}
