package main

import (
	"flag"
	"io/ioutil"
	"strings"

	"github.com/keelerh/omniscience/cmd/ingester/common"
	"github.com/keelerh/omniscience/cmd/ingester/github/lib"
	"github.com/pkg/errors"
)

var (
	fGithubAPITokenFilePath = flag.String(
		"github_api_token_file_path",
		"",
		"The path to the Github API token.")
	fGithubOrganization = flag.String(
		"github_organization",
		"",
		"The Github organization from which to index repo READMEs.")
)

func createGithubFetcher() (common.DocumentFetcher, error) {
	rawToken, err := ioutil.ReadFile(*fGithubAPITokenFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Github API token")
	}
	token := strings.TrimSpace(string(rawToken))
	return lib.NewGithub(token, *fGithubOrganization), nil
}

func main() {
	common.CreateIngesterCLI(createGithubFetcher)
}
