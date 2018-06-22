package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	dotGitignore string = ".gitignore"
)

type Reloader struct {
	// The path of the project directory. Defaults to"./" current working directory
	ProjectDir string
	//
	RunCmd string
	//
	ReloadCmd string
	//
	Pid int
}

func NewReloader() *Reloader {
	return &Reloader{
		ProjectDir: "./",
		RunCmd:     "echo alo",
		ReloadCmd:  "",
		Pid:        -1,
	}
}

func (self *Reloader) Bump() {
	fmt.Println("Bumping", self)

	if self.Pid > 0 {
		process, err := os.FindProcess(self.Pid)
		if err != nil {
			fmt.Println("Process already died")
		} else {
			fmt.Println("Killing", process.Pid)
			process.Kill()
		}
	}

	command := exec.Command(strings.Fields(self.RunCmd)[0], strings.Fields(self.RunCmd)[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Start()
	if err != nil {
		fmt.Println(err)
	}
	self.Pid = command.Process.Pid
}

func (self *Reloader) Run() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	bump := make(chan bool, 2)
	filter := NewFileFilter(self.ProjectDir)
	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev.Name)
				if !filter.Ignore(ev.Name) {
					bump <- true
				} else {
					log.Println("Ignoring")
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			case <-bump:
				go self.Bump()
			}
		}
	}()
	bump <- true

	err = watcher.Watch(self.ProjectDir)
	if err != nil {
		fmt.Println(err)
	}

	<-done

	watcher.Close()
}

// utils section
func loadGitIgnoreFileEx(path string) ([]string, error) {
	inFile, err := os.Open(path + dotGitignore)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	var ext []string

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || string(line[0]) == "#" {
			// ignoring empty lines and comments
			continue
		}
		if string(line[0]) == "*" {
			line = line[1:]
		}
		ext = append(ext, line)
	}
	return ext, nil
}

func createRegexpPatterns(expressins []string) []*regexp.Regexp {
	var patterns []*regexp.Regexp
	for _, ex := range expressins {
		r, err := regexp.Compile(regexp.QuoteMeta(ex))
		if err != nil {
			continue
		}
		patterns = append(patterns, r)
	}
	return patterns
}

type FileFilter struct {
	patterns []*regexp.Regexp
}

func (self *FileFilter) Ignore(fp string) bool {
	for _, rgx := range self.patterns {
		if rgx.MatchString(fp) {
			// fmt.Println("Matched ", rgx, fp)
			return true
		}
	}
	return false
}

func NewFileFilter(path string) *FileFilter {
	exs, _ := loadGitIgnoreFileEx(path)
	return &FileFilter{createRegexpPatterns(exs)}
}

func main() {
	app := cli.NewApp()
	app.Name = "gowatch"
	app.EnableBashCompletion = true
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:      "path",
			ShortName: "p",
			Usage:     "project path",
			Action: func(c *cli.Context) {
				println("completed task: ", c.Args().First())
			},
		},
	}
	app.Run(os.Args)
	reloader := NewReloader()
	reloader.Run()
}
