package main

import (
	"flag"
	"io/ioutil"
	"strings"

	"github.com/keelerh/omniscience/cmd/ingester/common"
	"github.com/keelerh/omniscience/cmd/ingester/confluence/lib"
	"github.com/pkg/errors"
)

var (
	fConfluenceHostname = flag.String(
		"confluence_domain",
		"",
		"A Confluence domain, e.g. your-domain.atlassian.net.")
	fConfluenceUsername = flag.String(
		"confluence_username",
		"",
		"The username associated with the generated Confluence API token.")
	fConfluenceAPITokenFilePath = flag.String(
		"confluence_api_token_file_path",
		"",
		"A Confluence API token; this can be generated at https://id.atlassian.com.")
)

func createConfluenceFetcher() (common.DocumentFetcher, error) {
	rawToken, err := ioutil.ReadFile(*fConfluenceAPITokenFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read Confluence API token")
	}

	token := strings.TrimSpace(string(rawToken))
	basicAuth := lib.BasicAuth(*fConfluenceUsername, token)
	confluenceSvc, err := lib.NewConfluence(*fConfluenceHostname, basicAuth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise Confluence service")
	}

	return confluenceSvc, nil
}

func main() {
	common.CreateIngesterCLI(createConfluenceFetcher)
}
