package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"

	"gopkg.in/yaml.v2"
)

var makeLabels = map[string]string{
	"bug":          "e11d21",
	"enhancement":  "009800",
	"question":     "0052cc",
	"API":          "bfe5bf",
	"architecture": "fef2c0",
	"chore":        "fbca04",
	"doc":          "fbca04",
	"feature":      "009800",
	"optimization": "009800",
	"framework":    "fef2c0",
	"production":   "d4c5f9",
	"refactoring":  "fbca04",
	"security":     "e11d21",
	"staging":      "d4c5f9",
	"TBD":          "ffffff",
	"test":         "fbca04",
	"URGENT":       "f7c6c7",
	"User Story":   "bfe5bf",
	"WIP":          "ffffff",
	"★★★":          "5319e7",
	"★★☆":          "5319e7",
	"★☆☆":          "5319e7",
	"◷0d1h":        "207de5",
	"◷0d3h":        "207de5",
	"◷1d":          "207de5",
	"◷2d":          "207de5",
	"◷4d":          "207de5",
}

var destructive = flag.Bool("d", false, "delete unknown labels")

func main() {
	owner, repo := parseArgs()
	user, token := loadHubCfg()
	ctx := context.Background()

	client := github.NewClient((&github.BasicAuthTransport{
		Username: user,
		Password: token,
	}).Client())

	labels, _, err := client.Issues.ListLabels(ctx, owner, repo, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, l := range labels {
		if makeLabels[*l.Name] == "" {
			if !*destructive {
				log.Printf("ignore %q", l)
			} else {
				log.Printf("delete %q", l)
				_, err := client.Issues.DeleteLabel(ctx, owner, repo, *l.Name)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else if makeLabels[*l.Name] != *l.Color {
			log.Printf("colour %q (%s to %s)", l, *l.Color, makeLabels[*l.Name])
			*l.Color = makeLabels[*l.Name]
			_, _, err := client.Issues.EditLabel(ctx, owner, repo, *l.Name, l)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("accept %q", l)
		}
		delete(makeLabels, *l.Name)
	}
	for name, color := range makeLabels {
		log.Printf("create %q (%s)", name, color)
		l := &github.Label{Name: &name, Color: &color}
		_, _, err := client.Issues.CreateLabel(ctx, owner, repo, l)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func parseArgs() (owner, repo string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] owner/repo\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 || strings.Count(flag.Arg(0), "/") != 1 {
		flag.Usage()
		os.Exit(2)
	}
	ownerRepo := strings.Split(flag.Arg(0), "/")
	return ownerRepo[0], ownerRepo[1]
}

func loadHubCfg() (user, token string) {
	var hubCfg = make(map[string][]struct {
		Protocol   string
		User       string
		OAuthToken string `yaml:"oauth_token"`
	})

	hubCfgData, err := ioutil.ReadFile(os.Getenv("HOME") + "/.config/hub")
	if os.IsNotExist(err) {
		log.Fatal("hub tool is not installed or configured")
	}
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal([]byte(hubCfgData), &hubCfg)
	if err != nil {
		log.Fatal(err)
	}
	if len(hubCfg["github.com"]) == 0 {
		log.Fatal("failed to detect configuration of hub tool")
	}

	return hubCfg["github.com"][0].User, hubCfg["github.com"][0].OAuthToken
}
