package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	OptVerbose bool
	OptDir     string
	Args       []string
)

// Number of concurrent clones at the same time.
var NumConcurrentClones = 12

// A mutex to protect output to the console so we don't get gobligook when multiple goroutines write at once
var OutMux = sync.Mutex{}

// Print help usage info
func Usage() {
	fmt.Println("multigit - run git commands against many git repositories at once")
	os.Exit(0)
}

func main() {
	flag.BoolVar(&OptVerbose, "verbose", false, "be verbose")
	flag.StringVar(&OptDir, "dir", "", "directory that contains your git repositories")
	flag.Usage = Usage
	flag.Parse()

	Args = flag.Args()
	if len(Args) == 0 {
		Usage()
		os.Exit(0)
	}

	if OptDir != "" {
		err := os.Chdir(OptDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Special-case for git clone
	if Args[0] == "clone" {
		Clone()
		return
	}

	// List all directories and determine which ones are a git directory
	basedir, err := os.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	subdirs, err := basedir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	gitdirs := make([]string, 0)
	maxnamelen := 0
	for _, fi := range subdirs {
		if fi.IsDir() {
			_, err := os.Stat(fi.Name() + "/.git")
			if err == nil {
				gitdirs = append(gitdirs, fi.Name())
				if len(fi.Name()) > maxnamelen {
					maxnamelen = len(fi.Name())
				}
			}
		}
	}

	// Run the commands
	var wg sync.WaitGroup
	results := make([]string, len(gitdirs))
	errmap := make(map[string]error)
	numerror := 0
	for i, dir := range gitdirs {
		wg.Add(1)
		go func(i int, dir string) {
			cmd := exec.Command("git", Args...)
			cmd.Dir = dir
			output, err := cmd.CombinedOutput()
			if err != nil {
				numerror++
			}
			errmap[dir] = err
			results[i] = string(output)
			wg.Done()
		}(i, dir)
	}
	wg.Wait()

	// Report the results
	for i, dir := range gitdirs {
		fmt.Print(dir + " ")
		if errmap[dir] != nil {
			fmt.Println("ERROR")
		}
		fmt.Println("")

		out := strings.Trim(results[i], "\n")
		if out != "" {
			outlines := strings.Split(out, "\n")
			for _, line := range outlines {
				fmt.Println("    " + line)
			}
		} else {
			if errmap[dir] != nil {
				fmt.Println("    " + errmap[dir].Error())
			}
		}
		fmt.Println("")
	}

	// Exit with correct code
	if numerror == len(gitdirs) {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Clone() {
	var directory, rawuri string
	var rest []string
	if len(Args) == 2 {
		rawuri = Args[1]
		rest = Args[:1]
	} else {
		if strings.Contains(Args[len(Args)-2], "--") {
			rawuri = Args[len(Args)-1]
			rest = Args[:len(Args)-1]
		} else {
			directory = Args[len(Args)-1]
			rawuri = Args[len(Args)-2]
			rest = Args[:len(Args)-2]
		}
	}

	// If the rawuri does not contain a *, then do a simple git clone
	if !strings.Contains(rawuri, "*") {
		cmd := exec.Command("git", Args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else { // The string contains a *, it's show time!
		// Parse the rawuri
		uri, err := url.Parse(rawuri)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if uri.Host == "" {
			// retry by munging an ssh address into a real URL
			uri, err = url.Parse("git://" + rawuri)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			hostparts := strings.Split(uri.Host, ":")
			if len(hostparts) > 1 {
				uri.Host = hostparts[0]
				uri.Path = hostparts[1] + uri.Path
			}
		}
		if uri.Host != "github.com" {
			fmt.Println("Sorry I don't know how to multi-clone on anything but github. Please open a feature-request if you would like to add support for another provider.")
			os.Exit(1)
		}
		if directory != "" {
			fmt.Println("Invalid directory for multi-clone")
			os.Exit(1)
		}

		pathParts := strings.Split(strings.Trim(uri.Path, "/"), "/")
		repos := GitHubRepos(pathParts[0], pathParts[1])

		if len(repos) == 0 {
			fmt.Println("No repositories found. Private repository cloning requires GITHUB_API_TOKEN")
		} else {
			fmt.Println("Cloning:")
			for _, repo := range repos {
				fmt.Println("    " + repo)
			}
			CloneRepositories("git@github.com:"+pathParts[0]+"/", repos, rest)
		}
	}
}

func GitHubRepos(user, repopattern string) []string {
	var client *github.Client
	token := os.Getenv("GITHUB_API_TOKEN")
	if token != "" {
		t := &oauth.Transport{
			Token: &oauth.Token{AccessToken: token},
		}
		client = github.NewClient(t.Client())
	} else {
		client = github.NewClient(nil)
	}

	// First get user repos
	allrepos := make([]github.Repository, 0)

	pnum := 1
	for {
		useropt := &github.RepositoryListOptions{
			ListOptions: github.ListOptions{PerPage: 100, Page: pnum},
		}
		userrepos, resp, err := client.Repositories.List(user, useropt)
		allrepos = append(allrepos, userrepos...)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if resp.NextPage != 0 {
			pnum = resp.NextPage
		} else {
			break
		}
	}

	// Next get organizational repos -- we don't care about errors
	pnum = 1
	for {
		orgopt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100, Page: pnum},
		}
		orgrepos, resp, _ := client.Repositories.ListByOrg(user, orgopt)
		allrepos = append(allrepos, orgrepos...)
		if resp.NextPage != 0 {
			pnum = resp.NextPage
		} else {
			break
		}
	}

	// Our list of repos
	repos := make([]string, 0)

	// Return everything if there is no repo name
	if repopattern == "" || repopattern == "*" {
		for _, repo := range allrepos {
			repos = append(repos, *repo.Name)
		}
		return repos
	}

	// For each repo, check if it's in the list of allowed repos as per the pattern
	for _, repo := range allrepos {
		matched, err := filepath.Match(repopattern, *repo.Name)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if matched {
			repos = append(repos, *repo.Name)
		}
	}
	return repos
}

func CloneRepositories(base string, repos []string, args []string) {
	// Run the commands
	ch := make(chan string)
	wg := sync.WaitGroup{}
	var numerror int64 = 0
	for n := 0; n < NumConcurrentClones; n++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case repo, ok := <-ch:
					if !ok {
						wg.Done()
						return
					}
					cmd := exec.Command("git", append(args, base+repo)...)
					output, err := cmd.CombinedOutput()
					if err != nil {
						atomic.AddInt64(&numerror, 1)
					}

					// Report output
					OutMux.Lock()
					fmt.Print(repo + " ")
					if err != nil {
						fmt.Print("ERROR")
					}
					fmt.Println("")
					outstr := strings.Trim(string(output), "\n")
					if outstr != "" {
						outlines := strings.Split(outstr, "\n")
						for _, line := range outlines {
							fmt.Println("    " + line)
						}
					} else {
						fmt.Println("    " + err.Error())
					}
					fmt.Println("")
					OutMux.Unlock()
				default:
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}

	// Feed values into channels
	for _, repo := range repos {
		ch <- repo
	}
	close(ch)
	wg.Wait()

	// Exit with correct code
	if numerror == int64(len(repos)) {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
