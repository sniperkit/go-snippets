package ipfs

import (
	"fmt"
	"github.com/datatogether/task_mgmt/tasks"
	"os"
	"path/filepath"
)

// TaskGithubAdd uses code hosted in a repository on github.com
// to coordinate adding a specified url
// TODO - Work in Progress.
type TaskGithubAdd struct {
	// url to where the code to execute lives
	// example: https://github.com/ipfs/ipfs-wiki/mirror
	RepoUrl string `json:"repoUrl"`
	// version control repoCommit to execute code from
	RepoCommit string `json:"repoCommit"`
	// url this code is to run against
	SourceUrl string `json:"sourceUrl"`
	// url of IPFS api server
	IpfsApiServerUrl string `json:"ipfsApiServerUrl"`
	// checksum of source resource
	// SourceChecksum string `json:"sourceChecksum"`
	// url of output
	// ResultUrl string `json:"resultUrl"`
	// multihash of output
	// ResultHash string `json:"resultHash"`
	// any message associated with this task (failure, info, etc.)
	// Message string `json:"message"`
}

func NewTaskGithubAdd() tasks.Taskable {
	return &TaskGithubAdd{}
}

func (t TaskGithubAdd) Valid() error {
	return fmt.Errorf("github add task is not yet finished")
}

func (t TaskGithubAdd) Do(updates chan tasks.Progress) {
	// this task should:
	// 1. Clone the repo from RepoUrl
	// 2. Download the zim file
	// 3. Exectute the code at the repo against
	// 4. Post the finalized output to an IPFS Node
	p := tasks.Progress{Percent: 0.0, Step: 1, Steps: 4, Status: "Cloning Repo"}
	updates <- p

	if err := cloneGitRepo(t.RepoUrl, filepath.Join(os.TempDir(), "repo")); err != nil {
		p.Error = fmt.Errorf("error cloning repo: %s", err.Error())
		updates <- p
		return
	}

	p.Steps++
	p.Percent = 0.25

}

func cloneGitRepo(url, filepath string) error {
	return fmt.Errorf("not yet finished")
}
