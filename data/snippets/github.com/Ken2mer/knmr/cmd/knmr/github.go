package main

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/Ken2mer/knmr/logger"
	"github.com/google/go-github/github"
	"github.com/tcnksm/go-gitconfig"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

var (
	username = "Ken2mer"
	reponame = "knmr"
)

var githubCommand = cli.Command{
	Name:   "github",
	Usage:  "github",
	Action: githubCmd,
}

func githubCmd(ctx *cli.Context) error {
	gh := newGithubClient()
	fns := []func(){
		gh.user,
		gh.follows,
		gh.code,
		gh.events,
		gh.subscription,
	}
	var wg sync.WaitGroup
	for _, fn := range fns {
		wg.Add(1)
		go func(f func()) {
			f()
			wg.Done()
		}(fn)
	}
	wg.Wait()
	return ghServe()
}

type ghClient struct {
	ctx    context.Context
	client *github.Client
}

func newGithubClient() ghClient {
	context := context.Background()
	return ghClient{
		ctx:    context,
		client: oauth2Client(context),
	}
}

func (c *ghClient) subscription() {
	subscription, err := getRepositorySubscription(c.ctx, c.client)
	if err != nil {
		return
	}
	logger.DumpJSON(subscription)
}

func (c *ghClient) events() {
	events, err := getEvents(c.ctx, c.client)
	if err != nil {
		return
	}
	dumpEvents(events)
}

func (c *ghClient) code() {
	code, err := getCodeSearchResult(c.ctx, c.client)
	if err != nil {
		return
	}
	dumpCodeSearchResult(code)
}

func (c *ghClient) follows() {
	follows, err := getFollowingUsers(c.ctx, c.client)
	if err != nil {
		return
	}
	dumpFollowUsers(follows)
}

func (c *ghClient) user() {
	user, err := getUser(c.ctx, c.client)
	if err != nil {
		return
	}
	dumpUser(user)
}

// cf. https://github.com/mackerelio/mkr/blob/af4f89ae6fac2290b9fe642de37f84a25de67d62/plugin/github.go

// Get github client having github token.
func oauth2Client(ctx context.Context) *github.Client {
	var oauthClient *http.Client
	if token := getGithubToken(); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: getGithubToken()},
		)
		oauthClient = oauth2.NewClient(ctx, ts)
	}
	return github.NewClient(oauthClient)
}

// Get github token from environment variables, or github.token in gitconfig file
func getGithubToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}
	token, _ = gitconfig.GithubToken()
	return token
}
