package github

import (
	"fmt"
	"github.com/google/go-github/github"
	"strings"
	"sync"
	"context"
)

var ReposNames []string

type Repo struct {
	owner string
	repo  string
	sha   string
	isIE  bool
}

func (repo Repo) String() string {
	return fmt.Sprintf("Repo{owner:%s, repo:%s, isIE:%v, sha:%s}", repo.owner, repo.repo, repo.isIE, repo.sha)
}

func createReposFromNames(names []string, xapOnly bool) []*Repo {
	fmt.Printf("creating repos from names %s\n", names)
	ret := make([]*Repo, 0)
	for _, name := range names {
		if strings.TrimSpace(name) != "" && strings.Contains(name, "/") {
			components := strings.Split(name, "/")
			isIE := components[0] == "InsightEdge"
			if !isIE || !xapOnly {
				ret = append(ret, &Repo{owner: components[0], repo: components[1], isIE: isIE})
			}
		}
	}
	fmt.Printf("repos are %s\nret is %s\n", names, ret)
	return ret
}

func fillSha(repos []*Repo, from string, client *github.Client) chan []*Repo {
	var wg sync.WaitGroup
	ret := make([]*Repo, 0)
	retChan := make(chan []*Repo, 1)
	for _, repo := range repos {
		wg.Add(1)
		go func(repo *Repo) {
			defer wg.Done()
			//fmt.Printf("working on repo %#v\n", repo)
			ref, _, err := client.Git.GetRef(context.Background(), repo.owner, repo.repo, "refs/heads/" + from)
			if err != nil {
				fmt.Printf("error: %#v\n", err)
			}else {
				//fmt.Printf("Ref: %s\nURL: %s\n", *ref.Ref, *ref.URL)
				//fmt.Printf("Type: %s\nURL: %s\nSHA: %s\n", *ref.Object.Type, *ref.Object.URL, *ref.Object.SHA)
				repo.sha = *ref.Object.SHA
				ret = append(ret, repo)
			}
		}(repo)
	}
	go func() {
		wg.Wait()
		retChan <- ret
	}()
	return retChan
}

func createBranchesWithProgress(repos []*Repo, branch string, client *github.Client) (progressChan chan RepoStatus, resChan chan map[string]interface{}) {
	progressChan = make(chan RepoStatus, len(repos))
	resChan = make(chan map[string]interface{})
	intermediateProgressChan := make(chan RepoStatus, len(repos))
	for _, repo := range repos {
		go func(repo *Repo) {
			_, _, err := client.Git.CreateRef(context.Background(), repo.owner, repo.repo, &github.Reference{
				Ref: github.String("refs/heads/" + branch),
				Object: &github.GitObject{
					SHA: &repo.sha,
				},
			})
			if err != nil {
				fmt.Printf("error createing branch %s in repo: %#v, error is %s:\n", branch, repo, err.Error())
				intermediateProgressChan <- RepoStatus{Name: repo.owner + "/" + repo.repo, Success: false}
			} else {
				fmt.Printf("branch %s created in repo %#v\n", branch, repo)
				intermediateProgressChan <- RepoStatus{Name: repo.owner + "/" + repo.repo, Success: true}
			}
		}(repo)
	}
	go func() {
		var statuses map[string]interface{} = make(map[string]interface{})
		for range repos {
			status := <-intermediateProgressChan
			statuses[status.Name] = status.Success
			progressChan <- status
		}
		resChan <- statuses
		close(progressChan)
		close(resChan)
		close(intermediateProgressChan)
	}()
	return progressChan, resChan
}

func deleteBranchesWithProgress(repos []*Repo, branch string, client *github.Client) (progressChan chan RepoStatus, resChan chan map[string]interface{}) {
	progressChan = make(chan RepoStatus, len(repos))
	resChan = make(chan map[string]interface{})
	intermediateProgressChan := make(chan RepoStatus, len(repos))
	for _, repo := range repos {
		go func(repo *Repo) {
			_, err := client.Git.DeleteRef(context.Background(), repo.owner, repo.repo, "refs/heads/" + branch)
			if err != nil {
				fmt.Printf("error deleting branch %s in repo: %#v, error is %s:\n", branch, repo, err.Error())
				intermediateProgressChan <- RepoStatus{Name: repo.repo, Success: false}
			} else {
				//fmt.Printf("branch %s deleted from repo %#v\n", branch, repo)
				intermediateProgressChan <- RepoStatus{Name: repo.repo, Success: true}
			}
		}(repo)
	}
	go func() {
		var statuses map[string]interface{} = make(map[string]interface{})
		for range repos {
			status := <-intermediateProgressChan
			statuses[status.Name] = status.Success
			progressChan <- status
		}
		resChan <- statuses
		close(progressChan)
		close(resChan)
		close(intermediateProgressChan)
	}()
	return progressChan, resChan
}

type BranchOnRepositories struct {
	Name     string          `json:"name"`
	Statuses map[string]bool `json:"statuses"`
}

func ListAllRefsAsMap(client *github.Client) chan map[string]map[string]string {
	type branch struct {
		repo   string
		branch string
	}
	branches := make(chan branch, 100)
	repos := createReposFromNames(ReposNames, false)
	fmt.Printf("ReposNames are: %s\n", ReposNames)
	var wg sync.WaitGroup
	for _, repo := range repos {
		fmt.Printf("ListAllRefsAsMap working on repo: %s\n", *repo)
		wg.Add(1)
		go func(repo *Repo) {
			defer func() {
				wg.Done()
			}()
			heads, err := ListRefs(client, repo.owner, repo.repo)
			if err == nil {
				for _, head := range heads {
					branchName := strings.TrimPrefix(*head.Ref, "refs/heads/")
					branches <- branch{repo.owner + "/" + repo.repo, branchName}
				}
			} else {
				fmt.Printf("List all ref error: %#v, on repo %#v\n", err, repo)
			}
		}(repo)
	}

	// closer
	go func() {
		wg.Wait()
		close(branches)
	}()

	var res map[string]map[string]string = make(map[string]map[string]string)
	resChan := make(chan map[string]map[string]string)
	go func() {
		for branch := range branches {
			if b, ok := res[branch.branch]; ok {
				b[branch.repo] = ""
			} else {
				res[branch.branch] = map[string]string{branch.repo: ""}
			}
		}
		resChan <- res
		close(resChan)
	}()
	return resChan
}

func DeleteBranchWithProgress(branch string, client *github.Client) (progressChan chan RepoStatus, resChan chan map[string]interface{}) {
	repos := createReposFromNames(ReposNames, false)
	return deleteBranchesWithProgress(repos, branch, client)
}

func CreateBranchsWithProgress(branch string, from string, xapOnly bool, client *github.Client) (progressChan chan RepoStatus, resChan chan map[string]interface{}) {
	repos := createReposFromNames(ReposNames, xapOnly)
	ch :=  fillSha(repos, from, client)
	var existingRepos []*Repo = <- ch
	return createBranchesWithProgress(existingRepos, branch, client)
}

func ListRefs(client *github.Client, owner string, repo string) ([]*github.Reference, error) {
	refs, _, err := client.Git.ListRefs(context.Background(), owner, repo, &github.ReferenceListOptions{Type: "heads/"})
	//refs, _, err := client.Git.ListRefs(owner, repo, &github.ReferenceListOptions{Type:"heads"});
	return refs, err
}

type RefList struct {
}

type UIBranch struct {
	Name  string       `json:"name"`
	Repos []RepoStatus `json:"repos"`
}

type RepoStatus struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
}

type repoWithRefs struct {
	name string
	refs []github.Reference
}
