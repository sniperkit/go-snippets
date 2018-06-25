package tools

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Git struct {
	Dir string
}

func CheckGitInstall() error {
	// Check if git is installed
	if _, err := exec.Command("git", "--version").Output(); err != nil {
		return errors.New("It seems like you do not have git installed and present in your PATH. Git is required for dashfiles to work")
	}
	return nil
}

func (git *Git) Version() (out string) {
	return git.gitPanic("--version")
}

func (git *Git) Clone(url string, targetFolder string) error {
	_, err := git.git("clone", parseGitUrl(url), targetFolder)
	return err
}

func (git *Git) Pull() error {
	_, err := git.git("pull")
	return err
}

func (git *Git) gitPanic(args ...string) string {
	out, err := git.git(args...)
	if err != nil {
		panic(err)
	}
	return out
}

func (git *Git) git(args ...string) (out string, err error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = git.Dir
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	bytes, err := cmd.Output()
	out = strings.TrimSpace(string(bytes))
	return
}

func parseGitUrl(url string) string {
	// Add shorthand for ssh based git repositories
	if match, _ := regexp.MatchString("^[\\w-]+/[\\w-]+$", url); match {
		return "git@github.com:" + url
	}
	return url
}
