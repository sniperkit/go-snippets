package hunter

import (
	"log"
	"net/url"
	"strings"
)

type GitHubContent struct {
	Username   string `"json:username"`
	Repository string `"json:repository"`
	Branch     string `"json:string"`
	Path       string `"json:string"`
}

func ParseGitHubURL(urlstr string) GitHubContent {
	var ghc GitHubContent

	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatalf("Error parsing: ", err)
	}

	//look into SplitN
	path := strings.Split(u.Path, "/")

	if len(path) > 3 {
		if !strings.Contains(u.Host, "raw.githubusercontent.com") {
			path = append(path[:3], path[4:]...)
		}
	}

	switch sl := len(path); {
	case sl > 4:
		ghc.Path = strings.Join(path[4:], "/")
		fallthrough
	case sl > 3:
		ghc.Branch = path[3]
		fallthrough
	case sl > 2:
		ghc.Repository = path[2]
		fallthrough
	case sl > 1:
		ghc.Username = path[1]
		fallthrough
	default:
	}

	return ghc

}
