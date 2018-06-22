package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/kr/fs"
)

func getAuthors(repoPath string) []string {

	// XXX stub
	authors := []string{}

	return authors

}

func getDiff(repoPath string, currentCommit string, previousCommit string) []byte {

	diffCmd := exec.Command("git", "--git-dir="+repoPath+"/.git", "diff", currentCommit, previousCommit)
	diffOut, diffErr := diffCmd.Output()
	check(diffErr)

	return diffOut

}

func getCommits(repoPath string) ([]string, map[string]map[string]string) {

	commits := map[string]map[string]string{}
	order := []string{}

	revCmd := exec.Command("git", "--git-dir="+repoPath+"/.git", "rev-list", "--all", "--no-merges", "--pretty=format:\"%H|%an|%at\"")
	revOut, revErr := revCmd.Output()
	check(revErr)
	lines := strings.Split(string(revOut), "\n")
	prevCommit := ""
	for i, val := range lines {
		if i % 2 != 0 {
			// XXX error check to ensure there are exactly 3 splits
			components := strings.Split(val[1:len(val)-1], "|")
			if _, ok := commits[prevCommit]; ok {
				commits[prevCommit]["prevCommit"] = components[0]
			}
			order = append(order, components[0])
			kvs := map[string]string{}
			kvs["prevCommit"] = components[0]
			kvs["author"] = components[1]
			kvs["timestamp"] = components[2]
			commits[components[0]] = kvs
			prevCommit = components[0]
		}
	}

	return order, commits

}

func getFiles(repoPath string) []string {

	files := []string{}

	walker := fs.Walk(repoPath)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if !strings.Contains(walker.Path(), ".git") && strings.Contains(walker.Path(), ".") {
			name := walker.Path()
			f, err := os.Open(name)
			if err == nil {
				defer f.Close()
				fi, err := f.Stat()
				if err == nil {
					switch mode := fi.Mode(); {
					case mode.IsDir():
						continue
					case mode.IsRegular():
						files = append(files, walker.Path())
					}
				}
			}
		}
	}

	return files

}

func getFileContents(filename string) []byte {

	data, err := ioutil.ReadFile(filename)
	check(err)

	return data

}

func checkoutCommit(repoPath string, commit string) {

	cwd, wdErr := os.Getwd()
	check(wdErr)
	chErr := os.Chdir(repoPath)
	check(chErr)
	checkoutCmd := exec.Command("git", "checkout", commit)
	_, checkoutErr := checkoutCmd.Output()
	check(checkoutErr)
	chErr = os.Chdir(cwd)
	check(chErr)

}

func cloneRepo(repo string) string {

	uuidRepo := uuid.New()

	// XXX don't forget to cleanup after we're finished
	err := os.Mkdir("/tmp/"+uuidRepo, 0644)
	check(err)

	if repo != "" {
		cloneCmd := exec.Command("git", "clone", repo, "/tmp/"+uuidRepo)
		cloneOut, cloneErr := cloneCmd.Output()
		check(cloneErr)
		fmt.Println(string(cloneOut))
	}

	return "/tmp/"+uuidRepo

}
