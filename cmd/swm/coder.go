package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/kalbasit/swm/code"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	cli "gopkg.in/urfave/cli.v2"
)

var coderCmd = &cli.Command{
	Name: "coder",
	Subcommands: []*cli.Command{
		// add project
		{
			Name:      "add-project",
			Usage:     "Clone a new project and places it in the selected profile and story",
			Action:    coderAddProject,
			ArgsUsage: "<url>",
		},
		// pull request
		{
			Name:    "pull-request",
			Usage:   "Pull request sub-command provides commands to interact with Github",
			Aliases: []string{"pr"},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "github.access_token", Usage: "The access token for accessing Github", EnvVars: []string{"GITHUB_ACCESS_TOKEN"}},
			},
			Before: createGithubClient,
			Subcommands: []*cli.Command{
				// list
				{
					Name:    "list",
					Usage:   "List the pull requests open for this repository over on Github",
					Aliases: []string{"ls"},
					Action:  coderPullRequestList,
				},
			},
		},
	},
}

func coderAddProject(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		log.Debug().Msgf("expecting one argument, the URL to clone. Got %d arguments", ctx.Args())
		return errors.New("expecting one argument as url, required")
	}
	// create a new coder
	c, err := newCoder(ctx)
	if err != nil {
		return err
	}
	if err = c.Scan(); err != nil {
		return err
	}
	// get the profile
	profile, err := c.Profile(ctx.String("profile"))
	if err != nil {
		log.Debug().Str("profile", ctx.String("profile")).Msg("profile not found")
		return err
	}
	// get the story
	story := profile.Story(ctx.String("story"))
	// add this project
	return story.AddProject(ctx.Args().First())
}

func coderPullRequestList(ctx *cli.Context) error {
	// create a new coder
	c, err := newCoder(ctx)
	if err != nil {
		return err
	}
	if err = c.Scan(); err != nil {
		return err
	}
	// get the project from the current PATH
	var prj code.Project
	var wd string
	wd, err = os.Getwd()
	if err != nil {
		return errors.Wrap(err, "error finding the current working directory")
	}
	prj, err = c.ProjectByAbsolutePath(wd)
	if err != nil {
		return errors.Wrap(err, "error finding the project for the current directory")
	}
	// get the list of prs
	var prs []*github.PullRequest
	prs, err = prj.ListPullRequests()
	if err != nil {
		return errors.Wrap(err, "error getting the list of the pull requests")
	}
	if len(prs) == 0 {
		fmt.Println("No pull requests found for the project.")
		return nil
	}
	// prepare the table writer and write down the PRs
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{"Number", "Title", "URL", "Created at"})
	for _, pr := range prs {
		table.Append([]string{strconv.Itoa(pr.GetNumber()), pr.GetTitle(), pr.GetHTMLURL(), pr.GetCreatedAt().String()})
	}
	table.Render()

	return nil
}
