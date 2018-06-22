package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

const (
	// Version of changelog
	Version = "0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = "changelog"
	app.Usage = "Generate changelog of gitea repository"
	app.Version = Version
	app.Commands = []cli.Command{
		cmdGenerate,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(4, "Failed to run app with %s: %v", os.Args, err)
	}
}

var cmdGenerate = cli.Command{
	Name:        "generate",
	Usage:       "generate changelog of gitea repository",
	Description: `generate changelog of gitea repository`,
	Action:      runGenerate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "milestone, m",
			Usage: "Generate which tag from",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Specify a config file",
		},
	},
}

type Config struct {
	Repo   string `yml:"repo"`
	Groups []struct {
		Name    string   `yml:"name"`
		Labels  []string `yml:"labels"`
		Default bool     `yml:"default"`
	}
}

var (
	defaultConfig = []byte(`repo: go-gitea/gitea
groups:
  - 
    name: BREAKING
    labels:
      - kind/breaking
  - 
    name: FEATURE
    labels:
      - kind/feature
  -
    name: BUGFIXES
    labels:
      - kind/bug
  - 
    name: ENHANCEMENT
    labels:
      - kind/enhancement
      - kind/refactor
  -
    name: SECURITY
    labels:
      - kind/security
  - 
    name: TESTING
    labels:
      - kind/testing
  - 
    name: TRANSLATION
    labels:
      - kind/translation
  - 
    name: BUILD
    labels:
      - kind/build
      - kind/lint
  - 
    name: DOCS
    labels:
    - kind/docs
  - 
    name: MISC
    default: true`)
)

func runGenerate(cmd *cli.Context) {
	milestone := cmd.String("milestone")
	if milestone == "" {
		fmt.Println("Please specify a milestone")
		return
	}

	var err error
	var configContent []byte
	if cmd.String("config") == "" {
		configContent = defaultConfig
	} else {
		configContent, err = ioutil.ReadFile(cmd.String("config"))
		if err != nil {
			fmt.Printf("Load config from file %s failed: %v\n", cmd.String("config"), err)
			return
		}
	}

	var config Config
	err = yaml.Unmarshal(configContent, &config)
	if err != nil {
		fmt.Printf("Unmarshal config content failed: %v\n", err)
		return
	}

	client := github.NewClient(nil)
	ctx := context.Background()

	var labels = make(map[string]string)
	var changelogs = make(map[string][]github.Issue)
	var defaultGroup string
	for _, g := range config.Groups {
		changelogs[g.Name] = []github.Issue{}
		for _, l := range g.Labels {
			labels[l] = g.Name
		}
		if g.Default {
			defaultGroup = g.Name
		}
	}

	if defaultGroup == "" {
		defaultGroup = config.Groups[len(config.Groups)-1].Name
	}

	var query = fmt.Sprintf(`repo:%s is:merged milestone:"%s"`, config.Repo, milestone)
	var p = 1
	var perPage = 100
	for {
		result, _, err := client.Search.Issues(ctx, query, &github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    p,
				PerPage: perPage,
			},
		})
		p++
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, pr := range result.Issues {
			var found bool
			for _, lb := range pr.Labels {
				if g, ok := labels[lb.GetName()]; ok {
					changelogs[g] = append(changelogs[g], pr)
					found = true
					break
				}
			}
			if !found {
				changelogs[defaultGroup] = append(changelogs[defaultGroup], pr)
			}
		}

		if len(result.Issues) != perPage {
			break
		}
	}

	fmt.Printf("## [%s](https://github.com/%s/releases/tag/v%s) - %s\n", milestone, config.Repo, milestone, time.Now().Format("2006-01-02"))
	for _, g := range config.Groups {
		if len(changelogs[g.Name]) == 0 {
			continue
		}

		fmt.Println("* " + g.Name)
		for _, pr := range changelogs[g.Name] {
			fmt.Printf("  * %s (#%d)\n", *pr.Title, *pr.Number)
		}
	}
}
