package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const SEPARATOR = string(os.PathSeparator)

func getGitArgs(githubName string, projectName string) []string {

	dir, _ := os.Getwd()

	target := strings.Join([]string{dir, projectName}, SEPARATOR)
	// https://codeload.github.com/divanvisagie/postl/zip/master

	arguments := []string{
		"clone",
		"--depth=1",
		fmt.Sprintf("https://github.com/%s.git", githubName),
		target,
	}

	return arguments
}

func runCommand(command string, arguments []string) string {
	commandOutput, err := exec.Command(command, arguments...).Output()
	if err != nil {
		return err.Error()
	}
	return string(commandOutput)
}

func removeGitInDirectory(directoryName string) {
	path := string(strings.Join([]string{directoryName, ".git"}, SEPARATOR))

	dir, _ := os.Getwd()
	target := strings.Join([]string{dir, path}, SEPARATOR)
	err := os.RemoveAll(target)

	if err != nil {
		log.Fatalln("Failed to delete .git directory")
	}
}

func main() {

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalln("You need to pass in parameters")
	}

	if len(args) == 1 {
		log.Fatalln("you need to provide a seed url")
	}

	projectName := args[0]
	githubName := args[1]

	commandArgs := getGitArgs(githubName, projectName)

	commandOutput := runCommand("git", commandArgs)

	fmt.Println(commandOutput)

	removeGitInDirectory(projectName)
}
