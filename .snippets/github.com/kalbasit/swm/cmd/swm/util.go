package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/google/go-github/github"
	"github.com/kalbasit/swm/code"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.i3wm.org/i3"
	"golang.org/x/oauth2"

	cli "gopkg.in/urfave/cli.v2"
)

func createLogger(ctx *cli.Context) error {
	// create the logger that pretty prints to the ctx.Writer
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: ctx.App.Writer}).
		With().
		Timestamp().
		Str("ignore-pattern", ctx.String("ignore-pattern")).
		Str("code-path", ctx.String("code-path")).
		Str("profile", ctx.String("profile")).
		Str("story", ctx.String("story")).
		Logger().
		Level(zerolog.InfoLevel)
	// handle debug
	if ctx.Bool("debug") {
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
	}

	return nil
}

func sortCommands(subCmds []*cli.Command) {
	for _, subCmd := range subCmds {
		sort.Sort(cli.FlagsByName(subCmd.Flags))
		sort.Sort(cli.CommandsByName(subCmd.Subcommands))
		if len(subCmd.Subcommands) > 0 {
			sortCommands(subCmd.Subcommands)
		}
	}
}

func createGithubClient(ctx *cli.Context) error {
	githubAccessToken := ctx.String("github.access_token")
	if githubAccessToken == "" {
		var err error
		githubTokenPath := path.Join(os.Getenv("HOME"), ".github_token")
		if _, err = os.Stat(githubTokenPath); err == nil {
			var con []byte
			con, err = ioutil.ReadFile(githubTokenPath)
			if err != nil {
				return errors.Wrap(err, "error reading the Github token from ~/.github_token")
			}
			githubAccessToken = string(bytes.TrimSpace(con))
		}
	}
	if githubAccessToken == "" {
		return errors.New("no Github token were provided")
	}
	log.Debug().Str("github.access_token", githubAccessToken).Msg("creating the Github client")
	code.GithubClient = github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)))

	return nil
}

func getDefaultProfile() string {
	var p string

	// try parsing it from the ACTIVE_PROFILE environment variable.
	p = os.Getenv("ACTIVE_PROFILE")
	// try parsing it from the TMUX environment variable (the session path).
	if p == "" {
		tmuxSocketPath := os.Getenv("TMUX")
		if tmuxSocketPath != "" {
			profileStory := strings.Split(path.Base(tmuxSocketPath), ",")[0]
			profileStoryArr := strings.Split(profileStory, "@")
			if len(profileStoryArr) == 2 {
				p = profileStoryArr[0]
			}
		}
	}
	// finally try parsing it from the i3 workspace
	if p == "" {
		i3Workspace, err := getActiveI3WorkspaceName()
		if err == nil && strings.Contains(i3Workspace, "@") {
			p = strings.Split(i3Workspace, "@")[0]
		}
	}

	return p
}

func getDefaultStory() string {
	var s string

	// try parsing it from the ACTIVE_STORY environment variable.
	s = os.Getenv("ACTIVE_STORY")
	// try parsing it from the TMUX environment variable (the session path).
	if s == "" {
		tmuxSocketPath := os.Getenv("TMUX")
		if tmuxSocketPath != "" {
			profileStory := strings.Split(path.Base(tmuxSocketPath), ",")[0]
			profileStoryArr := strings.Split(profileStory, "@")
			if len(profileStoryArr) == 2 {
				s = profileStoryArr[1]
			}
		}
	}
	// finally try parsing it from the i3 workspace
	if s == "" {
		i3Workspace, err := getActiveI3WorkspaceName()
		if err == nil && strings.Contains(i3Workspace, "@") {
			s = strings.Split(i3Workspace, "@")[1]
		}
	}

	return s
}

func getActiveI3WorkspaceName() (string, error) {
	// get the workspaces
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		return "", err
	}
	for _, workspace := range workspaces {
		if workspace.Focused {
			return workspace.Name, nil
		}
	}
	return "", errors.New("no active i3 workspace was found")
}

func newCoder(ctx *cli.Context) (code.Coder, error) {
	// parse the regex
	ignorePattern, err := regexp.Compile(ctx.String("ignore-pattern"))
	if err != nil {
		return nil, err
	}
	// create a new coder
	return code.New(ctx.String("code-path"), ignorePattern), nil
}
