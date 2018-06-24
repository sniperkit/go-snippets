package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/codeskyblue/go-sh"
	"github.com/git-hook/git-ledger"
	"github.com/urfave/cli"
)

// Get remote of the first remote listed by `git remote` for project.
func getRemote(project string) string {
	session := sh.NewSession()
	session.SetDir(project)
	remote, err := session.Command("git", "remote").Command("head", "-n", "1").Output()
	// pipe makes this code unreachable
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(remote))
}

// Get the slug of the first remote listed by `git remote` for project.
func getSlug(project string) (string, error) {
	if !ledger.IsGitProject(project) {
		return "", errors.New(fmt.Sprintf("Cannot add %s: not a git project", project))
	}
	session := sh.NewSession()
	session.SetDir(project)
	out, err := session.Command("git", "remote", "show", "-n", getRemote(project)).Command("grep", "Fetch").Output()
	if err != nil {
		return "", err
	}
	output := strings.TrimSpace(string(out))
	reg := regexp.MustCompile(`[^/:]*/[^/:]*$`)
	res := reg.FindStringSubmatch(output)
	return res[0], nil
}

// Add specified path to the ledger.
func add(c *cli.Context) error {
	project := c.Args().First()

	slug, err := getSlug(project)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err
	}
	absolutePath, _ := filepath.Abs(project)

	record := ledger.Record{Path: absolutePath, Slug: slug}
	record.RemoveFromLedger()
	record.AddToLedger()

	return nil
}
