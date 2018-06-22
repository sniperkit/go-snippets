package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "hot"
	app.Author = "Naoto Kaneko"
	app.Email = "naoty.k@gmail.com"
	app.Version = "0.1.0"
	app.Usage = "List files in order of the commited count"
	app.Action = hot
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "number, n",
			Usage: "the number of listed files",
		},
	}
	app.Run(os.Args)
}

func hot(context *cli.Context) {
	var files Files

	filenames, _ := gitLsFiles()

	if len(context.Args()) > 0 {
		pattern := context.Args()[0]
		filenames = Match(pattern, filenames)
	}

	for _, filename := range filenames {
		commitCount := getCommitTimes(filename)
		files = append(files, File{Name: filename, CommitCount: commitCount})
	}

	limit := context.Int("number")
	if limit == 0 {
		limit = len(files)
	}

	sort.Sort(files)
	for i, file := range files {
		if i >= limit {
			break
		}
		fmt.Printf("%d: %s\n", file.CommitCount, file.Name)
	}
}

func gitLsFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files")
	out, err := cmd.CombinedOutput()
	str := string(out)
	str = strings.Trim(str, "\n")
	return strings.Split(str, "\n"), err
}

func getCommitTimes(filename string) int {
	cmd := exec.Command("git", "log", "--oneline", "--", filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0
	}

	str := string(out)
	str = strings.Trim(str, "\n")
	return len(strings.Split(str, "\n"))
}
